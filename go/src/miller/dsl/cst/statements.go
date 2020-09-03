package cst

import (
	"miller/dsl"
)

// ================================================================
// CST build/execute for statements: assignments, bare booleans,
// break/continue/return, etc.
// ================================================================

// ----------------------------------------------------------------
func BuildSrecDirectFieldAssignmentNode(
	astNode *dsl.ASTNode,
) (*SrecDirectFieldAssignmentNode, error) {

	err := astNode.CheckArity(2)
	if err != nil {
		return nil, err
	}

	lhsASTNode := astNode.Children[0]
	rhsASTNode := astNode.Children[1]

	lhsFieldName := string(lhsASTNode.Token.Lit)
	rhs, err := BuildEvaluableNode(rhsASTNode)
	if err != nil {
		return nil, err
	}

	return &SrecDirectFieldAssignmentNode{
		lhsFieldName: lhsFieldName,
		rhs:          rhs,
	}, nil
}

func (this *SrecDirectFieldAssignmentNode) Execute(state *State) {
	value := this.rhs.Evaluate(state)
	if !value.IsAbsent() {
		state.Inrec.Put(&this.lhsFieldName, &value)
	}
}
