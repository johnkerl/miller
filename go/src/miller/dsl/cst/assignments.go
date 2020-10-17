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
func (this *RootNode) BuildAssignmentNode(
	astNode *dsl.ASTNode,
) (*AssignmentNode, error) {

	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeAssignment)
	err := astNode.CheckArity(2)
	if err != nil {
		return nil, err
	}

	lhsASTNode := astNode.Children[0]
	rhsASTNode := astNode.Children[1]

	lvalueNode, err := this.BuildAssignableNode(lhsASTNode)
	if err != nil {
		return nil, err
	}

	rvalueNode, err := this.BuildEvaluableNode(rhsASTNode)
	if err != nil {
		return nil, err
	}

	return &AssignmentNode{
		lvalueNode: lvalueNode,
		rvalueNode: rvalueNode,
	}, nil
}

// ----------------------------------------------------------------
type AssignmentNode struct {
	lvalueNode IAssignable
	rvalueNode IEvaluable
}

func NewAssignmentNode(
	lvalueNode IAssignable,
	rvalueNode IEvaluable,
) *AssignmentNode {
	return &AssignmentNode{
		lvalueNode: lvalueNode,
		rvalueNode: rvalueNode,
	}
}

func (this *AssignmentNode) Execute(state *State) (*BlockExitPayload, error) {
	rvalue := this.rvalueNode.Evaluate(state)
	if !rvalue.IsAbsent() {
		err := this.lvalueNode.Assign(&rvalue, state)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
