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
func NewArrayLiteralNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeArrayLiteral)

	// xxx temp
	return NewPanic(), nil

	return nil, errors.New("CST builder: unhandled AST array node " + string(astNode.Type))
}

//// ----------------------------------------------------------------
//type StringLiteral struct {
//	literal lib.Mlrval
//}
//
//func NewStringLiteral(literal string) *StringLiteral {
//	return &StringLiteral{
//		literal: lib.MlrvalFromString(literal),
//	}
//}
//func (this *StringLiteral) Evaluate(state *State) lib.Mlrval {
//	return this.literal
//}
