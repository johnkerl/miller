// ================================================================
// This is for awkish pattern-action-blocks, like mlr put 'NR > 10 { ... }'.
// Just shorthand for if-statements without elif/else.
// ================================================================

package cst

import (
	"errors"

	"mlr/internal/pkg/dsl"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/runtime"
	"mlr/internal/pkg/types"
)

type CondBlockNode struct {
	conditionNode      IEvaluable
	statementBlockNode *StatementBlockNode
}

// ----------------------------------------------------------------
// Sample AST:

func (root *RootNode) BuildCondBlockNode(astNode *dsl.ASTNode) (*CondBlockNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeCondBlock)
	lib.InternalCodingErrorIf(len(astNode.Children) != 2)

	conditionNode, err := root.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	statementBlockNode, err := root.BuildStatementBlockNode(astNode.Children[1])
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
func (node *CondBlockNode) Execute(
	state *runtime.State,
) (*BlockExitPayload, error) {
	condition := types.MLRVAL_TRUE
	if node.conditionNode != nil {
		condition = node.conditionNode.Evaluate(state)
	}
	boolValue, isBool := condition.GetBoolValue()

	// Data-heterogeneity case
	if condition.IsAbsent() {
		boolValue = false
	} else if !isBool {
		// TODO: line-number/token info for the DSL expression.
		return nil, errors.New("mlr: conditional expression did not evaluate to boolean.")
	}

	if boolValue == true {
		blockExitPayload, err := node.statementBlockNode.Execute(state)
		if err != nil {
			return nil, err
		}
		if blockExitPayload != nil {
			return blockExitPayload, nil
		}
	}
	return nil, nil
}
