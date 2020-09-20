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
func BuildFilterStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeFilterStatement)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)

	filterEvaluable, err := BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	return &FilterStatementNode{
		filterEvaluable: filterEvaluable,
	}, nil
}

func (this *FilterStatementNode) Execute(state *State) error {

	filterResult := this.filterEvaluable.Evaluate(state)

	if filterResult.IsAbsent() {
		return nil
	}

	boolValue, isBoolean := filterResult.GetBoolValue()
	if !isBoolean {
		return errors.New(
			"Miller: expression does not evaluate to boolean: got " +
				filterResult.GetTypeName() + ".",
		)
	}

	state.FilterResult = boolValue

	return nil
}
