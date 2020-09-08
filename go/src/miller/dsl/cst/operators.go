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
		return BuildUnaryPlusOperatorNode(cstChild), nil
		break
	case "-":
		return BuildUnaryMinusOperatorNode(cstChild), nil
		break
	case "~":
		return BuildBitwiseNOTOperatorNode(cstChild), nil
		break
	case "!":
		return BuildLogicalNOTOperatorNode(cstChild), nil
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
	case ".":
		return BuildDotOperatorNode(leftCSTChild, rightCSTChild), nil
		break

	case "+":
		return BuildPlusOperatorNode(leftCSTChild, rightCSTChild), nil
		break
	case "-":
		return BuildMinusOperatorNode(leftCSTChild, rightCSTChild), nil
		break
	case "*":
		return BuildTimesOperatorNode(leftCSTChild, rightCSTChild), nil
		break
	case "/":
		return BuildDivideOperatorNode(leftCSTChild, rightCSTChild), nil
		break
	case "//":
		return BuildIntDivideOperatorNode(leftCSTChild, rightCSTChild), nil
		break
	case "**":
		return BuildPowOperatorNode(leftCSTChild, rightCSTChild), nil
		break

	case ".+":
		return BuildDotPlusOperatorNode(leftCSTChild, rightCSTChild), nil
		break
	case ".-":
		return BuildDotMinusOperatorNode(leftCSTChild, rightCSTChild), nil
		break
	case ".*":
		return BuildDotTimesOperatorNode(leftCSTChild, rightCSTChild), nil
		break
	case "./":
		return BuildDotDivideOperatorNode(leftCSTChild, rightCSTChild), nil
		break

	case "%":
		return BuildModulusOperatorNode(leftCSTChild, rightCSTChild), nil
		break

	case "&":
		return BuildBitwiseANDOperatorNode(leftCSTChild, rightCSTChild), nil
		break
	case "|":
		return BuildBitwiseOROperatorNode(leftCSTChild, rightCSTChild), nil
		break
	case "^":
		return BuildBitwiseXOROperatorNode(leftCSTChild, rightCSTChild), nil
		break

	// TO DO: implement short-circuiting for these, as special cases.
	case "&&":
		return BuildLogicalANDOperatorNode(leftCSTChild, rightCSTChild), nil
		break
	case "||":
		return BuildLogicalOROperatorNode(leftCSTChild, rightCSTChild), nil
		break
	case "^^":
		return BuildLogicalXOROperatorNode(leftCSTChild, rightCSTChild), nil
		break

	case "==":
		return BuildEqualsOperatorNode(leftCSTChild, rightCSTChild), nil
		break
	case "!=":
		return BuildNotEqualsOperatorNode(leftCSTChild, rightCSTChild), nil
		break
	case ">":
		return BuildGreaterThanOperatorNode(leftCSTChild, rightCSTChild), nil
		break
	case ">=":
		return BuildGreaterThanOrEqualsOperatorNode(leftCSTChild, rightCSTChild), nil
		break
	case "<":
		return BuildLessThanOperatorNode(leftCSTChild, rightCSTChild), nil
		break
	case "<=":
		return BuildLessThanOrEqualsOperatorNode(leftCSTChild, rightCSTChild), nil
		break

		// xxx continue ...
	}

	return nil, errors.New(
		"CST BuildBinaryOperatorNode: unhandled AST node " + string(astNode.Type),
	)
}

// ----------------------------------------------------------------
// TODO: Look into short-circuiting
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
type UnaryPlusOperatorNode struct{ a IEvaluable }

func BuildUnaryPlusOperatorNode(a IEvaluable) *UnaryPlusOperatorNode {
	return &UnaryPlusOperatorNode{a: a}
}
func (this *UnaryPlusOperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	return lib.MlrvalUnaryPlus(&aout)
}

// ----------------------------------------------------------------
type UnaryMinusOperatorNode struct{ a IEvaluable }

func BuildUnaryMinusOperatorNode(a IEvaluable) *UnaryMinusOperatorNode {
	return &UnaryMinusOperatorNode{a: a}
}
func (this *UnaryMinusOperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	return lib.MlrvalUnaryMinus(&aout)
}

// ================================================================
type DotOperatorNode struct{ a, b IEvaluable }

func BuildDotOperatorNode(a, b IEvaluable) *DotOperatorNode {
	return &DotOperatorNode{a: a, b: b}
}
func (this *DotOperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalDot(&aout, &bout)
}

// ----------------------------------------------------------------
type PlusOperatorNode struct{ a, b IEvaluable }

func BuildPlusOperatorNode(a, b IEvaluable) *PlusOperatorNode {
	return &PlusOperatorNode{a: a, b: b}
}
func (this *PlusOperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalPlus(&aout, &bout)
}

// ----------------------------------------------------------------
type MinusOperatorNode struct{ a, b IEvaluable }

func BuildMinusOperatorNode(a, b IEvaluable) *MinusOperatorNode {
	return &MinusOperatorNode{a: a, b: b}
}
func (this *MinusOperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalMinus(&aout, &bout)
}

// ----------------------------------------------------------------
type TimesOperatorNode struct{ a, b IEvaluable }

func BuildTimesOperatorNode(a, b IEvaluable) *TimesOperatorNode {
	return &TimesOperatorNode{a: a, b: b}
}
func (this *TimesOperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalTimes(&aout, &bout)
}

// ----------------------------------------------------------------
type DivideOperatorNode struct{ a, b IEvaluable }

func BuildDivideOperatorNode(a, b IEvaluable) *DivideOperatorNode {
	return &DivideOperatorNode{a: a, b: b}
}
func (this *DivideOperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalDivide(&aout, &bout)
}

// ----------------------------------------------------------------
type IntDivideOperatorNode struct{ a, b IEvaluable }

func BuildIntDivideOperatorNode(a, b IEvaluable) *IntDivideOperatorNode {
	return &IntDivideOperatorNode{a: a, b: b}
}
func (this *IntDivideOperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalIntDivide(&aout, &bout)
}

// ----------------------------------------------------------------
type PowOperatorNode struct{ a, b IEvaluable }

func BuildPowOperatorNode(a, b IEvaluable) *PowOperatorNode {
	return &PowOperatorNode{a: a, b: b}
}
func (this *PowOperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalPow(&aout, &bout)
}

// ----------------------------------------------------------------
type DotPlusOperatorNode struct{ a, b IEvaluable }

func BuildDotPlusOperatorNode(a, b IEvaluable) *DotPlusOperatorNode {
	return &DotPlusOperatorNode{a: a, b: b}
}
func (this *DotPlusOperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalDotPlus(&aout, &bout)
}

// ----------------------------------------------------------------
type DotMinusOperatorNode struct{ a, b IEvaluable }

func BuildDotMinusOperatorNode(a, b IEvaluable) *DotMinusOperatorNode {
	return &DotMinusOperatorNode{a: a, b: b}
}
func (this *DotMinusOperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalDotMinus(&aout, &bout)
}

// ----------------------------------------------------------------
type DotTimesOperatorNode struct{ a, b IEvaluable }

func BuildDotTimesOperatorNode(a, b IEvaluable) *DotTimesOperatorNode {
	return &DotTimesOperatorNode{a: a, b: b}
}
func (this *DotTimesOperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalDotTimes(&aout, &bout)
}

// ----------------------------------------------------------------
type DotDivideOperatorNode struct{ a, b IEvaluable }

func BuildDotDivideOperatorNode(a, b IEvaluable) *DotDivideOperatorNode {
	return &DotDivideOperatorNode{a: a, b: b}
}
func (this *DotDivideOperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalDotDivide(&aout, &bout)
}

// ----------------------------------------------------------------
type ModulusOperatorNode struct{ a, b IEvaluable }

func BuildModulusOperatorNode(a, b IEvaluable) *ModulusOperatorNode {
	return &ModulusOperatorNode{a: a, b: b}
}
func (this *ModulusOperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalModulus(&aout, &bout)
}

// ----------------------------------------------------------------
type BitwiseANDOperatorNode struct{ a, b IEvaluable }

func BuildBitwiseANDOperatorNode(a, b IEvaluable) *BitwiseANDOperatorNode {
	return &BitwiseANDOperatorNode{a: a, b: b}
}
func (this *BitwiseANDOperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalBitwiseAND(&aout, &bout)
}

// ----------------------------------------------------------------
type BitwiseOROperatorNode struct{ a, b IEvaluable }

func BuildBitwiseOROperatorNode(a, b IEvaluable) *BitwiseOROperatorNode {
	return &BitwiseOROperatorNode{a: a, b: b}
}
func (this *BitwiseOROperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalBitwiseOR(&aout, &bout)
}

// ----------------------------------------------------------------
type BitwiseXOROperatorNode struct{ a, b IEvaluable }

func BuildBitwiseXOROperatorNode(a, b IEvaluable) *BitwiseXOROperatorNode {
	return &BitwiseXOROperatorNode{a: a, b: b}
}
func (this *BitwiseXOROperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalBitwiseXOR(&aout, &bout)
}

// ----------------------------------------------------------------
type BitwiseNOTOperatorNode struct{ a IEvaluable }

func BuildBitwiseNOTOperatorNode(a IEvaluable) *BitwiseNOTOperatorNode {
	return &BitwiseNOTOperatorNode{a: a}
}
func (this *BitwiseNOTOperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	return lib.MlrvalBitwiseNOT(&aout)
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

// ----------------------------------------------------------------
type LogicalXOROperatorNode struct{ a, b IEvaluable }

func BuildLogicalXOROperatorNode(a, b IEvaluable) *LogicalXOROperatorNode {
	return &LogicalXOROperatorNode{a: a, b: b}
}
func (this *LogicalXOROperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalLogicalXOR(&aout, &bout)
}

// ----------------------------------------------------------------
type LogicalNOTOperatorNode struct{ a IEvaluable }

func BuildLogicalNOTOperatorNode(a IEvaluable) *LogicalNOTOperatorNode {
	return &LogicalNOTOperatorNode{a: a}
}
func (this *LogicalNOTOperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	return lib.MlrvalLogicalNOT(&aout)
}

// ----------------------------------------------------------------
type EqualsOperatorNode struct{ a, b IEvaluable }

func BuildEqualsOperatorNode(a, b IEvaluable) *EqualsOperatorNode {
	return &EqualsOperatorNode{a: a, b: b}
}
func (this *EqualsOperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalEquals(&aout, &bout)
}

// ----------------------------------------------------------------
type NotEqualsOperatorNode struct{ a, b IEvaluable }

func BuildNotEqualsOperatorNode(a, b IEvaluable) *NotEqualsOperatorNode {
	return &NotEqualsOperatorNode{a: a, b: b}
}
func (this *NotEqualsOperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalNotEquals(&aout, &bout)
}

// ----------------------------------------------------------------
type GreaterThanOperatorNode struct{ a, b IEvaluable }

func BuildGreaterThanOperatorNode(a, b IEvaluable) *GreaterThanOperatorNode {
	return &GreaterThanOperatorNode{a: a, b: b}
}
func (this *GreaterThanOperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalGreaterThan(&aout, &bout)
}

// ----------------------------------------------------------------
type GreaterThanOrEqualsOperatorNode struct{ a, b IEvaluable }

func BuildGreaterThanOrEqualsOperatorNode(a, b IEvaluable) *GreaterThanOrEqualsOperatorNode {
	return &GreaterThanOrEqualsOperatorNode{a: a, b: b}
}
func (this *GreaterThanOrEqualsOperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalGreaterThanOrEquals(&aout, &bout)
}

// ----------------------------------------------------------------
type LessThanOperatorNode struct{ a, b IEvaluable }

func BuildLessThanOperatorNode(a, b IEvaluable) *LessThanOperatorNode {
	return &LessThanOperatorNode{a: a, b: b}
}
func (this *LessThanOperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalLessThan(&aout, &bout)
}

// ----------------------------------------------------------------
type LessThanOrEqualsOperatorNode struct{ a, b IEvaluable }

func BuildLessThanOrEqualsOperatorNode(a, b IEvaluable) *LessThanOrEqualsOperatorNode {
	return &LessThanOrEqualsOperatorNode{a: a, b: b}
}
func (this *LessThanOrEqualsOperatorNode) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalLessThanOrEquals(&aout, &bout)
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
