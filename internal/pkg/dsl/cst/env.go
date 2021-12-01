// ================================================================
// This handles ENV["FOO"] on the right-hand side of an assignment.  Note that
// environment variables aren't arbitrarily indexable like maps are -- they're
// only a single-level map from string to string, managed indirectly through
// library routines.
// ================================================================

package cst

import (
	"os"

	"github.com/johnkerl/miller/internal/pkg/dsl"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
	"github.com/johnkerl/miller/internal/pkg/runtime"
	"github.com/johnkerl/miller/internal/pkg/types"
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
		return types.MLRVAL_ABSENT
	}
	if !name.IsString() {
		return types.MLRVAL_ERROR
	}

	return mlrval.MlrvalFromString(os.Getenv(name.String()))
}
