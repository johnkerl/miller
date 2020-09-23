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

// mlr -n put -v for (k in $*) { emit { k : k } }
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
	lib.InternalCodingErrorIf(len(astNode.Children) != 3)

	// TODO

	// error if loop-over node isn't Mappable (inasmuch as can be detected at CST-build time)

	return NewForLoopKeyOnlyNode(), nil
}

// ----------------------------------------------------------------
func (this *ForLoopKeyOnlyNode) Execute(state *State) error {
	// TODO

	// absent/error handling
	// error if loop-over node isn't a mlrmap
	// loop over map keys
	// * make new stack frame, binding the localvar name to its .Copy()
	// * push the stack frame
	// * statement-block execute
	// * pop the stack frame

	return nil
}
