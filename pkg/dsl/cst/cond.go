// ================================================================
// This is for awkish pattern-action-blocks, like mlr put 'NR > 10 { ... }'.
// Just shorthand for if-statements without elif/else.
// ================================================================

package cst

import (
	"fmt"

	"github.com/johnkerl/miller/pkg/dsl"
	"github.com/johnkerl/miller/pkg/lib"
	"github.com/johnkerl/miller/pkg/mlrval"
	"github.com/johnkerl/miller/pkg/parsing/token"
	"github.com/johnkerl/miller/pkg/runtime"
)

type CondBlockNode struct {
	conditionNode      IEvaluable
	conditionToken     *token.Token
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
	conditionToken := astNode.Children[0].Token
	statementBlockNode, err := root.BuildStatementBlockNode(astNode.Children[1])
	if err != nil {
		return nil, err
	}
	condBlockNode := &CondBlockNode{
		conditionNode:      conditionNode,
		conditionToken:     conditionToken,
		statementBlockNode: statementBlockNode,
	}

	return condBlockNode, nil
}

// ----------------------------------------------------------------
func (node *CondBlockNode) Execute(
	state *runtime.State,
) (*BlockExitPayload, error) {
	condition := mlrval.TRUE
	if node.conditionNode != nil {
		condition = node.conditionNode.Evaluate(state)
	}
	boolValue, isBool := condition.GetBoolValue()

	// Data-heterogeneity case
	if condition.IsAbsent() {
		boolValue = false
	} else if !isBool {
		return nil, fmt.Errorf(
			"mlr: conditional expression did not evaluate to boolean%s.",
			dsl.TokenToLocationInfo(node.conditionToken),
		)
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
