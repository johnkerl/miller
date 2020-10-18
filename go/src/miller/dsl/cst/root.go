package cst

import (
	"container/list"
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
		beginBlocks:                 make([]*StatementBlockNode, 0),
		mainBlock:                   NewStatementBlockNode(),
		endBlocks:                   make([]*StatementBlockNode, 0),
		udfManager:                  NewUDFManager(),
		unresolvedFunctionCallsites: list.New(),
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

	err = cstRoot.resolveFunctionCallsites()
	if err != nil {
		return nil, err
	}

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

	// Example AST:
	//
	// $ mlr put -v 'begin{@a=1;@b=2} $x=3; $y=4' myfile.dkvp
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

		if astChild.Type == dsl.NodeTypeFunctionDefinition {
			err := this.BuildAndInstallUDF(astChild)
			if err != nil {
				return err
			}

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

// This is invoked within the buildMainPass call tree whenever a function is
// called before it's defined.
func (this *RootNode) rememberUnresolvedFunctionCallsite(udfCallsite *UDFCallsite) {
	this.unresolvedFunctionCallsites.PushBack(udfCallsite)
}

// After-pass after buildMainPass returns, in case a function was called before
// it was defined. It may be the case that:
//
// * A user-defined function was called before it was defined, and was actually defined.
// * A user-defined function was called before it was defined, and it was not actually defined.
// * The user misspelled the name of a built-in function.
//
// So, our error message should reflect all those options.

func (this *RootNode) resolveFunctionCallsites() error {
	for this.unresolvedFunctionCallsites.Len() > 0 {
		unresolvedFunctionCallsite := this.unresolvedFunctionCallsites.Remove(
			this.unresolvedFunctionCallsites.Front(),
		).(*UDFCallsite)

		functionName := unresolvedFunctionCallsite.udf.signature.functionName
		callsiteArity := unresolvedFunctionCallsite.udf.signature.arity

		udf, err := this.udfManager.LookUp(functionName, callsiteArity)
		if err != nil {
			return err
		}
		if udf == nil {
			return errors.New(
				"Miller: function name not found: " + functionName,
			)
		}

		unresolvedFunctionCallsite.udf = udf
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
