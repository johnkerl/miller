package dsl

import (
	"fmt"
	"miller/parsing/token"
)

// ----------------------------------------------------------------
// xxx comment interface{} everywhere vs. true types due to gocc polymorphism API.
// and, line-count for casts here vs in the BNF:
//
// Statement :
//   md_token_field_name md_token_assign md_token_number
//
// Statement :
//   md_token_field_name md_token_assign md_token_number
//     << dsl.NewASTNodeTernary("foo", $0, $1, $2) >> ;

// ----------------------------------------------------------------
type AST struct {
	Root *ASTNode
}

func NewAST(root interface{}) (*AST, error) {
	return &AST {
		root.(*ASTNode),
	}, nil
}

func (this *AST) Print() {
	this.Root.Print(0)
}

// ----------------------------------------------------------------
type ASTNode struct {
	Token token.Token
	// Type enum
	// There need to be tokenless nodes
	// text string
	Children []*ASTNode
}

func (this *ASTNode) Print(depth int) {
	for i := 0; i < depth; i++ {
		fmt.Print("  ")
	}
	tok := this.Token
	fmt.Printf("* Type=\"%s\" Literal=\"%s\"\n",
		token.TokMap.Id(tok.Type), string(tok.Lit))
	if this.Children != nil {
		for _, child := range this.Children {
			child.Print(depth + 1)
		}
	}
}

func NewASTNode(itok interface{}) (*ASTNode, error) {
	tok := itok.(*token.Token)
	return &ASTNode {
		*tok,
		// type
		nil,
	}, nil
}

// xxx temp
func Wrap(iparent interface{}) (*ASTNode, error) {
	parent := iparent.(*ASTNode)
	return parent, nil
}

func MakeZary(iparent interface{}) (*ASTNode, error) {
	parent := iparent.(*ASTNode)
	children := make([]*ASTNode, 0)
	parent.Children = children
	return parent, nil
}

func MakeUnary(iparent interface{}, childA interface{}) (*ASTNode, error) {
	parent := iparent.(*ASTNode)
	children := make([]*ASTNode, 1)
	children[0] = childA.(*ASTNode)
	parent.Children = children
	return parent, nil
}

func MakeBinary(iparent interface{}, childA interface{}, childB interface{}) (*ASTNode, error) {
	parent := iparent.(*ASTNode)
	children := make([]*ASTNode, 2)
	children[0] = childA.(*ASTNode)
	children[1] = childB.(*ASTNode)
	parent.Children = children
	return parent, nil
}
