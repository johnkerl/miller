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
func NewOperatorNode(astNode *dsl.ASTNode) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeOperator)

	arity := len(astNode.Children)
	switch arity {
	case 1:
		return NewUnaryOperatorNode(astNode)
		break
	case 2:
		return NewBinaryOperatorNode(astNode)
		break
	case 3:
		return NewTernaryOperatorNode(astNode)
		break
	}
	return nil, errors.New("CST build: AST operator node unhandled.")
}

// ----------------------------------------------------------------
func NewUnaryOperatorNode(astNode *dsl.ASTNode) (IEvaluable, error) {
	arity := len(astNode.Children)
	lib.InternalCodingErrorIf(arity != 1)
	astChild := astNode.Children[0]

	cstChild, err := NewEvaluable(astChild)
	if err != nil {
		return nil, err
	}

	sop := string(astNode.Token.Lit)
	switch sop {
	case "+":
		return NewUnaryPlusOperator(cstChild), nil
		break
	case "-":
		return NewUnaryMinusOperator(cstChild), nil
		break
	case "~":
		return NewBitwiseNOTOperator(cstChild), nil
		break
	case "!":
		return NewLogicalNOTOperator(cstChild), nil
		break
	}

	return nil, errors.New("CST build: AST unary operator node unhandled.")
}

// ----------------------------------------------------------------
func NewBinaryOperatorNode(astNode *dsl.ASTNode) (IEvaluable, error) {
	arity := len(astNode.Children)
	lib.InternalCodingErrorIf(arity != 2)

	leftASTChild := astNode.Children[0]
	rightASTChild := astNode.Children[1]

	leftCSTChild, err := NewEvaluable(leftASTChild)
	if err != nil {
		return nil, err
	}
	rightCSTChild, err := NewEvaluable(rightASTChild)
	if err != nil {
		return nil, err
	}

	sop := string(astNode.Token.Lit)
	switch sop {
	case ".":
		return NewDotOperator(leftCSTChild, rightCSTChild), nil
		break

	case "+":
		return NewPlusOperator(leftCSTChild, rightCSTChild), nil
		break
	case "-":
		return NewMinusOperator(leftCSTChild, rightCSTChild), nil
		break
	case "*":
		return NewTimesOperator(leftCSTChild, rightCSTChild), nil
		break
	case "/":
		return NewDivideOperator(leftCSTChild, rightCSTChild), nil
		break
	case "//":
		return NewIntDivideOperator(leftCSTChild, rightCSTChild), nil
		break
	case "**":
		return NewPowOperator(leftCSTChild, rightCSTChild), nil
		break

	case ".+":
		return NewDotPlusOperator(leftCSTChild, rightCSTChild), nil
		break
	case ".-":
		return NewDotMinusOperator(leftCSTChild, rightCSTChild), nil
		break
	case ".*":
		return NewDotTimesOperator(leftCSTChild, rightCSTChild), nil
		break
	case "./":
		return NewDotDivideOperator(leftCSTChild, rightCSTChild), nil
		break

	case "%":
		return NewModulusOperator(leftCSTChild, rightCSTChild), nil
		break

	case "&":
		return NewBitwiseANDOperator(leftCSTChild, rightCSTChild), nil
		break
	case "|":
		return NewBitwiseOROperator(leftCSTChild, rightCSTChild), nil
		break
	case "^":
		return NewBitwiseXOROperator(leftCSTChild, rightCSTChild), nil
		break

	// TO DO: implement short-circuiting for these, as special cases.
	case "&&":
		return NewLogicalANDOperator(leftCSTChild, rightCSTChild), nil
		break
	case "||":
		return NewLogicalOROperator(leftCSTChild, rightCSTChild), nil
		break
	case "^^":
		return NewLogicalXOROperator(leftCSTChild, rightCSTChild), nil
		break

	case "==":
		return NewEqualsOperator(leftCSTChild, rightCSTChild), nil
		break
	case "!=":
		return NewNotEqualsOperator(leftCSTChild, rightCSTChild), nil
		break
	case ">":
		return NewGreaterThanOperator(leftCSTChild, rightCSTChild), nil
		break
	case ">=":
		return NewGreaterThanOrEqualsOperator(leftCSTChild, rightCSTChild), nil
		break
	case "<":
		return NewLessThanOperator(leftCSTChild, rightCSTChild), nil
		break
	case "<=":
		return NewLessThanOrEqualsOperator(leftCSTChild, rightCSTChild), nil
		break

		// xxx continue ...
	}

	return nil, errors.New(
		"CST build: unandled AST binary operator node \"" + sop + "\"",
	)
}

// ----------------------------------------------------------------
// TODO: Look into short-circuiting
func NewTernaryOperatorNode(astNode *dsl.ASTNode) (IEvaluable, error) {
	arity := len(astNode.Children)
	lib.InternalCodingErrorIf(arity != 3)

	leftASTChild := astNode.Children[0]
	middleASTChild := astNode.Children[1]
	rightASTChild := astNode.Children[2]

	leftCSTChild, err := NewEvaluable(leftASTChild)
	if err != nil {
		return nil, err
	}
	middleCSTChild, err := NewEvaluable(middleASTChild)
	if err != nil {
		return nil, err
	}
	rightCSTChild, err := NewEvaluable(rightASTChild)
	if err != nil {
		return nil, err
	}

	sop := string(astNode.Token.Lit)
	switch sop {
	case "?:":
		return NewTernaryOperator(leftCSTChild, middleCSTChild, rightCSTChild), nil
		break
	}

	return nil, errors.New("CST build: AST ternary operator node unhandled.")
}

// ================================================================
type UnaryPlusOperator struct{ a IEvaluable }

func NewUnaryPlusOperator(a IEvaluable) *UnaryPlusOperator {
	return &UnaryPlusOperator{a: a}
}
func (this *UnaryPlusOperator) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	return lib.MlrvalUnaryPlus(&aout)
}

// ----------------------------------------------------------------
type UnaryMinusOperator struct{ a IEvaluable }

func NewUnaryMinusOperator(a IEvaluable) *UnaryMinusOperator {
	return &UnaryMinusOperator{a: a}
}
func (this *UnaryMinusOperator) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	return lib.MlrvalUnaryMinus(&aout)
}

// ================================================================
type DotOperator struct{ a, b IEvaluable }

func NewDotOperator(a, b IEvaluable) *DotOperator {
	return &DotOperator{a: a, b: b}
}
func (this *DotOperator) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalDot(&aout, &bout)
}

// ----------------------------------------------------------------
type PlusOperator struct{ a, b IEvaluable }

func NewPlusOperator(a, b IEvaluable) *PlusOperator {
	return &PlusOperator{a: a, b: b}
}
func (this *PlusOperator) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalPlus(&aout, &bout)
}

// ----------------------------------------------------------------
type MinusOperator struct{ a, b IEvaluable }

func NewMinusOperator(a, b IEvaluable) *MinusOperator {
	return &MinusOperator{a: a, b: b}
}
func (this *MinusOperator) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalMinus(&aout, &bout)
}

// ----------------------------------------------------------------
type TimesOperator struct{ a, b IEvaluable }

func NewTimesOperator(a, b IEvaluable) *TimesOperator {
	return &TimesOperator{a: a, b: b}
}
func (this *TimesOperator) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalTimes(&aout, &bout)
}

// ----------------------------------------------------------------
type DivideOperator struct{ a, b IEvaluable }

func NewDivideOperator(a, b IEvaluable) *DivideOperator {
	return &DivideOperator{a: a, b: b}
}
func (this *DivideOperator) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalDivide(&aout, &bout)
}

// ----------------------------------------------------------------
type IntDivideOperator struct{ a, b IEvaluable }

func NewIntDivideOperator(a, b IEvaluable) *IntDivideOperator {
	return &IntDivideOperator{a: a, b: b}
}
func (this *IntDivideOperator) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalIntDivide(&aout, &bout)
}

// ----------------------------------------------------------------
type PowOperator struct{ a, b IEvaluable }

func NewPowOperator(a, b IEvaluable) *PowOperator {
	return &PowOperator{a: a, b: b}
}
func (this *PowOperator) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalPow(&aout, &bout)
}

// ----------------------------------------------------------------
type DotPlusOperator struct{ a, b IEvaluable }

func NewDotPlusOperator(a, b IEvaluable) *DotPlusOperator {
	return &DotPlusOperator{a: a, b: b}
}
func (this *DotPlusOperator) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalDotPlus(&aout, &bout)
}

// ----------------------------------------------------------------
type DotMinusOperator struct{ a, b IEvaluable }

func NewDotMinusOperator(a, b IEvaluable) *DotMinusOperator {
	return &DotMinusOperator{a: a, b: b}
}
func (this *DotMinusOperator) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalDotMinus(&aout, &bout)
}

// ----------------------------------------------------------------
type DotTimesOperator struct{ a, b IEvaluable }

func NewDotTimesOperator(a, b IEvaluable) *DotTimesOperator {
	return &DotTimesOperator{a: a, b: b}
}
func (this *DotTimesOperator) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalDotTimes(&aout, &bout)
}

// ----------------------------------------------------------------
type DotDivideOperator struct{ a, b IEvaluable }

func NewDotDivideOperator(a, b IEvaluable) *DotDivideOperator {
	return &DotDivideOperator{a: a, b: b}
}
func (this *DotDivideOperator) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalDotDivide(&aout, &bout)
}

// ----------------------------------------------------------------
type ModulusOperator struct{ a, b IEvaluable }

func NewModulusOperator(a, b IEvaluable) *ModulusOperator {
	return &ModulusOperator{a: a, b: b}
}
func (this *ModulusOperator) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalModulus(&aout, &bout)
}

// ----------------------------------------------------------------
type BitwiseANDOperator struct{ a, b IEvaluable }

func NewBitwiseANDOperator(a, b IEvaluable) *BitwiseANDOperator {
	return &BitwiseANDOperator{a: a, b: b}
}
func (this *BitwiseANDOperator) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalBitwiseAND(&aout, &bout)
}

// ----------------------------------------------------------------
type BitwiseOROperator struct{ a, b IEvaluable }

func NewBitwiseOROperator(a, b IEvaluable) *BitwiseOROperator {
	return &BitwiseOROperator{a: a, b: b}
}
func (this *BitwiseOROperator) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalBitwiseOR(&aout, &bout)
}

// ----------------------------------------------------------------
type BitwiseXOROperator struct{ a, b IEvaluable }

func NewBitwiseXOROperator(a, b IEvaluable) *BitwiseXOROperator {
	return &BitwiseXOROperator{a: a, b: b}
}
func (this *BitwiseXOROperator) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalBitwiseXOR(&aout, &bout)
}

// ----------------------------------------------------------------
type BitwiseNOTOperator struct{ a IEvaluable }

func NewBitwiseNOTOperator(a IEvaluable) *BitwiseNOTOperator {
	return &BitwiseNOTOperator{a: a}
}
func (this *BitwiseNOTOperator) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	return lib.MlrvalBitwiseNOT(&aout)
}

// ----------------------------------------------------------------
type LogicalANDOperator struct{ a, b IEvaluable }

func NewLogicalANDOperator(a, b IEvaluable) *LogicalANDOperator {
	return &LogicalANDOperator{a: a, b: b}
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

func (this *LogicalANDOperator) Evaluate(state *State) lib.Mlrval {
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
type LogicalOROperator struct{ a, b IEvaluable }

func NewLogicalOROperator(a, b IEvaluable) *LogicalOROperator {
	return &LogicalOROperator{a: a, b: b}
}

// This is different from most of the evaluator functions in that it does
// short-circuiting: since is logical OR, the second argument is not evaluated
// if the first argumeent is false.
//
// See the disposition-matrix discussion for LogicalANDOperator.
func (this *LogicalOROperator) Evaluate(state *State) lib.Mlrval {
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
type LogicalXOROperator struct{ a, b IEvaluable }

func NewLogicalXOROperator(a, b IEvaluable) *LogicalXOROperator {
	return &LogicalXOROperator{a: a, b: b}
}
func (this *LogicalXOROperator) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalLogicalXOR(&aout, &bout)
}

// ----------------------------------------------------------------
type LogicalNOTOperator struct{ a IEvaluable }

func NewLogicalNOTOperator(a IEvaluable) *LogicalNOTOperator {
	return &LogicalNOTOperator{a: a}
}
func (this *LogicalNOTOperator) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	return lib.MlrvalLogicalNOT(&aout)
}

// ----------------------------------------------------------------
type EqualsOperator struct{ a, b IEvaluable }

func NewEqualsOperator(a, b IEvaluable) *EqualsOperator {
	return &EqualsOperator{a: a, b: b}
}
func (this *EqualsOperator) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalEquals(&aout, &bout)
}

// ----------------------------------------------------------------
type NotEqualsOperator struct{ a, b IEvaluable }

func NewNotEqualsOperator(a, b IEvaluable) *NotEqualsOperator {
	return &NotEqualsOperator{a: a, b: b}
}
func (this *NotEqualsOperator) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalNotEquals(&aout, &bout)
}

// ----------------------------------------------------------------
type GreaterThanOperator struct{ a, b IEvaluable }

func NewGreaterThanOperator(a, b IEvaluable) *GreaterThanOperator {
	return &GreaterThanOperator{a: a, b: b}
}
func (this *GreaterThanOperator) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalGreaterThan(&aout, &bout)
}

// ----------------------------------------------------------------
type GreaterThanOrEqualsOperator struct{ a, b IEvaluable }

func NewGreaterThanOrEqualsOperator(a, b IEvaluable) *GreaterThanOrEqualsOperator {
	return &GreaterThanOrEqualsOperator{a: a, b: b}
}
func (this *GreaterThanOrEqualsOperator) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalGreaterThanOrEquals(&aout, &bout)
}

// ----------------------------------------------------------------
type LessThanOperator struct{ a, b IEvaluable }

func NewLessThanOperator(a, b IEvaluable) *LessThanOperator {
	return &LessThanOperator{a: a, b: b}
}
func (this *LessThanOperator) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalLessThan(&aout, &bout)
}

// ----------------------------------------------------------------
type LessThanOrEqualsOperator struct{ a, b IEvaluable }

func NewLessThanOrEqualsOperator(a, b IEvaluable) *LessThanOrEqualsOperator {
	return &LessThanOrEqualsOperator{a: a, b: b}
}
func (this *LessThanOrEqualsOperator) Evaluate(state *State) lib.Mlrval {
	aout := this.a.Evaluate(state)
	bout := this.b.Evaluate(state)
	return lib.MlrvalLessThanOrEquals(&aout, &bout)
}

// ================================================================
type TernaryOperator struct{ a, b, c IEvaluable }

func NewTernaryOperator(a, b, c IEvaluable) *TernaryOperator {
	return &TernaryOperator{a: a, b: b, c: c}
}
func (this *TernaryOperator) Evaluate(state *State) lib.Mlrval {
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
