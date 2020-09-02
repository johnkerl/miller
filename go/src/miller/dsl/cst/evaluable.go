package cst

import (
	"errors"

	"miller/dsl"
)

// ================================================================
// This handles anything on the right-hand sides of assignment statements.
// (Also, computed field names on the left-hand sides of assignment
// statements.)
// ================================================================

// ----------------------------------------------------------------
func NewEvaluable(astNode *dsl.ASTNode) (IEvaluable, error) {
	if astNode.Children == nil {
		return NewEvaluableLeafNode(astNode)
	}

	if astNode.Type == dsl.NodeTypeOperator {
		return NewOperatorNode(astNode)
	}

	// xxx more

	return nil, errors.New("CST builder: unhandled AST node type " + string(astNode.Type))
}
