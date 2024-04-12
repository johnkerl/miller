// ================================================================
// Support for higher-order functions in Miller: select, apply, fold, reduce,
// sort, any, and every.
// ================================================================

package cst

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/facette/natsort"

	"github.com/johnkerl/miller/pkg/lib"
	"github.com/johnkerl/miller/pkg/mlrval"
	"github.com/johnkerl/miller/pkg/runtime"
)

// Most function types are in the github.com/johnkerl/miller/pkg/types package. These types, though,
// include functions which need to access CST state in order to call back to
// user-defined functions.  To avoid a package-cycle dependency, they are
// defined here.

// BinaryFuncWithState is for select, apply, and reduce.
type BinaryFuncWithState func(
	input1 *mlrval.Mlrval,
	input2 *mlrval.Mlrval,
	state *runtime.State,
) *mlrval.Mlrval

// TernaryFuncWithState is for fold.
type TernaryFuncWithState func(
	input1 *mlrval.Mlrval,
	input2 *mlrval.Mlrval,
	input3 *mlrval.Mlrval,
	state *runtime.State,
) *mlrval.Mlrval

// VariadicFuncWithState is for sort.
type VariadicFuncWithState func(
	inputs []*mlrval.Mlrval,
	state *runtime.State,
) *mlrval.Mlrval

// tHOFSpace is the datatype for the getHOFSpace cache-manager.
type tHOFSpace struct {
	udfCallsite *UDFCallsite
	argsArray   []*mlrval.Mlrval
}

// hofCache is persistent data for the getHOFSpace cache-manager.
var hofCache map[string]*tHOFSpace = make(map[string]*tHOFSpace)

// getHOFSpace manages a cache for the data needed by higher-order functions.
// Those functions may be invoked on every record of a big data file, so we try
// to cache data they need for UDF-callsite setup.
func getHOFSpace(
	funcVal *mlrval.Mlrval,
	arity int,
	hofName string,
	arrayOrMap string,
) *tHOFSpace {
	// At this callsite, localvars have been evaluated already -- so for 'y =
	// sort(x, f)' we have the *value* of f, not that variable name -- the
	// value will be a Mlrval of function type.
	//
	// If the func-type mlrval as the second argument points to a named UDF,
	// then funcVal.String() is its name -- e.g. in 'func f(a, b) { return b
	// <=> a }' and 'y = sort(x, f)', funcVal.String is "f".
	//
	// If the func-type mlrval as the second argument points to an unnamed UDF,
	// then funcVal.String() is a UUID -- e.g. in 'y = sortaf(x, func f(a, b) {
	// return b <=> a })', funcVal.String is something like
	// "function-literal-000052".
	udfName := funcVal.String()

	cacheKey := udfName + ":" + hofName + ":" + arrayOrMap

	// Cache hit, but check arity. Example: someone makes a correct-arity
	// callback for arrays, then re-uses it for maps where the arity needs to
	// be different. E.g.  apply([...], f(e) {...}) vs apply({...}, f(k,v) {...})
	entry := hofCache[cacheKey]
	if entry != nil {
		if entry.udfCallsite.arity != arity {
			fmt.Fprintf(
				os.Stderr,
				"mlr: %s: argument function \"%s\" has arity %d; needed %d for %s.\n",
				hofName,
				udfName,
				entry.udfCallsite.arity,
				arity,
				arrayOrMap,
			)
			os.Exit(1)
		}
		return entry
	}

	// Cache miss
	var udf *UDF = nil
	iUDF := funcVal.GetFunction()
	if iUDF == nil { // E.g. does not exist at all
		fmt.Fprintf(os.Stderr, "mlr: %s: argument function \"%s\" not found.\n", hofName, udfName)
		os.Exit(1)
	}
	udf = iUDF.(*UDF)

	if udf.signature.arity != arity { // Present, but with the wrong arity.
		fmt.Fprintf(
			os.Stderr,
			"mlr: %s: argument function \"%s\" has arity %d; needed %d for %s.\n",
			hofName,
			udfName,
			udf.signature.arity,
			arity,
			arrayOrMap,
		)
		os.Exit(1)
	}

	udfCallsite := NewUDFCallsiteForHigherOrderFunction(udf, arity)
	argsArray := make([]*mlrval.Mlrval, arity)
	entry = &tHOFSpace{
		udfCallsite: udfCallsite,
		argsArray:   argsArray,
	}
	// Remember for subsequent cache hit.
	hofCache[cacheKey] = entry
	return entry
}

// mustBeNonAbsent checks that a UDF for array reduce/fold/apply returned a value.
func isNonAbsentOrDie(mlrval *mlrval.Mlrval, hofName string) *mlrval.Mlrval {
	if mlrval.IsAbsent() {
		hofCheckDie(mlrval, hofName, "second-argument function must return a value")
	}
	return mlrval
}

// getKVPairForAccumulatorOrDie checks that a user-supplied accumulator value
// for a map fold is indeed a single-element map.
func getKVPairForAccumulatorOrDie(mlrval *mlrval.Mlrval, hofName string) *mlrval.Mlrmap {
	kvPair := getKVPair(mlrval)
	if kvPair == nil {
		hofCheckDie(mlrval, hofName, "accumulator value must be a single-element map")
	}
	return kvPair
}

// getKVPairForCallbackOrDie checks that a return value from a UDF for map
// reduce/fold/apply is indeed a single-element map.
func getKVPairForCallbackOrDie(mlrval *mlrval.Mlrval, hofName string) *mlrval.Mlrmap {
	kvPair := getKVPair(mlrval)
	if kvPair == nil {
		hofCheckDie(mlrval, hofName, "second-argument function must return single-element map")
	}
	return kvPair
}

// hofCheckDie is a helper function for HOFs on maps, to check that the
// user-supplied UDF returned a single-entry map.
func hofCheckDie(mlrval *mlrval.Mlrval, hofName string, message string) {
	fmt.Fprintf(
		os.Stderr,
		"mlr: %s: %s; got \"%s\".\n",
		hofName,
		message,
		mlrval.String(),
	)
	os.Exit(1)
}

// getKVPair is a helper function getKVPairOrDie.
func getKVPair(mlrval *mlrval.Mlrval) *mlrval.Mlrmap {
	mapval := mlrval.GetMap()
	if mapval == nil {
		return nil
	}
	if mapval == nil || mapval.FieldCount != 1 {
		return nil
	}
	return mapval
}

func isFunctionOrDie(mlrval *mlrval.Mlrval, hofName string) {
	if !mlrval.IsFunction() {
		fmt.Fprintf(os.Stderr, "mlr: %s: second argument must be a function; got %s.\n",
			hofName, mlrval.GetTypeName(),
		)
		os.Exit(1)
	}
}

// ================================================================
// SELECT HOF

func SelectHOF(
	input1 *mlrval.Mlrval,
	input2 *mlrval.Mlrval,
	state *runtime.State,
) *mlrval.Mlrval {
	if input1.IsArray() {
		return selectArray(input1, input2, state)
	} else if input1.IsMap() {
		return selectMap(input1, input2, state)
	} else {
		return mlrval.FromNotCollectionError("select", input1)
	}
}

func selectArray(
	input1 *mlrval.Mlrval,
	input2 *mlrval.Mlrval,
	state *runtime.State,
) *mlrval.Mlrval {
	inputArray, errVal := input1.GetArrayValueOrError("select")
	if inputArray == nil { // not an array
		return errVal
	}
	isFunctionOrDie(input2, "select")

	hofSpace := getHOFSpace(input2, 1, "select", "array")
	udfCallsite := hofSpace.udfCallsite
	argsArray := hofSpace.argsArray

	outputArray := make([]*mlrval.Mlrval, 0, len(inputArray))

	for i := range inputArray {
		argsArray[0] = inputArray[i]
		mret := udfCallsite.EvaluateWithArguments(state, udfCallsite.udf, argsArray)
		bret, ok := mret.GetBoolValue()
		if !ok {
			fmt.Fprintf(
				os.Stderr,
				"mlr: select: function returned non-boolean \"%s\".\n",
				mret.String(),
			)
			os.Exit(1)
		}
		if bret {
			outputArray = append(outputArray, inputArray[i].Copy())
		}
	}
	return mlrval.FromArray(outputArray)
}

func selectMap(
	input1 *mlrval.Mlrval,
	input2 *mlrval.Mlrval,
	state *runtime.State,
) *mlrval.Mlrval {
	inputMap, errVal := input1.GetMapValueOrError("select")
	if inputMap == nil { // not a map
		return errVal
	}
	isFunctionOrDie(input2, "select")

	hofSpace := getHOFSpace(input2, 2, "select", "map")
	udfCallsite := hofSpace.udfCallsite
	argsArray := hofSpace.argsArray

	outputMap := mlrval.NewMlrmap()

	for pe := inputMap.Head; pe != nil; pe = pe.Next {
		argsArray[0] = mlrval.FromString(pe.Key)
		argsArray[1] = pe.Value
		mret := udfCallsite.EvaluateWithArguments(state, udfCallsite.udf, argsArray)
		bret, ok := mret.GetBoolValue()
		if !ok {
			fmt.Fprintf(
				os.Stderr,
				"mlr: select: function returned non-boolean \"%s\".\n",
				mret.String(),
			)
			os.Exit(1)
		}
		if bret {
			outputMap.PutCopy(pe.Key, pe.Value)
		}
	}

	return mlrval.FromMap(outputMap)
}

// ================================================================
// APPLY HOF

func ApplyHOF(
	input1 *mlrval.Mlrval,
	input2 *mlrval.Mlrval,
	state *runtime.State,
) *mlrval.Mlrval {
	if input1.IsArray() {
		return applyArray(input1, input2, state)
	} else if input1.IsMap() {
		return applyMap(input1, input2, state)
	} else {
		return mlrval.FromNotCollectionError("apply", input1)
	}
}

func applyArray(
	input1 *mlrval.Mlrval,
	input2 *mlrval.Mlrval,
	state *runtime.State,
) *mlrval.Mlrval {
	inputArray, errVal := input1.GetArrayValueOrError("apply")
	if inputArray == nil {
		return errVal
	}
	isFunctionOrDie(input2, "apply")

	hofSpace := getHOFSpace(input2, 1, "apply", "array")
	udfCallsite := hofSpace.udfCallsite
	argsArray := hofSpace.argsArray

	outputArray := make([]*mlrval.Mlrval, len(inputArray))

	for i := range inputArray {
		argsArray[0] = inputArray[i]
		outputArray[i] = isNonAbsentOrDie(
			(udfCallsite.EvaluateWithArguments(state, udfCallsite.udf, argsArray)),
			"apply",
		)
	}
	return mlrval.FromArray(outputArray)
}

func applyMap(
	input1 *mlrval.Mlrval,
	input2 *mlrval.Mlrval,
	state *runtime.State,
) *mlrval.Mlrval {
	inputMap, errVal := input1.GetMapValueOrError("apply")
	if inputMap == nil { // not a map
		return errVal
	}
	isFunctionOrDie(input2, "apply")

	hofSpace := getHOFSpace(input2, 2, "apply", "map")
	udfCallsite := hofSpace.udfCallsite
	argsArray := hofSpace.argsArray

	outputMap := mlrval.NewMlrmap()

	for pe := inputMap.Head; pe != nil; pe = pe.Next {
		argsArray[0] = mlrval.FromString(pe.Key)
		argsArray[1] = pe.Value
		retval := udfCallsite.EvaluateWithArguments(state, udfCallsite.udf, argsArray)
		kvPair := getKVPairForCallbackOrDie(retval, "apply")
		outputMap.PutReference(kvPair.Head.Key, kvPair.Head.Value)
	}
	return mlrval.FromMap(outputMap)
}

// ================================================================
// REDUCE HOF

func ReduceHOF(
	input1 *mlrval.Mlrval,
	input2 *mlrval.Mlrval,
	state *runtime.State,
) *mlrval.Mlrval {
	if input1.IsArray() {
		return reduceArray(input1, input2, state)
	} else if input1.IsMap() {
		return reduceMap(input1, input2, state)
	} else {
		return mlrval.FromNotCollectionError("reduce", input1)
	}
}

func reduceArray(
	input1 *mlrval.Mlrval,
	input2 *mlrval.Mlrval,
	state *runtime.State,
) *mlrval.Mlrval {
	inputArray, errVal := input1.GetArrayValueOrError("reduce")
	if inputArray == nil {
		return errVal
	}
	isFunctionOrDie(input2, "reduce")

	hofSpace := getHOFSpace(input2, 2, "reduce", "array")
	udfCallsite := hofSpace.udfCallsite
	argsArray := hofSpace.argsArray

	n := len(inputArray)
	if n == 0 {
		return input1
	}
	accumulator := inputArray[0].Copy()

	for i := 1; i < n; i++ {
		argsArray[0] = accumulator
		argsArray[1] = inputArray[i]
		accumulator = (udfCallsite.EvaluateWithArguments(state, udfCallsite.udf, argsArray))
		isNonAbsentOrDie(accumulator, "apply")
	}
	return accumulator
}

func reduceMap(
	input1 *mlrval.Mlrval,
	input2 *mlrval.Mlrval,
	state *runtime.State,
) *mlrval.Mlrval {
	inputMap, errVal := input1.GetMapValueOrError("reduce")
	if inputMap == nil { // not a map
		return errVal
	}
	isFunctionOrDie(input2, "reduce")

	hofSpace := getHOFSpace(input2, 4, "reduce", "map")
	udfCallsite := hofSpace.udfCallsite
	argsArray := hofSpace.argsArray

	accumulator := inputMap.GetFirstPair()
	if accumulator == nil { // Input map is empty
		return input1
	}

	for pe := inputMap.Head.Next; pe != nil; pe = pe.Next {
		argsArray[0] = mlrval.FromString(accumulator.Head.Key)
		argsArray[1] = accumulator.Head.Value
		argsArray[2] = mlrval.FromString(pe.Key)
		argsArray[3] = pe.Value.Copy()
		retval := (udfCallsite.EvaluateWithArguments(state, udfCallsite.udf, argsArray))
		kvPair := getKVPairForCallbackOrDie(retval, "reduce")
		accumulator = kvPair
	}
	return mlrval.FromMap(accumulator)
}

// ================================================================
// FOLD HOF

func FoldHOF(
	input1 *mlrval.Mlrval,
	input2 *mlrval.Mlrval,
	input3 *mlrval.Mlrval,
	state *runtime.State,
) *mlrval.Mlrval {
	if input1.IsArray() {
		return foldArray(input1, input2, input3, state)
	} else if input1.IsMap() {
		return foldMap(input1, input2, input3, state)
	} else {
		return mlrval.FromNotCollectionError("fold", input1)
	}
}

func foldArray(
	input1 *mlrval.Mlrval,
	input2 *mlrval.Mlrval,
	input3 *mlrval.Mlrval,
	state *runtime.State,
) *mlrval.Mlrval {
	inputArray, errVal := input1.GetArrayValueOrError("fold")
	if inputArray == nil {
		return errVal
	}
	isFunctionOrDie(input2, "fold")

	hofSpace := getHOFSpace(input2, 2, "fold", "array")
	udfCallsite := hofSpace.udfCallsite
	argsArray := hofSpace.argsArray

	accumulator := input3.Copy()

	for i := range inputArray {
		argsArray[0] = accumulator
		argsArray[1] = inputArray[i]
		accumulator = (udfCallsite.EvaluateWithArguments(state, udfCallsite.udf, argsArray))
		isNonAbsentOrDie(accumulator, "apply")
	}
	return accumulator
}

func foldMap(
	input1 *mlrval.Mlrval,
	input2 *mlrval.Mlrval,
	input3 *mlrval.Mlrval,
	state *runtime.State,
) *mlrval.Mlrval {
	inputMap, errVal := input1.GetMapValueOrError("fold")
	if inputMap == nil { // not a map
		return errVal
	}
	isFunctionOrDie(input2, "fold")

	hofSpace := getHOFSpace(input2, 4, "fold", "map")
	udfCallsite := hofSpace.udfCallsite
	argsArray := hofSpace.argsArray

	if inputMap.IsEmpty() {
		return mlrval.ABSENT
	}

	accumulator := getKVPairForAccumulatorOrDie(input3, "reduce").Copy()

	for pe := inputMap.Head; pe != nil; pe = pe.Next {
		argsArray[0] = mlrval.FromString(accumulator.Head.Key)
		argsArray[1] = accumulator.Head.Value
		argsArray[2] = mlrval.FromString(pe.Key)
		argsArray[3] = pe.Value.Copy()
		retval := (udfCallsite.EvaluateWithArguments(state, udfCallsite.udf, argsArray))
		kvPair := getKVPairForCallbackOrDie(retval, "reduce")
		accumulator = kvPair
	}
	return mlrval.FromMap(accumulator)
}

// ================================================================
// SORT HOF

func SortHOF(
	inputs []*mlrval.Mlrval,
	state *runtime.State,
) *mlrval.Mlrval {

	if len(inputs) == 1 {
		if inputs[0].IsArray() {
			return sortA(inputs[0], "")
		} else if inputs[0].IsMap() {
			return sortM(inputs[0], "")
		} else {
			return mlrval.FromNotCollectionError("sort", inputs[0])
		}

	} else if inputs[1].IsStringOrVoid() {
		if inputs[0].IsArray() {
			return sortA(inputs[0], inputs[1].String())
		} else if inputs[0].IsMap() {
			return sortM(inputs[0], inputs[1].String())
		} else {
			return mlrval.FromNotCollectionError("sort", inputs[0])
		}

	} else if inputs[1].IsFunction() {
		if inputs[0].IsArray() {
			return sortAF(inputs[0], inputs[1], state)
		} else if inputs[0].IsMap() {
			return sortMF(inputs[0], inputs[1], state)
		} else {
			return mlrval.FromNotCollectionError("sort", inputs[0])
		}

	} else {
		fmt.Fprintf(os.Stderr, "mlr: sort: second argument must be a string or function; got %s.\n",
			inputs[1].GetTypeName(),
		)
		os.Exit(1)
	}
	// Not reached
	lib.InternalCodingErrorIf(true)
	return nil
}

// ----------------------------------------------------------------
// Helpers for sort with string flags in place of callback UDF.

type tSortType int

const (
	sortTypeLexical tSortType = iota
	sortTypeCaseFold
	sortTypeNumerical
	sortTypeNatural
)

// decodeSortFlags maps strings like "cr" in the second argument to sort
// into sortType=sortTypeCaseFold and reverse=true, etc.
func decodeSortFlags(flags string) (sortType tSortType, reverse bool, byMapValue bool) {
	sortType = sortTypeNumerical
	reverse = false
	byMapValue = false
	for _, c := range flags {
		switch c {
		case 'n':
			sortType = sortTypeNumerical
		case 'f':
			sortType = sortTypeLexical
		case 'c':
			sortType = sortTypeCaseFold
		case 't':
			sortType = sortTypeNatural
		case 'r':
			reverse = true
		case 'v':
			byMapValue = true
		}
	}
	return sortType, reverse, byMapValue
}

// sortA implements sort on array, with string flags rather than callback UDF.
func sortA(
	input1 *mlrval.Mlrval,
	flags string,
) *mlrval.Mlrval {
	temp, errVal := input1.GetArrayValueOrError("sort")
	if temp == nil { // not an array
		return errVal
	}
	output := input1.Copy()

	// byMapValue is ignored for sorting arrays
	sortType, reverse, _ := decodeSortFlags(flags)

	a := output.GetArray()
	switch sortType {
	case sortTypeNumerical:
		sortANumerical(a, reverse)
	case sortTypeLexical:
		sortALexical(a, reverse)
	case sortTypeCaseFold:
		sortACaseFold(a, reverse)
	case sortTypeNatural:
		sortANatural(a, reverse)
	}

	return output
}

func sortANumerical(array []*mlrval.Mlrval, reverse bool) {
	if !reverse {
		sort.Slice(array, func(i, j int) bool {
			return mlrval.LessThan(array[i], array[j])
		})
	} else {
		sort.Slice(array, func(i, j int) bool {
			return mlrval.GreaterThan(array[i], array[j])
		})
	}
}

func sortALexical(array []*mlrval.Mlrval, reverse bool) {
	if !reverse {
		sort.Slice(array, func(i, j int) bool {
			return array[i].String() < array[j].String()
		})
	} else {
		sort.Slice(array, func(i, j int) bool {
			return array[i].String() > array[j].String()
		})
	}
}

func sortACaseFold(array []*mlrval.Mlrval, reverse bool) {
	if !reverse {
		sort.Slice(array, func(i, j int) bool {
			return strings.ToLower(array[i].String()) < strings.ToLower(array[j].String())
		})
	} else {
		sort.Slice(array, func(i, j int) bool {
			return strings.ToLower(array[i].String()) > strings.ToLower(array[j].String())
		})
	}
}

func sortANatural(array []*mlrval.Mlrval, reverse bool) {
	if !reverse {
		sort.Slice(array, func(i, j int) bool {
			return natsort.Compare(array[i].String(), array[j].String())
		})
	} else {
		sort.Slice(array, func(i, j int) bool {
			return natsort.Compare(array[j].String(), array[i].String())
		})
	}
}

// sortM implements sort on map, with string flags rather than callback UDF.
func sortM(
	input1 *mlrval.Mlrval,
	flags string,
) *mlrval.Mlrval {
	inputMap, errVal := input1.GetMapValueOrError("sort")
	if inputMap == nil { // not a map
		return errVal
	}

	// Get sort-flags, if provided
	sortType, reverse, byMapValue := decodeSortFlags(flags)

	// Copy the entries to an array for sorting.
	n := inputMap.FieldCount
	entries := make([]mlrval.MlrmapEntryForArray, n)
	i := 0
	for pe := inputMap.Head; pe != nil; pe = pe.Next {
		entries[i].Key = pe.Key
		entries[i].Value = pe.Value // pointer alias for now until new map at end of this function
		i++
	}

	// Do the key-sort
	switch sortType {
	case sortTypeNumerical:
		sortMNumerical(entries, reverse, byMapValue)
	case sortTypeLexical:
		sortMLexical(entries, reverse, byMapValue)
	case sortTypeCaseFold:
		sortMCaseFold(entries, reverse, byMapValue)
	case sortTypeNatural:
		sortMNatural(entries, reverse, byMapValue)
	}

	// Make a new map with entries in the new sort order.
	outmap := mlrval.NewMlrmap()
	for i := int64(0); i < n; i++ {
		entry := entries[i]
		outmap.PutCopy(entry.Key, entry.Value)
	}

	return mlrval.FromMap(outmap)
}

func sortMNumerical(array []mlrval.MlrmapEntryForArray, reverse bool, byMapValue bool) {
	if !byMapValue {
		if !reverse {
			sort.Slice(array, func(i, j int) bool {
				na, erra := strconv.ParseFloat(array[i].Key, 64)
				nb, errb := strconv.ParseFloat(array[j].Key, 64)
				if erra == nil && errb == nil {
					return na < nb
				} else {
					return array[i].Key < array[j].Key
				}
			})
		} else {
			sort.Slice(array, func(i, j int) bool {
				na, erra := strconv.ParseFloat(array[i].Key, 64)
				nb, errb := strconv.ParseFloat(array[j].Key, 64)
				if erra == nil && errb == nil {
					return na > nb
				} else {
					return array[i].Key > array[j].Key
				}
			})
		}
	} else {
		if !reverse {
			sort.Slice(array, func(i, j int) bool {
				return mlrval.LessThan(array[i].Value, array[j].Value)
			})
		} else {
			sort.Slice(array, func(i, j int) bool {
				return mlrval.LessThan(array[j].Value, array[i].Value)
			})
		}
	}
}

func sortMLexical(array []mlrval.MlrmapEntryForArray, reverse bool, byMapValue bool) {
	if !byMapValue {
		if !reverse {
			// Or sort.Strings(keys) would work here as well.
			sort.Slice(array, func(i, j int) bool {
				return array[i].Key < array[j].Key
			})
		} else {
			sort.Slice(array, func(i, j int) bool {
				return array[i].Key > array[j].Key
			})
		}
	} else {
		if !reverse {
			sort.Slice(array, func(i, j int) bool {
				return array[i].Value.String() < array[j].Value.String()
			})
		} else {
			sort.Slice(array, func(i, j int) bool {
				return array[i].Value.String() > array[j].Value.String()
			})
		}
	}
}

func sortMCaseFold(array []mlrval.MlrmapEntryForArray, reverse bool, byMapValue bool) {
	if !byMapValue {
		if !reverse {
			sort.Slice(array, func(i, j int) bool {
				return strings.ToLower(array[i].Key) < strings.ToLower(array[j].Key)
			})
		} else {
			sort.Slice(array, func(i, j int) bool {
				return strings.ToLower(array[i].Key) > strings.ToLower(array[j].Key)
			})
		}
	} else {
		if !reverse {
			sort.Slice(array, func(i, j int) bool {
				return strings.ToLower(array[i].Value.String()) < strings.ToLower(array[j].Value.String())
			})
		} else {
			sort.Slice(array, func(i, j int) bool {
				return strings.ToLower(array[i].Value.String()) > strings.ToLower(array[j].Value.String())
			})
		}
	}
}

func sortMNatural(array []mlrval.MlrmapEntryForArray, reverse bool, byMapValue bool) {
	if !byMapValue {
		if !reverse {
			sort.Slice(array, func(i, j int) bool {
				return natsort.Compare(strings.ToLower(array[i].Key), strings.ToLower(array[j].Key))
			})
		} else {
			sort.Slice(array, func(i, j int) bool {
				return natsort.Compare(strings.ToLower(array[j].Key), strings.ToLower(array[i].Key))
			})
		}
	} else {
		if !reverse {
			sort.Slice(array, func(i, j int) bool {
				return natsort.Compare(
					strings.ToLower(array[i].Value.String()),
					strings.ToLower(array[j].Value.String()),
				)
			})
		} else {
			sort.Slice(array, func(i, j int) bool {
				return natsort.Compare(
					strings.ToLower(array[j].Value.String()),
					strings.ToLower(array[i].Value.String()),
				)
			})
		}
	}
}

// sortAF implements sort on arrays with callback UDF.
func sortAF(
	input1 *mlrval.Mlrval,
	input2 *mlrval.Mlrval,
	state *runtime.State,
) *mlrval.Mlrval {
	inputArray, errVal := input1.GetArrayValueOrError("select")
	if inputArray == nil { // not an array
		return errVal
	}
	isFunctionOrDie(input2, "sort")

	hofSpace := getHOFSpace(input2, 2, "sort", "array")
	udfCallsite := hofSpace.udfCallsite
	argsArray := hofSpace.argsArray

	outputArray := mlrval.CopyMlrvalArray(inputArray)

	sort.Slice(outputArray, func(i, j int) bool {
		argsArray[0] = outputArray[i]
		argsArray[1] = outputArray[j]
		// Call the user's comparator function.
		mret := udfCallsite.EvaluateWithArguments(state, udfCallsite.udf, argsArray)
		// Unpack the mlrval.Mlrval return value into a number.
		nret, ok := mret.GetNumericToFloatValue()
		if !ok {
			fmt.Fprintf(
				os.Stderr,
				"mlr: sort: comparator function \"%s\" returned non-number \"%s\".\n",
				input2.String(),
				mret.String(),
			)
			os.Exit(1)
		}
		lib.InternalCodingErrorIf(!ok)
		// Go sort-callback conventions: true if a < b, false otherwise.
		return nret < 0
	})
	return mlrval.FromArray(outputArray)
}

// sortMF implements sort on arrays with callback UDF.
func sortMF(
	input1 *mlrval.Mlrval,
	input2 *mlrval.Mlrval,
	state *runtime.State,
) *mlrval.Mlrval {
	inputMap, errVal := input1.GetMapValueOrError("sort")
	if inputMap == nil { // not a map
		return errVal
	}
	isFunctionOrDie(input2, "sort")

	pairsArray := inputMap.ToPairsArray()

	hofSpace := getHOFSpace(input2, 4, "sort", "map")
	udfCallsite := hofSpace.udfCallsite
	argsArray := hofSpace.argsArray

	sort.Slice(pairsArray, func(i, j int) bool {
		argsArray[0] = mlrval.FromString(pairsArray[i].Key)
		argsArray[1] = pairsArray[i].Value
		argsArray[2] = mlrval.FromString(pairsArray[j].Key)
		argsArray[3] = pairsArray[j].Value

		// Call the user's comparator function.
		mret := udfCallsite.EvaluateWithArguments(state, udfCallsite.udf, argsArray)
		// Unpack the mlrval.Mlrval return value into a number.
		nret, ok := mret.GetNumericToFloatValue()
		if !ok {
			fmt.Fprintf(
				os.Stderr,
				"mlr: sort: comparator function \"%s\" returned non-number \"%s\".\n",
				input2.String(),
				mret.String(),
			)
			os.Exit(1)
		}
		lib.InternalCodingErrorIf(!ok)
		// Go sort-callback conventions: true if a < b, false otherwise.
		return nret < 0
	})

	sortedMap := mlrval.MlrmapFromPairsArray(pairsArray)
	return mlrval.FromMap(sortedMap)
}

// ================================================================
// ANY HOF

func AnyHOF(
	input1 *mlrval.Mlrval,
	input2 *mlrval.Mlrval,
	state *runtime.State,
) *mlrval.Mlrval {
	if input1.IsArray() {
		return anyArray(input1, input2, state)
	} else if input1.IsMap() {
		return anyMap(input1, input2, state)
	} else {
		return mlrval.FromNotCollectionError("any", input1)
	}
}

func anyArray(
	input1 *mlrval.Mlrval,
	input2 *mlrval.Mlrval,
	state *runtime.State,
) *mlrval.Mlrval {
	inputArray, errVal := input1.GetArrayValueOrError("any")
	if inputArray == nil { // not an array
		return errVal
	}
	isFunctionOrDie(input2, "any")

	hofSpace := getHOFSpace(input2, 1, "any", "array")
	udfCallsite := hofSpace.udfCallsite
	argsArray := hofSpace.argsArray

	boolAny := false
	for i := range inputArray {
		argsArray[0] = inputArray[i]
		mret := udfCallsite.EvaluateWithArguments(state, udfCallsite.udf, argsArray)
		bret, ok := mret.GetBoolValue()
		if !ok {
			fmt.Fprintf(
				os.Stderr,
				"mlr: any: function returned non-boolean \"%s\".\n",
				mret.String(),
			)
			os.Exit(1)
		}
		if bret {
			boolAny = true
			break
		}
	}
	return mlrval.FromBool(boolAny)
}

func anyMap(
	input1 *mlrval.Mlrval,
	input2 *mlrval.Mlrval,
	state *runtime.State,
) *mlrval.Mlrval {
	inputMap, errVal := input1.GetMapValueOrError("any")
	if inputMap == nil { // not a map
		return errVal
	}
	isFunctionOrDie(input2, "any")

	hofSpace := getHOFSpace(input2, 2, "any", "map")
	udfCallsite := hofSpace.udfCallsite
	argsArray := hofSpace.argsArray

	boolAny := false

	for pe := inputMap.Head; pe != nil; pe = pe.Next {
		argsArray[0] = mlrval.FromString(pe.Key)
		argsArray[1] = pe.Value
		mret := udfCallsite.EvaluateWithArguments(state, udfCallsite.udf, argsArray)
		bret, ok := mret.GetBoolValue()
		if !ok {
			fmt.Fprintf(
				os.Stderr,
				"mlr: any: function returned non-boolean \"%s\".\n",
				mret.String(),
			)
			os.Exit(1)
		}
		if bret {
			boolAny = true
			break
		}
	}

	return mlrval.FromBool(boolAny)
}

// ================================================================
// EVERY HOF

func EveryHOF(
	input1 *mlrval.Mlrval,
	input2 *mlrval.Mlrval,
	state *runtime.State,
) *mlrval.Mlrval {
	if input1.IsArray() {
		return everyArray(input1, input2, state)
	} else if input1.IsMap() {
		return everyMap(input1, input2, state)
	} else {
		return mlrval.FromNotCollectionError("every", input1)
	}
}

func everyArray(
	input1 *mlrval.Mlrval,
	input2 *mlrval.Mlrval,
	state *runtime.State,
) *mlrval.Mlrval {
	inputArray, errVal := input1.GetArrayValueOrError("every")
	if inputArray == nil { // not an array
		return errVal
	}
	isFunctionOrDie(input2, "every")

	hofSpace := getHOFSpace(input2, 1, "every", "array")
	udfCallsite := hofSpace.udfCallsite
	argsArray := hofSpace.argsArray

	boolEvery := true
	for i := range inputArray {
		argsArray[0] = inputArray[i]
		mret := udfCallsite.EvaluateWithArguments(state, udfCallsite.udf, argsArray)
		bret, ok := mret.GetBoolValue()
		if !ok {
			fmt.Fprintf(
				os.Stderr,
				"mlr: every: function returned non-boolean \"%s\".\n",
				mret.String(),
			)
			os.Exit(1)
		}
		if !bret {
			boolEvery = false
			break
		}
	}
	return mlrval.FromBool(boolEvery)
}

func everyMap(
	input1 *mlrval.Mlrval,
	input2 *mlrval.Mlrval,
	state *runtime.State,
) *mlrval.Mlrval {
	inputMap, errVal := input1.GetMapValueOrError("every")
	if inputMap == nil { // not a map
		return errVal
	}
	isFunctionOrDie(input2, "every")

	hofSpace := getHOFSpace(input2, 2, "every", "map")
	udfCallsite := hofSpace.udfCallsite
	argsArray := hofSpace.argsArray

	boolEvery := true

	for pe := inputMap.Head; pe != nil; pe = pe.Next {
		argsArray[0] = mlrval.FromString(pe.Key)
		argsArray[1] = pe.Value
		mret := udfCallsite.EvaluateWithArguments(state, udfCallsite.udf, argsArray)
		bret, ok := mret.GetBoolValue()
		if !ok {
			fmt.Fprintf(
				os.Stderr,
				"mlr: every: function returned non-boolean \"%s\".\n",
				mret.String(),
			)
			os.Exit(1)
		}
		if !bret {
			boolEvery = false
			break
		}
	}

	return mlrval.FromBool(boolEvery)
}
