// ================================================================
// This handles ENV["FOO"] on the right-hand side of an assignment.  Note that
// environment variables aren't arbitrarily indexable like maps are -- they're
// only a single-level map from string to string, managed indirectly through
// library routines.
// ================================================================

package cst

import (
	"os"

	"miller/src/dsl"
	"miller/src/lib"
	"miller/src/runtime"
	"miller/src/types"
)

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

func (this *EnvironmentVariableNode) Evaluate(
	output *types.Mlrval,
	state *runtime.State,
) {
	var name types.Mlrval
	this.nameEvaluable.Evaluate(&name, state)
	if name.IsAbsent() {
		output.SetFromAbsent()
		return
	}
	if !name.IsString() {
		output.SetFromError()
		return
	}

	output.SetFromString(os.Getenv(name.String()))
}
