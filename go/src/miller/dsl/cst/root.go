package cst

import (
	"errors"
	"miller/dsl"
	"miller/types"
)

// ================================================================
// Top-level entry point for building a CST from an AST at parse time, and for
// executing the CST at runtime.
// ================================================================

// ----------------------------------------------------------------
func NewEmptyRoot() *RootNode {
	return &RootNode{
		beginBlocks: make([]*StatementBlockNode, 0),
		mainBlock:   NewStatementBlockNode(),
		endBlocks:   make([]*StatementBlockNode, 0),
		udfManager:  NewUDFManager(),
	}
}

// ----------------------------------------------------------------
func Build(ast *dsl.AST) (*RootNode, error) {
	if ast.RootNode == nil {
		return nil, errors.New("Cannot build CST from nil AST root")
	}

	cstRoot := NewEmptyRoot()

	err := cstRoot.buildMainPass(ast)

	if err != nil {
		return nil, err
	}

	// TODO: UDF-resolver after-pass

	return cstRoot, nil
}

// ----------------------------------------------------------------
// This builds the CST almost entirely. The only afterwork is that user-defined
// functions may be called before they are defined, so a follow-up pass will
// need to resolve those callsites.

func (this *RootNode) buildMainPass(ast *dsl.AST) error {

	// They can do mlr put '': there are simply zero statements.
	if ast.RootNode.Type == dsl.NodeTypeEmptyStatement {
		return nil
	}

	if ast.RootNode.Type != dsl.NodeTypeStatementBlock {
		return errors.New(
			"CST root build: non-statement-block AST root node unhandled",
		)
	}
	astChildren := ast.RootNode.Children

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Pass 2 to handle everyting else besides functions definitions, ignoring
	// them when they are encountered.

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

		// TODO: fill out
		if astChild.Type == dsl.NodeTypeFunctionDefinition {
			// TODO
			//fmt.Printf("UDF stub: found function %s\n", string(astChild.Token.Lit))
			// UDFManager.Install(...)

		} else if astChild.Type == dsl.NodeTypeBeginBlock || astChild.Type == dsl.NodeTypeEndBlock {
			statementBlockNode, err := this.BuildStatementBlockNodeFromBeginOrEnd(astChild)
			if err != nil {
				return err
			}

			if astChild.Type == dsl.NodeTypeBeginBlock {
				this.beginBlocks = append(this.beginBlocks, statementBlockNode)
			} else {
				this.endBlocks = append(this.endBlocks, statementBlockNode)
			}
		} else {
			statementNode, err := this.BuildStatementNode(astChild)
			if err != nil {
				return err
			}
			this.mainBlock.AppendStatementNode(statementNode)
		}
	}

	return nil
}

// ----------------------------------------------------------------
func (this *RootNode) ExecuteBeginBlocks(state *State) error {
	for _, beginBlock := range this.beginBlocks {
		_, err := beginBlock.Execute(state)
		if err != nil {
			return err
		}
	}
	return nil
}

// ----------------------------------------------------------------
func (this *RootNode) ExecuteMainBlock(state *State) (outrec *types.Mlrmap, err error) {
	_, err = this.mainBlock.Execute(state)
	return state.Inrec, err
}

// ----------------------------------------------------------------
func (this *RootNode) ExecuteEndBlocks(state *State) error {
	for _, endBlock := range this.endBlocks {
		_, err := endBlock.Execute(state)
		if err != nil {
			return err
		}
	}
	return nil
}
