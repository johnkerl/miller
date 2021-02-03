// ================================================================
// Top-level entry point for building a CST from an AST at parse time, and for
// executing the CST at runtime.
// ================================================================

package cst

import (
	"container/list"
	"errors"
	"fmt"
	"os"

	"miller/cliutil"
	"miller/dsl"
	"miller/lib"
	"miller/output"
	"miller/runtime"
	"miller/types"
)

// ----------------------------------------------------------------
func NewEmptyRoot(
	recordWriterOptions *cliutil.TWriterOptions,
) *RootNode {
	return &RootNode{
		beginBlocks:                   make([]*StatementBlockNode, 0),
		mainBlock:                     NewStatementBlockNode(),
		endBlocks:                     make([]*StatementBlockNode, 0),
		udfManager:                    NewUDFManager(),
		udsManager:                    NewUDSManager(),
		allowUDFSRedefinitions:        false,
		unresolvedFunctionCallsites:   list.New(),
		unresolvedSubroutineCallsites: list.New(),
		outputHandlerManagers:         list.New(),
		recordWriterOptions:           recordWriterOptions,
	}
}

// Nominally for mlr put/filter we want to flag overwritten UDFs/UDSs as an
// error.  But in the REPL, which is interactive, people should be able to
// redefine.  This method allows the latter use-case.
func (this *RootNode) WithRedefinableUDFS() *RootNode {
	this.allowUDFSRedefinitions = true
	return this
}

// ----------------------------------------------------------------
// If the user has multiple put -f / put -e pieces, we can AST-parse each
// separately and build them. However we cannot resolve UDF/UDS references
// until after they're all ingested -- e.g. first piece calls a function which
// the second defines, or mutual recursion across pieces, etc.
func (this *RootNode) IngestAST(
	ast *dsl.AST,
	isFilter bool, // false for 'mlr put', true for 'mlr filter'
) error {
	if ast.RootNode == nil {
		return errors.New("Cannot build CST from nil AST root")
	}

	// Check for things that are syntax errors but not done in the AST for
	// pragmatic reasons. For example, $anything in begin/end blocks;
	// begin/end/func not at top level; etc.
	err := ValidateAST(ast, isFilter)
	if err != nil {
		return err
	}

	err = this.buildMainPass(ast)
	if err != nil {
		return err
	}

	return nil
}

func (this *RootNode) Resolve() error {

	err := this.resolveFunctionCallsites()
	if err != nil {
		return err
	}

	err = this.resolveSubroutineCallsites()
	if err != nil {
		return err
	}

	return nil
}

// ----------------------------------------------------------------
// This builds the CST almost entirely. The only afterwork is that user-defined
// functions may be called before they are defined, so a follow-up pass will
// need to resolve those callsites.

func (this *RootNode) buildMainPass(ast *dsl.AST) error {

	if ast.RootNode.Type != dsl.NodeTypeStatementBlock {
		return errors.New(
			"CST root build: non-statement-block AST root node unhandled",
		)
	}
	astChildren := ast.RootNode.Children

	// Example AST:
	//
	// $ mlr put -v 'begin{@a=1;@b=2} $x=3; $y=4' myfile.dkvp
	// DSL EXPRESSION:
	// begin{@a=1;@b=2} $x=3; $y=4
	// RAW AST:
	// * StatementBlock
	//     * BeginBlock
	//         * StatementBlock
	//             * Assignment "="
	//                 * DirectOosvarValue "a"
	//                 * IntLiteral "1"
	//             * Assignment "="
	//                 * DirectOosvarValue "b"
	//                 * IntLiteral "2"
	//     * Assignment "="
	//         * DirectFieldValue "x"
	//         * IntLiteral "3"
	//     * Assignment "="
	//         * DirectFieldValue "y"
	//         * IntLiteral "4"

	for _, astChild := range astChildren {

		if astChild.Type == dsl.NodeTypeFunctionDefinition {
			err := this.BuildAndInstallUDF(astChild)
			if err != nil {
				return err
			}

		} else if astChild.Type == dsl.NodeTypeSubroutineDefinition {
			err := this.BuildAndInstallUDS(astChild)
			if err != nil {
				return err
			}

		} else if astChild.Type == dsl.NodeTypeBeginBlock || astChild.Type == dsl.NodeTypeEndBlock {
			statementBlockNode, err := this.BuildStatementBlockNodeFromBeginOrEnd(astChild)
			if err != nil {
				return err
			}

			if astChild.Type == dsl.NodeTypeBeginBlock {
				this.beginBlocks = append(this.beginBlocks, statementBlockNode)
			} else {
				this.endBlocks = append(this.endBlocks, statementBlockNode)
			}
		} else {
			statementNode, err := this.BuildStatementNode(astChild)
			if err != nil {
				return err
			}
			this.mainBlock.AppendStatementNode(statementNode)
		}
	}

	return nil
}

// This is invoked within the buildMainPass call tree whenever a function is
// called before it's defined.
func (this *RootNode) rememberUnresolvedFunctionCallsite(udfCallsite *UDFCallsite) {
	this.unresolvedFunctionCallsites.PushBack(udfCallsite)
}

func (this *RootNode) rememberUnresolvedSubroutineCallsite(udsCallsite *UDSCallsite) {
	this.unresolvedSubroutineCallsites.PushBack(udsCallsite)
}

// After-pass after buildMainPass returns, in case a function was called before
// it was defined. It may be the case that:
//
// * A user-defined function was called before it was defined, and was actually defined.
// * A user-defined function was called before it was defined, and it was not actually defined.
// * The user misspelled the name of a built-in function.
//
// So, our error message should reflect all those options.

func (this *RootNode) resolveFunctionCallsites() error {
	for this.unresolvedFunctionCallsites.Len() > 0 {
		unresolvedFunctionCallsite := this.unresolvedFunctionCallsites.Remove(
			this.unresolvedFunctionCallsites.Front(),
		).(*UDFCallsite)

		functionName := unresolvedFunctionCallsite.udf.signature.funcOrSubrName
		callsiteArity := unresolvedFunctionCallsite.udf.signature.arity

		udf, err := this.udfManager.LookUp(functionName, callsiteArity)
		if err != nil {
			return err
		}
		if udf == nil {
			return errors.New(
				"Miller: function name not found: " + functionName,
			)
		}

		unresolvedFunctionCallsite.udf = udf
	}
	return nil
}

func (this *RootNode) resolveSubroutineCallsites() error {
	for this.unresolvedSubroutineCallsites.Len() > 0 {
		unresolvedSubroutineCallsite := this.unresolvedSubroutineCallsites.Remove(
			this.unresolvedSubroutineCallsites.Front(),
		).(*UDSCallsite)

		subroutineName := unresolvedSubroutineCallsite.uds.signature.funcOrSubrName
		callsiteArity := unresolvedSubroutineCallsite.uds.signature.arity

		uds, err := this.udsManager.LookUp(subroutineName, callsiteArity)
		if err != nil {
			return err
		}
		if uds == nil {
			return errors.New(
				"Miller: subroutine name not found: " + subroutineName,
			)
		}

		unresolvedSubroutineCallsite.uds = uds
	}
	return nil
}

// ----------------------------------------------------------------
// Various 'tee > $hostname . ".dat", $*' statements will have
// OutputHandlerManager instances to track file-descriptors for all unique
// values of $hostname in the input stream.
//
// At CST-build time, the builders are expected to call this so we can put
// OutputHandlerManager instances on a list. Then, at end of stream, we
// can close all the descriptors, flush the record-output streams, etc.

func (this *RootNode) RegisterOutputHandlerManager(
	outputHandlerManager output.OutputHandlerManager,
) {
	this.outputHandlerManagers.PushBack(outputHandlerManager)
}

func (this *RootNode) ProcessEndOfStream() {
	for entry := this.outputHandlerManagers.Front(); entry != nil; entry = entry.Next() {
		outputHandlerManager := entry.Value.(output.OutputHandlerManager)
		errs := outputHandlerManager.Close()
		if len(errs) != 0 {
			for _, err := range errs {
				fmt.Fprintf(
					os.Stderr,
					"%s: error on end-of-stream close: %v\n",
					lib.MlrExeName(),
					err,
				)
			}
			os.Exit(1)
		}
	}
}

// ----------------------------------------------------------------
func (this *RootNode) ExecuteBeginBlocks(state *runtime.State) error {
	for _, beginBlock := range this.beginBlocks {
		_, err := beginBlock.Execute(state)
		if err != nil {
			return err
		}
	}
	return nil
}

// ----------------------------------------------------------------
func (this *RootNode) ExecuteMainBlock(state *runtime.State) (outrec *types.Mlrmap, err error) {
	_, err = this.mainBlock.Execute(state)
	return state.Inrec, err
}

// ----------------------------------------------------------------
func (this *RootNode) ExecuteEndBlocks(state *runtime.State) error {
	for _, endBlock := range this.endBlocks {
		_, err := endBlock.Execute(state)
		if err != nil {
			return err
		}
	}
	return nil
}

// ----------------------------------------------------------------
// src/miller/auxents/repl
// TODO: more comments
func (this *RootNode) ExecuteREPLExperimental(state *runtime.State) (outrec *types.Mlrmap, err error) {
	_, err = this.mainBlock.ExecuteFrameless(state)
	return state.Inrec, err
}

func (this *RootNode) ResetForREPL() {
	this.beginBlocks = make([]*StatementBlockNode, 0)
	this.mainBlock = NewStatementBlockNode()
	this.endBlocks = make([]*StatementBlockNode, 0)
	this.unresolvedFunctionCallsites = list.New()
	this.unresolvedSubroutineCallsites = list.New()
}
