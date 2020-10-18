package cst

import (
	"fmt"
	"os"

	"miller/dsl"
	"miller/lib"
)

// ================================================================
// This handles print and dump statements.
// ================================================================

// ================================================================
type DumpStatementNode struct {
	// TODO: redirect options
	ostream    *os.File
	expression IEvaluable
}

// ----------------------------------------------------------------
func (this *RootNode) BuildDumpStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeDumpStatement)
	return this.BuildDumpxStatementNode(astNode, os.Stdout)
}

func (this *RootNode) BuildEdumpStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeEdumpStatement)
	return this.BuildDumpxStatementNode(astNode, os.Stderr)
}

// Common code for building dump/edump nodes
func (this *RootNode) BuildDumpxStatementNode(
	astNode *dsl.ASTNode,
	ostream *os.File,
) (IExecutable, error) {
	var expression IEvaluable = nil
	var err error = nil

	if len(astNode.Children) == 0 {
		// OK
	} else if len(astNode.Children) == 1 {
		expression, err = this.BuildEvaluableNode(astNode.Children[0])
		if err != nil {
			return nil, err
		}
	} else {
		// Should not have been allowed by the BNF grammar
		lib.InternalCodingErrorIf(true)
	}

	return &DumpStatementNode{
		ostream,
		expression,
	}, nil
}

// ----------------------------------------------------------------
func (this *DumpStatementNode) Execute(state *State) (*BlockExitPayload, error) {
	if this.expression == nil { // 'dump' without argument means 'dump @*'
		// Not Fprintln since JSON output is LF-terminated already
		fmt.Fprint(this.ostream, state.Oosvars.String())
	} else {
		evaluation := this.expression.Evaluate(state)
		fmt.Fprintln(this.ostream, evaluation.String())
	}

	return nil, nil
}

// ================================================================
type PrintStatementNode struct {
	// TODO: redirect options
	ostream    *os.File
	terminator string
	expression IEvaluable
}

// ----------------------------------------------------------------
func (this *RootNode) BuildPrintStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypePrintStatement)
	return this.BuildPrintxStatementNode(
		astNode,
		os.Stdout,
		"\n",
	)
}

func (this *RootNode) BuildEprintStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeEprintStatement)
	return this.BuildPrintxStatementNode(
		astNode,
		os.Stderr,
		"\n",
	)
}

func (this *RootNode) BuildPrintnStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypePrintnStatement)
	return this.BuildPrintxStatementNode(
		astNode,
		os.Stdout,
		"",
	)
}

func (this *RootNode) BuildEprintnStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeEprintnStatement)
	return this.BuildPrintxStatementNode(
		astNode,
		os.Stderr,
		"",
	)
}

// Common code for building print/eprint/printn/eprintn nodes
func (this *RootNode) BuildPrintxStatementNode(
	astNode *dsl.ASTNode,
	ostream *os.File,
	terminator string,
) (IExecutable, error) {
	var expression IEvaluable = nil
	var err error = nil
	if len(astNode.Children) == 0 {
		// OK
	} else if len(astNode.Children) == 1 {
		expression, err = this.BuildEvaluableNode(astNode.Children[0])
		if err != nil {
			return nil, err
		}
	} else {
		// Should not have been allowed by the BNF grammar
		lib.InternalCodingErrorIf(true)
	}

	return &PrintStatementNode{
		ostream,
		terminator,
		expression,
	}, nil
}

// ----------------------------------------------------------------
func (this *PrintStatementNode) Execute(state *State) (*BlockExitPayload, error) {
	if this.expression != nil {
		evaluation := this.expression.Evaluate(state)
		fmt.Fprint(this.ostream, evaluation.String())
	}
	fmt.Fprintf(this.ostream, this.terminator)
	return nil, nil
}
