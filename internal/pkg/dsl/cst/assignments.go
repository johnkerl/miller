// ================================================================
// CST build/execute for assignment and unset statements.
// ================================================================

package cst

import (
	"mlr/internal/pkg/dsl"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/runtime"
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

	lvalueNodes := make([]IAssignable, len(astNode.Children))

	for i, lhsASTNode := range astNode.Children {
		lvalueNode, err := root.BuildAssignableNode(lhsASTNode)
		if err != nil {
			return nil, err
		}
		lvalueNodes[i] = lvalueNode
	}

	return &UnsetNode{
		lvalueNodes: lvalueNodes,
	}, nil
}

// ----------------------------------------------------------------
type UnsetNode struct {
	lvalueNodes []IAssignable
}

func NewUnsetNode(
	lvalueNodes []IAssignable,
) *UnsetNode {
	return &UnsetNode{
		lvalueNodes: lvalueNodes,
	}
}

func (node *UnsetNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	for _, lvalueNode := range node.lvalueNodes {
		lvalueNode.Unassign(state)
	}
	return nil, nil
}
