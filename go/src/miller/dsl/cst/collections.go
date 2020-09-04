package cst

import (
	"miller/dsl"
	"miller/lib"
)

// ================================================================
// CST build/execute for AST array-literal, map-literal, index-access, and
// slice-access nodes
// ================================================================

// ----------------------------------------------------------------
type ArrayLiteralNode struct {
	evaluables []IEvaluable
}

func BuildArrayLiteralNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeArrayLiteral)
	// An empty array should have non-nil zero-length children, not nil
	// children
	lib.InternalCodingErrorIf(astNode.Children == nil)

	evaluables := make([]IEvaluable, 0)

	for _, astChild := range astNode.Children {
		element, err := BuildEvaluableNode(astChild)
		if err != nil {
			return nil, err
		}
		evaluables = append(evaluables, element)
	}

	return &ArrayLiteralNode{evaluables: evaluables}, nil
}

func (this *ArrayLiteralNode) Evaluate(state *State) lib.Mlrval {
	mlrvals := make([]lib.Mlrval, 0)
	for _, evaluable := range this.evaluables {
		mlrval := evaluable.Evaluate(state)
		mlrvals = append(mlrvals, mlrval)
	}
	return lib.MlrvalFromArrayLiteral(mlrvals)
}

// ----------------------------------------------------------------
type ArrayOrMapIndexAccessNode struct {
	baseEvaluable  IEvaluable
	indexEvaluable IEvaluable
}

func BuildArrayOrMapIndexAccessNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeArrayOrMapIndexAccess)
	lib.InternalCodingErrorIf(len(astNode.Children) != 2)

	baseASTNode := astNode.Children[0]
	indexASTNode := astNode.Children[1]

	baseEvaluable, err := BuildEvaluableNode(baseASTNode)
	if err != nil {
		return nil, err
	}
	indexEvaluable, err := BuildEvaluableNode(indexASTNode)
	if err != nil {
		return nil, err
	}

	return &ArrayOrMapIndexAccessNode{
		baseEvaluable:  baseEvaluable,
		indexEvaluable: indexEvaluable,
	}, nil
}

func (this *ArrayOrMapIndexAccessNode) Evaluate(state *State) lib.Mlrval {
	baseMlrval := this.baseEvaluable.Evaluate(state)
	indexMlrval := this.indexEvaluable.Evaluate(state)
	// Base-is-array and index-is-int will be checked there
	return baseMlrval.GetArrayIndex(&indexMlrval)
}

// ----------------------------------------------------------------
func BuildArraySliceAccessNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeArraySliceAccess)

	// TODO

	return BuildPanicNode(), nil
}

//	if astNode.Type == dsl.NodeTypeArraySliceEmptyLowerIndex {
//		return BuildPanicNode(), nil // xxx temp
//	}
//	if astNode.Type == dsl.NodeTypeArraySliceEmptyUpperIndex {
//		return BuildPanicNode(), nil // xxx temp
//	}

// ----------------------------------------------------------------
type MapLiteralNode struct {
	// needs array of key/value Mlrval pairs
}

func BuildMapLiteralNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeMapLiteral)
	// An empty array should have non-nil zero-length children, not nil
	// children
	lib.InternalCodingErrorIf(astNode.Children == nil)

	// TODO

	return &MapLiteralNode{
	}, nil
}

func (this *MapLiteralNode) Evaluate(state *State) lib.Mlrval {
	return lib.MlrvalFromError() // xxx temp
}
