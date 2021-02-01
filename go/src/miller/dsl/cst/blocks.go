// ================================================================
// This is for begin and end blocks, but not the main block which is direct
// from the CST root.
// ================================================================

package cst

import (
	"miller/dsl"
	"miller/lib"
	"miller/runtime"
)

// ----------------------------------------------------------------
func NewStatementBlockNode() *StatementBlockNode {
	return &StatementBlockNode{
		executables: make([]IExecutable, 0),
	}
}

// ----------------------------------------------------------------
func (this *StatementBlockNode) AppendStatementNode(executable IExecutable) {
	this.executables = append(this.executables, executable)
}

// ----------------------------------------------------------------
func (this *RootNode) BuildStatementBlockNodeFromBeginOrEnd(
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
	// RAW AST:
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
	statementBlockNode, err := this.BuildStatementBlockNode(astStatementBlockNode)
	if err != nil {
		return nil, err
	} else {
		return statementBlockNode, nil
	}
}

// ----------------------------------------------------------------
func (this *RootNode) BuildStatementBlockNode(
	astNode *dsl.ASTNode,
) (*StatementBlockNode, error) {

	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeStatementBlock)

	statementBlockNode := NewStatementBlockNode()

	astChildren := astNode.Children

	for _, astChild := range astChildren {
		statement, err := this.BuildStatementNode(astChild)
		if err != nil {
			return nil, err
		}
		statementBlockNode.AppendStatementNode(statement)
	}
	return statementBlockNode, nil
}

// ----------------------------------------------------------------
func (this *StatementBlockNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	state.Stack.PushStackFrame()
	defer state.Stack.PopStackFrame()
	for _, statement := range this.executables {
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

func (this *StatementBlockNode) ExecuteFrameless(state *runtime.State) (*BlockExitPayload, error) {
	for _, statement := range this.executables {
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
