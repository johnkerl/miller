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
	mlrvals    []types.Mlrval
}

func (this *RootNode) BuildArrayLiteralNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeArrayLiteral)
	// An empty array should have non-nil zero-length children, not nil
	// children
	lib.InternalCodingErrorIf(astNode.Children == nil)

	evaluables := make([]IEvaluable, len(astNode.Children))
	mlrvals := make([]types.Mlrval, len(astNode.Children))

	for i, astChild := range astNode.Children {
		element, err := this.BuildEvaluableNode(astChild)
		if err != nil {
			return nil, err
		}
		evaluables[i] = element
		mlrvals[i] = types.MlrvalFromError()
	}

	return &ArrayLiteralNode{
		evaluables: evaluables,
		mlrvals:    mlrvals,
	}, nil
}

func (this *ArrayLiteralNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	for i, _ := range this.evaluables {
		this.evaluables[i].Evaluate(&this.mlrvals[i], state)
	}
	output.SetFromArrayLiteralReference(this.mlrvals)
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

func (this *ArrayOrMapIndexAccessNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	var baseMlrval, indexMlrval types.Mlrval
	this.baseEvaluable.Evaluate(&baseMlrval, state)
	this.indexEvaluable.Evaluate(&indexMlrval, state)

	// Base-is-array and index-is-int will be checked there
	if baseMlrval.IsArray() {
		*output = baseMlrval.ArrayGet(&indexMlrval)
	} else if baseMlrval.IsMap() {
		*output = baseMlrval.MapGet(&indexMlrval)
	} else if baseMlrval.IsAbsent() {
		output.SetFromAbsent()
	} else {
		output.SetFromError()
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

func (this *ArraySliceAccessNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	var baseMlrval, lowerIndexMlrval, upperIndexMlrval types.Mlrval
	this.baseEvaluable.Evaluate(&baseMlrval, state)
	this.lowerIndexEvaluable.Evaluate(&lowerIndexMlrval, state)
	this.upperIndexEvaluable.Evaluate(&upperIndexMlrval, state)

	if baseMlrval.IsAbsent() {
		output.SetFromAbsent()
	}
	if baseMlrval.IsString() {
		types.MlrvalSubstr(output, &baseMlrval, &lowerIndexMlrval, &upperIndexMlrval)
		return
	}
	array := baseMlrval.GetArray()
	if array == nil {
		output.SetFromError()
		return
	}
	n := len(array)

	if lowerIndexMlrval.IsAbsent() {
		output.SetFromAbsent()
		return
	}
	if upperIndexMlrval.IsAbsent() {
		output.SetFromAbsent()
		return
	}

	lowerIndex, ok := lowerIndexMlrval.GetIntValue()
	if !ok {
		if lowerIndexMlrval.IsEmpty() {
			lowerIndex = 1
		} else {
			output.SetFromError()
			return
		}
	}
	upperIndex, ok := upperIndexMlrval.GetIntValue()
	if !ok {
		if upperIndexMlrval.IsEmpty() {
			upperIndex = n
		} else {
			output.SetFromError()
			return
		}
	}

	// UnaliasArrayIndex returns a boolean second return value to indicate
	// whether the index is in range. But here, for the slicing operation, we
	// inspect the in-range-ness ourselves so we discard that 2nd return value.
	lowerZindex, _ := types.UnaliasArrayIndex(&array, lowerIndex)
	upperZindex, _ := types.UnaliasArrayIndex(&array, upperIndex)

	if lowerZindex > upperZindex {
		output.SetFromArrayLiteralReference(make([]types.Mlrval, 0))
		return
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

	output.SetFromArrayLiteralReference(retval)
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
func (this *PositionalFieldNameNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	var indexMlrval types.Mlrval
	this.indexEvaluable.Evaluate(&indexMlrval, state)
	if indexMlrval.IsAbsent() {
		output.SetFromAbsent()
		return
	}

	index, ok := indexMlrval.GetIntValue()
	if !ok {
		output.SetFromError()
		return
	}

	name, ok := state.Inrec.GetNameAtPositionalIndex(index)
	if !ok {
		output.SetFromAbsent()
		return
	}

	output.SetFromString(name)
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

func (this *PositionalFieldValueNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	var indexMlrval types.Mlrval
	this.indexEvaluable.Evaluate(&indexMlrval, state)
	if indexMlrval.IsAbsent() {
		output.SetFromAbsent()
		return
	}

	index, ok := indexMlrval.GetIntValue()
	if !ok {
		output.SetFromError()
		return
	}

	retval := state.Inrec.GetWithPositionalIndex(index)
	if retval == nil {
		output.SetFromAbsent()
		return
	}

	output.CopyFrom(retval)
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

func (this *ArrayOrMapPositionalNameAccessNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	var baseMlrval, indexMlrval types.Mlrval
	this.baseEvaluable.Evaluate(&baseMlrval, state)
	this.indexEvaluable.Evaluate(&indexMlrval, state)

	if indexMlrval.IsAbsent() {
		output.SetFromAbsent()
		return
	}

	index, ok := indexMlrval.GetIntValue()
	if !ok {
		output.SetFromError()
		return
	}

	if baseMlrval.IsArray() {
		n, _ := baseMlrval.GetArrayLength()
		zindex, ok := types.UnaliasArrayLengthIndex(int(n), index)
		if ok {
			output.SetFromInt(zindex + 1) // Miller user-space indices are 1-up
			return
		} else {
			output.SetFromAbsent()
			return
		}

	} else if baseMlrval.IsMap() {
		name, ok := baseMlrval.GetMap().GetNameAtPositionalIndex(index)
		if !ok {
			output.SetFromAbsent()
			return
		} else {
			output.SetFromString(name)
			return
		}

	} else if baseMlrval.IsAbsent() {
		output.SetFromAbsent()
		return

	} else {
		output.SetFromError()
		return
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

func (this *ArrayOrMapPositionalValueAccessNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	var baseMlrval, indexMlrval types.Mlrval
	this.baseEvaluable.Evaluate(&baseMlrval, state)
	this.indexEvaluable.Evaluate(&indexMlrval, state)

	if indexMlrval.IsAbsent() {
		output.SetFromAbsent()
		return
	}

	index, ok := indexMlrval.GetIntValue()
	if !ok {
		output.SetFromError()
		return
	}

	if baseMlrval.IsArray() {
		// xxx pending pointer-output refactor
		element := baseMlrval.ArrayGet(&indexMlrval)
		output.CopyFrom(&element)

	} else if baseMlrval.IsMap() {
		value := baseMlrval.GetMap().GetWithPositionalIndex(index)
		if value == nil {
			output.SetFromAbsent()
			return
		}

		output.CopyFrom(value)
		return

	} else if baseMlrval.IsAbsent() {
		output.CopyFrom(&baseMlrval)
		return

	} else {
		output.SetFromError()
		return
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
	evaluablePairs []*EvaluablePair
}

func (this *RootNode) BuildMapLiteralNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeMapLiteral)
	// An empty array should have non-nil zero-length children, not nil
	// children
	lib.InternalCodingErrorIf(astNode.Children == nil)

	evaluablePairs := make([]*EvaluablePair, len(astNode.Children))
	for i, astChild := range astNode.Children {
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

		evaluablePairs[i] = NewEvaluablePair(cstKey, cstValue)
	}
	return &MapLiteralNode{
		evaluablePairs: evaluablePairs,
	}, nil
}

func (this *MapLiteralNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	output.SetFromEmptyMap()

	for i, _ := range this.evaluablePairs {
		var key, value types.Mlrval
		this.evaluablePairs[i].Key.Evaluate(&key, state)
		this.evaluablePairs[i].Value.Evaluate(&value, state)

		if !value.IsAbsent() {
			output.MapPut(&key, &value)
		}
	}
}
