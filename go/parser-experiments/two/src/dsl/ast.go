package dsl

import (
	"fmt"

	"experimental/token"
)

type TNodeType string

const (
	NodeTypeEmptyStatement TNodeType = "empty statement"

	NodeTypeStatementBlock = "statement block"
	NodeTypeStatement      = "statement"
)

// ----------------------------------------------------------------
type AST struct {
	RootNode *ASTNode
}

// This is for the GOCC/BNF parser, which produces an AST
func NewAST(iroot interface{}) (*AST, error) {
	return &AST{
		RootNode: iroot.(*ASTNode),
	}, nil
}

func (this *AST) Print() {
	this.RootNode.Print()
}

// ----------------------------------------------------------------
type ASTNode struct {
	Token    *token.Token // Nil for tokenless/structural nodes
	Type     TNodeType
	Children []*ASTNode
}

func (this *ASTNode) Print() {
	this.PrintAux(0)
}
func (this *ASTNode) PrintAux(depth int) {
	for i := 0; i < depth; i++ {
		fmt.Print("    ")
	}
	tok := this.Token
	fmt.Print("* " + this.Type)

	if tok != nil {
		//fmt.Printf(" \"%s\" \"%s\"",
		//	token.TokMap.Id(tok.Type), string(tok.Lit))
		fmt.Printf(" \"%s\"", string(tok.Lit))
	}
	fmt.Println()
	if this.Children != nil {
		for _, child := range this.Children {
			child.PrintAux(depth + 1)
		}
	}
}

// For handling empty expressions.
func NewASTNodeEmptyNestable(nodeType TNodeType) *ASTNode {
	return &ASTNode{
		Token:    nil,
		Type:     nodeType,
		Children: nil,
	}
}

// xxx comment why grammar use
func NewASTNodeNestable(itok interface{}, nodeType TNodeType) *ASTNode {
	var tok *token.Token = nil
	if itok != nil {
		tok = itok.(*token.Token)
	}
	return &ASTNode{
		Token:    tok,
		Type:     nodeType,
		Children: nil,
	}
}

// Signature: Token Type
func NewASTNodeZary(itok interface{}, nodeType TNodeType) (*ASTNode, error) {
	parent := NewASTNodeNestable(itok, nodeType)
	convertToZary(parent)
	return parent, nil
}

// Signature: Token Node Type
func NewASTNodeUnary(itok, childA interface{}, nodeType TNodeType) (*ASTNode, error) {
	parent := NewASTNodeNestable(itok, nodeType)
	convertToUnary(parent, childA)
	return parent, nil
}

// Signature: Token Node Node Type
func NewASTNodeBinaryNestable(itok, childA, childB interface{}, nodeType TNodeType) *ASTNode {
	parent := NewASTNodeNestable(itok, nodeType)
	convertToBinary(parent, childA, childB)
	return parent
}

// Signature: Token Node Node Type
func NewASTNodeBinary(
	itok, childA, childB interface{}, nodeType TNodeType,
) (*ASTNode, error) {
	return NewASTNodeBinaryNestable(itok, childA, childB, nodeType), nil
}

func convertToZary(iparent interface{}) {
	parent := iparent.(*ASTNode)
	children := make([]*ASTNode, 0)
	parent.Children = children
}

// xxx inline this. can be a one-liner.
func convertToUnary(iparent interface{}, childA interface{}) {
	parent := iparent.(*ASTNode)
	children := make([]*ASTNode, 1)
	children[0] = childA.(*ASTNode)
	parent.Children = children
}

func convertToBinary(iparent interface{}, childA, childB interface{}) {
	parent := iparent.(*ASTNode)
	children := make([]*ASTNode, 2)
	children[0] = childA.(*ASTNode)
	children[1] = childB.(*ASTNode)
	parent.Children = children
}

func PrependChild(iparent interface{}, ichild interface{}) (*ASTNode, error) {
	parent := iparent.(*ASTNode)
	child := ichild.(*ASTNode)
	if parent.Children == nil {
		convertToUnary(iparent, ichild)
	} else {
		parent.Children = append([]*ASTNode{child}, parent.Children...)
	}
	return parent, nil
}

// TODO: comment
func Wrap(inode interface{}) (*ASTNode, error) {
	node := inode.(*ASTNode)
	return node, nil
}
