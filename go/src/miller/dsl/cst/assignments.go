package cst

import (
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

func (this *AssignmentNode) Execute(state *State) error {
	rvalue := this.rvalue.Evaluate(state)
	if !rvalue.IsAbsent() {
		err := this.lvalue.Assign(&rvalue, state)
		if err != nil {
			return err
		}
	}
	return nil
}
