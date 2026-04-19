// This is for begin and end blocks, but not the main block which is direct
// from the CST root.

package cst

import (
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/runtime"

	"github.com/johnkerl/pgpg/go/lib/pkg/asts"
)

func NewStatementBlockNode() *StatementBlockNode {
	return &StatementBlockNode{
		executables: []IExecutable{},
	}
}

func (node *StatementBlockNode) AppendStatementNode(executable IExecutable) {
	node.executables = append(node.executables, executable)
}

func (root *RootNode) BuildStatementBlockNodeFromBeginOrEnd(
	astBeginOrEndNode *asts.ASTNode,
) (*StatementBlockNode, error) {

	lib.InternalCodingErrorIf(
		astBeginOrEndNode.Type != asts.NodeType(NodeTypeBeginBlock) &&
			astBeginOrEndNode.Type != asts.NodeType(NodeTypeEndBlock),
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
	// PGPG: BeginBlock/EndBlock have StatementBlockInBraces as child.
	// With "parent":1,"children":[1], StatementBlockInBraces.Children[0] is the StatementBlock.
	// Unwrap so we pass StatementBlock to BuildStatementBlockNode.
	if astStatementBlockNode.Type == asts.NodeType(NodeTypeStatementBlockInBraces) {
		lib.InternalCodingErrorIf(astStatementBlockNode.Children == nil || len(astStatementBlockNode.Children) < 1)
		astStatementBlockNode = astStatementBlockNode.Children[0]
	}
	statementBlockNode, err := root.BuildStatementBlockNode(astStatementBlockNode)
	if err != nil {
		return nil, err
	}
	return statementBlockNode, nil
}

func (root *RootNode) BuildStatementBlockNode(
	astNode *asts.ASTNode,
) (*StatementBlockNode, error) {
	// PGPG: StatementBlockInBraces has "children":[1] so its child is StatementBlock; unwrap.
	if astNode.Type == asts.NodeType(NodeTypeStatementBlockInBraces) &&
		astNode.Children != nil && len(astNode.Children) == 1 &&
		astNode.Children[0].Type == asts.NodeType(NodeTypeStatementBlock) {
		astNode = astNode.Children[0]
	}
	lib.InternalCodingErrorIf(astNode.Type != asts.NodeType(NodeTypeStatementBlock))

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
