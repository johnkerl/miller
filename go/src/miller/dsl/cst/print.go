// ================================================================
// This handles print and dump statements.
// ================================================================

package cst

import (
	"fmt"
	"os"

	"miller/dsl"
	"miller/lib"
)

// ================================================================
type PrintStatementNode struct {
	ostream     *os.File
	terminator  string
	expressions []IEvaluable
	// xxx redirect ...
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

func (this *RootNode) BuildPrintnStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypePrintnStatement)
	return this.BuildPrintxStatementNode(
		astNode,
		os.Stdout,
		"",
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
	lib.InternalCodingErrorIf(len(astNode.Children) != 2)

	expressionsNode := astNode.Children[0]
	_ /*redirectNode*/ = astNode.Children[1]

	expressions := make([]IEvaluable, len(expressionsNode.Children))
	for i, childNode := range expressionsNode.Children {
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
			if !evaluation.IsAbsent() {
				fmt.Fprint(this.ostream, evaluation.String())
			}
		}
		fmt.Fprintf(this.ostream, this.terminator)
	}
	return nil, nil
}
