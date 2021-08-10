// ================================================================
// This handles bare booleans and filter statements.
//
// Example: 'NR > 10' without if or '{...}' body.  For put, these are no-ops
// except side-effects (like regex-captures); for filter, they set the filter
// condition only if they're the last statement in the main block.
// ================================================================

package cst

import (
	"mlr/src/dsl"
	"mlr/src/lib"
	"mlr/src/runtime"
)

// ----------------------------------------------------------------
type BareBooleanStatementNode struct {
	bareBooleanEvaluable IEvaluable
}
type FilterStatementNode struct {
	filterEvaluable IEvaluable
}

// ----------------------------------------------------------------
func (root *RootNode) BuildBareBooleanOrFilterStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	if root.dslInstanceType == DSLInstanceTypePut {
		return root.BuildBareBooleanStatementNode(astNode)
	} else {
		// More to do -- only if final
		return root.BuildFilterStatementNode(astNode)
	}
}

// ----------------------------------------------------------------
func (root *RootNode) BuildBareBooleanStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(
		astNode.Type != dsl.NodeTypeBareBoolean &&
			astNode.Type != dsl.NodeTypeFilterStatement)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)

	bareBooleanEvaluable, err := root.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	return &BareBooleanStatementNode{
		bareBooleanEvaluable: bareBooleanEvaluable,
	}, nil
}

func (node *BareBooleanStatementNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	return nil, nil
}

// ----------------------------------------------------------------
func (root *RootNode) BuildFilterStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(
		astNode.Type != dsl.NodeTypeBareBoolean &&
			astNode.Type != dsl.NodeTypeFilterStatement)
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
