// ================================================================
// This handles tee statements. This produces new records (in addition to $*)
// into th output record stream.
// ================================================================

package cst

import (
	"fmt"

	"miller/dsl"
	"miller/lib"
)

// ================================================================
// Examples:
//   tee @a
//   tee @a, @b
//
// Each argument must be a non-indexed oosvar/localvar/fieldname, so we can use
// their names as keys in the emitted record.  These restrictions are enforced
// in the CST logic, to keep this parser/AST logic simpler.

type TeeStatementNode struct {
	teeEvaluable IEvaluable
	// xxx redirect
}

// ----------------------------------------------------------------
// Example:
//   'tee > "foo.dat", $*'
// Only $* can be the expression for tee. (This is a syntactic special case of emit.)

func (this *RootNode) BuildTeeStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeTeeStatement)
	lib.InternalCodingErrorIf(len(astNode.Children) != 2)

	expressionNode := astNode.Children[0]

	teeEvaluable, err := this.BuildEvaluableNode(expressionNode)
	if err != nil {
		return nil, err
	}
	return &TeeStatementNode{
		teeEvaluable: teeEvaluable,
	}, nil
}

func (this *TeeStatementNode) Execute(state *State) (*BlockExitPayload, error) {
	teeValue := this.teeEvaluable.Evaluate(state)
	if !teeValue.IsAbsent() {
		// xxx temp
		fmt.Println(teeValue.String())
	}

	return nil, nil
}
