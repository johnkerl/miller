// ================================================================
// This is for things that get us out of statement blocks: break, continue,
// return.
// ================================================================

package cst

import (
	"errors"

	"miller/dsl"
	"miller/lib"
	"miller/runtime"
)

// ----------------------------------------------------------------
type BreakNode struct {
}

func (this *RootNode) BuildBreakNode(astNode *dsl.ASTNode) (*BreakNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeBreak)
	lib.InternalCodingErrorIf(astNode.Children == nil)
	lib.InternalCodingErrorIf(len(astNode.Children) != 0)

	return &BreakNode{}, nil
}

func (this *BreakNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	return &BlockExitPayload{
		BLOCK_EXIT_BREAK,
		nil,
	}, nil
}

// ----------------------------------------------------------------
type ContinueNode struct {
}

func (this *RootNode) BuildContinueNode(astNode *dsl.ASTNode) (*ContinueNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeContinue)
	lib.InternalCodingErrorIf(astNode.Children == nil)
	lib.InternalCodingErrorIf(len(astNode.Children) != 0)

	return &ContinueNode{}, nil
}

func (this *ContinueNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	return &BlockExitPayload{
		BLOCK_EXIT_CONTINUE,
		nil,
	}, nil
}

// ----------------------------------------------------------------
type ReturnNode struct {
	returnValueExpression IEvaluable
}

func (this *RootNode) BuildReturnNode(astNode *dsl.ASTNode) (*ReturnNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeReturn)
	lib.InternalCodingErrorIf(astNode.Children == nil)
	if len(astNode.Children) == 0 {
		return &ReturnNode{nil}, nil
	} else if len(astNode.Children) == 1 {
		returnValueExpression, err := this.BuildEvaluableNode(astNode.Children[0])
		if err != nil {
			return nil, err
		}
		return &ReturnNode{returnValueExpression}, nil
	} else {
		lib.InternalCodingErrorIf(true)
	}
	return nil, errors.New("Internal coding error: Statement should not be reached.")
}

func (this *ReturnNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	if this.returnValueExpression == nil {
		return &BlockExitPayload{
			BLOCK_EXIT_RETURN_VOID,
			nil,
		}, nil
	} else {
		// This can be of type MT_ERROR but we do not use Go-level error return here
		returnValue := this.returnValueExpression.Evaluate(state) // TODO: TypeGatedMlrval
		return &BlockExitPayload{
			BLOCK_EXIT_RETURN_VALUE,
			&returnValue,
		}, nil
	}
}
