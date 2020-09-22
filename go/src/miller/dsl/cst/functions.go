package cst

import (
	"errors"
	"fmt"

	"miller/dsl"
	"miller/lib"
	"miller/types"
)

// ================================================================
// CST build/execute for AST operator/function nodes
//
// Operators and functions are semantically the same thing -- they differ only
// syntactically. Binary operators are infix, like '1+2', while functions are
// prefix, like 'max(1,2)'. Both parse to the same AST shape.
// ================================================================

func BuildFunctionCallsiteNode(astNode *dsl.ASTNode) (IEvaluable, error) {
	lib.InternalCodingErrorIf(
		astNode.Type != dsl.NodeTypeFunctionCallsite &&
			astNode.Type != dsl.NodeTypeOperator,
	)
	lib.InternalCodingErrorIf(astNode.Token == nil)
	lib.InternalCodingErrorIf(astNode.Children == nil)

	functionName := string(astNode.Token.Lit)

	functionInfo := BuiltinFunctionManager.LookUp(functionName)
	if functionInfo != nil {
		if functionInfo.hasMultipleArities { // E.g. "+" and "-"
			return BuildMultipleArityFunctionCallsiteNode(astNode, functionInfo)
		} else if functionInfo.zaryFunc != nil {
			return BuildZaryFunctionCallsiteNode(astNode, functionInfo)
		} else if functionInfo.unaryFunc != nil {
			return BuildUnaryFunctionCallsiteNode(astNode, functionInfo)
		} else if functionInfo.binaryFunc != nil {
			return BuildBinaryFunctionCallsiteNode(astNode, functionInfo)
		} else if functionInfo.ternaryFunc != nil {
			return BuildTernaryFunctionCallsiteNode(astNode, functionInfo)
		} else if functionInfo.variadicFunc != nil {
			return BuildVariadicFunctionCallsiteNode(astNode, functionInfo)
		} else {
			return nil, errors.New(
				"CST BuildFunctionCallsiteNode: function not implemented yet: " +
					functionName,
			)
		}
	} else {
		return nil, errors.New(
			"CST BuildFunctionCallsiteNode: function name not found: " +
				functionName,
		)
	}
}

// ----------------------------------------------------------------
func BuildMultipleArityFunctionCallsiteNode(
	astNode *dsl.ASTNode,
	functionInfo *FunctionInfo,
) (IEvaluable, error) {
	callsiteArity := len(astNode.Children)
	if callsiteArity == 1 && functionInfo.unaryFunc != nil {
		return BuildUnaryFunctionCallsiteNode(astNode, functionInfo)
	}
	if callsiteArity == 2 && functionInfo.binaryFunc != nil {
		return BuildBinaryFunctionCallsiteNode(astNode, functionInfo)
	}

	return nil, errors.New(
		fmt.Sprintf(
			"CST BuildMultipleArityFunctionCallsiteNode: function name not found: " +
				functionInfo.name,
		),
	)
}

// ----------------------------------------------------------------
type ZaryFunctionCallsiteNode struct {
	zaryFunc types.ZaryFunc
}

func BuildZaryFunctionCallsiteNode(
	astNode *dsl.ASTNode,
	functionInfo *FunctionInfo,
) (IEvaluable, error) {
	callsiteArity := len(astNode.Children)
	expectedArity := 0
	if callsiteArity != expectedArity {
		return nil, errors.New(
			fmt.Sprintf(
				"Miller: function %s invoked with %d argument%s; expected %d",
				functionInfo.name,
				callsiteArity,
				lib.Plural(callsiteArity),
				expectedArity,
			),
		)
	}

	return &ZaryFunctionCallsiteNode{zaryFunc: functionInfo.zaryFunc}, nil
}

func (this *ZaryFunctionCallsiteNode) Evaluate(state *State) types.Mlrval {
	return this.zaryFunc()
}

// ----------------------------------------------------------------
type UnaryFunctionCallsiteNode struct {
	unaryFunc  types.UnaryFunc
	evaluable1 IEvaluable
}

func BuildUnaryFunctionCallsiteNode(
	astNode *dsl.ASTNode,
	functionInfo *FunctionInfo,
) (IEvaluable, error) {
	callsiteArity := len(astNode.Children)
	expectedArity := 1
	if callsiteArity != expectedArity {
		return nil, errors.New(
			fmt.Sprintf(
				"Miller: function %s invoked with %d argument%s; expected %d",
				functionInfo.name,
				callsiteArity,
				lib.Plural(callsiteArity),
				expectedArity,
			),
		)
	}

	evaluable1, err := BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}

	return &UnaryFunctionCallsiteNode{
		unaryFunc:  functionInfo.unaryFunc,
		evaluable1: evaluable1,
	}, nil
}

func (this *UnaryFunctionCallsiteNode) Evaluate(state *State) types.Mlrval {
	arg1 := this.evaluable1.Evaluate(state)
	return this.unaryFunc(&arg1)
}

// ----------------------------------------------------------------
type BinaryFunctionCallsiteNode struct {
	binaryFunc types.BinaryFunc
	evaluable1 IEvaluable
	evaluable2 IEvaluable
}

func BuildBinaryFunctionCallsiteNode(
	astNode *dsl.ASTNode,
	functionInfo *FunctionInfo,
) (IEvaluable, error) {
	callsiteArity := len(astNode.Children)
	expectedArity := 2
	if callsiteArity != expectedArity {
		return nil, errors.New(
			fmt.Sprintf(
				"Miller: function %s invoked with %d argument%s; expected %d",
				functionInfo.name,
				callsiteArity,
				lib.Plural(callsiteArity),
				expectedArity,
			),
		)
	}

	evaluable1, err := BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	evaluable2, err := BuildEvaluableNode(astNode.Children[1])
	if err != nil {
		return nil, err
	}

	// Special short-circuiting cases
	if functionInfo.name == "&&" {
		return BuildLogicalANDOperatorNode(
			evaluable1,
			evaluable2,
		), nil
	}
	if functionInfo.name == "||" {
		return BuildLogicalOROperatorNode(
			evaluable1,
			evaluable2,
		), nil
	}

	return &BinaryFunctionCallsiteNode{
		binaryFunc: functionInfo.binaryFunc,
		evaluable1: evaluable1,
		evaluable2: evaluable2,
	}, nil
}

func (this *BinaryFunctionCallsiteNode) Evaluate(state *State) types.Mlrval {
	arg1 := this.evaluable1.Evaluate(state)
	arg2 := this.evaluable2.Evaluate(state)
	return this.binaryFunc(&arg1, &arg2)
}

// ----------------------------------------------------------------
type TernaryFunctionCallsiteNode struct {
	ternaryFunc types.TernaryFunc
	evaluable1  IEvaluable
	evaluable2  IEvaluable
	evaluable3  IEvaluable
}

func BuildTernaryFunctionCallsiteNode(
	astNode *dsl.ASTNode,
	functionInfo *FunctionInfo,
) (IEvaluable, error) {
	callsiteArity := len(astNode.Children)
	expectedArity := 3
	if callsiteArity != expectedArity {
		return nil, errors.New(
			fmt.Sprintf(
				"Miller: function %s invoked with %d argument%s; expected %d",
				functionInfo.name,
				callsiteArity,
				lib.Plural(callsiteArity),
				expectedArity,
			),
		)
	}

	evaluable1, err := BuildEvaluableNode(astNode.Children[0])
	evaluable2, err := BuildEvaluableNode(astNode.Children[1])
	evaluable3, err := BuildEvaluableNode(astNode.Children[2])
	if err != nil {
		return nil, err
	}

	// Special short-circuiting case
	if functionInfo.name == "?:" {
		return BuildStandardTernaryOperatorNode(
			evaluable1,
			evaluable2,
			evaluable3,
		), nil
	}

	return &TernaryFunctionCallsiteNode{
		ternaryFunc: functionInfo.ternaryFunc,
		evaluable1:  evaluable1,
		evaluable2:  evaluable2,
		evaluable3:  evaluable3,
	}, nil
}

func (this *TernaryFunctionCallsiteNode) Evaluate(state *State) types.Mlrval {
	arg1 := this.evaluable1.Evaluate(state)
	arg2 := this.evaluable2.Evaluate(state)
	arg3 := this.evaluable3.Evaluate(state)
	return this.ternaryFunc(&arg1, &arg2, &arg3)
}

// ----------------------------------------------------------------
type VariadicFunctionCallsiteNode struct {
	variadicFunc types.VariadicFunc
	evaluables   []IEvaluable
}

func BuildVariadicFunctionCallsiteNode(
	astNode *dsl.ASTNode,
	functionInfo *FunctionInfo,
) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Children == nil)
	evaluables := make([]IEvaluable, len(astNode.Children))

	var err error = nil
	for i, astChildNode := range astNode.Children {
		evaluables[i], err = BuildEvaluableNode(astChildNode)
		if err != nil {
			return nil, err
		}
	}
	return &VariadicFunctionCallsiteNode{
		variadicFunc: functionInfo.variadicFunc,
		evaluables:   evaluables,
	}, nil
}

func (this *VariadicFunctionCallsiteNode) Evaluate(state *State) types.Mlrval {
	args := make([]*types.Mlrval, len(this.evaluables))
	for i, evaluable := range this.evaluables {
		arg := evaluable.Evaluate(state)
		args[i] = &arg
	}
	return this.variadicFunc(args)
}

// ================================================================
type LogicalANDOperatorNode struct{ a, b IEvaluable }

func BuildLogicalANDOperatorNode(a, b IEvaluable) *LogicalANDOperatorNode {
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

func (this *LogicalANDOperatorNode) Evaluate(state *State) types.Mlrval {
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
	return types.MlrvalLogicalAND(&aout, &bout)
}

// ================================================================
type LogicalOROperatorNode struct{ a, b IEvaluable }

func BuildLogicalOROperatorNode(a, b IEvaluable) *LogicalOROperatorNode {
	return &LogicalOROperatorNode{a: a, b: b}
}

// This is different from most of the evaluator functions in that it does
// short-circuiting: since is logical OR, the second argument is not evaluated
// if the first argumeent is false.
//
// See the disposition-matrix discussion for LogicalANDOperator.
func (this *LogicalOROperatorNode) Evaluate(state *State) types.Mlrval {
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
	return types.MlrvalLogicalOR(&aout, &bout)
}

// ================================================================
type StandardTernaryOperatorNode struct{ a, b, c IEvaluable }

func BuildStandardTernaryOperatorNode(a, b, c IEvaluable) *StandardTernaryOperatorNode {
	return &StandardTernaryOperatorNode{a: a, b: b, c: c}
}
func (this *StandardTernaryOperatorNode) Evaluate(state *State) types.Mlrval {
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

func BinaryShortCircuitPlaceholder(a, b *types.Mlrval) types.Mlrval {
	lib.InternalCodingErrorPanic("Short-circuting was not correctly implemented")
	return types.MlrvalFromError() // not reached
}

func TernaryShortCircuitPlaceholder(a, b, c *types.Mlrval) types.Mlrval {
	lib.InternalCodingErrorPanic("Short-circuting was not correctly implemented")
	return types.MlrvalFromError() // not reached
}
