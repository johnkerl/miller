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
func BuildEvaluableNode(astNode *dsl.ASTNode) (IEvaluable, error) {

	if astNode.Children == nil {
		return BuildLeafNode(astNode)
	}

	switch astNode.Type {

	case dsl.NodeTypeOperator:
		return BuildOperatorNode(astNode)

	case dsl.NodeTypeArrayLiteral:
		return BuildArrayLiteralNode(astNode)

	case dsl.NodeTypeMapLiteral:
		return BuildMapLiteralNode(astNode)

	case dsl.NodeTypeArrayOrMapIndexAccess:
		return BuildArrayOrMapIndexAccessNode(astNode)

	case dsl.NodeTypeArraySliceAccess:
		return BuildArraySliceAccessNode(astNode)

	}

	// xxx if/while/etc
	// xxx function
	// xxx more

	return nil, errors.New(
		"CST BuildEvaluableNode: unhandled AST node type " + string(astNode.Type),
	)
}
