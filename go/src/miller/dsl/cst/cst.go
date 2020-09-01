package cst

import (
	"miller/containers"
	"miller/lib"
)

// ----------------------------------------------------------------
// There are three CST roots: begin-block, body-block, and end-block.
//
// Next-level items are:
// * srec assignments
// * oosvar assignments
// * localvar assignments
// * emit et al.
// * bare-boolean
// * break/continue/return
// * statement block (if-body, for-body, etc)
// ----------------------------------------------------------------

// ----------------------------------------------------------------
// AST nodes (TNodeType) at the moment:
//
// NodeTypeStringLiteral
// NodeTypeIntLiteral
// NodeTypeFloatLiteral
// NodeTypeBoolLiteral
//
// NodeTypeDirectFieldName
// NodeTypeIndirectFieldName
//
// NodeTypeStatementBlock
// NodeTypeAssignment
// NodeTypeOperator
// NodeTypeContextVariable
// ----------------------------------------------------------------

// ----------------------------------------------------------------
type State struct {
	Inrec   *containers.Lrec
	Context *containers.Context
	// oosvars too
	// stack frames will go into individual statement-block nodes
}

func NewState(
	inrec *containers.Lrec,
	context *containers.Context,
) *State {
	return &State{
		Inrec:   inrec,
		Context: context,
	}
}

// ----------------------------------------------------------------
type IExecutable interface {
	Execute(state *State)
}

// ----------------------------------------------------------------
type Root struct {
	// array of statements/blocks
	executables []IExecutable
}

type DirectSrecFieldAssignment struct {
	lhsFieldName string
	rhs          IEvaluable
}
// xxx implement IExecutable

type IndirectSrecFieldAssignment struct {
	lhsFieldName IEvaluable
	rhs          IEvaluable
}
// xxx implement IExecutable

type StatementBlock struct {
	// list of statement
}
// xxx implement IExecutable

// ================================================================
type IEvaluable interface {
	Evaluate(state *State) lib.Mlrval
	// Needs an Evaluate which takes context and produces a mlrval
}

// ----------------------------------------------------------------
type StringLiteral struct {
	literal lib.Mlrval
}

func NewStringLiteral(literal string) *StringLiteral {
	return &StringLiteral{
		literal: lib.MlrvalFromString(literal),
	}
}
func (this *StringLiteral) Evaluate(state *State) lib.Mlrval {
	return this.literal
}

// ----------------------------------------------------------------
type IntLiteral struct {
	literal lib.Mlrval
}

func NewIntLiteral(literal int64) *IntLiteral {
	return &IntLiteral{
		literal: lib.MlrvalFromInt64(literal),
	}
}
func (this *IntLiteral) Evaluate(state *State) lib.Mlrval {
	return this.literal
}

// ----------------------------------------------------------------
type FloatLiteral struct {
	literal lib.Mlrval
}

func NewFloatLiteral(literal float64) *FloatLiteral {
	return &FloatLiteral{
		literal: lib.MlrvalFromFloat64(literal),
	}
}
func (this *FloatLiteral) Evaluate(state *State) lib.Mlrval {
	return this.literal
}

// ----------------------------------------------------------------
type BoolLiteral struct {
	literal lib.Mlrval
}

func NewBoolLiteral(literal bool) *BoolLiteral {
	return &BoolLiteral{
		literal: lib.MlrvalFromBool(literal),
	}
}
func (this *BoolLiteral) Evaluate(state *State) lib.Mlrval {
	return this.literal
}

// ----------------------------------------------------------------
// context variable by name
// srec direct field name
// srec indirect field name
// unary operator
//   "-"
// binary operator
//   . + - * / //

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
