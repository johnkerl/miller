package cst

import (
	"fmt"
	"os"

	"miller/dsl"
	"miller/lib"
)

// ================================================================
// CST build/execute for statements: assignments, bare booleans,
// break/continue/return, etc.
// ================================================================

// ----------------------------------------------------------------
func BuildAssignmentNode(
	astNode *dsl.ASTNode,
) (*AssignmentNode, error) {

	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeAssignment)
	err := astNode.CheckArity(2)
	if err != nil {
		return nil, err
	}

	lhsASTNode := astNode.Children[0]
	rhsASTNode := astNode.Children[1]

	lvalue, err := BuildAssignableNode(lhsASTNode)
	if err != nil {
		return nil, err
	}

	rvalue, err := BuildEvaluableNode(rhsASTNode)
	if err != nil {
		return nil, err
	}

	return &AssignmentNode{
		lvalue: lvalue,
		rvalue: rvalue,
	}, nil
}

// ----------------------------------------------------------------
type AssignmentNode struct {
	lvalue IAssignable
	rvalue IEvaluable
}

func NewAssignmentNode(
	lvalue IAssignable,
	rvalue IEvaluable,
) *AssignmentNode {
	return &AssignmentNode{
		lvalue: lvalue,
		rvalue: rvalue,
	}
}

func (this *AssignmentNode) Execute(state *State) {
	rvalue := this.rvalue.Evaluate(state)
	if !rvalue.IsAbsent() {
		// xxx need to propagate the error coming back in the Execute() API
		err := this.lvalue.Assign(&rvalue, state)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
