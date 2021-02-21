// ================================================================
// Methods for built-in functions
// ================================================================

package cst

import (
	"errors"
	"fmt"

	"miller/src/dsl"
	"miller/src/lib"
	"miller/src/runtime"
	"miller/src/types"
)

// ----------------------------------------------------------------
func (this *RootNode) BuildBuiltinFunctionCallsiteNode(
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
			return this.BuildMultipleArityFunctionCallsiteNode(astNode, builtinFunctionInfo)
		} else if builtinFunctionInfo.zaryFunc != nil {
			return this.BuildZaryFunctionCallsiteNode(astNode, builtinFunctionInfo)
		} else if builtinFunctionInfo.unaryFunc != nil {
			return this.BuildUnaryFunctionCallsiteNode(astNode, builtinFunctionInfo)
		} else if builtinFunctionInfo.contextualUnaryFunc != nil {
			return this.BuildContextualUnaryFunctionCallsiteNode(astNode, builtinFunctionInfo)
		} else if builtinFunctionInfo.binaryFunc != nil {
			return this.BuildBinaryFunctionCallsiteNode(astNode, builtinFunctionInfo)
		} else if builtinFunctionInfo.ternaryFunc != nil {
			return this.BuildTernaryFunctionCallsiteNode(astNode, builtinFunctionInfo)
		} else if builtinFunctionInfo.variadicFunc != nil {
			return this.BuildVariadicFunctionCallsiteNode(astNode, builtinFunctionInfo)
		} else {
			return nil, errors.New(
				"CST BuildFunctionCallsiteNode: builtin function not implemented yet: " +
					functionName,
			)
		}
	}

	return nil, nil // not found
}

// ----------------------------------------------------------------
func (this *RootNode) BuildMultipleArityFunctionCallsiteNode(
	astNode *dsl.ASTNode,
	builtinFunctionInfo *BuiltinFunctionInfo,
) (IEvaluable, error) {
	callsiteArity := len(astNode.Children)
	if callsiteArity == 1 && builtinFunctionInfo.unaryFunc != nil {
		return this.BuildUnaryFunctionCallsiteNode(astNode, builtinFunctionInfo)
	}
	if callsiteArity == 2 && builtinFunctionInfo.binaryFunc != nil {
		return this.BuildBinaryFunctionCallsiteNode(astNode, builtinFunctionInfo)
	}
	if callsiteArity == 3 && builtinFunctionInfo.ternaryFunc != nil {
		return this.BuildTernaryFunctionCallsiteNode(astNode, builtinFunctionInfo)
	}

	return nil, errors.New(
		fmt.Sprintf(
			"CST BuildMultipleArityFunctionCallsiteNode: function name not found: " +
				builtinFunctionInfo.name,
		),
	)
}

// ----------------------------------------------------------------
type ZaryFunctionCallsiteNode struct {
	zaryFunc types.ZaryFunc
}

func (this *RootNode) BuildZaryFunctionCallsiteNode(
	astNode *dsl.ASTNode,
	builtinFunctionInfo *BuiltinFunctionInfo,
) (IEvaluable, error) {
	callsiteArity := len(astNode.Children)
	expectedArity := 0
	if callsiteArity != expectedArity {
		return nil, errors.New(
			fmt.Sprintf(
				"Miller: function %s invoked with %d argument%s; expected %d",
				builtinFunctionInfo.name,
				callsiteArity,
				lib.Plural(callsiteArity),
				expectedArity,
			),
		)
	}

	return &ZaryFunctionCallsiteNode{zaryFunc: builtinFunctionInfo.zaryFunc}, nil
}

func (this *ZaryFunctionCallsiteNode) Evaluate(state *runtime.State) types.Mlrval {
	// xxx temp
	output := types.MlrvalFromAbsent()
	this.zaryFunc(&output)
	return output
}

// ----------------------------------------------------------------
type UnaryFunctionCallsiteNode struct {
	unaryFunc  types.UnaryFunc
	evaluable1 IEvaluable
}

func (this *RootNode) BuildUnaryFunctionCallsiteNode(
	astNode *dsl.ASTNode,
	builtinFunctionInfo *BuiltinFunctionInfo,
) (IEvaluable, error) {
	callsiteArity := len(astNode.Children)
	expectedArity := 1
	if callsiteArity != expectedArity {
		return nil, errors.New(
			fmt.Sprintf(
				"Miller: function %s invoked with %d argument%s; expected %d",
				builtinFunctionInfo.name,
				callsiteArity,
				lib.Plural(callsiteArity),
				expectedArity,
			),
		)
	}

	evaluable1, err := this.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}

	return &UnaryFunctionCallsiteNode{
		unaryFunc:  builtinFunctionInfo.unaryFunc,
		evaluable1: evaluable1,
	}, nil
}

func (this *UnaryFunctionCallsiteNode) Evaluate(state *runtime.State) types.Mlrval {
	arg1 := this.evaluable1.Evaluate(state)
	// xxx temp
	output := types.MlrvalFromAbsent()
	this.unaryFunc(&output, &arg1)
	return output
}

// ----------------------------------------------------------------
type ContextualUnaryFunctionCallsiteNode struct {
	contextualUnaryFunc types.ContextualUnaryFunc
	evaluable1          IEvaluable
}

func (this *RootNode) BuildContextualUnaryFunctionCallsiteNode(
	astNode *dsl.ASTNode,
	builtinFunctionInfo *BuiltinFunctionInfo,
) (IEvaluable, error) {
	callsiteArity := len(astNode.Children)
	expectedArity := 1
	if callsiteArity != expectedArity {
		return nil, errors.New(
			fmt.Sprintf(
				"Miller: function %s invoked with %d argument%s; expected %d",
				builtinFunctionInfo.name,
				callsiteArity,
				lib.Plural(callsiteArity),
				expectedArity,
			),
		)
	}

	evaluable1, err := this.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}

	return &ContextualUnaryFunctionCallsiteNode{
		contextualUnaryFunc: builtinFunctionInfo.contextualUnaryFunc,
		evaluable1:          evaluable1,
	}, nil
}

func (this *ContextualUnaryFunctionCallsiteNode) Evaluate(state *runtime.State) types.Mlrval {
	arg1 := this.evaluable1.Evaluate(state)
	// xxx temp
	output := types.MlrvalFromAbsent()
	this.contextualUnaryFunc(&output, &arg1, state.Context)
	return output
}

// ----------------------------------------------------------------
type BinaryFunctionCallsiteNode struct {
	binaryFunc types.BinaryFunc
	evaluable1 IEvaluable
	evaluable2 IEvaluable
}

func (this *RootNode) BuildBinaryFunctionCallsiteNode(
	astNode *dsl.ASTNode,
	builtinFunctionInfo *BuiltinFunctionInfo,
) (IEvaluable, error) {
	callsiteArity := len(astNode.Children)
	expectedArity := 2
	if callsiteArity != expectedArity {
		return nil, errors.New(
			fmt.Sprintf(
				"Miller: function %s invoked with %d argument%s; expected %d",
				builtinFunctionInfo.name,
				callsiteArity,
				lib.Plural(callsiteArity),
				expectedArity,
			),
		)
	}

	evaluable1, err := this.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	evaluable2, err := this.BuildEvaluableNode(astNode.Children[1])
	if err != nil {
		return nil, err
	}

	// Special short-circuiting cases
	if builtinFunctionInfo.name == "&&" {
		return this.BuildLogicalANDOperatorNode(
			evaluable1,
			evaluable2,
		), nil
	}
	if builtinFunctionInfo.name == "||" {
		return this.BuildLogicalOROperatorNode(
			evaluable1,
			evaluable2,
		), nil
	}
	if builtinFunctionInfo.name == "??" {
		return this.BuildAbsentCoalesceOperatorNode(
			evaluable1,
			evaluable2,
		), nil
	}
	if builtinFunctionInfo.name == "???" {
		return this.BuildEmptyCoalesceOperatorNode(
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

func (this *BinaryFunctionCallsiteNode) Evaluate(state *runtime.State) types.Mlrval {
	// xxx temp
	output := types.MlrvalFromError()
	arg1 := this.evaluable1.Evaluate(state)
	arg2 := this.evaluable2.Evaluate(state)
	this.binaryFunc(&output, &arg1, &arg2)
	return output
}

// ----------------------------------------------------------------
type TernaryFunctionCallsiteNode struct {
	ternaryFunc types.TernaryFunc
	evaluable1  IEvaluable
	evaluable2  IEvaluable
	evaluable3  IEvaluable
}

func (this *RootNode) BuildTernaryFunctionCallsiteNode(
	astNode *dsl.ASTNode,
	builtinFunctionInfo *BuiltinFunctionInfo,
) (IEvaluable, error) {
	callsiteArity := len(astNode.Children)
	expectedArity := 3
	if callsiteArity != expectedArity {
		return nil, errors.New(
			fmt.Sprintf(
				"Miller: function %s invoked with %d argument%s; expected %d",
				builtinFunctionInfo.name,
				callsiteArity,
				lib.Plural(callsiteArity),
				expectedArity,
			),
		)
	}

	evaluable1, err := this.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	evaluable2, err := this.BuildEvaluableNode(astNode.Children[1])
	if err != nil {
		return nil, err
	}
	evaluable3, err := this.BuildEvaluableNode(astNode.Children[2])
	if err != nil {
		return nil, err
	}

	// Special short-circuiting case
	if builtinFunctionInfo.name == "?:" {
		return this.BuildStandardTernaryOperatorNode(
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

func (this *TernaryFunctionCallsiteNode) Evaluate(state *runtime.State) types.Mlrval {
	arg1 := this.evaluable1.Evaluate(state)
	arg2 := this.evaluable2.Evaluate(state)
	arg3 := this.evaluable3.Evaluate(state)
	// xxx temp
	output := types.MlrvalFromAbsent()
	this.ternaryFunc(&output, &arg1, &arg2, &arg3)
	return output
}

// ----------------------------------------------------------------
type VariadicFunctionCallsiteNode struct {
	variadicFunc types.VariadicFunc
	evaluables   []IEvaluable
}

func (this *RootNode) BuildVariadicFunctionCallsiteNode(
	astNode *dsl.ASTNode,
	builtinFunctionInfo *BuiltinFunctionInfo,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Children == nil)
	evaluables := make([]IEvaluable, len(astNode.Children))

	callsiteArity := len(astNode.Children)
	if callsiteArity < builtinFunctionInfo.minimumVariadicArity {
		return nil, errors.New(
			fmt.Sprintf(
				"Miller: function %s takes minimum argument count %d; got %d.\n",
				builtinFunctionInfo.minimumVariadicArity,
				callsiteArity,
			),
		)
	}

	var err error = nil
	for i, astChildNode := range astNode.Children {
		evaluables[i], err = this.BuildEvaluableNode(astChildNode)
		if err != nil {
			return nil, err
		}
	}
	return &VariadicFunctionCallsiteNode{
		variadicFunc: builtinFunctionInfo.variadicFunc,
		evaluables:   evaluables,
	}, nil
}

func (this *VariadicFunctionCallsiteNode) Evaluate(state *runtime.State) types.Mlrval {
	args := make([]*types.Mlrval, len(this.evaluables))
	for i, evaluable := range this.evaluables {
		arg := evaluable.Evaluate(state)
		args[i] = &arg
	}
	// xxx temp
	output := types.MlrvalFromError()
	this.variadicFunc(&output, args)
	return output
}

// ================================================================
type LogicalANDOperatorNode struct{ a, b IEvaluable }

func (this *RootNode) BuildLogicalANDOperatorNode(a, b IEvaluable) *LogicalANDOperatorNode {
	return &LogicalANDOperatorNode{a: a, b: b}
}

// This is different from most of the evaluator functions in that it does
// short-circuiting: since is logical AND, the second argument is not evaluated
// if the first argument is false.
//
// Disposition matrix:
//
//       {
//a      b  ERROR   ABSENT  EMPTY  STRING INT    FLOAT  BOOL
//ERROR  :  {ERROR, ERROR,  ERROR, ERROR, ERROR, ERROR, ERROR},
//ABSENT :  {ERROR, absent, ERROR, ERROR, ERROR, ERROR, absent},
//EMPTY  :  {ERROR, ERROR,  ERROR, ERROR, ERROR, ERROR, ERROR},
//STRING :  {ERROR, ERROR,  ERROR, ERROR, ERROR, ERROR, ERROR},
//INT    :  {ERROR, ERROR,  ERROR, ERROR, ERROR, ERROR, ERROR},
//FLOAT  :  {ERROR, ERROR,  ERROR, ERROR, ERROR, ERROR, ERROR},
//BOOL   :  {ERROR, absent, ERROR, ERROR, ERROR, ERROR, a&&b},
//       }
//
// which without the all-error rows/columns reduces to
//
//       {
//a      b  ABSENT   BOOL
//ABSENT :  {absent, absent},
//BOOL   :  {absent, a&&b},
//       }
//
// So:
// * Evaluate a
// * If a is not absent or bool: return error
// * If a is absent: return absent
// * If a is false: return a
// * Now a is boolean true
// * Evaluate b
// * If b is not absent or bool: return error
// * If b is absent: return absent
// * Return a && b

func (this *LogicalANDOperatorNode) Evaluate(state *runtime.State) types.Mlrval {
	aout := this.a.Evaluate(state)
	atype := aout.GetType()
	if !(atype == types.MT_ABSENT || atype == types.MT_BOOL) {
		return types.MlrvalFromError()
	}
	if atype == types.MT_ABSENT {
		return aout
	}
	if aout.IsFalse() {
		// This means false && bogus type evaluates to true, which is sad but
		// which we MUST do in order to not violate the short-circuiting
		// property.  We would have to evaluate b to know if it were error or
		// not.
		return aout
	}

	bout := this.b.Evaluate(state)
	btype := bout.GetType()
	if !(btype == types.MT_ABSENT || btype == types.MT_BOOL) {
		return types.MlrvalFromError()
	}
	if btype == types.MT_ABSENT {
		return bout
	}
	// xxx temp
	output := types.MlrvalFromError()
	types.MlrvalLogicalAND(&output, &aout, &bout)
	return output
}

// ================================================================
type LogicalOROperatorNode struct{ a, b IEvaluable }

func (this *RootNode) BuildLogicalOROperatorNode(a, b IEvaluable) *LogicalOROperatorNode {
	return &LogicalOROperatorNode{a: a, b: b}
}

// This is different from most of the evaluator functions in that it does
// short-circuiting: since is logical OR, the second argument is not evaluated
// if the first argument is false.
//
// See the disposition-matrix discussion for LogicalANDOperator.
func (this *LogicalOROperatorNode) Evaluate(state *runtime.State) types.Mlrval {
	aout := this.a.Evaluate(state)
	atype := aout.GetType()
	if !(atype == types.MT_ABSENT || atype == types.MT_BOOL) {
		return types.MlrvalFromError()
	}
	if atype == types.MT_ABSENT {
		return aout
	}
	if aout.IsTrue() {
		// This means true || bogus type evaluates to true, which is sad but
		// which we MUST do in order to not violate the short-circuiting
		// property.  We would have to evaluate b to know if it were error or
		// not.
		return aout
	}

	bout := this.b.Evaluate(state)
	btype := bout.GetType()
	if !(btype == types.MT_ABSENT || btype == types.MT_BOOL) {
		return types.MlrvalFromError()
	}
	if btype == types.MT_ABSENT {
		return bout
	}
	// xxx temp
	output := types.MlrvalFromError()
	types.MlrvalLogicalOR(&output, &aout, &bout)
	return output
}

// ================================================================
// a ?? b evaluates to b only when a is absent. Example: '$foo ?? 0' when the
// current record has no field $foo.
type AbsentCoalesceOperatorNode struct{ a, b IEvaluable }

func (this *RootNode) BuildAbsentCoalesceOperatorNode(a, b IEvaluable) *AbsentCoalesceOperatorNode {
	return &AbsentCoalesceOperatorNode{a: a, b: b}
}

// This is different from most of the evaluator functions in that it does
// short-circuiting: the second argument is not evaluated if the first
// argument is not absent.
func (this *AbsentCoalesceOperatorNode) Evaluate(state *runtime.State) types.Mlrval {
	aout := this.a.Evaluate(state)
	atype := aout.GetType()
	if atype != types.MT_ABSENT {
		return aout
	}

	bout := this.b.Evaluate(state)
	return bout
}

// ================================================================
// a ?? b evaluates to b only when a is absent or empty. Example: '$foo ?? 0'
// when the current record has no field $foo, or when $foo is empty..
type EmptyCoalesceOperatorNode struct{ a, b IEvaluable }

func (this *RootNode) BuildEmptyCoalesceOperatorNode(a, b IEvaluable) *EmptyCoalesceOperatorNode {
	return &EmptyCoalesceOperatorNode{a: a, b: b}
}

// This is different from most of the evaluator functions in that it does
// short-circuiting: the second argument is not evaluated if the first
// argument is not absent.
func (this *EmptyCoalesceOperatorNode) Evaluate(state *runtime.State) types.Mlrval {
	aout := this.a.Evaluate(state)
	atype := aout.GetType()
	if atype == types.MT_ABSENT || atype == types.MT_VOID || (atype == types.MT_STRING && aout.String() == "") {
		bout := this.b.Evaluate(state)
		return bout
	} else {
		return aout
	}
}

// ================================================================
type StandardTernaryOperatorNode struct{ a, b, c IEvaluable }

func (this *RootNode) BuildStandardTernaryOperatorNode(a, b, c IEvaluable) *StandardTernaryOperatorNode {
	return &StandardTernaryOperatorNode{a: a, b: b, c: c}
}
func (this *StandardTernaryOperatorNode) Evaluate(state *runtime.State) types.Mlrval {
	aout := this.a.Evaluate(state)

	boolValue, isBool := aout.GetBoolValue()
	if !isBool {
		return types.MlrvalFromError()
	}

	// Short-circuit: defer evaluation unless needed
	if boolValue == true {
		return this.b.Evaluate(state)
	} else {
		return this.c.Evaluate(state)
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

func BinaryShortCircuitPlaceholder(output, input1, input2 *types.Mlrval) {
	lib.InternalCodingErrorPanic("Short-circuting was not correctly implemented")
	output.SetFromError() // not reached
}

func TernaryShortCircuitPlaceholder(output, input1, input2, input3 *types.Mlrval) {
	lib.InternalCodingErrorPanic("Short-circuting was not correctly implemented")
	output.SetFromError() // not reached
}
