// ================================================================
// This handles filter statements.
// ================================================================

package cst

import (
	"miller/src/dsl"
	"miller/src/lib"
	"miller/src/runtime"
)

// ----------------------------------------------------------------
type FilterStatementNode struct {
	filterEvaluable IEvaluable
}

// ----------------------------------------------------------------
func (root *RootNode) BuildFilterStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(
		astNode.Type != dsl.NodeTypeFilterStatement &&
			astNode.Type != dsl.NodeTypeBareBoolean)
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
