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

func (this *AssignmentNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	rvalue := this.rvalueNode.Evaluate(state)
	if !rvalue.IsAbsent() {
		err := this.lvalueNode.Assign(&rvalue, state)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// ================================================================
func (this *RootNode) BuildUnsetNode(
	astNode *dsl.ASTNode,
) (*UnsetNode, error) {

	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeUnset)
	err := astNode.CheckArity(1)
	if err != nil {
		return nil, err
	}

	lhsASTNode := astNode.Children[0]

	lvalueNode, err := this.BuildAssignableNode(lhsASTNode)
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

func (this *UnsetNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	this.lvalueNode.Unset(state)
	return nil, nil
}
