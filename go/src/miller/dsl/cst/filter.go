// ================================================================
// This handles filter statements.
// ================================================================

package cst

import (
	"errors"

	"miller/dsl"
	"miller/lib"
	"miller/runtime"
)

// ----------------------------------------------------------------
type FilterStatementNode struct {
	filterEvaluable IEvaluable
}

// ----------------------------------------------------------------
func (this *RootNode) BuildFilterStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(
		astNode.Type != dsl.NodeTypeFilterStatement &&
			astNode.Type != dsl.NodeTypeBareBoolean)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)

	filterEvaluable, err := this.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	return &FilterStatementNode{
		filterEvaluable: filterEvaluable,
	}, nil
}

func (this *FilterStatementNode) Execute(state *runtime.State) (*BlockExitPayload, error) {

	filterResult := this.filterEvaluable.Evaluate(state)

	if filterResult.IsAbsent() {
		state.LastFilterResultDefined = false
		return nil, nil
	}

	boolValue, isBool := filterResult.GetBoolValue()
	if !isBool {
		return nil, errors.New(
			"Miller: expression does not evaluate to boolean: got " +
				filterResult.GetTypeName() + ".",
		)
	}

	state.FilterResult = boolValue
	state.LastFilterResultDefined = true

	return nil, nil
}
