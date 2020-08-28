package ast // TODO: move up one level

import (
	"miller/parsing/token"
)

// ----------------------------------------------------------------
// AST type:
// root *Node

// Node type:
// * text string
// * type enum
// * children []Node

// ----------------------------------------------------------------
type AST struct {
	Root *ASTNode
}

func NewAST(root *ASTNode) *AST {
	return &AST {
		root,
	}
}

// ----------------------------------------------------------------
type ASTNode struct {
	Text string
	// Type enum
	Children []*ASTNode
}

func NewASTNode(text string) *ASTNode {
	return &ASTNode {
		text,
		// type
		nil,
	}
}

func NewASTNodeZary(text string) *ASTNode {
	children := make([]*ASTNode, 0)
	return &ASTNode {
		text,
		// type
		children,
	}
}

func NewASTNodeUnary(text string, childA *ASTNode) *ASTNode {
	children := make([]*ASTNode, 1)
	children[0] = childA
	return &ASTNode {
		text,
		// type
		children,
	}
}

func NewASTNodeBinary(text string, childA *ASTNode, childB *ASTNode) *ASTNode {
	children := make([]*ASTNode, 1)
	children[0] = childA
	children[1] = childB
	return &ASTNode {
		text,
		// type
		children,
	}
}

// ----------------------------------------------------------------
// prototype stuff from gocc example
type (
	StatementList []Statement
	Statement     string
)

func NewStatementList(statement interface{}) (StatementList, error) {
	return StatementList{statement.(Statement)}, nil
}

func AppendStatement(statementList, statement interface{}) (StatementList, error) {
	return append(statementList.(StatementList), statement.(Statement)), nil
}

func NewStatement(statementList interface{}) (Statement, error) {
	return Statement(statementList.(*token.Token).Lit), nil
}
