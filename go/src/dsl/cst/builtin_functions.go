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

	return &ZaryFunctionCallsiteNode{
		zaryFunc: builtinFunctionInfo.zaryFunc,
	}, nil
}

func (this *ZaryFunctionCallsiteNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	this.zaryFunc(output)
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

func (this *UnaryFunctionCallsiteNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	var arg1 types.Mlrval
	this.evaluable1.Evaluate(&arg1, state)
	this.unaryFunc(output, &arg1)
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

func (this *ContextualUnaryFunctionCallsiteNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	var arg1 types.Mlrval
	this.evaluable1.Evaluate(&arg1, state)
	this.contextualUnaryFunc(output, &arg1, state.Context)
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

func (this *BinaryFunctionCallsiteNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	var arg1, arg2 types.Mlrval
	this.evaluable1.Evaluate(&arg1, state)
	this.evaluable2.Evaluate(&arg2, state)
	this.binaryFunc(output, &arg1, &arg2)
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

func (this *TernaryFunctionCallsiteNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	var arg1, arg2, arg3 types.Mlrval
	this.evaluable1.Evaluate(&arg1, state)
	this.evaluable2.Evaluate(&arg2, state)
	this.evaluable3.Evaluate(&arg3, state)
	this.ternaryFunc(output, &arg1, &arg2, &arg3)
}

// ----------------------------------------------------------------
type VariadicFunctionCallsiteNode struct {
	variadicFunc types.VariadicFunc
	evaluables   []IEvaluable
	args         []*types.Mlrval
}

func (this *RootNode) BuildVariadicFunctionCallsiteNode(
	astNode *dsl.ASTNode,
	builtinFunctionInfo *BuiltinFunctionInfo,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Children == nil)
	evaluables := make([]IEvaluable, len(astNode.Children))
	args := make([]*types.Mlrval, len(astNode.Children))

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
		args[i] = types.MlrvalPointerFromError()
	}
	return &VariadicFunctionCallsiteNode{
		variadicFunc: builtinFunctionInfo.variadicFunc,
		evaluables:   evaluables,
		args:         args,
	}, nil
}

func (this *VariadicFunctionCallsiteNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	for i, _ := range this.evaluables {
		this.evaluables[i].Evaluate(this.args[i], state)
	}
	this.variadicFunc(output, this.args)
}

// ================================================================
type LogicalANDOperatorNode struct {
	a, b IEvaluable
}

func (this *RootNode) BuildLogicalANDOperatorNode(a, b IEvaluable) *LogicalANDOperatorNode {
	return &LogicalANDOperatorNode{
		a: a,
		b: b,
	}
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

func (this *LogicalANDOperatorNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	var aout, bout types.Mlrval
	this.a.Evaluate(&aout, state)
	atype := aout.GetType()
	if !(atype == types.MT_ABSENT || atype == types.MT_BOOL) {
		output.SetFromError()
		return
	}
	if atype == types.MT_ABSENT {
		output.SetFromAbsent()
		return
	}
	if aout.IsFalse() {
		// This means false && bogus type evaluates to true, which is sad but
		// which we MUST do in order to not violate the short-circuiting
		// property.  We would have to evaluate b to know if it were error or
		// not.
		output.CopyFrom(&aout)
		return
	}

	this.b.Evaluate(&bout, state)
	btype := bout.GetType()
	if !(btype == types.MT_ABSENT || btype == types.MT_BOOL) {
		output.SetFromError()
		return
	}
	if btype == types.MT_ABSENT {
		output.SetFromAbsent()
		return
	}
	types.MlrvalLogicalAND(output, &aout, &bout)
}

// ================================================================
type LogicalOROperatorNode struct {
	a, b IEvaluable
}

func (this *RootNode) BuildLogicalOROperatorNode(a, b IEvaluable) *LogicalOROperatorNode {
	return &LogicalOROperatorNode{
		a: a,
		b: b,
	}
}

// This is different from most of the evaluator functions in that it does
// short-circuiting: since is logical OR, the second argument is not evaluated
// if the first argument is false.
//
// See the disposition-matrix discussion for LogicalANDOperator.
func (this *LogicalOROperatorNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	var aout, bout types.Mlrval
	this.a.Evaluate(&aout, state)
	atype := aout.GetType()
	if !(atype == types.MT_ABSENT || atype == types.MT_BOOL) {
		output.SetFromError()
		return
	}
	if atype == types.MT_ABSENT {
		output.SetFromAbsent()
		return
	}
	if aout.IsTrue() {
		// This means true || bogus type evaluates to true, which is sad but
		// which we MUST do in order to not violate the short-circuiting
		// property.  We would have to evaluate b to know if it were error or
		// not.
		output.CopyFrom(&aout)
		return
	}

	this.b.Evaluate(&bout, state)
	btype := bout.GetType()
	if !(btype == types.MT_ABSENT || btype == types.MT_BOOL) {
		output.SetFromError()
		return
	}
	if btype == types.MT_ABSENT {
		output.SetFromAbsent()
		return
	}
	types.MlrvalLogicalOR(output, &aout, &bout)
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
func (this *AbsentCoalesceOperatorNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	this.a.Evaluate(output, state)
	if output.GetType() != types.MT_ABSENT {
		return
	}

	this.b.Evaluate(output, state)
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
func (this *EmptyCoalesceOperatorNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	this.a.Evaluate(output, state)
	atype := output.GetType()
	if atype == types.MT_ABSENT || atype == types.MT_VOID || (atype == types.MT_STRING && output.String() == "") {
		this.b.Evaluate(output, state)
	}
}

// ================================================================
type StandardTernaryOperatorNode struct{ a, b, c IEvaluable }

func (this *RootNode) BuildStandardTernaryOperatorNode(a, b, c IEvaluable) *StandardTernaryOperatorNode {
	return &StandardTernaryOperatorNode{a: a, b: b, c: c}
}
func (this *StandardTernaryOperatorNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	this.a.Evaluate(output, state)

	boolValue, isBool := output.GetBoolValue()
	if !isBool {
		output.SetFromError()
		return
	}

	// Short-circuit: defer evaluation unless needed
	if boolValue == true {
		this.b.Evaluate(output, state)
	} else {
		this.c.Evaluate(output, state)
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
