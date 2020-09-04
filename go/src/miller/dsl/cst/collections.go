package cst

import (
	"errors"

	"miller/dsl"
	"miller/lib"
)

// ================================================================
// CST build/execute for AST array-literal, map-literal, index-access, and
// slice-access nodes
// ================================================================

// ----------------------------------------------------------------
func BuildArrayLiteralNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeArrayLiteral)

	// xxx assert children not nil (0-length non-nil ok, nil not ok)

	// xxx temp
	return BuildPanicNode(), nil
}

// ----------------------------------------------------------------
func BuildArraySliceAccessNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeArraySliceAccess)

	// xxx temp
	return BuildPanicNode(), nil
}

//// ----------------------------------------------------------------
//type ArrayLiteralNode struct {
//	elements []IEvaluable
//}
//
//func BuildArrayLiteralNode(astChildren []*dsl.ASTNode) *ArrayLiteralNode {
//	...
//	return &ArrayLiteralNode{
//		elements: ...
//	}
//}
//func (this *ArrayLiteralNode) Evaluate(state *State) lib.Mlrval {
//	return ...
//}

////	if astNode.Type == dsl.NodeTypeArraySliceEmptyLowerIndex {
////		return BuildPanicNode(), nil // xxx temp
////	}
////	if astNode.Type == dsl.NodeTypeArraySliceEmptyUpperIndex {
////		return BuildPanicNode(), nil // xxx temp
////	}

// ----------------------------------------------------------------
func BuildMapLiteralNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeMapLiteral)

	// xxx temp
	return BuildPanicNode(), nil

	return nil, errors.New("CST builder: unhandled AST array node " + string(astNode.Type))
}
