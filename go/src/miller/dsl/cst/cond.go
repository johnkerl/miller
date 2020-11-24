// ================================================================
// This is for awkish pattern-action-blocks, like mlr put 'NR > 10 { ... }'.
// Just shorthand for if-statements without elif/else.
// ================================================================

package cst

import (
	"errors"

	"miller/dsl"
	"miller/lib"
	"miller/types"
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
func (this *CondBlockNode) Execute(state *State) (*BlockExitPayload, error) {
	condition := types.MlrvalFromTrue()
	if this.conditionNode != nil {
		condition = this.conditionNode.Evaluate(state)
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
