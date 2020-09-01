package cst

import (
	"miller/containers"
	"miller/lib"
)

// ----------------------------------------------------------------
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
