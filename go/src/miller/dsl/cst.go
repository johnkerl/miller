// maybe subpackage to avoid names all starting with "CST" ?
package dsl

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
type CSTState struct {
	Inrec *containers.Lrec
	Context *containers.Context
}
func NewCSTState(
	inrec *containers.Lrec,
	context *containers.Context,
) *CSTState {
	return &CSTState {
		Inrec: inrec,
		Context: context,
	}
}

// ----------------------------------------------------------------
type CSTRoot struct {
	// array of statements/blocks
}

type ICSTExecutable interface {
	Execute(cstState *CSTState)
}

// "Parent" class (this is Go, so composition, but ...) with an Execute method
// taking state: inrec, context, oosvars, ...
//
// maybe an IExecutable interface?
type CSTDirectSrecFieldAssignment struct {
	lhsFieldName string
	rhs          ICSTEvaluable
}
// xxx implement ICSTExecutable

type CSTIndirectSrecFieldAssignment struct {
	lhsFieldName ICSTEvaluable
	rhs          ICSTEvaluable
}
// xxx implement ICSTExecutable

type CSTStatementBlock struct {
	// list of statement
}
// xxx implement ICSTExecutable

// ================================================================
type ICSTEvaluable interface {
	Evaluate(cstState *CSTState) lib.Mlrval
	// Needs an Evaluate which takes context and produces a mlrval
}

// ----------------------------------------------------------------
type CSTStringLiteral struct {
	literal lib.Mlrval
}
func NewCSTStringLiteral(literal string) *CSTStringLiteral {
	return &CSTStringLiteral{
		literal: lib.MlrvalFromString(literal),
	}
}
func (this *CSTStringLiteral) Evaluate(cstState *CSTState) lib.Mlrval {
	return this.literal
}

// ----------------------------------------------------------------
type CSTIntLiteral struct {
	literal lib.Mlrval
}
func NewCSTIntLiteral(literal int64) *CSTIntLiteral {
	return &CSTIntLiteral{
		literal: lib.MlrvalFromInt64(literal),
	}
}
func (this *CSTIntLiteral) Evaluate(cstState *CSTState) lib.Mlrval {
	return this.literal
}

// ----------------------------------------------------------------
type CSTFloatLiteral struct {
	literal lib.Mlrval
}
func NewCSTFloatLiteral(literal float64) *CSTFloatLiteral {
	return &CSTFloatLiteral{
		literal: lib.MlrvalFromFloat64(literal),
	}
}
func (this *CSTFloatLiteral) Evaluate(cstState *CSTState) lib.Mlrval {
	return this.literal
}

// ----------------------------------------------------------------
type CSTBoolLiteral struct {
	literal lib.Mlrval
}
func NewCSTBoolLiteral(literal bool) *CSTBoolLiteral {
	return &CSTBoolLiteral{
		literal: lib.MlrvalFromBool(literal),
	}
}
func (this *CSTBoolLiteral) Evaluate(cstState *CSTState) lib.Mlrval {
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

// xxx package miller/dst/cst
// xxx package miller/dst/cst/literals
// xxx package miller/dst/cst/operators
// etc -- ???

// ----------------------------------------------------------------
type CSTDotOperator struct {
	a ICSTEvaluable
	b ICSTEvaluable
}
func NewCSTDotOperator(a, b ICSTEvaluable) *CSTDotOperator {
	return &CSTDotOperator{
		a:a, b:b,
	}
}
func (this *CSTDotOperator) Evaluate(cstState *CSTState) lib.Mlrval {
	aout := this.a.Evaluate(cstState)
	bout := this.b.Evaluate(cstState)
	return lib.MlrvalDot(&aout, &bout)
}

