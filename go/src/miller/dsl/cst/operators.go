package cst

import (
	"errors"

	"miller/dsl"
	"miller/lib"
)

// ================================================================
// CST build/execute for AST operator nodes
// ================================================================

// ================================================================
func BuildOperatorNode(astNode *dsl.ASTNode) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeOperator)

	arity := len(astNode.Children)
	switch arity {
	case 1:
		return BuildUnaryOperatorNode(astNode)
		break
	case 2:
		return BuildBinaryOperatorNode(astNode)
		break
	case 3:
		return BuildTernaryOperatorNode(astNode)
		break
	}
	return nil, errors.New(
		"CST BuildOperatorNode: unhandled AST node " + string(astNode.Type),
	)
}

// ----------------------------------------------------------------
func BuildUnaryOperatorNode(astNode *dsl.ASTNode) (IEvaluable, error) {
	arity := len(astNode.Children)
	lib.InternalCodingErrorIf(arity != 1)
	astChild := astNode.Children[0]

	cstChild, err := BuildEvaluableNode(astChild)
	if err != nil {
		return nil, err
	}

	sop := string(astNode.Token.Lit)
	switch sop {
	case "+":
		return BuildUnaryFunctionNode(cstChild, lib.MlrvalUnaryPlus), nil
		break
	case "-":
		return BuildUnaryFunctionNode(cstChild, lib.MlrvalUnaryMinus), nil
		break
	case "~":
		return BuildUnaryFunctionNode(cstChild, lib.MlrvalBitwiseNOT), nil
		break
	case "!":
		return BuildUnaryFunctionNode(cstChild, lib.MlrvalLogicalNOT), nil
		break
	}

	return nil, errors.New(
		"CST BuildUnaryOperatorNode: unhandled AST node " + string(astNode.Type),
	)
}

// ----------------------------------------------------------------
func BuildBinaryOperatorNode(astNode *dsl.ASTNode) (IEvaluable, error) {
	arity := len(astNode.Children)
	lib.InternalCodingErrorIf(arity != 2)

	leftASTChild := astNode.Children[0]
	rightASTChild := astNode.Children[1]

	leftCSTChild, err := BuildEvaluableNode(leftASTChild)
	if err != nil {
		return nil, err
	}
	rightCSTChild, err := BuildEvaluableNode(rightASTChild)
	if err != nil {
		return nil, err
	}

	sop := string(astNode.Token.Lit)
	switch sop {

	// xxx lookup table in function_manager:
	// name:
	// * help string
	// * binaryFunc ptr

	case ".":
		return BuildBinaryFunctionNode(leftCSTChild, rightCSTChild, lib.MlrvalDot), nil
		break
	case "+":
		return BuildBinaryFunctionNode(leftCSTChild, rightCSTChild, lib.MlrvalPlus), nil
		break
	case "-":
		return BuildBinaryFunctionNode(leftCSTChild, rightCSTChild, lib.MlrvalMinus), nil
		break
	case "*":
		return BuildBinaryFunctionNode(leftCSTChild, rightCSTChild, lib.MlrvalTimes), nil
		break
	case "/":
		return BuildBinaryFunctionNode(leftCSTChild, rightCSTChild, lib.MlrvalDivide), nil
		break
	case "//":
		return BuildBinaryFunctionNode(leftCSTChild, rightCSTChild, lib.MlrvalIntDivide), nil
		break
	case "**":
		return BuildBinaryFunctionNode(leftCSTChild, rightCSTChild, lib.MlrvalPow), nil
		break
	case ".+":
		return BuildBinaryFunctionNode(leftCSTChild, rightCSTChild, lib.MlrvalDotPlus), nil
		break
	case ".-":
		return BuildBinaryFunctionNode(leftCSTChild, rightCSTChild, lib.MlrvalDotMinus), nil
		break
	case ".*":
		return BuildBinaryFunctionNode(leftCSTChild, rightCSTChild, lib.MlrvalDotTimes), nil
		break
	case "./":
		return BuildBinaryFunctionNode(leftCSTChild, rightCSTChild, lib.MlrvalDotDivide), nil
		break
	case "%":
		return BuildBinaryFunctionNode(leftCSTChild, rightCSTChild, lib.MlrvalModulus), nil
		break
	case "&":
		return BuildBinaryFunctionNode(leftCSTChild, rightCSTChild, lib.MlrvalBitwiseAND), nil
		break
	case "|":
		return BuildBinaryFunctionNode(leftCSTChild, rightCSTChild, lib.MlrvalBitwiseOR), nil
		break
	case "^":
		return BuildBinaryFunctionNode(leftCSTChild, rightCSTChild, lib.MlrvalBitwiseXOR), nil
		break
	case "<<":
		return BuildBinaryFunctionNode(leftCSTChild, rightCSTChild, lib.MlrvalLeftShift), nil
		break
	case ">>":
		return BuildBinaryFunctionNode(leftCSTChild, rightCSTChild, lib.MlrvalSignedRightShift), nil
		break
	case ">>>":
		return BuildBinaryFunctionNode(leftCSTChild, rightCSTChild, lib.MlrvalUnsignedRightShift), nil
		break

	case "&&":
		return BuildLogicalANDOperatorNode(leftCSTChild, rightCSTChild), nil
		break
	case "||":
		return BuildLogicalOROperatorNode(leftCSTChild, rightCSTChild), nil
		break
	case "^^":
		return BuildBinaryFunctionNode(leftCSTChild, rightCSTChild, lib.MlrvalLogicalXOR), nil
		break
	case "==":
		return BuildBinaryFunctionNode(leftCSTChild, rightCSTChild, lib.MlrvalEquals), nil
		break
	case "!=":
		return BuildBinaryFunctionNode(leftCSTChild, rightCSTChild, lib.MlrvalNotEquals), nil
		break
	case ">":
		return BuildBinaryFunctionNode(leftCSTChild, rightCSTChild, lib.MlrvalGreaterThan), nil
		break
	case ">=":
		return BuildBinaryFunctionNode(leftCSTChild, rightCSTChild, lib.MlrvalGreaterThanOrEquals), nil
		break
	case "<":
		return BuildBinaryFunctionNode(leftCSTChild, rightCSTChild, lib.MlrvalLessThan), nil
		break
	case "<=":
		return BuildBinaryFunctionNode(leftCSTChild, rightCSTChild, lib.MlrvalLessThanOrEquals), nil
		break

		// xxx continue ...
	}

	return nil, errors.New(
		"CST BuildBinaryOperatorNode: unhandled AST node " + string(astNode.Type),
	)
}

// ----------------------------------------------------------------
func BuildTernaryOperatorNode(astNode *dsl.ASTNode) (IEvaluable, error) {
	arity := len(astNode.Children)
	lib.InternalCodingErrorIf(arity != 3)

	leftASTChild := astNode.Children[0]
	middleASTChild := astNode.Children[1]
	rightASTChild := astNode.Children[2]

	leftCSTChild, err := BuildEvaluableNode(leftASTChild)
	if err != nil {
		return nil, err
	}
	middleCSTChild, err := BuildEvaluableNode(middleASTChild)
	if err != nil {
		return nil, err
	}
	rightCSTChild, err := BuildEvaluableNode(rightASTChild)
	if err != nil {
		return nil, err
	}

	sop := string(astNode.Token.Lit)
	switch sop {
	case "?:":
		return BuildStandardTernaryOperatorNode(leftCSTChild, middleCSTChild, rightCSTChild), nil
		break
	}

	return nil, errors.New(
		"CST BuildTernnaryOperatorNode: unhandled AST node " + string(astNode.Type),
	)
}

// ================================================================
type UnaryFunctionNode struct {
	a         IEvaluable
	unaryFunc lib.UnaryFunc
}

func BuildUnaryFunctionNode(
	a IEvaluable,
	unaryFunc lib.UnaryFunc,
) *UnaryFunctionNode {
	return &UnaryFunctionNode{a: a, unaryFunc: unaryFunc}
}

func (this *UnaryFunctionNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	return this.unaryFunc(&aout)
}

// ================================================================
type BinaryFunctionNode struct {
	a, b       IEvaluable
	binaryFunc lib.BinaryFunc
}

func BuildBinaryFunctionNode(
	a, b IEvaluable,
	binaryFunc lib.BinaryFunc,
) *BinaryFunctionNode {
	return &BinaryFunctionNode{a: a, b: b, binaryFunc: binaryFunc}
}

func (this *BinaryFunctionNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return this.binaryFunc(&aout, &bout)
}

// ----------------------------------------------------------------
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

func (this *LogicalANDOperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	atype := aout.GetType()
	if !(atype == lib.MT_ABSENT || atype == lib.MT_BOOL) {
		return lib.MlrvalFromError()
	}
	if atype == lib.MT_ABSENT {
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
	if !(btype == lib.MT_ABSENT || btype == lib.MT_BOOL) {
		return lib.MlrvalFromError()
	}
	if btype == lib.MT_ABSENT {
		return bout
	}
	return lib.MlrvalLogicalAND(&aout, &bout)
}

// ----------------------------------------------------------------
type LogicalOROperatorNode struct{ a, b IEvaluable }

func BuildLogicalOROperatorNode(a, b IEvaluable) *LogicalOROperatorNode {
	return &LogicalOROperatorNode{a: a, b: b}
}

// This is different from most of the evaluator functions in that it does
// short-circuiting: since is logical OR, the second argument is not evaluated
// if the first argumeent is false.
//
// See the disposition-matrix discussion for LogicalANDOperator.
func (this *LogicalOROperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	atype := aout.GetType()
	if !(atype == lib.MT_ABSENT || atype == lib.MT_BOOL) {
		return lib.MlrvalFromError()
	}
	if atype == lib.MT_ABSENT {
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
	if !(btype == lib.MT_ABSENT || btype == lib.MT_BOOL) {
		return lib.MlrvalFromError()
	}
	if btype == lib.MT_ABSENT {
		return bout
	}
	return lib.MlrvalLogicalOR(&aout, &bout)
}

// ================================================================
type StandardTernaryOperatorNode struct{ a, b, c IEvaluable }

func BuildStandardTernaryOperatorNode(a, b, c IEvaluable) *StandardTernaryOperatorNode {
	return &StandardTernaryOperatorNode{a: a, b: b, c: c}
}
func (this *StandardTernaryOperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)

	boolValue, isBoolean := aout.GetBoolValue()
	if !isBoolean {
		return lib.MlrvalFromError()
	}

	// Short-circuit: defer evaluation unless needed
	if boolValue == true {
		return this.b.Evaluate(state)
	} else {
		return this.c.Evaluate(state)
	}
}
