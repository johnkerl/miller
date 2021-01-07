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
type DumpStatementNode struct {
	// TODO: redirect options
	ostream     *os.File
	expressions []IEvaluable
	// xxx redirect
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
	lib.InternalCodingErrorIf(len(astNode.Children) != 2)
	expressionsNode := astNode.Children[0]

	expressions := make([]IEvaluable, len(expressionsNode.Children))
	for i, childNode := range expressionsNode.Children {
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
			if !evaluation.IsAbsent() {
				fmt.Fprintln(this.ostream, evaluation.String())
			}
		}
	}

	return nil, nil
}
