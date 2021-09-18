// ================================================================
// CST build/execute for AST operator/function nodes.
//
// Operators and functions are semantically the same thing -- they differ only
// syntactically. Binary operators are infix, like '1+2', while functions are
// prefix, like 'max(1,2)'. Both parse to the same AST shape.
// ================================================================

package cst

import (
	"fmt"
	"os"
	"sort"

	"mlr/src/dsl"
	"mlr/src/lib"
	"mlr/src/runtime"
	"mlr/src/types"
)

// ----------------------------------------------------------------
// BinaryFuncWithState is for the sortaf and sortmf functions.  Most function
// types are in the mlr/src/types packae. This type, though, includes functions
// which need to access CST state in order to call back to user-defined
// functions.  To avoid a package-cycle dependency, they are defined here.
type BinaryFuncWithState func(
	input1 *types.Mlrval,
	input2 *types.Mlrval,
	state *runtime.State,
	udfManager *UDFManager,
) *types.Mlrval

// ----------------------------------------------------------------
// Function lookup:
//
// * Try builtins first
// * Absent a match there, try UDF lookup (i.e. the UDF has been defined before being called)
// * Absent a match there:
//   o Make a UDF-placeholder node with present signature but nil function-pointer
//   o Append that node to CST to-be-resolved list
//   o On a next pass, we will walk that list resolving against all encountered
//     UDF definitions. (It will be an error then if it's still unresolvable.)

func (root *RootNode) BuildFunctionCallsiteNode(astNode *dsl.ASTNode) (IEvaluable, error) {
	lib.InternalCodingErrorIf(
		astNode.Type != dsl.NodeTypeFunctionCallsite &&
			astNode.Type != dsl.NodeTypeOperator,
	)
	lib.InternalCodingErrorIf(astNode.Token == nil)
	lib.InternalCodingErrorIf(astNode.Children == nil)

	functionName := string(astNode.Token.Lit)

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Special-case the dot operator, which is:
	// * string + string, with coercion to string if either side is int/float/bool/etc.;
	// * map attribute access, if the left-hand side is a map.

	if functionName == "." {
		dotCallsiteNode, err := root.BuildDotCallsiteNode(astNode)
		if err != nil {
			return nil, err
		}
		return dotCallsiteNode, nil
	}

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Look for a builtin function with the given name.

	builtinFunctionCallsiteNode, err := root.BuildBuiltinFunctionCallsiteNode(astNode)
	if err != nil {
		return nil, err
	}
	if builtinFunctionCallsiteNode != nil {
		return builtinFunctionCallsiteNode, nil
	}

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Look for a user-defined function with the given name.

	callsiteArity := len(astNode.Children)
	udf, err := root.udfManager.LookUp(functionName, callsiteArity)
	if err != nil {
		return nil, err
	}

	// AST snippet for '$z = f($x, $y)':
	// * Assignment "="
	//     * DirectFieldValue "z"
	//     * FunctionCallsite "f"
	//         * DirectFieldValue "x"
	//         * DirectFieldValue "y"
	//
	// Here we need to make an array of our arguments at the callsite, to be
	// paired up with the parameters within he function definition at runtime.
	argumentNodes := make([]IEvaluable, callsiteArity)
	for i, argumentASTNode := range astNode.Children {
		argumentNode, err := root.BuildEvaluableNode(argumentASTNode)
		if err != nil {
			return nil, err
		}
		argumentNodes[i] = argumentNode
	}

	if udf == nil {
		// Mark this as unresolved for an after-pass to see if a UDF with this
		// name/arity has been defined farther down in the DSL expression after
		// this callsite. This happens example when a function is called before
		// it's defined.
		udf = NewUnresolvedUDF(functionName, callsiteArity)
		udfCallsiteNode := NewUDFCallsite(argumentNodes, udf)
		root.rememberUnresolvedFunctionCallsite(udfCallsiteNode)
		return udfCallsiteNode, nil
	} else {
		udfCallsiteNode := NewUDFCallsite(argumentNodes, udf)
		return udfCallsiteNode, nil
	}
}

// ================================================================
// Most DSL functions are implemented in the types package. But these call UDFs
// which are here in the dsl/cst package, so they can't be in the types package

// tSortXFSpace is the datatype for the getSortXFSpace cache-manager.
type tSortXFSpace struct {
	udfCallsite *UDFCallsite
	argsArray   []*types.Mlrval
}

var sortXFCache map[string]*tSortXFSpace = make(map[string]*tSortXFSpace)

// getSortXFSpace manages a cache for the data needed by sortaf/sortmf.  Those
// functions may be invoked on every record of a big data file, so we try to
// cache data they need for UDF-callsite setup.
func getSortXFSpace(
	udfName string,
	udfManager *UDFManager,
	arity int, // 2 for sortaf, 4 for sortmf
) *tSortXFSpace {
	entry := sortXFCache[udfName]
	if entry != nil {
		return entry
	}

	udf, err := udfManager.LookUp(udfName, arity)
	if err != nil { // e.g. exists with wrong arity
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if udf == nil { // e.g. does not exist at all
		fmt.Fprintf(os.Stderr, "mlr: sortaf: comparator function \"%s\" not found.\n", udfName)
		os.Exit(1)
	}

	udfCallsite := NewUDFCallsiteForSortF(udf)
	argsArray := make([]*types.Mlrval, arity)

	entry = &tSortXFSpace{
		udfCallsite: udfCallsite,
		argsArray:   argsArray,
	}

	sortXFCache[udfName] = entry
	return entry
}

// ----------------------------------------------------------------

// SortAF implements the sortaf function, which takes an array as first
// argument and string UDF-name as second argument. It sorts the array using
// the UDF as the comparator.
//
// * Forward sort: func f(a,b) { return a <=> b }
// * Reverse sort: func f(a,b) { return b <=> a }
func SortAF(
	input1 *types.Mlrval,
	input2 *types.Mlrval,
	state *runtime.State,
	udfManager *UDFManager,
) *types.Mlrval {
	inputArray := input1.GetArray()
	if inputArray == nil { // not an array
		return input1
	}
	if !input2.IsString() {
		return types.MLRVAL_ERROR
	}
	udfName := input2.String()

	sortAFSpace := getSortXFSpace(udfName, udfManager, 2)
	udfCallsite := sortAFSpace.udfCallsite
	argsArray := sortAFSpace.argsArray

	outputArray := types.CopyMlrvalArray(inputArray)

	sort.Slice(outputArray, func(i, j int) bool {
		argsArray[0] = &outputArray[i]
		argsArray[1] = &outputArray[j]
		// Call the user's comparator function.
		mret := udfCallsite.EvaluateWithArguments(state, argsArray)
		// Unpack the types.Mlrval return value into a number.
		nret, ok := mret.GetNumericToFloatValue()
		if !ok {
			fmt.Fprintf(
				os.Stderr,
				"mlr: sortaf: comparator function \"%s\" returned non-number \"%s\".\n",
				udfName,
				mret.String(),
			)
			os.Exit(1)
		}
		lib.InternalCodingErrorIf(!ok)
		// Go sort-callback conventions: true if a < b, false otherwise.
		return nret < 0
	})
	return types.MlrvalPointerFromArrayReference(outputArray)
}

// ----------------------------------------------------------------

// SortMF implements the sortmf function, which takes a map as first argument
// and string UDF-name as second argument. It sorts the map using the UDF as
// the comparator.
//
// * Forward sort by key: func f(ak,av,bk,bv) { return ak <=> bk }
// * Reverse sort by key: func f(ak,av,bk,bv) { return bk <=> ak }
func SortMF(
	input1 *types.Mlrval,
	input2 *types.Mlrval,
	state *runtime.State,
	udfManager *UDFManager,
) *types.Mlrval {
	inputMap := input1.GetMap()
	if inputMap == nil { // not a map
		return input1
	}
	if !input2.IsString() {
		return types.MLRVAL_ERROR
	}

	pairsArray := inputMap.ToPairsArray()

	udfName := input2.String()
	sortAFSpace := getSortXFSpace(udfName, udfManager, 4)
	udfCallsite := sortAFSpace.udfCallsite
	argsArray := sortAFSpace.argsArray

	sort.Slice(pairsArray, func(i, j int) bool {
		argsArray[0] = types.MlrvalPointerFromString(pairsArray[i].Key)
		argsArray[1] = pairsArray[i].Value
		argsArray[2] = types.MlrvalPointerFromString(pairsArray[j].Key)
		argsArray[3] = pairsArray[j].Value

		// Call the user's comparator function.
		mret := udfCallsite.EvaluateWithArguments(state, argsArray)
		// Unpack the types.Mlrval return value into a number.
		nret, ok := mret.GetNumericToFloatValue()
		if !ok {
			fmt.Fprintf(
				os.Stderr,
				"mlr: sortaf: comparator function \"%s\" returned non-number \"%s\".\n",
				udfName,
				mret.String(),
			)
			os.Exit(1)
		}
		lib.InternalCodingErrorIf(!ok)
		// Go sort-callback conventions: true if a < b, false otherwise.
		return nret < 0
	})

	sortedMap := types.MlrmapFromPairsArray(pairsArray)
	return types.MlrvalPointerFromMapReferenced(sortedMap)
}
