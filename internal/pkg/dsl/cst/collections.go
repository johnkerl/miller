// ================================================================
// CST build/execute for AST array-literal, map-literal, index-access, and
// slice-access nodes
// ================================================================

package cst

import (
	"github.com/johnkerl/miller/internal/pkg/bifs"
	"github.com/johnkerl/miller/internal/pkg/dsl"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
	"github.com/johnkerl/miller/internal/pkg/runtime"
)

// ----------------------------------------------------------------
type ArrayLiteralNode struct {
	evaluables []IEvaluable
}

func (node *RootNode) BuildArrayLiteralNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeArrayLiteral)
	// An empty array should have non-nil zero-length children, not nil
	// children
	lib.InternalCodingErrorIf(astNode.Children == nil)

	evaluables := make([]IEvaluable, len(astNode.Children))

	for i, astChild := range astNode.Children {
		element, err := node.BuildEvaluableNode(astChild)
		if err != nil {
			return nil, err
		}
		evaluables[i] = element
	}

	return &ArrayLiteralNode{
		evaluables: evaluables,
	}, nil
}

func (node *ArrayLiteralNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	mlrvals := make([]*mlrval.Mlrval, len(node.evaluables))
	for i := range node.evaluables {
		mlrvals[i] = node.evaluables[i].Evaluate(state)
	}
	return mlrval.FromArray(mlrvals)
}

// ----------------------------------------------------------------
type CollectionIndexAccessNode struct {
	baseEvaluable  IEvaluable
	indexEvaluable IEvaluable
}

func (node *RootNode) BuildArrayOrMapIndexAccessNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeArrayOrMapIndexAccess)
	lib.InternalCodingErrorIf(len(astNode.Children) != 2)

	baseASTNode := astNode.Children[0]
	indexASTNode := astNode.Children[1]

	baseEvaluable, err := node.BuildEvaluableNode(baseASTNode)
	if err != nil {
		return nil, err
	}
	indexEvaluable, err := node.BuildEvaluableNode(indexASTNode)
	if err != nil {
		return nil, err
	}

	return &CollectionIndexAccessNode{
		baseEvaluable:  baseEvaluable,
		indexEvaluable: indexEvaluable,
	}, nil
}

func (node *CollectionIndexAccessNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	baseMlrval := node.baseEvaluable.Evaluate(state)
	indexMlrval := node.indexEvaluable.Evaluate(state)

	// Base-is-array and index-is-int will be checked there
	if baseMlrval.IsArray() {
		output := baseMlrval.ArrayGet(indexMlrval)
		return &output
	} else if baseMlrval.IsMap() {
		output := baseMlrval.MapGet(indexMlrval)
		return &output
	} else if baseMlrval.IsStringOrVoid() {
		mindex, isInt := indexMlrval.GetIntValue()
		if !isInt {
			return mlrval.ERROR
		}
		// Handle UTF-8 correctly: len(input1.printrep) will count bytes, not runes.
		runes := []rune(baseMlrval.String())
		// Miller uses 1-up, and negatively aliased, indexing for strings and arrays.
		zindex, inBounds := mlrval.UnaliasArrayLengthIndex(len(runes), int(mindex))
		if !inBounds {
			return mlrval.ERROR
		}
		return mlrval.FromString(string(runes[zindex]))

	} else if baseMlrval.IsAbsent() {
		return mlrval.ABSENT
	} else {
		return mlrval.ERROR
	}
}

// ----------------------------------------------------------------
type ArraySliceAccessNode struct {
	baseEvaluable       IEvaluable
	lowerIndexEvaluable IEvaluable
	upperIndexEvaluable IEvaluable
}

func (node *RootNode) BuildArraySliceAccessNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeArraySliceAccess)
	lib.InternalCodingErrorIf(len(astNode.Children) != 3)

	baseASTNode := astNode.Children[0]
	lowerIndexASTNode := astNode.Children[1]
	upperIndexASTNode := astNode.Children[2]

	baseEvaluable, err := node.BuildEvaluableNode(baseASTNode)
	if err != nil {
		return nil, err
	}

	lowerIndexEvaluable, err := node.BuildEvaluableNode(lowerIndexASTNode)
	if err != nil {
		return nil, err
	}

	upperIndexEvaluable, err := node.BuildEvaluableNode(upperIndexASTNode)
	if err != nil {
		return nil, err
	}

	return &ArraySliceAccessNode{
		baseEvaluable:       baseEvaluable,
		lowerIndexEvaluable: lowerIndexEvaluable,
		upperIndexEvaluable: upperIndexEvaluable,
	}, nil
}

func (node *ArraySliceAccessNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	baseMlrval := node.baseEvaluable.Evaluate(state)
	lowerIndexMlrval := node.lowerIndexEvaluable.Evaluate(state)
	upperIndexMlrval := node.upperIndexEvaluable.Evaluate(state)

	if baseMlrval.IsAbsent() {
		return mlrval.ABSENT
	}
	if baseMlrval.IsString() {
		return bifs.BIF_substr_1_up(baseMlrval, lowerIndexMlrval, upperIndexMlrval)
	}
	array := baseMlrval.GetArray()
	if array == nil {
		return mlrval.ERROR
	}
	n := len(array)

	if lowerIndexMlrval.IsAbsent() {
		return mlrval.ABSENT
	}
	if upperIndexMlrval.IsAbsent() {
		return mlrval.ABSENT
	}

	lowerIndex, ok := lowerIndexMlrval.GetIntValue()
	if !ok {
		if lowerIndexMlrval.IsVoid() {
			lowerIndex = 1
		} else {
			return mlrval.ERROR
		}
	}
	upperIndex, ok := upperIndexMlrval.GetIntValue()
	if !ok {
		if upperIndexMlrval.IsVoid() {
			upperIndex = int64(n)
		} else {
			return mlrval.ERROR
		}
	}

	// UnaliasArrayIndex returns a boolean second return value to indicate
	// whether the index is in range. But here, for the slicing operation, we
	// inspect the in-range-ness ourselves so we discard that 2nd return value.
	lowerZindex, _ := mlrval.UnaliasArrayIndex(&array, int(lowerIndex))
	upperZindex, _ := mlrval.UnaliasArrayIndex(&array, int(upperIndex))

	if lowerZindex > upperZindex {
		return mlrval.FromEmptyArray()
	}

	// Semantics: say x=[1,2,3,4,5]. Then x[3:10] is [3,4,5].
	//
	// Cases:
	//      [* * * * *]              actual data
	//  [o o]                        1. attempted indexing: lo, hi both out of bounds
	//  [o o o o o o ]               2. attempted indexing: hi in bounds, lo out
	//  [o o o o o o o o o o o o]    3. attempted indexing: lo, hi both out of bounds
	//        [o o o]                4. attempted indexing: lo, hi in bounds
	//        [o o o o o o ]         5. attempted indexing: lo in bounds, hi out
	//                  [o o o o]    6. attempted indexing: lo, hi both out of bounds

	if lowerZindex < 0 {
		lowerZindex = 0
		if lowerZindex > upperZindex {
			return mlrval.FromEmptyArray()
		}
	}
	if upperZindex > n-1 {
		upperZindex = n - 1
		if lowerZindex > upperZindex {
			return mlrval.FromEmptyArray()
		}
	}

	// Go     slices have inclusive lower bound, exclusive upper bound.
	// Miller slices have inclusive lower bound, inclusive upper bound.
	var m = upperZindex - lowerZindex + 1

	retval := make([]*mlrval.Mlrval, m)

	di := 0
	for si := lowerZindex; si <= upperZindex; si++ {
		retval[di] = array[si].Copy()
		di++
	}

	return mlrval.FromArray(retval)
}

// ================================================================
// For input record 'a=7,b=8,c=9',  $[[2]] = "b"

type PositionalFieldNameNode struct {
	indexEvaluable IEvaluable
}

func (node *RootNode) BuildPositionalFieldNameNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypePositionalFieldName)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)

	indexASTNode := astNode.Children[0]

	indexEvaluable, err := node.BuildEvaluableNode(indexASTNode)
	if err != nil {
		return nil, err
	}

	return &PositionalFieldNameNode{
		indexEvaluable: indexEvaluable,
	}, nil
}

// TODO: code-dedupe these next four Evaluate methods
func (node *PositionalFieldNameNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	indexMlrval := node.indexEvaluable.Evaluate(state)
	if indexMlrval.IsAbsent() {
		return mlrval.ABSENT
	}

	index, ok := indexMlrval.GetIntValue()
	if !ok {
		return mlrval.ERROR
	}

	name, ok := state.Inrec.GetNameAtPositionalIndex(index)
	if !ok {
		return mlrval.ABSENT
	}

	return mlrval.FromString(name)
}

// ================================================================
// For input record 'a=7,b=8,c=9',  $[[2]] = 8

type PositionalFieldValueNode struct {
	indexEvaluable IEvaluable
}

func (node *RootNode) BuildPositionalFieldValueNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypePositionalFieldValue)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)

	indexASTNode := astNode.Children[0]

	indexEvaluable, err := node.BuildEvaluableNode(indexASTNode)
	if err != nil {
		return nil, err
	}

	return &PositionalFieldValueNode{
		indexEvaluable: indexEvaluable,
	}, nil
}

func (node *PositionalFieldValueNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	indexMlrval := node.indexEvaluable.Evaluate(state)
	if indexMlrval.IsAbsent() {
		return mlrval.ABSENT
	}

	index, ok := indexMlrval.GetIntValue()
	if !ok {
		return mlrval.ERROR
	}

	retval := state.Inrec.GetWithPositionalIndex(index)
	if retval == nil {
		return mlrval.ABSENT
	}

	return retval
}

// ================================================================
// For x = [7,8,9], x[[2]] = 2
// For y = {"a":7,"b":8,"c":9}, y[[2]] = "b"
type ArrayOrMapPositionalNameAccessNode struct {
	baseEvaluable  IEvaluable
	indexEvaluable IEvaluable
}

func (node *RootNode) BuildArrayOrMapPositionalNameAccessNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeArrayOrMapPositionalNameAccess)
	lib.InternalCodingErrorIf(len(astNode.Children) != 2)

	baseASTNode := astNode.Children[0]
	indexASTNode := astNode.Children[1]

	baseEvaluable, err := node.BuildEvaluableNode(baseASTNode)
	if err != nil {
		return nil, err
	}
	indexEvaluable, err := node.BuildEvaluableNode(indexASTNode)
	if err != nil {
		return nil, err
	}

	return &ArrayOrMapPositionalNameAccessNode{
		baseEvaluable:  baseEvaluable,
		indexEvaluable: indexEvaluable,
	}, nil
}

func (node *ArrayOrMapPositionalNameAccessNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	baseMlrval := node.baseEvaluable.Evaluate(state)
	indexMlrval := node.indexEvaluable.Evaluate(state)

	if indexMlrval.IsAbsent() {
		return mlrval.ABSENT
	}

	index, ok := indexMlrval.GetIntValue()
	if !ok {
		return mlrval.ERROR
	}

	if baseMlrval.IsArray() {
		n, _ := baseMlrval.GetArrayLength()
		zindex, ok := mlrval.UnaliasArrayLengthIndex(int(n), int(index))
		if ok {
			return mlrval.FromInt(int64(zindex + 1)) // Miller user-space indices are 1-up
		} else {
			return mlrval.ABSENT
		}

	} else if baseMlrval.IsMap() {
		name, ok := baseMlrval.GetMap().GetNameAtPositionalIndex(index)
		if !ok {
			return mlrval.ABSENT
		} else {
			return mlrval.FromString(name)
		}

	} else if baseMlrval.IsAbsent() {
		return mlrval.ABSENT

	} else {
		return mlrval.ERROR
	}
}

// ================================================================
// For x = [7,8,9], x[[2]] = 8
// For y = {"a":7,"b":8,"c":9}, y[[2]] = 8
type ArrayOrMapPositionalValueAccessNode struct {
	baseEvaluable  IEvaluable
	indexEvaluable IEvaluable
}

func (node *RootNode) BuildArrayOrMapPositionalValueAccessNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeArrayOrMapPositionalValueAccess)
	lib.InternalCodingErrorIf(len(astNode.Children) != 2)

	baseASTNode := astNode.Children[0]
	indexASTNode := astNode.Children[1]

	baseEvaluable, err := node.BuildEvaluableNode(baseASTNode)
	if err != nil {
		return nil, err
	}
	indexEvaluable, err := node.BuildEvaluableNode(indexASTNode)
	if err != nil {
		return nil, err
	}

	return &ArrayOrMapPositionalValueAccessNode{
		baseEvaluable:  baseEvaluable,
		indexEvaluable: indexEvaluable,
	}, nil
}

func (node *ArrayOrMapPositionalValueAccessNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	baseMlrval := node.baseEvaluable.Evaluate(state)
	indexMlrval := node.indexEvaluable.Evaluate(state)

	if indexMlrval.IsAbsent() {
		return mlrval.ABSENT
	}

	index, ok := indexMlrval.GetIntValue()
	if !ok {
		return mlrval.ERROR
	}

	if baseMlrval.IsArray() {
		// xxx pending pointer-output refactor
		retval := baseMlrval.ArrayGet(indexMlrval)
		return &retval

	} else if baseMlrval.IsMap() {
		value := baseMlrval.GetMap().GetWithPositionalIndex(index)
		if value == nil {
			return mlrval.ABSENT
		}

		return value

	} else if baseMlrval.IsAbsent() {
		return mlrval.ABSENT

	} else {
		return mlrval.ERROR
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

func (node *RootNode) BuildMapLiteralNode(
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

		cstKey, err := node.BuildEvaluableNode(astKey)
		if err != nil {
			return nil, err
		}
		cstValue, err := node.BuildEvaluableNode(astValue)
		if err != nil {
			return nil, err
		}

		evaluablePairs[i] = NewEvaluablePair(cstKey, cstValue)
	}
	return &MapLiteralNode{
		evaluablePairs: evaluablePairs,
	}, nil
}

func (node *MapLiteralNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	output := mlrval.FromEmptyMap()

	for i := range node.evaluablePairs {
		key := node.evaluablePairs[i].Key.Evaluate(state)
		value := node.evaluablePairs[i].Value.Evaluate(state)

		if !value.IsAbsent() {
			output.MapPut(key, value)
		}
	}

	return output
}
