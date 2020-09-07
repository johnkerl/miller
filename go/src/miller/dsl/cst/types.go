package cst

import (
	"miller/lib"
)

// ================================================================
// Main type definitions for CST build/execute
// ================================================================

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
// When we do mlr put '...DSL expression here...', this state is what is needed
// to execute the expression. That includes the current record, AWK-like variables
// such as FILENAME and NR, and out-of-stream variables.
type State struct {
	Inrec   *lib.Mlrmap
	Context *lib.Context
	// TODO: oosvars too
	// TODO: stack frames will go into individual statement-block nodes
}

func NewState(
	inrec *lib.Mlrmap,
	context *lib.Context,
) *State {
	return &State{
		Inrec:   inrec,
		Context: context,
	}
}

// ----------------------------------------------------------------
// This is for all statements and statemnt blocks within the CST.
type IExecutable interface {
	Execute(state *State)
}

// ----------------------------------------------------------------
type Root struct {
	// Statements/blocks
	executables []IExecutable
}

// ----------------------------------------------------------------
type SrecDirectFieldAssignmentNode struct {
	lhsFieldName string
	rhs          IEvaluable
}

type IndirectSrecFieldAssignmentNode struct {
	lhsFieldName IEvaluable
	rhs          IEvaluable
}

type StatementBlockNode struct {
	// TODO: list of statement
}

// ================================================================
// This is for any right-hand side (RHS) of an assignment statement.  Also, for
// computed field names on the left-hand side, like '$a . $b' in mlr put '$[$a
// . $b]' = $x + $y'. Also known as an "Rvalue".
type IEvaluable interface {
	Evaluate(state *State) lib.Mlrval
}

// This is for computing map entries at runtime. For example, in mlr put 'mymap
// = {"sum": $x + $y, "diff": $x - $y}; ...', the first pair would have key
// being string-literal "sum" and value being the evaluable expression '$x + $y'.
type EvaluablePair struct {
	Key   IEvaluable
	Value IEvaluable
}

func NewEvaluablePair(key IEvaluable, value IEvaluable) *EvaluablePair {
	return &EvaluablePair{
		Key:   key,
		Value: value,
	}
}

// ================================================================
// This is for any left-hand side (LHS) of an assignment statement.
// Also known as an "Lvalue".
type IAssignable interface {
	Assign(state *State, mlrval lib.Mlrval) error
}
