package cst

import (
	"miller/dsl"
	"miller/lib"
)

// This is for begin and end blocks, but not the main block which is direct
// from the CST root.

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
func BuildStatementBlockNodeFromBeginOrEnd(
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
	statementBlockNode, err := BuildStatementBlockNode(astStatementBlockNode)
	if err != nil {
		return nil, err
	} else {
		return statementBlockNode, nil
	}
}

// ----------------------------------------------------------------
func BuildStatementBlockNode(astNode *dsl.ASTNode) (*StatementBlockNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeStatementBlock)

	statementBlockNode := NewStatementBlockNode()

	astChildren := astNode.Children

	for _, astChild := range astChildren {
		statement, err := BuildStatementNode(astChild)
		if err != nil {
			return nil, err
		}
		statementBlockNode.AppendStatementNode(statement)
	}
	return statementBlockNode, nil
}

// ----------------------------------------------------------------
func (this *StatementBlockNode) Execute(state *State) error {
	state.stack.PushStackFrame()
	defer state.stack.PopStackFrame()
	for _, statement := range this.executables {
		err := statement.Execute(state)
		if err != nil {
			return err
		}
	}
	return nil
}

// Assumes the caller has wrapped PushStackFrame() / PopStackFrame().  That
// could be done here, but is instead done in the caller to simplify the
// binding of for-loop variables. In particular, in
//   'for (i = 0; i < 10; i += 1) {...}'
// the 'i = 0' and 'i += 1' are StatementBlocks and if they pushed their
// own stack frame then the 'i=0' would be in an evanescent, isolated frame.

func (this *StatementBlockNode) ExecuteFrameless(state *State) error {
	for _, statement := range this.executables {
		err := statement.Execute(state)
		if err != nil {
			return err
		}
	}
	return nil
}
