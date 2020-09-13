package cst

import (
	"errors"

	"miller/dsl"
	"miller/lib"
)

// ================================================================
// Top-level entry point for building a CST from an AST at parse time, and for
// executing the CST at runtime.
// ================================================================

// ----------------------------------------------------------------
func BuildEmptyRoot() *RootNode {
	return &RootNode{
		beginBlocks: make([]*StatementBlockNode, 0),
		mainBlock:   NewStatementBlockNode(),
		endBlocks:   make([]*StatementBlockNode, 0),
	}
}

// ----------------------------------------------------------------
func Build(ast *dsl.AST) (*RootNode, error) {
	if ast.RootNode == nil {
		return nil, errors.New("Cannot build CST from nil AST root")
	}

	cstRoot := BuildEmptyRoot()

	// They can do mlr put '': there are simply zero statements.
	if ast.RootNode.Type == dsl.NodeTypeEmptyStatement {
		return cstRoot, nil
	}

	if ast.RootNode.Type != dsl.NodeTypeStatementBlock {
		return nil, errors.New(
			"CST root build: non-statement-block AST root node unhandled",
		)
	}
	astChildren := ast.RootNode.Children

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

	for _, astChild := range astChildren {
		if astChild.Type == dsl.NodeTypeBeginBlock || astChild.Type == dsl.NodeTypeBeginBlock {
			statementBlockNode, err := BuildStatementBlockNodeFromBeginOrEnd(astChild)
			if err != nil {
				return nil, err
			}

			if astChild.Type == dsl.NodeTypeBeginBlock {
				cstRoot.beginBlocks = append(cstRoot.beginBlocks, statementBlockNode)
			} else {
				cstRoot.endBlocks = append(cstRoot.endBlocks, statementBlockNode)
			}
		} else {
			statementNode, err := BuildStatementNode(astChild)
			if err != nil {
				return nil, err
			}
			cstRoot.mainBlock.AppendStatementNode(statementNode)
		}
	}
	return cstRoot, nil
}

// ----------------------------------------------------------------
func (this *RootNode) ExecuteBeginBlocks(state *State) error {
	for _, beginBlock := range this.beginBlocks {
		err := beginBlock.Execute(state)
		if err != nil {
			return err
		}
	}
	return nil
}

// ----------------------------------------------------------------
func (this *RootNode) ExecuteMainBlock(state *State) (outrec *types.Mlrmap, err error) {
	err = this.mainBlock.Execute(state)

	return state.Inrec, err
}

// ----------------------------------------------------------------
func (this *RootNode) ExecuteEndBlocks(state *State) error {
	for _, endBlock := range this.endBlocks {
		err := endBlock.Execute(state)
		if err != nil {
			return err
		}
	}
	return nil
}
