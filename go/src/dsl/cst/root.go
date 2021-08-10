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

	"mlr/src/cliutil"
	"mlr/src/dsl"
	"mlr/src/output"
	"mlr/src/runtime"
	"mlr/src/types"
)

// ----------------------------------------------------------------
func NewEmptyRoot(
	recordWriterOptions *cliutil.TWriterOptions,
	dslInstanceType DSLInstanceType,
) *RootNode {
	return &RootNode{
		beginBlocks:                   make([]*StatementBlockNode, 0),
		mainBlock:                     NewStatementBlockNode(),
		replImmediateBlock:            NewStatementBlockNode(),
		endBlocks:                     make([]*StatementBlockNode, 0),
		udfManager:                    NewUDFManager(),
		udsManager:                    NewUDSManager(),
		allowUDFUDSRedefinitions:      false,
		unresolvedFunctionCallsites:   list.New(),
		unresolvedSubroutineCallsites: list.New(),
		outputHandlerManagers:         list.New(),
		recordWriterOptions:           recordWriterOptions,
		dslInstanceType:               dslInstanceType,
	}
}

// Nominally for mlr put/filter we want to flag overwritten UDFs/UDSs as an
// error.  But in the REPL, which is interactive, people should be able to
// redefine.  This method allows the latter use-case.
func (root *RootNode) WithRedefinableUDFUDS() *RootNode {
	root.allowUDFUDSRedefinitions = true
	return root
}

// ----------------------------------------------------------------
// If the user has multiple put -f / put -e pieces, we can AST-parse each
// separately and build them. However we cannot resolve UDF/UDS references
// until after they're all ingested -- e.g. first piece calls a function which
// the second defines, or mutual recursion across pieces, etc.
func (root *RootNode) IngestAST(
	ast *dsl.AST,
	// False for non-REPL use. Also false for bulk-load REPL use.  True for
	// interactive REPL statements which are intended to be executed once
	// (immediately) but not retained.
	isReplImmediate bool,
	doWarnings bool,
	warningsAreFatal bool,
) error {
	if ast.RootNode == nil {
		return errors.New("Cannot build CST from nil AST root")
	}

	// Check for things that are syntax errors but not done in the AST for
	// pragmatic reasons. For example, $anything in begin/end blocks;
	// begin/end/func not at top level; etc.
	err := ValidateAST(ast, root.dslInstanceType)
	if err != nil {
		return err
	}

	if doWarnings {
		ok := WarnOnAST(ast)
		if !ok {
			// Messages already printed out
			if warningsAreFatal {
				fmt.Printf(
					"%s: Exiting due to warnings treated as fatal.\n",
					"mlr",
				)
				os.Exit(1)
			}
		}
	}

	err = root.buildMainPass(ast, isReplImmediate)
	if err != nil {
		return err
	}

	return nil
}

func (root *RootNode) Resolve() error {

	err := root.resolveFunctionCallsites()
	if err != nil {
		return err
	}

	err = root.resolveSubroutineCallsites()
	if err != nil {
		return err
	}

	return nil
}

// ----------------------------------------------------------------
// This builds the CST almost entirely. The only afterwork is that user-defined
// functions may be called before they are defined, so a follow-up pass will
// need to resolve those callsites.

func (root *RootNode) buildMainPass(ast *dsl.AST, isReplImmediate bool) error {

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
	// AST:
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
			err := root.BuildAndInstallUDF(astChild)
			if err != nil {
				return err
			}

		} else if astChild.Type == dsl.NodeTypeSubroutineDefinition {
			err := root.BuildAndInstallUDS(astChild)
			if err != nil {
				return err
			}

		} else if astChild.Type == dsl.NodeTypeBeginBlock || astChild.Type == dsl.NodeTypeEndBlock {
			statementBlockNode, err := root.BuildStatementBlockNodeFromBeginOrEnd(astChild)
			if err != nil {
				return err
			}

			if astChild.Type == dsl.NodeTypeBeginBlock {
				root.beginBlocks = append(root.beginBlocks, statementBlockNode)
			} else {
				root.endBlocks = append(root.endBlocks, statementBlockNode)
			}
		} else if isReplImmediate {
			statementNode, err := root.BuildStatementNode(astChild)
			if err != nil {
				return err
			}
			root.replImmediateBlock.AppendStatementNode(statementNode)
		} else {
			statementNode, err := root.BuildStatementNode(astChild)
			if err != nil {
				return err
			}
			root.mainBlock.AppendStatementNode(statementNode)
		}
	}

	return nil
}

// This is invoked within the buildMainPass call tree whenever a function is
// called before it's defined.
func (root *RootNode) rememberUnresolvedFunctionCallsite(udfCallsite *UDFCallsite) {
	root.unresolvedFunctionCallsites.PushBack(udfCallsite)
}

func (root *RootNode) rememberUnresolvedSubroutineCallsite(udsCallsite *UDSCallsite) {
	root.unresolvedSubroutineCallsites.PushBack(udsCallsite)
}

// After-pass after buildMainPass returns, in case a function was called before
// it was defined. It may be the case that:
//
// * A user-defined function was called before it was defined, and was actually defined.
// * A user-defined function was called before it was defined, and it was not actually defined.
// * The user misspelled the name of a built-in function.
//
// So, our error message should reflect all those options.

func (root *RootNode) resolveFunctionCallsites() error {
	for root.unresolvedFunctionCallsites.Len() > 0 {
		unresolvedFunctionCallsite := root.unresolvedFunctionCallsites.Remove(
			root.unresolvedFunctionCallsites.Front(),
		).(*UDFCallsite)

		functionName := unresolvedFunctionCallsite.udf.signature.funcOrSubrName
		callsiteArity := unresolvedFunctionCallsite.udf.signature.arity

		udf, err := root.udfManager.LookUp(functionName, callsiteArity)
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

func (root *RootNode) resolveSubroutineCallsites() error {
	for root.unresolvedSubroutineCallsites.Len() > 0 {
		unresolvedSubroutineCallsite := root.unresolvedSubroutineCallsites.Remove(
			root.unresolvedSubroutineCallsites.Front(),
		).(*UDSCallsite)

		subroutineName := unresolvedSubroutineCallsite.uds.signature.funcOrSubrName
		callsiteArity := unresolvedSubroutineCallsite.uds.signature.arity

		uds, err := root.udsManager.LookUp(subroutineName, callsiteArity)
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

func (root *RootNode) RegisterOutputHandlerManager(
	outputHandlerManager output.OutputHandlerManager,
) {
	root.outputHandlerManagers.PushBack(outputHandlerManager)
}

func (root *RootNode) ProcessEndOfStream() {
	for entry := root.outputHandlerManagers.Front(); entry != nil; entry = entry.Next() {
		outputHandlerManager := entry.Value.(output.OutputHandlerManager)
		errs := outputHandlerManager.Close()
		if len(errs) != 0 {
			for _, err := range errs {
				fmt.Fprintf(
					os.Stderr,
					"%s: error on end-of-stream close: %v\n",
					"mlr",
					err,
				)
			}
			os.Exit(1)
		}
	}
}

// ----------------------------------------------------------------
func (root *RootNode) ExecuteBeginBlocks(state *runtime.State) error {
	for _, beginBlock := range root.beginBlocks {
		_, err := beginBlock.Execute(state)
		if err != nil {
			return err
		}
	}
	return nil
}

// ----------------------------------------------------------------
func (root *RootNode) ExecuteMainBlock(state *runtime.State) (outrec *types.Mlrmap, err error) {
	_, err = root.mainBlock.Execute(state)
	return state.Inrec, err
}

// ----------------------------------------------------------------
func (root *RootNode) ExecuteEndBlocks(state *runtime.State) error {
	for _, endBlock := range root.endBlocks {
		_, err := endBlock.Execute(state)
		if err != nil {
			return err
		}
	}
	return nil
}

// ----------------------------------------------------------------
// These are for the Miller REPL.

// If a DSL string was parsed into an AST and ingested in 'immediate' mode and
// build into the CST, it's not populated into the main-statements block for
// remembered execution on every record. Rather, it's just stored once, to be
// executed once, and then discarded.

// This is the 'execute once' part of that.
func (root *RootNode) ExecuteREPLImmediate(state *runtime.State) (outrec *types.Mlrmap, err error) {
	_, err = root.replImmediateBlock.ExecuteFrameless(state)
	return state.Inrec, err
}

// This is the 'and then discarded' part of that.
func (root *RootNode) ResetForREPL() {
	root.replImmediateBlock = NewStatementBlockNode()
	root.unresolvedFunctionCallsites = list.New()
	root.unresolvedSubroutineCallsites = list.New()
}

// This is for the REPL's context-printer command.
func (root *RootNode) ShowBlockReport() {
	fmt.Printf("#begin %d\n", len(root.beginBlocks))
	fmt.Printf("#main  %d\n", len(root.mainBlock.executables))
	fmt.Printf("#end   %d\n", len(root.endBlocks))
}
