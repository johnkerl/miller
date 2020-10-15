package cst

import (
	"errors"

	"miller/dsl"
	"miller/lib"
)

// ================================================================
// This is for while/do-while loops
// ================================================================

// ================================================================
type WhileLoopNode struct {
	conditionNode      IEvaluable
	statementBlockNode *StatementBlockNode
}

func NewWhileLoopNode(
	conditionNode IEvaluable,
	statementBlockNode *StatementBlockNode,
) *WhileLoopNode {
	return &WhileLoopNode{
		conditionNode:      conditionNode,
		statementBlockNode: statementBlockNode,
	}
}

func BuildWhileLoopNode(astNode *dsl.ASTNode) (*WhileLoopNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeWhileLoop)
	lib.InternalCodingErrorIf(len(astNode.Children) != 2)

	conditionNode, err := BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	statementBlockNode, err := BuildStatementBlockNode(astNode.Children[1])
	if err != nil {
		return nil, err
	}

	return NewWhileLoopNode(
		conditionNode,
		statementBlockNode,
	), nil
}

// ----------------------------------------------------------------
func (this *WhileLoopNode) Execute(state *State) (*BlockExitStatus, error) {
	var blockExitStatus *BlockExitStatus = nil
	var err error = nil
	for {
		condition := this.conditionNode.Evaluate(state)
		boolValue, isBool := condition.GetBoolValue()
		if !isBool {
			// TODO: line-number/token info for the DSL expression.
			return nil, errors.New("Miller: conditional expression did not evaluate to boolean.")
		}
		if boolValue != true {
			break
		}
		blockExitStatus, err = this.statementBlockNode.Execute(state)
		if err != nil {
			return nil, err
		}
	}
	return blockExitStatus, nil
}

// ================================================================
type DoWhileLoopNode struct {
	statementBlockNode *StatementBlockNode
	conditionNode      IEvaluable
}

func NewDoWhileLoopNode(
	statementBlockNode *StatementBlockNode,
	conditionNode IEvaluable,
) *DoWhileLoopNode {
	return &DoWhileLoopNode{
		statementBlockNode: statementBlockNode,
		conditionNode:      conditionNode,
	}
}

func BuildDoWhileLoopNode(astNode *dsl.ASTNode) (*DoWhileLoopNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeDoWhileLoop)
	lib.InternalCodingErrorIf(len(astNode.Children) != 2)

	statementBlockNode, err := BuildStatementBlockNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	conditionNode, err := BuildEvaluableNode(astNode.Children[1])
	if err != nil {
		return nil, err
	}

	return NewDoWhileLoopNode(
		statementBlockNode,
		conditionNode,
	), nil
}

// ----------------------------------------------------------------
func (this *DoWhileLoopNode) Execute(state *State) (*BlockExitStatus, error) {
	var blockExitStatus *BlockExitStatus = nil
	var err error = nil
	for {
		blockExitStatus, err = this.statementBlockNode.Execute(state)
		if err != nil {
			return nil, err
		}

		condition := this.conditionNode.Evaluate(state)
		boolValue, isBool := condition.GetBoolValue()
		if !isBool {
			// TODO: line-number/token info for the DSL expression.
			return nil, errors.New("Miller: conditional expression did not evaluate to boolean.")
		}
		if boolValue == false {
			break
		}
	}
	return blockExitStatus, nil
}
