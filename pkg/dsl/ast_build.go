// AST-build methods, for use by callbacks within the GOCC/BNF Miller
// DSL grammar in mlr.bnf.

package dsl

import (
	"fmt"

	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/parsing/token"
)

// This is for the GOCC/BNF parser, which produces an AST
func NewASTWithErrorReturn(iroot interface{}) (*AST, error) {
	return &AST{
		RootNode: iroot.(*ASTNode),
	}, nil
}

func NewASTNodeTerminal(itok interface{}, nodeType TNodeType) *ASTNode {
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

// NewASTNode is the ASTNode constructor. If children is non-nil and length 0, a
// zary node is created. (Example: a function call with zero arguments.) If
// children is nil, a terminal node is created. (Example: a string or integer
// literal.)
func NewASTNode(
	itok interface{},
	nodeType TNodeType,
	children []interface{},
) *ASTNode {

	var tok *token.Token = nil
	if itok != nil {
		tok = itok.(*token.Token)
	}

	node := &ASTNode{
		Token:    tok,
		Type:     nodeType,
		Children: nil,
	}

	if children == nil {
		return node
	}

	n := len(children)
	node.Children = make([]*ASTNode, n)
	for i, child := range children {
		node.Children[i] = child.(*ASTNode)
	}
	return node
}

// Pass-through expressions in the grammar sometimes need to be turned from
// (ASTNode) to (ASTNode, error). This is for GOCC.
func WithErrorReturn(iparent interface{}) (*ASTNode, error) {
	return iparent.(*ASTNode), nil
}

func NewASTNodeWithErrorReturn(
	itok interface{},
	nodeType TNodeType,
	children []interface{},
) (*ASTNode, error) {
	return WithErrorReturn(NewASTNode(itok, nodeType, children))
}

func NewASTNodeTerminalWithErrorReturn(itok interface{}, nodeType TNodeType) (*ASTNode, error) {
	return WithErrorReturn(NewASTNode(itok, nodeType, nil))
}

// Strips the leading '$' from field names, or '@' from oosvar names. Not done
// in the parser itself due to LR-1 conflicts.
func NewASTNodeStripDollarOrAtSign(itok interface{}, nodeType TNodeType) (*ASTNode, error) {
	oldToken := itok.(*token.Token)
	lib.InternalCodingErrorIf(len(oldToken.Lit) < 2)
	lib.InternalCodingErrorIf(oldToken.Lit[0] != '$' && oldToken.Lit[0] != '@')
	newToken := &token.Token{
		Type: oldToken.Type,
		Lit:  oldToken.Lit[1:],
		Pos:  oldToken.Pos,
	}
	return NewASTNodeTerminal(newToken, nodeType), nil
}

// Strips the leading '${' and trailing '}' from braced field names, or '@{'
// and '}' from oosvar names. Not done in the parser itself due to LR-1
// conflicts.
func NewASTNodeStripDollarOrAtSignAndCurlyBraces(
	itok interface{},
	nodeType TNodeType,
) (*ASTNode, error) {
	oldToken := itok.(*token.Token)
	n := len(oldToken.Lit)
	lib.InternalCodingErrorIf(n < 4)
	lib.InternalCodingErrorIf(oldToken.Lit[0] != '$' && oldToken.Lit[0] != '@')
	lib.InternalCodingErrorIf(oldToken.Lit[1] != '{')
	lib.InternalCodingErrorIf(oldToken.Lit[n-1] != '}')
	newToken := &token.Token{
		Type: oldToken.Type,
		Lit:  oldToken.Lit[2 : n-1],
		Pos:  oldToken.Pos,
	}
	return NewASTNodeTerminal(newToken, nodeType), nil
}

// Likewise for the leading/trailing double quotes on string literals.  Also,
// since string literals can have backslash-escaped double-quotes like
// "...\"...\"...", we also unbackslash here.
func NewASTNodeStripDoubleQuotePair(
	itok interface{},
	nodeType TNodeType,
) (*ASTNode, error) {
	oldToken := itok.(*token.Token)
	n := len(oldToken.Lit)
	contents := string(oldToken.Lit[1 : n-1])

	newToken := &token.Token{
		Type: oldToken.Type,
		Lit:  []byte(contents),
		Pos:  oldToken.Pos,
	}
	return NewASTNodeTerminal(newToken, nodeType), nil
}

func WithChildPrepended(iparent interface{}, ichild interface{}) (*ASTNode, error) {
	parent := iparent.(*ASTNode)
	child := ichild.(*ASTNode)
	if parent.Children == nil {
		parent.Children = []*ASTNode{child}
	} else {
		parent.Children = append([]*ASTNode{child}, parent.Children...)
	}
	return parent, nil
}

func WithTwoChildrenPreprended(iparent interface{}, ichildA, ichildB interface{}) (*ASTNode, error) {
	parent := iparent.(*ASTNode)
	childA := ichildA.(*ASTNode)
	childB := ichildB.(*ASTNode)
	if parent.Children == nil {
		parent.Children = []*ASTNode{childA, childB}
	} else {
		parent.Children = append([]*ASTNode{childA, childB}, parent.Children...)
	}
	return parent, nil
}

func WithChildAppended(iparent interface{}, child interface{}) (*ASTNode, error) {
	parent := iparent.(*ASTNode)
	if parent.Children == nil {
		parent.Children = []*ASTNode{child.(*ASTNode)}
	} else {
		parent.Children = append(parent.Children, child.(*ASTNode))
	}
	return parent, nil
}

func WithChildrenAdopted(iparent interface{}, ichild interface{}) (*ASTNode, error) {
	parent := iparent.(*ASTNode)
	child := ichild.(*ASTNode)
	parent.Children = child.Children
	child.Children = nil
	return parent, nil
}

func (node *ASTNode) CheckArity(
	arity int,
) error {
	if len(node.Children) != arity {
		return fmt.Errorf("expected AST node arity %d, got %d", arity, len(node.Children))
	}
	return nil
}

// Tokens are produced by GOCC. However there is an exception: for the ternary
// operator I want the AST to have a "?:" token, which GOCC doesn't produce
// since nothing is actually spelled like that in the DSL.
func NewASTToken(iliteral interface{}, iclonee interface{}) *token.Token {
	literal := iliteral.(string)
	clonee := iclonee.(*token.Token)
	return &token.Token{
		Type: clonee.Type,
		Lit:  []byte(literal),
		Pos:  clonee.Pos,
	}
}
