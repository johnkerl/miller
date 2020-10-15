package cst

import (
	"errors"

	"miller/dsl"
	"miller/lib"
)

// ================================================================
// This handles filter statements.
// ================================================================

// ----------------------------------------------------------------
type FilterStatementNode struct {
	filterEvaluable IEvaluable
}

// ----------------------------------------------------------------
// TODO: disallow bare boolean except for final statement in 'mlr filter' ...
func BuildFilterStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(
		astNode.Type != dsl.NodeTypeFilterStatement &&
			astNode.Type != dsl.NodeTypeBareBoolean)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)

	filterEvaluable, err := BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	return &FilterStatementNode{
		filterEvaluable: filterEvaluable,
	}, nil
}

func (this *FilterStatementNode) Execute(state *State) (*BlockExitPayload, error) {

	filterResult := this.filterEvaluable.Evaluate(state)

	if filterResult.IsAbsent() {
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

	return nil, nil
}
