package cst

import (
	"os"

	"miller/dsl"
	"miller/lib"
	"miller/types"
)

// ================================================================
// This handles ENV["FOO"] on the right-hand side of an assignment.  Note that
// environment variables aren't arbitrarily indexable like maps are -- they're
// only a single-level map from string to string, managed indirectly through
// library routines.
// ================================================================

type EnvironmentVariableNode struct {
	nameEvaluable IEvaluable
}

func (this *RootNode) BuildEnvironmentVariableNode(astNode *dsl.ASTNode) (*EnvironmentVariableNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeEnvironmentVariable)
	lib.InternalCodingErrorIf(astNode.Children == nil)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)
	nameEvaluable, err := this.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	return &EnvironmentVariableNode{
		nameEvaluable: nameEvaluable,
	}, nil
}

func (this *EnvironmentVariableNode) Evaluate(state *State) types.Mlrval {
	name := this.nameEvaluable.Evaluate(state)
	if name.IsAbsent() {
		return types.MlrvalFromAbsent()
	}
	if !name.IsString() {
		return types.MlrvalFromError()
	}

	return types.MlrvalFromString(os.Getenv(name.String()))
}
