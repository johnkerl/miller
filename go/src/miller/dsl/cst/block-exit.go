package cst

import (
	"errors"

	"miller/dsl"
	"miller/lib"
)

// ================================================================
// This is for things that get us out of statement blocks: break, continue,
// return.
// ================================================================

// ----------------------------------------------------------------
type BreakNode struct {
}

func BuildBreakNode(astNode *dsl.ASTNode) (*BreakNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeBreak)
	lib.InternalCodingErrorIf(astNode.Children != nil)

	return &BreakNode{}, nil
}

func (this *BreakNode) Execute(state *State) (*BlockExitPayload, error) {
	return &BlockExitPayload {
		BLOCK_EXIT_BREAK,
		nil,
	}, nil
}

// ----------------------------------------------------------------
type ContinueNode struct {
}

func BuildContinueNode(astNode *dsl.ASTNode) (*ContinueNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeContinue)
	lib.InternalCodingErrorIf(astNode.Children != nil)

	return &ContinueNode{}, nil
}

func (this *ContinueNode) Execute(state *State) (*BlockExitPayload, error) {
	return &BlockExitPayload {
		BLOCK_EXIT_CONTINUE,
		nil,
	}, nil
}

// ----------------------------------------------------------------
type ReturnNode struct {
	returnValueExpression IEvaluable
}

func BuildReturnNode(astNode *dsl.ASTNode) (*ReturnNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeReturn)
	if astNode.Children == nil {
		return &ReturnNode{nil}, nil
	} else if len(astNode.Children) == 1 {
		return &ReturnNode{nil}, nil
	} else {
		lib.InternalCodingErrorIf(true)
	}
	return nil, errors.New("Internal coding error: Statement should not be reached.")
}

func (this *ReturnNode) Execute(state *State) (*BlockExitPayload, error) {
	if this.returnValueExpression == nil {
		return &BlockExitPayload {
			BLOCK_EXIT_RETURN_VOID,
			nil,
		}, nil
	} else {
		returnValue := this.returnValueExpression.Evaluate(state)
		return &BlockExitPayload {
			BLOCK_EXIT_RETURN_VALUE,
			&returnValue,
		}, nil
	}
}
