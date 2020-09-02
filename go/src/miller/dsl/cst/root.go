package cst

import (
	"errors"

	"miller/containers"
	"miller/dsl"
)

// ----------------------------------------------------------------
func NewRoot() *Root {
	return &Root{
		make([]IExecutable, 0),
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
	if ast.Root.Type != dsl.NodeTypeStatementBlock {
		return nil, errors.New("Non-statement-block AST root node unhandled")
	}
	astChildren := ast.Root.Children

	cstRoot := NewRoot()

	// For this very early stub, only process statement nodes (which is all the
	// grammar produces at this point ...)
	for _, astChild := range astChildren {
		if astChild.Type != dsl.NodeTypeSrecDirectAssignment {
			return nil, errors.New("Non-assignment AST node unhandled")
		}

		statement, err := NewSrecDirectFieldAssignment(astChild)
		if err != nil {
			return nil, err
		}

		cstRoot.AppendStatement(statement)
	}

	return cstRoot, nil
}

// ----------------------------------------------------------------
func (this *Root) Execute(state *State) *containers.Lrec {

	for _, statement := range(this.executables) {
		statement.Execute(state)
	}

	return state.Inrec
}
