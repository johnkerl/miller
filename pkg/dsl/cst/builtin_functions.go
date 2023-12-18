// ================================================================
// Methods for built-in functions
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
func (root *RootNode) BuildBuiltinFunctionCallsiteNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(
		astNode.Type != dsl.NodeTypeFunctionCallsite &&
			astNode.Type != dsl.NodeTypeOperator,
	)
	lib.InternalCodingErrorIf(astNode.Token == nil)
	lib.InternalCodingErrorIf(astNode.Children == nil)

	functionName := string(astNode.Token.Lit)

	builtinFunctionInfo := BuiltinFunctionManagerInstance.LookUp(functionName)
	if builtinFunctionInfo != nil {
		if builtinFunctionInfo.hasMultipleArities { // E.g. "+" and "-"
			return root.BuildMultipleArityFunctionCallsiteNode(astNode, builtinFunctionInfo)
		} else if builtinFunctionInfo.zaryFunc != nil {
			return BuildZaryFunctionCallsiteNode(astNode, builtinFunctionInfo)
		} else if builtinFunctionInfo.unaryFunc != nil {
			return root.BuildUnaryFunctionCallsiteNode(astNode, builtinFunctionInfo)
		} else if builtinFunctionInfo.unaryFuncWithContext != nil {
			return root.BuildUnaryFunctionWithContextCallsiteNode(astNode, builtinFunctionInfo)
		} else if builtinFunctionInfo.binaryFunc != nil {
			return root.BuildBinaryFunctionCallsiteNode(astNode, builtinFunctionInfo)
		} else if builtinFunctionInfo.binaryFuncWithState != nil {
			return root.BuildBinaryFunctionWithStateCallsiteNode(astNode, builtinFunctionInfo)
		} else if builtinFunctionInfo.regexCaptureBinaryFunc != nil {
			return root.BuildRegexCaptureBinaryFunctionCallsiteNode(astNode, builtinFunctionInfo)
		} else if builtinFunctionInfo.ternaryFunc != nil {
			return root.BuildTernaryFunctionCallsiteNode(astNode, builtinFunctionInfo)
		} else if builtinFunctionInfo.ternaryFuncWithState != nil {
			return root.BuildTernaryFunctionWithStateCallsiteNode(astNode, builtinFunctionInfo)
		} else if builtinFunctionInfo.variadicFunc != nil {
			return root.BuildVariadicFunctionCallsiteNode(astNode, builtinFunctionInfo)
		} else if builtinFunctionInfo.variadicFuncWithState != nil {
			return root.BuildVariadicFunctionWithStateCallsiteNode(astNode, builtinFunctionInfo)
		} else {
			return nil, fmt.Errorf(
				"at CST BuildFunctionCallsiteNode: builtin function not implemented yet: %s",
				functionName,
			)
		}
	}

	return nil, nil // not found
}

// ----------------------------------------------------------------
func (root *RootNode) BuildMultipleArityFunctionCallsiteNode(
	astNode *dsl.ASTNode,
	builtinFunctionInfo *BuiltinFunctionInfo,
) (IEvaluable, error) {
	callsiteArity := len(astNode.Children)
	if callsiteArity == 1 && builtinFunctionInfo.unaryFunc != nil {
		return root.BuildUnaryFunctionCallsiteNode(astNode, builtinFunctionInfo)
	}
	if callsiteArity == 2 && builtinFunctionInfo.binaryFunc != nil {
		return root.BuildBinaryFunctionCallsiteNode(astNode, builtinFunctionInfo)
	}
	if callsiteArity == 3 && builtinFunctionInfo.ternaryFunc != nil {
		return root.BuildTernaryFunctionCallsiteNode(astNode, builtinFunctionInfo)
	}

	return nil, fmt.Errorf(
		"at CST BuildMultipleArityFunctionCallsiteNode: function name not found: " +
			builtinFunctionInfo.name,
	)
}

// ----------------------------------------------------------------
type ZaryFunctionCallsiteNode struct {
	zaryFunc bifs.ZaryFunc
}

func BuildZaryFunctionCallsiteNode(
	astNode *dsl.ASTNode,
	builtinFunctionInfo *BuiltinFunctionInfo,
) (IEvaluable, error) {
	callsiteArity := len(astNode.Children)
	expectedArity := 0
	if callsiteArity != expectedArity {
		return nil, fmt.Errorf(
			"mlr: function %s invoked with %d argument%s; expected %d",
			builtinFunctionInfo.name,
			callsiteArity,
			lib.Plural(callsiteArity),
			expectedArity,
		)
	}

	return &ZaryFunctionCallsiteNode{
		zaryFunc: builtinFunctionInfo.zaryFunc,
	}, nil
}

func (node *ZaryFunctionCallsiteNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return node.zaryFunc()
}

// ----------------------------------------------------------------
type UnaryFunctionCallsiteNode struct {
	unaryFunc  bifs.UnaryFunc
	evaluable1 IEvaluable
}

func (root *RootNode) BuildUnaryFunctionCallsiteNode(
	astNode *dsl.ASTNode,
	builtinFunctionInfo *BuiltinFunctionInfo,
) (IEvaluable, error) {
	callsiteArity := len(astNode.Children)
	expectedArity := 1
	if callsiteArity != expectedArity {
		return nil, fmt.Errorf(
			"mlr: function %s invoked with %d argument%s; expected %d",
			builtinFunctionInfo.name,
			callsiteArity,
			lib.Plural(callsiteArity),
			expectedArity,
		)
	}

	evaluable1, err := root.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}

	return &UnaryFunctionCallsiteNode{
		unaryFunc:  builtinFunctionInfo.unaryFunc,
		evaluable1: evaluable1,
	}, nil
}

func (node *UnaryFunctionCallsiteNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return node.unaryFunc(node.evaluable1.Evaluate(state))
}

// ----------------------------------------------------------------
type UnaryFunctionWithContextCallsiteNode struct {
	unaryFuncWithContext bifs.UnaryFuncWithContext
	evaluable1           IEvaluable
}

func (root *RootNode) BuildUnaryFunctionWithContextCallsiteNode(
	astNode *dsl.ASTNode,
	builtinFunctionInfo *BuiltinFunctionInfo,
) (IEvaluable, error) {
	callsiteArity := len(astNode.Children)
	expectedArity := 1
	if callsiteArity != expectedArity {
		return nil, fmt.Errorf(
			"mlr: function %s invoked with %d argument%s; expected %d",
			builtinFunctionInfo.name,
			callsiteArity,
			lib.Plural(callsiteArity),
			expectedArity,
		)
	}

	evaluable1, err := root.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}

	return &UnaryFunctionWithContextCallsiteNode{
		unaryFuncWithContext: builtinFunctionInfo.unaryFuncWithContext,
		evaluable1:           evaluable1,
	}, nil
}

func (node *UnaryFunctionWithContextCallsiteNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return node.unaryFuncWithContext(node.evaluable1.Evaluate(state), state.Context)
}

// ----------------------------------------------------------------
type BinaryFunctionCallsiteNode struct {
	binaryFunc bifs.BinaryFunc
	evaluable1 IEvaluable
	evaluable2 IEvaluable
}

func (root *RootNode) BuildBinaryFunctionCallsiteNode(
	astNode *dsl.ASTNode,
	builtinFunctionInfo *BuiltinFunctionInfo,
) (IEvaluable, error) {
	callsiteArity := len(astNode.Children)
	expectedArity := 2
	if callsiteArity != expectedArity {
		return nil, fmt.Errorf(
			"mlr: function %s invoked with %d argument%s; expected %d",
			builtinFunctionInfo.name,
			callsiteArity,
			lib.Plural(callsiteArity),
			expectedArity,
		)
	}

	evaluable1, err := root.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	evaluable2, err := root.BuildEvaluableNode(astNode.Children[1])
	if err != nil {
		return nil, err
	}

	// Special short-circuiting cases
	if builtinFunctionInfo.name == "&&" {
		return BuildLogicalANDOperatorNode(
			evaluable1,
			evaluable2,
		), nil
	}
	if builtinFunctionInfo.name == "||" {
		return BuildLogicalOROperatorNode(
			evaluable1,
			evaluable2,
		), nil
	}
	if builtinFunctionInfo.name == "??" {
		return BuildAbsentCoalesceOperatorNode(
			evaluable1,
			evaluable2,
		), nil
	}
	if builtinFunctionInfo.name == "???" {
		return BuildEmptyCoalesceOperatorNode(
			evaluable1,
			evaluable2,
		), nil
	}

	return &BinaryFunctionCallsiteNode{
		binaryFunc: builtinFunctionInfo.binaryFunc,
		evaluable1: evaluable1,
		evaluable2: evaluable2,
	}, nil
}

func (node *BinaryFunctionCallsiteNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return node.binaryFunc(
		node.evaluable1.Evaluate(state),
		node.evaluable2.Evaluate(state),
	)
}

// ----------------------------------------------------------------
type BinaryFunctionWithStateCallsiteNode struct {
	binaryFuncWithState BinaryFuncWithState
	evaluable1          IEvaluable
	evaluable2          IEvaluable
}

func (root *RootNode) BuildBinaryFunctionWithStateCallsiteNode(
	astNode *dsl.ASTNode,
	builtinFunctionInfo *BuiltinFunctionInfo,
) (IEvaluable, error) {
	callsiteArity := len(astNode.Children)
	expectedArity := 2
	if callsiteArity != expectedArity {
		return nil, fmt.Errorf(
			"mlr: function %s invoked with %d argument%s; expected %d",
			builtinFunctionInfo.name,
			callsiteArity,
			lib.Plural(callsiteArity),
			expectedArity,
		)
	}

	evaluable1, err := root.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	evaluable2, err := root.BuildEvaluableNode(astNode.Children[1])
	if err != nil {
		return nil, err
	}

	return &BinaryFunctionWithStateCallsiteNode{
		binaryFuncWithState: builtinFunctionInfo.binaryFuncWithState,
		evaluable1:          evaluable1,
		evaluable2:          evaluable2,
	}, nil
}

func (node *BinaryFunctionWithStateCallsiteNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return node.binaryFuncWithState(
		node.evaluable1.Evaluate(state),
		node.evaluable2.Evaluate(state),
		state,
	)
}

// ----------------------------------------------------------------
type TernaryFunctionWithStateCallsiteNode struct {
	ternaryFuncWithState TernaryFuncWithState
	evaluable1           IEvaluable
	evaluable2           IEvaluable
	evaluable3           IEvaluable
}

func (root *RootNode) BuildTernaryFunctionWithStateCallsiteNode(
	astNode *dsl.ASTNode,
	builtinFunctionInfo *BuiltinFunctionInfo,
) (IEvaluable, error) {
	callsiteArity := len(astNode.Children)
	expectedArity := 3
	if callsiteArity != expectedArity {
		return nil, fmt.Errorf(
			"mlr: function %s invoked with %d argument%s; expected %d",
			builtinFunctionInfo.name,
			callsiteArity,
			lib.Plural(callsiteArity),
			expectedArity,
		)
	}

	evaluable1, err := root.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	evaluable2, err := root.BuildEvaluableNode(astNode.Children[1])
	if err != nil {
		return nil, err
	}
	evaluable3, err := root.BuildEvaluableNode(astNode.Children[2])
	if err != nil {
		return nil, err
	}

	return &TernaryFunctionWithStateCallsiteNode{
		ternaryFuncWithState: builtinFunctionInfo.ternaryFuncWithState,
		evaluable1:           evaluable1,
		evaluable2:           evaluable2,
		evaluable3:           evaluable3,
	}, nil
}

func (node *TernaryFunctionWithStateCallsiteNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return node.ternaryFuncWithState(
		node.evaluable1.Evaluate(state),
		node.evaluable2.Evaluate(state),
		node.evaluable3.Evaluate(state),
		state,
	)
}

// ----------------------------------------------------------------
// RegexCaptureBinaryFunctionCallsiteNode special-cases the =~ and !=~
// operators which set the CST State object's captures array for "\1".."\9".
// This is identical to BinaryFunctionCallsite except that
// BinaryFunctionCallsite's impl function takes two *mlrval.Mlrval arguments and
// returns a *mlrval.Mlrval, whereas RegexCaptureBinaryFunctionCallsiteNode's
// impl function takes two *mlrval.Mlrval arguments but returns *mlrval.Mlrval
// along with a []string captures array. The captures are stored in the State
// object for use in subsequent statements.
//
// Note the use of "capture" is ambiguous:
//
//   - There is the regex-match part which captures submatches out
//     of a full match expression, and saves them.
//
// * Then there is the part which inserts these captures into another string.
//
//   - For sub/gsub, the former and latter are both within the sub/gsub routine.
//     E.g. with
//     $y = sub($x, "(..)_(...)", "\2:\1"
//     and $x being "ab_cde", $y will be "cde:ab".
//
//   - For =~ and !=~, the former are right there, but the latter can be several
//     lines later. E.g.
//     if ($x =~ "(..)_(...)") {
//     ... other lines of code ...
//     $y = "\2:\1";
//     }
//
// So: this RegexCaptureBinaryFunctionCallsiteNode only refers to the =~ and
// !=~ callsites only -- not sub/gsub, and not the capture-using replacement
// statements like '$y = "\2:\1".
type RegexCaptureBinaryFunctionCallsiteNode struct {
	regexCaptureBinaryFunc bifs.RegexCaptureBinaryFunc
	evaluable1             IEvaluable
	evaluable2             IEvaluable
}

func (root *RootNode) BuildRegexCaptureBinaryFunctionCallsiteNode(
	astNode *dsl.ASTNode,
	builtinFunctionInfo *BuiltinFunctionInfo,
) (IEvaluable, error) {
	callsiteArity := len(astNode.Children)
	expectedArity := 2
	if callsiteArity != expectedArity {
		return nil, fmt.Errorf(
			"mlr: function %s invoked with %d argument%s; expected %d",
			builtinFunctionInfo.name,
			callsiteArity,
			lib.Plural(callsiteArity),
			expectedArity,
		)
	}

	evaluable1, err := root.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	evaluable2, err := root.BuildEvaluableNode(astNode.Children[1])
	if err != nil {
		return nil, err
	}

	return &RegexCaptureBinaryFunctionCallsiteNode{
		regexCaptureBinaryFunc: builtinFunctionInfo.regexCaptureBinaryFunc,
		evaluable1:             evaluable1,
		evaluable2:             evaluable2,
	}, nil
}

func (node *RegexCaptureBinaryFunctionCallsiteNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	output, captures := node.regexCaptureBinaryFunc(
		node.evaluable1.Evaluate(state),
		node.evaluable2.Evaluate(state),
	)
	state.SetRegexCaptures(captures)
	return output
}

// ----------------------------------------------------------------
// DotCallsiteNode special-cases the dot operator, which is:
// * string + string, with coercion to string if either side is int/float/bool/etc.
// * map attribute access, if the left-hand side is a map.
type DotCallsiteNode struct {
	evaluable1 IEvaluable
	evaluable2 IEvaluable
	string2    string
}

func (root *RootNode) BuildDotCallsiteNode(
	astNode *dsl.ASTNode,
) (IEvaluable, error) {
	callsiteArity := len(astNode.Children)
	expectedArity := 2
	if callsiteArity != expectedArity {
		return nil, fmt.Errorf(
			"mlr: function %s invoked with %d argument%s; expected %d",
			".",
			callsiteArity,
			lib.Plural(callsiteArity),
			expectedArity,
		)
	}

	evaluable1, err := root.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	evaluable2, err := root.BuildEvaluableNode(astNode.Children[1])
	if err != nil {
		return nil, err
	}

	return &DotCallsiteNode{
		evaluable1: evaluable1,
		evaluable2: evaluable2,
		string2:    string(astNode.Children[1].Token.Lit),
	}, nil
}

func (node *DotCallsiteNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	// For strict mode, absence should be detected on the node.evaluable1 evaluator.
	value1 := node.evaluable1.Evaluate(state)

	mapvalue1 := value1.GetMap()

	if mapvalue1 != nil {
		// Case 1: map.attribute as shorthand for map["attribute"]
		value2 := mapvalue1.Get(node.string2)
		if value2 == nil {
			return mlrval.ABSENT.StrictModeCheck(state.StrictMode, "map access ["+node.string2+"]")
		} else {
			return value2
		}
	} else {
		// Case 2: string concatenation
		value2 := node.evaluable2.Evaluate(state)
		return bifs.BIF_dot(value1, value2)
	}
}

// ----------------------------------------------------------------
type TernaryFunctionCallsiteNode struct {
	ternaryFunc bifs.TernaryFunc
	evaluable1  IEvaluable
	evaluable2  IEvaluable
	evaluable3  IEvaluable
}

func (root *RootNode) BuildTernaryFunctionCallsiteNode(
	astNode *dsl.ASTNode,
	builtinFunctionInfo *BuiltinFunctionInfo,
) (IEvaluable, error) {
	callsiteArity := len(astNode.Children)
	expectedArity := 3
	if callsiteArity != expectedArity {
		return nil, fmt.Errorf(
			"mlr: function %s invoked with %d argument%s; expected %d",
			builtinFunctionInfo.name,
			callsiteArity,
			lib.Plural(callsiteArity),
			expectedArity,
		)
	}

	evaluable1, err := root.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	evaluable2, err := root.BuildEvaluableNode(astNode.Children[1])
	if err != nil {
		return nil, err
	}
	evaluable3, err := root.BuildEvaluableNode(astNode.Children[2])
	if err != nil {
		return nil, err
	}

	// Special short-circuiting case
	if builtinFunctionInfo.name == "?:" {
		return BuildStandardTernaryOperatorNode(
			evaluable1,
			evaluable2,
			evaluable3,
		), nil
	}

	return &TernaryFunctionCallsiteNode{
		ternaryFunc: builtinFunctionInfo.ternaryFunc,
		evaluable1:  evaluable1,
		evaluable2:  evaluable2,
		evaluable3:  evaluable3,
	}, nil
}

func (node *TernaryFunctionCallsiteNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	return node.ternaryFunc(
		node.evaluable1.Evaluate(state),
		node.evaluable2.Evaluate(state),
		node.evaluable3.Evaluate(state),
	)
}

// ----------------------------------------------------------------
type VariadicFunctionCallsiteNode struct {
	variadicFunc bifs.VariadicFunc
	evaluables   []IEvaluable
}

func (root *RootNode) BuildVariadicFunctionCallsiteNode(
	astNode *dsl.ASTNode,
	builtinFunctionInfo *BuiltinFunctionInfo,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Children == nil)
	evaluables := make([]IEvaluable, len(astNode.Children))

	callsiteArity := len(astNode.Children)

	if callsiteArity < builtinFunctionInfo.minimumVariadicArity {
		return nil, fmt.Errorf(
			"mlr: function %s takes minimum argument count %d; got %d.\n",
			builtinFunctionInfo.name,
			builtinFunctionInfo.minimumVariadicArity,
			callsiteArity,
		)
	}

	if builtinFunctionInfo.maximumVariadicArity != 0 {
		if callsiteArity > builtinFunctionInfo.maximumVariadicArity {
			return nil, fmt.Errorf(
				"mlr: function %s takes maximum argument count %d; got %d.\n",
				builtinFunctionInfo.name,
				builtinFunctionInfo.maximumVariadicArity,
				callsiteArity,
			)
		}
	}

	var err error = nil
	for i, astChildNode := range astNode.Children {
		evaluables[i], err = root.BuildEvaluableNode(astChildNode)
		if err != nil {
			return nil, err
		}
	}
	return &VariadicFunctionCallsiteNode{
		variadicFunc: builtinFunctionInfo.variadicFunc,
		evaluables:   evaluables,
	}, nil
}

func (node *VariadicFunctionCallsiteNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	args := make([]*mlrval.Mlrval, len(node.evaluables))
	for i := range node.evaluables {
		args[i] = node.evaluables[i].Evaluate(state)
	}
	return node.variadicFunc(args)
}

// ----------------------------------------------------------------
type VariadicFunctionWithStateCallsiteNode struct {
	variadicFuncWithState VariadicFuncWithState
	evaluables            []IEvaluable
}

func (root *RootNode) BuildVariadicFunctionWithStateCallsiteNode(
	astNode *dsl.ASTNode,
	builtinFunctionInfo *BuiltinFunctionInfo,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Children == nil)
	evaluables := make([]IEvaluable, len(astNode.Children))

	callsiteArity := len(astNode.Children)

	if callsiteArity < builtinFunctionInfo.minimumVariadicArity {
		return nil, fmt.Errorf(
			"mlr: function %s takes minimum argument count %d; got %d.\n",
			builtinFunctionInfo.name,
			builtinFunctionInfo.minimumVariadicArity,
			callsiteArity,
		)
	}

	if builtinFunctionInfo.maximumVariadicArity != 0 {
		if callsiteArity > builtinFunctionInfo.maximumVariadicArity {
			return nil, fmt.Errorf(
				"mlr: function %s takes maximum argument count %d; got %d.\n",
				builtinFunctionInfo.name,
				builtinFunctionInfo.maximumVariadicArity,
				callsiteArity,
			)
		}
	}

	var err error = nil
	for i, astChildNode := range astNode.Children {
		evaluables[i], err = root.BuildEvaluableNode(astChildNode)
		if err != nil {
			return nil, err
		}
	}
	return &VariadicFunctionWithStateCallsiteNode{
		variadicFuncWithState: builtinFunctionInfo.variadicFuncWithState,
		evaluables:            evaluables,
	}, nil
}

func (node *VariadicFunctionWithStateCallsiteNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	args := make([]*mlrval.Mlrval, len(node.evaluables))
	for i := range node.evaluables {
		args[i] = node.evaluables[i].Evaluate(state)
	}
	return node.variadicFuncWithState(args, state)
}

// ================================================================
type LogicalANDOperatorNode struct {
	a, b IEvaluable
}

func BuildLogicalANDOperatorNode(a, b IEvaluable) *LogicalANDOperatorNode {
	return &LogicalANDOperatorNode{
		a: a,
		b: b,
	}
}

// This is different from most of the evaluator functions in that it does
// short-circuiting: since is logical AND, the second argument is not evaluated
// if the first argument is false. Thus we cannot use disposition matrices.
//
// * evaluate a
// * if   a is error:
//   * return a
// * elif a is absent:
//   * Evaluate b
//   * if   b is error: return error
//   * elif b is empty or absent: return absent
//   * elif b is empty or absent: return absent
//   * else: return b
// * elif a is empty:
//   * evaluate b
//   * if   b is error: return error
//   * elif b is empty: return empty
//   * elif b is absent: return absent
//   * else: return b
// * else:
//   * return the BIF (using its disposition matrix)

// mlr help type-arithmetic-info-extended | lumin -c red .error. | lumin -c blue .absent. | lumin -c green .empty.

func (node *LogicalANDOperatorNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	aout := node.a.Evaluate(state)
	atype := aout.Type()

	if atype == mlrval.MT_ERROR {
		return aout
	}

	if atype == mlrval.MT_ABSENT {
		bout := node.b.Evaluate(state)
		btype := bout.Type()
		if btype == mlrval.MT_ERROR {
			return bout
		}
		if btype == mlrval.MT_VOID || btype == mlrval.MT_ABSENT {
			return mlrval.ABSENT
		}
		if btype != mlrval.MT_BOOL {
			return mlrval.FromNotNamedTypeError("&&", bout, "absent or boolean")
		}
		return bout
	}

	if atype == mlrval.MT_VOID {
		bout := node.b.Evaluate(state)
		btype := bout.Type()
		if btype == mlrval.MT_ERROR {
			return bout
		}
		if btype == mlrval.MT_VOID {
			return mlrval.FromNotNamedTypeError("&&", bout, "absent or boolean")
		}
		if btype == mlrval.MT_ABSENT {
			return mlrval.ABSENT
		}
		if btype != mlrval.MT_BOOL {
			return mlrval.FromNotNamedTypeError("&&", bout, "absent or boolean")
		}
		return bout
	}

	// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	if aout.IsFalse() {
		// This means false && bogus type evaluates to false, which is sad but
		// which we MUST do in order to not violate the short-circuiting
		// property.  We would have to evaluate b to know if it were error or
		// not.
		return aout
	}

	bout := node.b.Evaluate(state)
	btype := bout.Type()
	if !(btype == mlrval.MT_ABSENT || btype == mlrval.MT_BOOL) {
		return mlrval.FromNotNamedTypeError("&&", bout, "absent or boolean")
	}
	if btype == mlrval.MT_ABSENT {
		return mlrval.ABSENT
	}

	return bifs.BIF_logical_AND(aout, bout)
}

// ================================================================
type LogicalOROperatorNode struct {
	a, b IEvaluable
}

func BuildLogicalOROperatorNode(a, b IEvaluable) *LogicalOROperatorNode {
	return &LogicalOROperatorNode{
		a: a,
		b: b,
	}
}

// This is different from most of the evaluator functions in that it does
// short-circuiting: since is logical OR, the second argument is not evaluated
// if the first argument is false.

func (node *LogicalOROperatorNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	aout := node.a.Evaluate(state)
	atype := aout.Type()

	if atype == mlrval.MT_ERROR {
		return aout
	}

	if atype == mlrval.MT_ABSENT {
		bout := node.b.Evaluate(state)
		btype := bout.Type()
		if btype == mlrval.MT_ERROR {
			return bout
		}
		if btype == mlrval.MT_VOID || btype == mlrval.MT_ABSENT {
			return mlrval.ABSENT
		}
		if btype == mlrval.MT_VOID {
			return mlrval.FromNotNamedTypeError("||", bout, "absent or boolean")
		}
		if btype != mlrval.MT_BOOL {
			return mlrval.FromNotNamedTypeError("||", bout, "absent or boolean")
		}
		return bout
	}

	if atype == mlrval.MT_VOID {
		bout := node.b.Evaluate(state)
		btype := bout.Type()
		if btype == mlrval.MT_ERROR {
			return bout
		}
		if btype == mlrval.MT_VOID {
			return mlrval.FromNotNamedTypeError("||", bout, "absent or boolean")
		}
		if btype == mlrval.MT_ABSENT {
			return mlrval.ABSENT
		}
		if btype != mlrval.MT_BOOL {
			return mlrval.FromNotNamedTypeError("||", bout, "absent or boolean")
		}
		return bout
	}

	// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	if aout.IsTrue() {
		// This means true || bogus type evaluates to true, which is sad but
		// which we MUST do in order to not violate the short-circuiting
		// property.  We would have to evaluate b to know if it were error or
		// not.
		return aout
	}

	bout := node.b.Evaluate(state)
	btype := bout.Type()
	if !(btype == mlrval.MT_ABSENT || btype == mlrval.MT_BOOL) {
		return mlrval.FromNotNamedTypeError("||", bout, "absent or boolean")
	}
	if btype == mlrval.MT_ABSENT {
		return mlrval.ABSENT
	}

	return bifs.BIF_logical_OR(aout, bout)
}

// ================================================================
// a ?? b evaluates to b only when a is absent. Example: '$foo ?? 0' when the
// current record has no field $foo.
type AbsentCoalesceOperatorNode struct{ a, b IEvaluable }

func BuildAbsentCoalesceOperatorNode(a, b IEvaluable) *AbsentCoalesceOperatorNode {
	return &AbsentCoalesceOperatorNode{a: a, b: b}
}

// This is different from most of the evaluator functions in that it does
// short-circuiting: the second argument is not evaluated if the first
// argument is not absent.
func (node *AbsentCoalesceOperatorNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	aout := node.a.Evaluate(state)
	if aout.Type() != mlrval.MT_ABSENT {
		return aout
	}

	return node.b.Evaluate(state)
}

// ================================================================
// a ?? b evaluates to b only when a is absent or empty. Example: '$foo ?? 0'
// when the current record has no field $foo, or when $foo is empty..
type EmptyCoalesceOperatorNode struct{ a, b IEvaluable }

func BuildEmptyCoalesceOperatorNode(a, b IEvaluable) *EmptyCoalesceOperatorNode {
	return &EmptyCoalesceOperatorNode{a: a, b: b}
}

// This is different from most of the evaluator functions in that it does
// short-circuiting: the second argument is not evaluated if the first
// argument is not absent.
func (node *EmptyCoalesceOperatorNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	aout := node.a.Evaluate(state)
	atype := aout.Type()
	if atype == mlrval.MT_ABSENT || atype == mlrval.MT_VOID || (atype == mlrval.MT_STRING && aout.String() == "") {
		return node.b.Evaluate(state)
	} else {
		return aout
	}
}

// ================================================================
type StandardTernaryOperatorNode struct{ a, b, c IEvaluable }

func BuildStandardTernaryOperatorNode(a, b, c IEvaluable) *StandardTernaryOperatorNode {
	return &StandardTernaryOperatorNode{a: a, b: b, c: c}
}
func (node *StandardTernaryOperatorNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	aout := node.a.Evaluate(state)

	boolValue, isBool := aout.GetBoolValue()
	if !isBool {
		return mlrval.FromNotBooleanError("?:", aout)
	}

	// Short-circuit: defer evaluation unless needed
	if boolValue == true {
		return node.b.Evaluate(state)
	} else {
		return node.c.Evaluate(state)
	}
}

// ================================================================
// The function-manager logic is designed to make it easy to implement a large
// number of functions/operators with a small number of keystrokes. The general
// paradigm is evaluate the arguments, then invoke the function/operator.
//
// For some, such as the binary operators "&&" and "||", and the ternary
// operator "?:", there is short-circuiting logic wherein one argument may not
// be evaluated depending on another's value. These functions are placeholders
// for the function-manager lookup table to indicate the arity of the function,
// even though at runtime these functions should not get invoked.

func BinaryShortCircuitPlaceholder(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	lib.InternalCodingErrorPanic("Short-circuting was not correctly implemented")
	return nil // not reached
}

func TernaryShortCircuitPlaceholder(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval {
	lib.InternalCodingErrorPanic("Short-circuting was not correctly implemented")
	return nil // not reached
}
