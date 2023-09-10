// ================================================================
// CST build/execute for AST array-literal, map-literal, index-access, and
// slice-access nodes
// ================================================================

package cst

import (
	"fmt"

	"github.com/johnkerl/miller/pkg/bifs"
	"github.com/johnkerl/miller/pkg/dsl"
	"github.com/johnkerl/miller/pkg/lib"
	"github.com/johnkerl/miller/pkg/mlrval"
	"github.com/johnkerl/miller/pkg/runtime"
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
type ArrayOrMapIndexAccessNode struct {
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

	return &ArrayOrMapIndexAccessNode{
		baseEvaluable:  baseEvaluable,
		indexEvaluable: indexEvaluable,
	}, nil
}

func (node *ArrayOrMapIndexAccessNode) Evaluate(
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
			return mlrval.FromError(
				fmt.Errorf(
					"unacceptable non-int index value %s of type %s on base value %s",
					indexMlrval.StringMaybeQuoted(),
					indexMlrval.GetTypeName(),
					baseMlrval.StringMaybeQuoted(),
				),
			)
		}
		// Handle UTF-8 correctly: len(input1.printrep) will count bytes, not runes.
		runes := []rune(baseMlrval.String())
		// Miller uses 1-up, and negatively aliased, indexing for strings and arrays.
		zindex, inBounds := mlrval.UnaliasArrayLengthIndex(len(runes), int(mindex))
		if !inBounds {
			return mlrval.FromError(
				fmt.Errorf(
					"cannot index base string %s of length %d with out-of-bounds index %d",
					baseMlrval.StringMaybeQuoted(),
					len(runes),
					int(mindex),
				),
			)
		}
		return mlrval.FromString(string(runes[zindex]))

	} else if baseMlrval.IsAbsent() {
		// For strict mode, absence should be detected on the baseMlrval and indexMlrval evaluators.
		return mlrval.ABSENT
	} else {
		return mlrval.FromError(
			fmt.Errorf(
				"cannot index base value %s of type %s, which is not array, map, or string",
				baseMlrval.StringMaybeQuoted(),
				baseMlrval.GetTypeName(),
			),
		)
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
		// For strict mode, absence should be detected on the baseMlrval and indexMlrval evaluators.
		return mlrval.ABSENT
	}
	if baseMlrval.IsString() {
		return bifs.BIF_substr_1_up(baseMlrval, lowerIndexMlrval, upperIndexMlrval)
	}
	array := baseMlrval.GetArray()
	if array == nil {
		return mlrval.FromError(
			fmt.Errorf(
				"cannot slice base value %s with non-array type %s",
				baseMlrval.StringMaybeQuoted(),
				baseMlrval.GetTypeName(),
			),
		)
	}
	n := len(array)

	sliceIsEmpty, absentOrError, lowerZindex, upperZindex :=
		bifs.MillerSliceAccess(lowerIndexMlrval, upperIndexMlrval, n, false)

	if sliceIsEmpty {
		return mlrval.FromEmptyArray()
	}
	if absentOrError != nil {
		return absentOrError
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
		return mlrval.ABSENT.StrictModeCheck(state.StrictMode, "$[[(absent)]]")
	}

	index, ok := indexMlrval.GetIntValue()
	if !ok {
		return mlrval.FromNotIntError("$[[...]]", indexMlrval)
	}

	name, ok := state.Inrec.GetNameAtPositionalIndex(index)
	if !ok {
		return mlrval.ABSENT.StrictModeCheck(state.StrictMode, "$[["+indexMlrval.String()+"]]")
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
		return mlrval.ABSENT.StrictModeCheck(state.StrictMode, "$[[[(absent)]]]")
	}

	index, ok := indexMlrval.GetIntValue()
	if !ok {
		return mlrval.FromNotIntError("$[[...]]", indexMlrval)
	}

	retval := state.Inrec.GetWithPositionalIndex(index)
	if retval == nil {
		return mlrval.ABSENT.StrictModeCheck(state.StrictMode, "$[[["+indexMlrval.String()+"]]]")
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
		// For strict mode, absence should be detected on the baseMlrval and indexMlrval evaluators.
		return mlrval.ABSENT
	}

	index, ok := indexMlrval.GetIntValue()
	if !ok {
		return mlrval.FromNotIntError("$[[...]]", indexMlrval)
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
		// For strict mode, absence should be detected on the baseMlrval and indexMlrval evaluators.
		return mlrval.ABSENT

	} else {
		return mlrval.FromError(
			fmt.Errorf(
				"cannot index base value %s of type %s, which is not array, map, or string",
				baseMlrval.StringMaybeQuoted(),
				baseMlrval.GetTypeName(),
			),
		)
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
		// For strict mode, absence should be detected on the baseMlrval and indexMlrval evaluators.
		return mlrval.ABSENT
	}

	index, ok := indexMlrval.GetIntValue()
	if !ok {
		return mlrval.FromNotIntError("$[[...]]", indexMlrval)
	}

	if baseMlrval.IsArray() {
		// xxx pending pointer-output refactor
		retval := baseMlrval.ArrayGet(indexMlrval)
		return &retval

	} else if baseMlrval.IsMap() {
		value := baseMlrval.GetMap().GetWithPositionalIndex(index)
		if value == nil {
			// For strict mode, absence should be detected on the baseMlrval and indexMlrval evaluators.
			return mlrval.ABSENT
		}

		return value

	} else if baseMlrval.IsAbsent() {
		// For strict mode, absence should be detected on the baseMlrval and indexMlrval evaluators.
		return mlrval.ABSENT

	} else {
		return mlrval.FromError(
			fmt.Errorf(
				"cannot index base value %s of type %s, which is not array, map, or string",
				baseMlrval.StringMaybeQuoted(),
				baseMlrval.GetTypeName(),
			),
		)
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
