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
	if astNode.Type != dsl.NodeTypeOperator {
		return nil, errors.New("Internal coding error detected") // xxx libify
	}

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
	if arity != 1 {
		return nil, errors.New("Internal coding error detected") // xxx libify
	}
	//astChild := astNode.Children[0]

	return nil, errors.New("CST build: AST unary operator node unhandled.")
}

// ----------------------------------------------------------------
func NewBinaryOperatorNode(astNode *dsl.ASTNode) (IEvaluable, error) {
	arity := len(astNode.Children)
	if arity != 2 {
		return nil, errors.New("Internal coding error detected") // xxx libify
	}

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

		// xxx continue ...
	}

	return nil, errors.New(
		"CST build: unandled AST binary operator node \"" + sop + "\"",
	)
}

// ----------------------------------------------------------------
func NewTernaryOperatorNode(astNode *dsl.ASTNode) (IEvaluable, error) {
	arity := len(astNode.Children)
	if arity != 3 {
		return nil, errors.New("Internal coding error detected") // xxx libify
	}

	//leftASTChild := astNode.Children[0]
	//middleASTChild := astNode.Children[1]
	//rightASTChild := astNode.Children[2]

	return nil, errors.New("CST build: AST ternary operator node unhandled.")
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
