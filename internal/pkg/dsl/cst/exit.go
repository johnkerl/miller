// ================================================================
// This handles exit statements.
// ================================================================

package cst

import (
	"fmt"

	"github.com/johnkerl/miller/internal/pkg/dsl"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/runtime"
)

// ================================================================
type ExitStatementNode struct {
	exitCodeEvaluable IEvaluable
}

func (root *RootNode) BuildExitStatementNode(
	astNode *dsl.ASTNode,
) (IExecutable, error) {
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)
	exitCodeNode := astNode.Children[0]

	exitCodeEvaluable, err := root.BuildEvaluableNode(exitCodeNode)
	if err != nil {
		return nil, err
	}

	retval := &ExitStatementNode{
		exitCodeEvaluable: exitCodeEvaluable,
	}

	return retval, nil
}

// ----------------------------------------------------------------
func (node *ExitStatementNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	exitCodeMlrval := node.exitCodeEvaluable.Evaluate(state)

	intValue, isInt := exitCodeMlrval.GetIntValue()
	if !isInt {
		return nil, fmt.Errorf("expected exit statement int argument; got %d", exitCodeMlrval.GetTypeName())
	}

	state.ExitInfo.HasExitCode = true
	state.ExitInfo.ExitCode = int(intValue)

	return nil, nil
}

