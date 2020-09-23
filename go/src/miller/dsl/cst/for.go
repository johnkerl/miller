package cst

import (
	//"errors"

	"miller/dsl"
	"miller/lib"
	//"miller/types"
)

// ================================================================
// This is for various flavors of for-loop
// ================================================================

// ----------------------------------------------------------------
type ForLoopKeyOnlyNode struct {
}

func NewForLoopKeyOnlyNode(
) *ForLoopKeyOnlyNode {
	return &ForLoopKeyOnlyNode{
	}
}

// ----------------------------------------------------------------
// Sample AST:

// mlr -n put -v for (k in $*) { emit { k : k} }
// DSL EXPRESSION:
// for (k in $*) { emit { k : k} }
// RAW AST:
// * StatementBlock
//     * ForLoopKeyOnly "for"
//         * LocalVariable "k"
//         * FullSrec "$*"
//         * StatementBlock
//             * EmitStatement "emit"
//                 * MapLiteral "{}"
//                     * MapLiteralKeyValuePair ":"
//                         * LocalVariable "k"
//                         * LocalVariable "k"

func BuildForLoopKeyOnlyNode(astNode *dsl.ASTNode) (*ForLoopKeyOnlyNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeForLoopKeyOnly)

	return NewForLoopKeyOnlyNode(), nil
}

// ----------------------------------------------------------------
func (this *ForLoopKeyOnlyNode) Execute(state *State) error {
	return nil
}
