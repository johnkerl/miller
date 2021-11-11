// ================================================================
// This handles bare booleans and filter statements.
//
// Example of the former: 'NR > 10' or '$x =~ "(..)_(...)"' without if or '{...}' body.
//
// Example of the latter: mlr put 'filter NR > 10'.
//
// For mlr put, these are no-ops except side-effects (like regex-captures).
// (Unless the user uses the filter keyword, like mlr put 'filter NR < 10'.)
//
// For mlr filter, they set the filter condition.
//
// For mlr repl, they become the expression to be printed after evaluation --
// e.g. if the user types '1+2' then the repl prints '3' and that expression is
// stored as a bare-boolean evaluable. Which is a misnomer (sorry!) since 3 is
// not a boolean.
// ================================================================

package cst

import (
	"mlr/internal/pkg/dsl"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/runtime"
)

// ----------------------------------------------------------------
// BareBooleanStatementNode is for implicit filter statements such as mlr
// filter 'NR < 10' -- "implicit" since the word "filter" doesn't appear within
// the single quotes as part of the DSL expression per se. Or (as noted above)
// mlr put 'NR < 10' is a no-op.
type BareBooleanStatementNode struct {
	bareBooleanEvaluable IEvaluable
	writeToState         bool
}

// ----------------------------------------------------------------
func (root *RootNode) BuildBareBooleanStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeBareBoolean)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)

	bareBooleanEvaluable, err := root.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	return &BareBooleanStatementNode{
		bareBooleanEvaluable: bareBooleanEvaluable,
		writeToState:         root.dslInstanceType != DSLInstanceTypePut,
	}, nil
}

func (node *BareBooleanStatementNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	// Evaluate always -- along with any side effects such as regex captures
	// from things like '$x =~ "(..)_(...)"' -- but for mlr put do not store
	// the result as state.FilterExpression. (For mlr put, bare booleans are
	// no-ops, outside of any side effects they may have.)
	result := node.bareBooleanEvaluable.Evaluate(state)
	if node.writeToState {
		state.FilterExpression = result
	}
	return nil, nil
}

// ----------------------------------------------------------------
// FilterStatementNode is for explicit filter statements such as mlr put
// 'filter NR < 10', where the word "filter" appears within the single quotes
// and is part of the DSL expression per se.

type FilterStatementNode struct {
	filterEvaluable IEvaluable
}

func (root *RootNode) BuildFilterStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeFilterStatement)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)

	filterEvaluable, err := root.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	return &FilterStatementNode{
		filterEvaluable: filterEvaluable,
	}, nil
}

func (node *FilterStatementNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	state.FilterExpression = node.filterEvaluable.Evaluate(state).Copy()
	return nil, nil
}
