// ================================================================
// This is for awkish pattern-action-blocks, like mlr put 'NR > 10 { ... }'.
// Just shorthand for if-statements without elif/else.
// ================================================================

package cst

import (
	"errors"

	"miller/src/dsl"
	"miller/src/lib"
	"miller/src/runtime"
	"miller/src/types"
)

type CondBlockNode struct {
	conditionNode      IEvaluable
	statementBlockNode *StatementBlockNode
}

// ----------------------------------------------------------------
// Sample AST:

func (this *RootNode) BuildCondBlockNode(astNode *dsl.ASTNode) (*CondBlockNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeCondBlock)
	lib.InternalCodingErrorIf(len(astNode.Children) != 2)

	conditionNode, err := this.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	statementBlockNode, err := this.BuildStatementBlockNode(astNode.Children[1])
	if err != nil {
		return nil, err
	}
	condBlockNode := &CondBlockNode{
		conditionNode:      conditionNode,
		statementBlockNode: statementBlockNode,
	}

	return condBlockNode, nil
}

// ----------------------------------------------------------------
func (this *CondBlockNode) Execute(
	state *runtime.State,
) (*BlockExitPayload, error) {
	var condition types.Mlrval
	condition.SetFromTrue()
	if this.conditionNode != nil {
		this.conditionNode.Evaluate(&condition, state)
	}
	boolValue, isBool := condition.GetBoolValue()

	// Data-heterogeneity case
	if condition.IsAbsent() {
		boolValue = false
	} else if !isBool {
		// TODO: line-number/token info for the DSL expression.
		return nil, errors.New("Miller: conditional expression did not evaluate to boolean.")
	}

	if boolValue == true {
		blockExitPayload, err := this.statementBlockNode.Execute(state)
		if err != nil {
			return nil, err
		}
		if blockExitPayload != nil {
			return blockExitPayload, nil
		}
	}
	return nil, nil
}
