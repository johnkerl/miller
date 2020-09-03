package cst

import (
	"errors"

	"miller/dsl"
	"miller/lib"
)

// ================================================================
// CST build/execute for AST array nodes
// ================================================================

// ----------------------------------------------------------------
func BuildArrayLiteralNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeArrayLiteral)

	// xxx temp
	return BuildPanicNode(), nil

	return nil, errors.New("CST builder: unhandled AST array node " + string(astNode.Type))
}

//// ----------------------------------------------------------------
//type ArrayLiteralNode struct {
//	literal lib.Mlrval
//}
//
//func BuildArrayLiteralNode(literal string) *ArrayLiteralNode {
//	return &ArrayLiteral{
//		literal: lib.MlrvalFromString(literal),
//	}
//}
//func (this *ArrayLiteralNode) Evaluate(state *State) lib.Mlrval {
//	return this.literal
//}
