// This handles ENV["FOO"] on the right-hand side of an assignment.  Note that
// environment variables aren't arbitrarily indexable like maps are -- they're
// only a single-level map from string to string, managed indirectly through
// library routines.

package cst

import (
	"os"

	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/runtime"
	"github.com/johnkerl/pgpg/go/lib/pkg/asts"
)

type EnvironmentVariableNode struct {
	nameEvaluable IEvaluable
}

func (root *RootNode) BuildEnvironmentVariableNode(astNode *asts.ASTNode) (*EnvironmentVariableNode, error) {
	lib.InternalCodingErrorWithMessageIf(
		astNode.Type != asts.NodeType(NodeTypeEnvironmentVariable) && string(astNode.Type) != NodeTypeEnvironmentVariable,
		"expected EnvironmentVariable node",
	)
	lib.InternalCodingErrorIf(astNode.Children == nil)
	lib.InternalCodingErrorIf(len(astNode.Children) != 1)
	child := astNode.Children[0]
	// Unwrap Parenthesized if present: ENV[("FOO")] or similar.
	if string(child.Type) == NodeTypeParenthesized && child.Children != nil && len(child.Children) == 1 {
		child = child.Children[0]
	}
	// ENV.FOO: child is non_sigil_name (identifier FOO) -> use as env var name directly
	if string(child.Type) == "non_sigil_name" || string(child.Type) == NodeTypeLocalVariable {
		sval := tokenLit(child)
		if sval != "" {
			return &EnvironmentVariableNode{
				nameEvaluable: root.BuildStringLiteralNode(sval),
			}, nil
		}
	}
	// PGPG may put the string literal's text on the node or its child; extract it directly
	// so we get the correct env var name even when BuildEvaluableNode would yield empty.
	if string(child.Type) == NodeTypeStringLiteral || string(child.Type) == "string_literal" {
		sval := tokenLit(child)
		if sval == "" && child.Children != nil && len(child.Children) == 1 {
			sval = tokenLit(child.Children[0])
		}
		if sval != "" {
			// PGPG lexer may include surrounding quotes in the lexeme; strip them.
			if len(sval) >= 2 && sval[0] == '"' && sval[len(sval)-1] == '"' {
				sval = sval[1 : len(sval)-1]
			}
			return &EnvironmentVariableNode{
				nameEvaluable: root.BuildStringLiteralNode(lib.UnbackslashStringLiteral(sval)),
			}, nil
		}
	}
	nameEvaluable, err := root.BuildEvaluableNode(child)
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
