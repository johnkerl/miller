package cst

import (
	"errors"

	"miller/dsl"
	"miller/lib"
)

// ================================================================
// CST build/execute for AST map nodes
// ================================================================

// ----------------------------------------------------------------
func BuildMapLiteralNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeMapLiteral)

	// xxx temp
	return BuildPanicNode(), nil

	return nil, errors.New("CST builder: unhandled AST array node " + string(astNode.Type))
}
