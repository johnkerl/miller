package cst

import (
	"errors"

	"miller/dsl"
)

// ----------------------------------------------------------------
func NewEvaluable(astNode *dsl.ASTNode) (IEvaluable, error) {
	if astNode.Children == nil {
		return NewEvaluableLeafNode(astNode)
	}

	if astNode.Type == dsl.NodeTypeOperator {
		return NewOperatorNode(astNode)
	}

	// xxx more

	return nil, errors.New("CST builder: unhandled AST node type" + string(astNode.Type))
}
