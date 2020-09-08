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
func BuildRoot() *Root {
	return &Root{
		executables: make([]IExecutable, 0),
	}
}

// ----------------------------------------------------------------
func (this *Root) AppendStatement(executable IExecutable) {
	this.executables = append(this.executables, executable)
}

// ----------------------------------------------------------------
func Build(ast *dsl.AST) (*Root, error) {
	if ast.Root == nil {
		return nil, errors.New("Cannot build CST from nil AST root")
	}

	cstRoot := BuildRoot()

	// They can do mlr put '': there are simply zero statements.
	if ast.Root.Type == dsl.NodeTypeEmptyStatement {
		return cstRoot, nil
	}

	if ast.Root.Type != dsl.NodeTypeStatementBlock {
		return nil, errors.New(
			"CST root build: non-statement-block AST root node unhandled",
		)
	}
	astChildren := ast.Root.Children

	for _, astChild := range astChildren {
		statement, err := BuildStatementNode(astChild)
		if err != nil {
			return nil, err
		}
		cstRoot.AppendStatement(statement)
	}
	return cstRoot, nil
}

// ----------------------------------------------------------------
func (this *Root) Execute(state *State) (outrec *lib.Mlrmap, err error) {

	for _, statement := range this.executables {
		err := statement.Execute(state)
		if err != nil {
			return nil, err
		}
	}

	return state.Inrec, nil
}
