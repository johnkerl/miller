// ================================================================
// CST build/execute for assignment and unset statements.
// ================================================================

package cst

import (
	"miller/src/dsl"
	"miller/src/lib"
	"miller/src/runtime"
)

// ================================================================
func (root *RootNode) BuildAssignmentNode(
	astNode *dsl.ASTNode,
) (*AssignmentNode, error) {

	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeAssignment)
	err := astNode.CheckArity(2)
	if err != nil {
		return nil, err
	}

	lhsASTNode := astNode.Children[0]
	rhsASTNode := astNode.Children[1]

	lvalueNode, err := root.BuildAssignableNode(lhsASTNode)
	if err != nil {
		return nil, err
	}

	rvalueNode, err := root.BuildEvaluableNode(rhsASTNode)
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

func (node *AssignmentNode) Execute(
	state *runtime.State,
) (*BlockExitPayload, error) {
	rvalue := node.rvalueNode.Evaluate(state)
	if !rvalue.IsAbsent() {
		err := node.lvalueNode.Assign(rvalue, state)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// ================================================================
func (root *RootNode) BuildUnsetNode(
	astNode *dsl.ASTNode,
) (*UnsetNode, error) {

	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeUnset)
	err := astNode.CheckArity(1)
	if err != nil {
		return nil, err
	}

	lhsASTNode := astNode.Children[0]

	lvalueNode, err := root.BuildAssignableNode(lhsASTNode)
	if err != nil {
		return nil, err
	}

	return &UnsetNode{
		lvalueNode: lvalueNode,
	}, nil
}

// ----------------------------------------------------------------
type UnsetNode struct {
	lvalueNode IAssignable
}

func NewUnsetNode(
	lvalueNode IAssignable,
) *UnsetNode {
	return &UnsetNode{
		lvalueNode: lvalueNode,
	}
}

func (node *UnsetNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	node.lvalueNode.Unassign(state)
	return nil, nil
}
