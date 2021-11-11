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
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/dsl"
	"mlr/internal/pkg/output"
	"mlr/internal/pkg/parsing/lexer"
	"mlr/internal/pkg/parsing/parser"
	"mlr/internal/pkg/runtime"
	"mlr/internal/pkg/types"
)

// NewEmptyRoot sets up an empty CST, before ingesting any DSL strings.  For
// mlr put and mlr filter, CSTs are constructed, then ASTs are ingested from
// DSL strings. For the REPL, a CST is constructed once, then an AST is
// ingested on every line of input from the REPL.
func NewEmptyRoot(
	recordWriterOptions *cli.TWriterOptions,
	dslInstanceType DSLInstanceType, // mlr put, mlr filter, or mlr repl
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

// ASTBuildVisitorFunc is a callback, used by RootNode's Build method, which
// CST-builder callsites can use to visit parse-to-AST result of multi-string
// DSL inputs. Nominal use: mlr put -v, mlr put -d, etc.
type ASTBuildVisitorFunc func(dslString string, astNode *dsl.AST)

// Used by DSL -> AST -> CST callsites including mlr put, mlr filter, and mlr
// repl. The RootNode must be separately instantiated (e.g. NewEmptyRoot())
// since the CST is partially reset on every line of input from the REPL
// prompt.
func (root *RootNode) Build(
	dslStrings []string,
	dslInstanceType DSLInstanceType,
	isReplImmediate bool,
	doWarnings bool,
	warningsAreFatal bool,
	astBuildVisitorFunc ASTBuildVisitorFunc,
) error {

	for _, dslString := range dslStrings {
		astRootNode, err := buildASTFromStringWithMessage(dslString)
		if err != nil {
			// Error message already printed out
			return err
		}

		// E.g. mlr put -v -- let it print out what it needs to.
		if astBuildVisitorFunc != nil {
			astBuildVisitorFunc(dslString, astRootNode)
		}

		err = root.IngestAST(
			astRootNode,
			isReplImmediate,
			doWarnings,
			warningsAreFatal,
		)
		if err != nil {
			return err
		}
	}

	err := root.Resolve()
	if err != nil {
		return err
	}

	return err
}

func buildASTFromStringWithMessage(dslString string) (*dsl.AST, error) {
	astRootNode, err := buildASTFromString(dslString)
	if err != nil {
		// Leave this out until we get better control over the error-messaging.
		// At present it's overly parser-internal, and confusing. :(
		fmt.Fprintln(os.Stderr, "mlr: cannot parse DSL expression.")
		return nil, err
	} else {
		return astRootNode, nil
	}
}

func buildASTFromString(dslString string) (*dsl.AST, error) {
	// For non-Windows, already stripped by the shell; helpful here for Windows.
	if strings.HasPrefix(dslString, "'") && strings.HasSuffix(dslString, "'") {
		dslString = dslString[1 : len(dslString)-1]
	}

	theLexer := lexer.NewLexer([]byte(dslString))
	theParser := parser.NewParser()
	interfaceAST, err := theParser.Parse(theLexer)
	if err != nil {
		return nil, err
	}
	astRootNode := interfaceAST.(*dsl.AST)
	return astRootNode, nil
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

	// For debug:
	// fmt.Println("PRE")
	// ast.Print()
	root.regexProtectPrePass(ast)
	// fmt.Println("POST")
	// ast.Print()

	err = root.buildMainPass(ast, isReplImmediate)
	if err != nil {
		return err
	}

	return nil
}

// Resolve is called after IngestAST has been called one or more times.
// See comments above IngestAST.
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
// regexProtectPrePass rewrites string-literal nodes in regex position (e.g.
// second arg to gsub) to have regex node type. This is so we can have "\t" be
// a tab character for string literals generally, but remain backslash-t for
// regex literals.
//
// Callsites to have regexes protected:
// * sub/gsub second argument;
// * regextract/regextract_or_else second argument;
// * =~ and !=~ right-hand side -- since these are infix operators, this means
//   (in the AST point of view) second argument.
//
// Sample ASTs:
//
// $ mlr -n put -v '$y =~ "\t"'
// AST:
// * statement block
//     * bare boolean
//         * operator "=~"
//             * direct field value "y"
//             * string literal "	"
//
// $ mlr -n put -v '$y = sub($x, "\t", "TAB")'
// AST:
// * statement block
//     * assignment "="
//         * direct field value "y"
//         * function callsite "sub"
//             * direct field value "x"
//             * string literal "	"
//             * string literal "TAB"

func (root *RootNode) regexProtectPrePass(ast *dsl.AST) {
	root.regexProtectPrePassAux(ast.RootNode)
}

func (root *RootNode) regexProtectPrePassAux(astNode *dsl.ASTNode) {

	if astNode.Children == nil || len(astNode.Children) == 0 {
		return
	}

	isCallsiteOfInterest := false
	if astNode.Type == dsl.NodeTypeOperator {
		if astNode.Token != nil {
			nodeName := string(astNode.Token.Lit)
			if nodeName == "=~" || nodeName == "!=~" {
				isCallsiteOfInterest = true
			}
		}
	} else if astNode.Type == dsl.NodeTypeFunctionCallsite {
		if astNode.Token != nil {
			nodeName := string(astNode.Token.Lit)
			if nodeName == "sub" || nodeName == "gsub" || nodeName == "regextract" || nodeName == "regextract_or_else" {
				isCallsiteOfInterest = true
			}
		}
	}

	for i, astChild := range astNode.Children {
		if isCallsiteOfInterest && i == 1 {
			if astChild.Type == dsl.NodeTypeStringLiteral {
				astChild.Type = dsl.NodeTypeRegex
			}
		}
		root.regexProtectPrePassAux(astChild)
	}

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

		if astChild.Type == dsl.NodeTypeNamedFunctionDefinition {
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
			// Unresolvable at CST-build time but perhaps a local variable. For example,
			// the UDF callsite '$z = f($x, $y)', and supposing
			// there will be 'f = func(a, b) { return a*b }' in scope at runtime.
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
				"mlr: subroutine name not found: " + subroutineName,
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

func (root *RootNode) ExecuteBeginBlocks(state *runtime.State) error {
	for _, beginBlock := range root.beginBlocks {
		_, err := beginBlock.Execute(state)
		if err != nil {
			return err
		}
	}
	return nil
}

func (root *RootNode) ExecuteMainBlock(state *runtime.State) (outrec *types.Mlrmap, err error) {
	_, err = root.mainBlock.Execute(state)
	return state.Inrec, err
}

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
