// ================================================================
// This is for begin and end blocks, but not the main block which is direct
// from the CST root.
// ================================================================

package cst

import (
	"mlr/internal/pkg/dsl"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/runtime"
)

// ----------------------------------------------------------------
func NewStatementBlockNode() *StatementBlockNode {
	return &StatementBlockNode{
		executables: make([]IExecutable, 0),
	}
}

// ----------------------------------------------------------------
func (node *StatementBlockNode) AppendStatementNode(executable IExecutable) {
	node.executables = append(node.executables, executable)
}

// ----------------------------------------------------------------
func (root *RootNode) BuildStatementBlockNodeFromBeginOrEnd(
	astBeginOrEndNode *dsl.ASTNode,
) (*StatementBlockNode, error) {

	lib.InternalCodingErrorIf(
		astBeginOrEndNode.Type != dsl.NodeTypeBeginBlock &&
			astBeginOrEndNode.Type != dsl.NodeTypeEndBlock,
	)
	lib.InternalCodingErrorIf(astBeginOrEndNode.Children == nil)
	// TODO: change the BNF to make it always 1 in the AST
	lib.InternalCodingErrorIf(len(astBeginOrEndNode.Children) > 1)

	if len(astBeginOrEndNode.Children) == 0 {
		return NewStatementBlockNode(), nil
	}

	// Example AST:
	//
	// $ mlr put -v 'begin{@a=1;@b=2} $x=3; $y=4' s
	// DSL EXPRESSION:
	// begin{@a=1;@b=2} $x=3; $y=4
	// AST:
	// * StatementBlock
	//     * BeginBlock
	//         * StatementBlock
	//             * Assignment "="
	//                 * DirectOosvarValue "a"
	//                 * IntLiteral "1"
	//             * Assignment "="
	//                 * DirectOosvarValue "b"
	//                 * IntLiteral "2"
	//     * Assignment "="
	//         * DirectFieldValue "x"
	//         * IntLiteral "3"
	//     * Assignment "="
	//         * DirectFieldValue "y"
	//         * IntLiteral "4"

	astStatementBlockNode := astBeginOrEndNode.Children[0]
	lib.InternalCodingErrorIf(astStatementBlockNode.Type != dsl.NodeTypeStatementBlock)
	statementBlockNode, err := root.BuildStatementBlockNode(astStatementBlockNode)
	if err != nil {
		return nil, err
	} else {
		return statementBlockNode, nil
	}
}

// ----------------------------------------------------------------
func (root *RootNode) BuildStatementBlockNode(
	astNode *dsl.ASTNode,
) (*StatementBlockNode, error) {

	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeStatementBlock)

	statementBlockNode := NewStatementBlockNode()

	astChildren := astNode.Children

	for _, astChild := range astChildren {
		statement, err := root.BuildStatementNode(astChild)
		if err != nil {
			return nil, err
		}
		statementBlockNode.AppendStatementNode(statement)
	}
	return statementBlockNode, nil
}

// ----------------------------------------------------------------
func (node *StatementBlockNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	state.Stack.PushStackFrame()
	defer state.Stack.PopStackFrame()
	for _, statement := range node.executables {
		blockExitPayload, err := statement.Execute(state)
		if err != nil {
			return nil, err
		}
		if blockExitPayload != nil {
			return blockExitPayload, nil
		}
	}

	return nil, nil
}

// ----------------------------------------------------------------
// Assumes the caller has wrapped PushStackFrame() / PopStackFrame().  That
// could be done here, but is instead done in the caller to simplify the
// binding of for-loop variables. In particular, in
//
//   'for (i = 0; i < 10; i += 1) {...}'
//
// the 'i = 0' and 'i += 1' are StatementBlocks and if they pushed their
// own stack frame then the 'i=0' would be in an evanescent, isolated frame.

func (node *StatementBlockNode) ExecuteFrameless(state *runtime.State) (*BlockExitPayload, error) {
	for _, statement := range node.executables {
		blockExitPayload, err := statement.Execute(state)
		if err != nil {
			return nil, err
		}
		if blockExitPayload != nil {
			return blockExitPayload, nil
		}
	}

	return nil, nil
}
