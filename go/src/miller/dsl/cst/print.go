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
	ostream     *os.File
	expressions []IEvaluable
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
	expressions := make([]IEvaluable, len(astNode.Children))
	for i, childNode := range astNode.Children {
		expression, err := this.BuildEvaluableNode(childNode)
		if err != nil {
			return nil, err
		}
		expressions[i] = expression
	}

	return &DumpStatementNode{
		ostream,
		expressions,
	}, nil
}

// ----------------------------------------------------------------
func (this *DumpStatementNode) Execute(state *State) (*BlockExitPayload, error) {
	if len(this.expressions) == 0 { // 'dump' without argument means 'dump @*'
		// Not Fprintln since JSON output is LF-terminated already
		fmt.Fprint(this.ostream, state.Oosvars.String())
	} else {
		for _, expression := range this.expressions {
			evaluation := expression.Evaluate(state)
			fmt.Fprintln(this.ostream, evaluation.String())
		}
	}

	return nil, nil
}

// ================================================================
type PrintStatementNode struct {
	// TODO: redirect options
	ostream     *os.File
	terminator  string
	expressions []IEvaluable
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
	expressions := make([]IEvaluable, len(astNode.Children))
	for i, childNode := range astNode.Children {
		expression, err := this.BuildEvaluableNode(childNode)
		if err != nil {
			return nil, err
		}
		expressions[i] = expression
	}

	return &PrintStatementNode{
		ostream,
		terminator,
		expressions,
	}, nil
}

// ----------------------------------------------------------------
func (this *PrintStatementNode) Execute(state *State) (*BlockExitPayload, error) {
	if len(this.expressions) == 0 {
		fmt.Fprintf(this.ostream, this.terminator)
	} else {
		for i, expression := range this.expressions {
			if i > 0 {
				fmt.Fprint(this.ostream, " ")
			}
			evaluation := expression.Evaluate(state)
			fmt.Fprint(this.ostream, evaluation.String())
		}
		fmt.Fprintf(this.ostream, this.terminator)
	}
	return nil, nil
}
