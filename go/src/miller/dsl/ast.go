package dsl

import (
	"fmt"
	"miller/parsing/token"
)

type TNodeType string

const (
	NodeTypeStatementBlock = "StatementBlock"
	NodeTypeStatement      = "Statement"
	NodeTypeToken          = "Token"
	NodeTypeOperator       = "Operator"
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
	return &AST{
		root.(*ASTNode),
	}, nil
}

func (this *AST) Print() {
	this.Root.Print(0)
}

// ----------------------------------------------------------------
type ASTNode struct {
	Token   *token.Token // Nil for tokenless/structural nodes
	NodeType TNodeType
	Children []*ASTNode
}

func (this *ASTNode) Print(depth int) {
	for i := 0; i < depth; i++ {
		fmt.Print("    ")
	}
	tok := this.Token
	fmt.Print("* " + this.NodeType)

	if tok != nil {
		fmt.Printf(" \"%s\" \"%s\"",
			token.TokMap.Id(tok.Type), string(tok.Lit))
	}
	fmt.Println()
	if this.Children != nil {
		for _, child := range this.Children {
			child.Print(depth + 1)
		}
	}
}

func NewASTNode(itok interface{}, nodeType TNodeType) (*ASTNode, error) {
	var tok *token.Token = nil
	if itok != nil {
		tok = itok.(*token.Token)
	}
	return &ASTNode{
		tok,
		nodeType,
		nil,
	}, nil
}

func NewASTNodeUnary(itok,  childA interface{}, nodeType TNodeType) (*ASTNode, error) {
	parent, err := NewASTNode(itok, nodeType)
	if err != nil {
		return nil, err
	}
	_, err = MakeUnary(parent, childA)
	if err != nil {
		return nil, err
	}
	return parent, nil
}

func NewASTNodeBinary(itok, childA, childB interface{}, nodeType TNodeType) (*ASTNode, error) {
	parent, err := NewASTNode(itok, nodeType)
	if err != nil {
		return nil, err
	}
	_, err = MakeBinary(parent, childA, childB)
	if err != nil {
		return nil, err
	}
	return parent, nil
}

func NewASTNodeTernary(itok, childA, childB, childC interface{}, nodeType TNodeType) (*ASTNode, error) {
	parent, err := NewASTNode(itok, nodeType)
	if err != nil {
		return nil, err
	}
	_, err = MakeTernary(parent, childA, childB, childC)
	if err != nil {
		return nil, err
	}
	return parent, nil
}

// Pass-through expressions in the grammar sometimes need to be turned from
// (ASTNode) to (ASTNode, error)
func PairNil(iparent interface{}) (*ASTNode, error) {
	return iparent.(*ASTNode), nil
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

func MakeBinary(iparent interface{}, childA, childB interface{}) (*ASTNode, error) {
	parent := iparent.(*ASTNode)
	children := make([]*ASTNode, 2)
	children[0] = childA.(*ASTNode)
	children[1] = childB.(*ASTNode)
	parent.Children = children
	return parent, nil
}

func MakeTernary(iparent interface{}, childA, childB, childC interface{}) (*ASTNode, error) {
	parent := iparent.(*ASTNode)
	children := make([]*ASTNode, 3)
	children[0] = childA.(*ASTNode)
	children[1] = childB.(*ASTNode)
	children[2] = childC.(*ASTNode)
	parent.Children = children
	return parent, nil
}

func AppendChild(iparent interface{}, child interface{}) (*ASTNode, error) {
	parent := iparent.(*ASTNode)
	if parent.Children == nil {
		MakeUnary(iparent, child)
	} else {
		parent.Children = append(parent.Children, child.(*ASTNode))
	}
	return parent, nil
}
