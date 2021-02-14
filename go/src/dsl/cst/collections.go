// ================================================================
// CST build/execute for AST array-literal, map-literal, index-access, and
// slice-access nodes
// ================================================================

package cst

import (
	"miller/src/dsl"
	"miller/src/lib"
	"miller/src/runtime"
	"miller/src/types"
)

// ----------------------------------------------------------------
type ArrayLiteralNode struct {
	evaluables []IEvaluable
}

func (this *RootNode) BuildArrayLiteralNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeArrayLiteral)
	// An empty array should have non-nil zero-length children, not nil
	// children
	lib.InternalCodingErrorIf(astNode.Children == nil)

	evaluables := make([]IEvaluable, 0)

	for _, astChild := range astNode.Children {
		element, err := this.BuildEvaluableNode(astChild)
		if err != nil {
			return nil, err
		}
		evaluables = append(evaluables, element)
	}

	return &ArrayLiteralNode{evaluables: evaluables}, nil
}

func (this *ArrayLiteralNode) Evaluate(state *runtime.State) types.Mlrval {
	mlrvals := make([]types.Mlrval, 0)
	for _, evaluable := range this.evaluables {
		mlrval := evaluable.Evaluate(state)
		mlrvals = append(mlrvals, mlrval)
	}
	return types.MlrvalFromArrayLiteralReference(mlrvals)
}

// ----------------------------------------------------------------
type ArrayOrMapIndexAccessNode struct {
	baseEvaluable  IEvaluable
	indexEvaluable IEvaluable
}

func (this *RootNode) BuildArrayOrMapIndexAccessNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeArrayOrMapIndexAccess)
	lib.InternalCodingErrorIf(len(astNode.Children) != 2)

	baseASTNode := astNode.Children[0]
	indexASTNode := astNode.Children[1]

	baseEvaluable, err := this.BuildEvaluableNode(baseASTNode)
	if err != nil {
		return nil, err
	}
	indexEvaluable, err := this.BuildEvaluableNode(indexASTNode)
	if err != nil {
		return nil, err
	}

	return &ArrayOrMapIndexAccessNode{
		baseEvaluable:  baseEvaluable,
		indexEvaluable: indexEvaluable,
	}, nil
}

func (this *ArrayOrMapIndexAccessNode) Evaluate(state *runtime.State) types.Mlrval {
	baseMlrval := this.baseEvaluable.Evaluate(state)
	indexMlrval := this.indexEvaluable.Evaluate(state)

	// Base-is-array and index-is-int will be checked there
	if baseMlrval.IsArray() {
		return baseMlrval.ArrayGet(&indexMlrval)
	} else if baseMlrval.IsMap() {
		return baseMlrval.MapGet(&indexMlrval)
	} else if baseMlrval.IsAbsent() {
		return baseMlrval
	} else {
		return types.MlrvalFromError()
	}
}

// ----------------------------------------------------------------
type ArraySliceAccessNode struct {
	baseEvaluable       IEvaluable
	lowerIndexEvaluable IEvaluable
	upperIndexEvaluable IEvaluable
}

func (this *RootNode) BuildArraySliceAccessNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeArraySliceAccess)
	lib.InternalCodingErrorIf(len(astNode.Children) != 3)

	baseASTNode := astNode.Children[0]
	lowerIndexASTNode := astNode.Children[1]
	upperIndexASTNode := astNode.Children[2]

	baseEvaluable, err := this.BuildEvaluableNode(baseASTNode)
	if err != nil {
		return nil, err
	}

	lowerIndexEvaluable, err := this.BuildEvaluableNode(lowerIndexASTNode)
	if err != nil {
		return nil, err
	}

	upperIndexEvaluable, err := this.BuildEvaluableNode(upperIndexASTNode)
	if err != nil {
		return nil, err
	}

	return &ArraySliceAccessNode{
		baseEvaluable:       baseEvaluable,
		lowerIndexEvaluable: lowerIndexEvaluable,
		upperIndexEvaluable: upperIndexEvaluable,
	}, nil
}

func (this *ArraySliceAccessNode) Evaluate(state *runtime.State) types.Mlrval {
	baseMlrval := this.baseEvaluable.Evaluate(state)
	lowerIndexMlrval := this.lowerIndexEvaluable.Evaluate(state)
	upperIndexMlrval := this.upperIndexEvaluable.Evaluate(state)

	if baseMlrval.IsAbsent() {
		return types.MlrvalFromAbsent()
	}
	if baseMlrval.IsString() {
		return types.MlrvalSubstr(&baseMlrval, &lowerIndexMlrval, &upperIndexMlrval)
	}
	array := baseMlrval.GetArray()
	if array == nil {
		return types.MlrvalFromError()
	}
	n := len(array)

	if lowerIndexMlrval.IsAbsent() {
		return types.MlrvalFromAbsent()
	}
	if upperIndexMlrval.IsAbsent() {
		return types.MlrvalFromAbsent()
	}

	lowerIndex, ok := lowerIndexMlrval.GetIntValue()
	if !ok {
		if lowerIndexMlrval.IsEmpty() {
			lowerIndex = 1
		} else {
			return types.MlrvalFromError()
		}
	}
	upperIndex, ok := upperIndexMlrval.GetIntValue()
	if !ok {
		if upperIndexMlrval.IsEmpty() {
			upperIndex = n
		} else {
			return types.MlrvalFromError()
		}
	}

	// UnaliasArrayIndex returns a boolean second return value to indicate
	// whether the index is in range. But here, for the slicing operation, we
	// inspect the in-range-ness ourselves so we discard that 2nd return value.
	lowerZindex, _ := types.UnaliasArrayIndex(&array, lowerIndex)
	upperZindex, _ := types.UnaliasArrayIndex(&array, upperIndex)

	if lowerZindex > upperZindex {
		return types.MlrvalFromArrayLiteralReference(make([]types.Mlrval, 0))
	}

	// Say x=[1,2,3,4,5]. Then x[3:10] is [3,4,5].
	if lowerZindex < 0 {
		lowerZindex = 0
	}
	if upperZindex > n-1 {
		upperZindex = n - 1
	}

	// Go     slices have inclusive lower bound, exclusive upper bound.
	// Miller slices have inclusive lower bound, inclusive upper bound.
	var m = upperZindex - lowerZindex + 1
	retval := make([]types.Mlrval, m)
	di := 0
	for si := lowerZindex; si <= upperZindex; si++ {
		retval[di] = *array[si].Copy()
		di++
	}

	return types.MlrvalFromArrayLiteralReference(retval)
}

// ================================================================
// For input record 'a=7,b=8,c=9',  $[[2]] = "b"

type PositionalFieldNameNode struct {
	indexEvaluable IEvaluable
}

func (this *RootNode) BuildPositionalFieldNameNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypePositionalFieldName)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)

	indexASTNode := astNode.Children[0]

	indexEvaluable, err := this.BuildEvaluableNode(indexASTNode)
	if err != nil {
		return nil, err
	}

	return &PositionalFieldNameNode{
		indexEvaluable: indexEvaluable,
	}, nil
}

// TODO: code-dedupe these next four Evaluate methods
func (this *PositionalFieldNameNode) Evaluate(state *runtime.State) types.Mlrval {
	indexMlrval := this.indexEvaluable.Evaluate(state)
	if indexMlrval.IsAbsent() {
		return types.MlrvalFromAbsent()
	}

	index, ok := indexMlrval.GetIntValue()
	if !ok {
		return types.MlrvalFromError()
	}

	name, ok := state.Inrec.GetNameAtPositionalIndex(index)
	if !ok {
		return types.MlrvalFromAbsent()
	}

	return types.MlrvalFromString(name)
}

// ================================================================
// For input record 'a=7,b=8,c=9',  $[[2]] = 8

type PositionalFieldValueNode struct {
	indexEvaluable IEvaluable
}

func (this *RootNode) BuildPositionalFieldValueNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypePositionalFieldValue)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)

	indexASTNode := astNode.Children[0]

	indexEvaluable, err := this.BuildEvaluableNode(indexASTNode)
	if err != nil {
		return nil, err
	}

	return &PositionalFieldValueNode{
		indexEvaluable: indexEvaluable,
	}, nil
}

func (this *PositionalFieldValueNode) Evaluate(state *runtime.State) types.Mlrval {
	indexMlrval := this.indexEvaluable.Evaluate(state)
	if indexMlrval.IsAbsent() {
		return types.MlrvalFromAbsent()
	}

	index, ok := indexMlrval.GetIntValue()
	if !ok {
		return types.MlrvalFromError()
	}

	retval := state.Inrec.GetWithPositionalIndex(index)
	if retval == nil {
		return types.MlrvalFromAbsent()
	}

	return *retval
}

// ================================================================
// For x = [7,8,9], x[[2]] = 2
// For y = {"a":7,"b":8,"c":9}, y[[2]] = "b"
type ArrayOrMapPositionalNameAccessNode struct {
	baseEvaluable  IEvaluable
	indexEvaluable IEvaluable
}

func (this *RootNode) BuildArrayOrMapPositionalNameAccessNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeArrayOrMapPositionalNameAccess)
	lib.InternalCodingErrorIf(len(astNode.Children) != 2)

	baseASTNode := astNode.Children[0]
	indexASTNode := astNode.Children[1]

	baseEvaluable, err := this.BuildEvaluableNode(baseASTNode)
	if err != nil {
		return nil, err
	}
	indexEvaluable, err := this.BuildEvaluableNode(indexASTNode)
	if err != nil {
		return nil, err
	}

	return &ArrayOrMapPositionalNameAccessNode{
		baseEvaluable:  baseEvaluable,
		indexEvaluable: indexEvaluable,
	}, nil
}

func (this *ArrayOrMapPositionalNameAccessNode) Evaluate(state *runtime.State) types.Mlrval {
	baseMlrval := this.baseEvaluable.Evaluate(state)
	indexMlrval := this.indexEvaluable.Evaluate(state)

	if indexMlrval.IsAbsent() {
		return types.MlrvalFromAbsent()
	}

	index, ok := indexMlrval.GetIntValue()
	if !ok {
		return types.MlrvalFromError()
	}

	if baseMlrval.IsArray() {
		n, _ := baseMlrval.GetArrayLength()
		zindex, ok := types.UnaliasArrayLengthIndex(int(n), index)
		if ok {
			return types.MlrvalFromInt(zindex + 1) // Miller user-space indices are 1-up
		} else {
			return types.MlrvalFromAbsent()
		}

	} else if baseMlrval.IsMap() {
		name, ok := baseMlrval.GetMap().GetNameAtPositionalIndex(index)
		if !ok {
			return types.MlrvalFromAbsent()
		} else {
			return types.MlrvalFromString(name)
		}

	} else if baseMlrval.IsAbsent() {
		return baseMlrval

	} else {
		return types.MlrvalFromError()
	}
}

// ================================================================
// For x = [7,8,9], x[[2]] = 8
// For y = {"a":7,"b":8,"c":9}, y[[2]] = 8
type ArrayOrMapPositionalValueAccessNode struct {
	baseEvaluable  IEvaluable
	indexEvaluable IEvaluable
}

func (this *RootNode) BuildArrayOrMapPositionalValueAccessNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeArrayOrMapPositionalValueAccess)
	lib.InternalCodingErrorIf(len(astNode.Children) != 2)

	baseASTNode := astNode.Children[0]
	indexASTNode := astNode.Children[1]

	baseEvaluable, err := this.BuildEvaluableNode(baseASTNode)
	if err != nil {
		return nil, err
	}
	indexEvaluable, err := this.BuildEvaluableNode(indexASTNode)
	if err != nil {
		return nil, err
	}

	return &ArrayOrMapPositionalValueAccessNode{
		baseEvaluable:  baseEvaluable,
		indexEvaluable: indexEvaluable,
	}, nil
}

func (this *ArrayOrMapPositionalValueAccessNode) Evaluate(state *runtime.State) types.Mlrval {
	baseMlrval := this.baseEvaluable.Evaluate(state)
	indexMlrval := this.indexEvaluable.Evaluate(state)

	if indexMlrval.IsAbsent() {
		return types.MlrvalFromAbsent()
	}

	index, ok := indexMlrval.GetIntValue()
	if !ok {
		return types.MlrvalFromError()
	}

	if baseMlrval.IsArray() {
		return baseMlrval.ArrayGet(&indexMlrval)

	} else if baseMlrval.IsMap() {
		value := baseMlrval.GetMap().GetWithPositionalIndex(index)
		if value == nil {
			return types.MlrvalFromAbsent()
		}

		retval := value.Copy()
		return *retval

	} else if baseMlrval.IsAbsent() {
		return baseMlrval

	} else {
		return types.MlrvalFromError()
	}
}

// ================================================================
// This is for computing map entries at runtime. For example, in
//
//   mlr put 'mymap = {"sum": $x + $y, "diff": $x - $y}; ...'
//
// the first pair would have key being string-literal "sum" and value being the
// evaluable expression '$x + $y'.

type EvaluablePair struct {
	Key   IEvaluable
	Value IEvaluable
}

func NewEvaluablePair(key IEvaluable, value IEvaluable) *EvaluablePair {
	return &EvaluablePair{
		Key:   key,
		Value: value,
	}
}

// ----------------------------------------------------------------
type MapLiteralNode struct {
	evaluablePairs []EvaluablePair
	// needs array of key/value Mlrval pairs
}

func (this *RootNode) BuildMapLiteralNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeMapLiteral)
	// An empty array should have non-nil zero-length children, not nil
	// children
	lib.InternalCodingErrorIf(astNode.Children == nil)

	evaluablePairs := make([]EvaluablePair, 0)
	for _, astChild := range astNode.Children {
		lib.InternalCodingErrorIf(astChild.Children == nil)
		lib.InternalCodingErrorIf(len(astChild.Children) != 2)
		astKey := astChild.Children[0]
		astValue := astChild.Children[1]

		cstKey, err := this.BuildEvaluableNode(astKey)
		if err != nil {
			return nil, err
		}
		cstValue, err := this.BuildEvaluableNode(astValue)
		if err != nil {
			return nil, err
		}

		evaluablePair := NewEvaluablePair(cstKey, cstValue)
		evaluablePairs = append(evaluablePairs, *evaluablePair)
	}
	return &MapLiteralNode{evaluablePairs: evaluablePairs}, nil
}

func (this *MapLiteralNode) Evaluate(state *runtime.State) types.Mlrval {
	mlrval := types.MlrvalEmptyMap()

	for _, evaluablePair := range this.evaluablePairs {
		mkey := evaluablePair.Key.Evaluate(state)
		mvalue := evaluablePair.Value.Evaluate(state)

		if !mvalue.IsAbsent() {
			mlrval.MapPut(&mkey, &mvalue)
		}
	}

	return mlrval
}
