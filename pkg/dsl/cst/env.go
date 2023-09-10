// ================================================================
// This handles ENV["FOO"] on the right-hand side of an assignment.  Note that
// environment variables aren't arbitrarily indexable like maps are -- they're
// only a single-level map from string to string, managed indirectly through
// library routines.
// ================================================================

package cst

import (
	"os"

	"github.com/johnkerl/miller/pkg/dsl"
	"github.com/johnkerl/miller/pkg/lib"
	"github.com/johnkerl/miller/pkg/mlrval"
	"github.com/johnkerl/miller/pkg/runtime"
)

type EnvironmentVariableNode struct {
	nameEvaluable IEvaluable
}

func (root *RootNode) BuildEnvironmentVariableNode(astNode *dsl.ASTNode) (*EnvironmentVariableNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeEnvironmentVariable)
	lib.InternalCodingErrorIf(astNode.Children == nil)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)
	nameEvaluable, err := root.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	return &EnvironmentVariableNode{
		nameEvaluable: nameEvaluable,
	}, nil
}

func (node *EnvironmentVariableNode) Evaluate(
	state *runtime.State,
) *mlrval.Mlrval {
	name := node.nameEvaluable.Evaluate(state)
	if name.IsAbsent() {
		return mlrval.ABSENT.StrictModeCheck(state.StrictMode, "ENV[(absent)]")
	}
	if !name.IsString() {
		return mlrval.FromTypeErrorUnary("ENV[]", name)
	}

	return mlrval.FromString(os.Getenv(name.String()))
}
