package dsl

import (
	"fmt"

	"experimental/token"
)

type TNodeType string

const (
	NodeTypeEmptyStatement TNodeType = "empty statement"

	NodeTypeStatementBlock = "statement block"
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

// func NewASTNode(itok interface{}, nodeType TNodeType) (*ASTNode, error) {
// 	return NewASTNodeNestable(itok, nodeType), nil
// }
// 
// // For handling empty expressions.
// func NewASTNodeEmptyNestable(nodeType TNodeType) *ASTNode {
// 	return &ASTNode{
// 		Token:    nil,
// 		Type:     nodeType,
// 		Children: nil,
// 	}
// }
// 
// // For handling empty expressions.
// func NewASTNodeEmpty(nodeType TNodeType) (*ASTNode, error) {
// 	return NewASTNodeEmptyNestable(nodeType), nil
// }
// 
// // Strips the leading '$' from field names, or '@' from oosvar names. Not done
// // in the parser itself due to LR-1 conflicts.
// func NewASTNodeStripDollarOrAtSign(itok interface{}, nodeType TNodeType) (*ASTNode, error) {
// 	oldToken := itok.(*token.Token)
// 	lib.InternalCodingErrorIf(len(oldToken.Lit) < 2)
// 	lib.InternalCodingErrorIf(oldToken.Lit[0] != '$' && oldToken.Lit[0] != '@')
// 	newToken := &token.Token{
// 		Type: oldToken.Type,
// 		Lit:  oldToken.Lit[1:],
// 		Pos:  oldToken.Pos,
// 	}
// 	return NewASTNodeNestable(newToken, nodeType), nil
// }
// 
// // Strips the leading '${' and trailing '}' from braced field names, or '@{'
// // and '}' from oosvar names. Not done in the parser itself due to LR-1
// // conflicts.
// func NewASTNodeStripDollarOrAtSignAndCurlyBraces(
// 	itok interface{},
// 	nodeType TNodeType,
// ) (*ASTNode, error) {
// 	oldToken := itok.(*token.Token)
// 	n := len(oldToken.Lit)
// 	lib.InternalCodingErrorIf(n < 4)
// 	lib.InternalCodingErrorIf(oldToken.Lit[0] != '$' && oldToken.Lit[0] != '@')
// 	lib.InternalCodingErrorIf(oldToken.Lit[1] != '{')
// 	lib.InternalCodingErrorIf(oldToken.Lit[n-1] != '}')
// 	newToken := &token.Token{
// 		Type: oldToken.Type,
// 		Lit:  oldToken.Lit[2 : n-1],
// 		Pos:  oldToken.Pos,
// 	}
// 	return NewASTNodeNestable(newToken, nodeType), nil
// }
// 
// // Likewise for the leading/trailing double quotes on string literals.  Also,
// // since string literals can have backslash-escaped double-quotes like
// // "...\"...\"...", we also unbackslash here.
// func NewASTNodeStripDoubleQuotePair(
// 	itok interface{},
// 	nodeType TNodeType,
// ) (*ASTNode, error) {
// 	oldToken := itok.(*token.Token)
// 	n := len(oldToken.Lit)
// 	contents := string(oldToken.Lit[1 : n-1])
// 	contents = strings.Replace(contents, "\\\"", "\"", -1)
// 	newToken := &token.Token{
// 		Type: oldToken.Type,
// 		Lit:  []byte(contents),
// 		Pos:  oldToken.Pos,
// 	}
// 	return NewASTNodeNestable(newToken, nodeType), nil
// }
// 
// // xxx comment why grammar use
// func NewASTNodeNestable(itok interface{}, nodeType TNodeType) *ASTNode {
// 	var tok *token.Token = nil
// 	if itok != nil {
// 		tok = itok.(*token.Token)
// 	}
// 	return &ASTNode{
// 		Token:    tok,
// 		Type:     nodeType,
// 		Children: nil,
// 	}
// }
// 
// func NewASTNodeZary(itok interface{}, nodeType TNodeType) (*ASTNode, error) {
// 	parent := NewASTNodeNestable(itok, nodeType)
// 	convertToZary(parent)
// 	return parent, nil
// }
// 
// func NewASTNodeUnary(itok, childA interface{}, nodeType TNodeType) (*ASTNode, error) {
// 	parent := NewASTNodeNestable(itok, nodeType)
// 	convertToUnary(parent, childA)
// 	return parent, nil
// }
// 
// // Signature: Token Node Node Type
// func NewASTNodeBinaryNestable(itok, childA, childB interface{}, nodeType TNodeType) *ASTNode {
// 	parent := NewASTNodeNestable(itok, nodeType)
// 	convertToBinary(parent, childA, childB)
// 	return parent
// }
// 
// // Signature: Token Node Node Type
// func NewASTNodeBinary(
// 	itok, childA, childB interface{}, nodeType TNodeType,
// ) (*ASTNode, error) {
// 	return NewASTNodeBinaryNestable(itok, childA, childB, nodeType), nil
// }
// 
// func NewASTNodeTernary(itok, childA, childB, childC interface{}, nodeType TNodeType) (*ASTNode, error) {
// 	parent := NewASTNodeNestable(itok, nodeType)
// 	convertToTernary(parent, childA, childB, childC)
// 	return parent, nil
// }
// 
// func NewASTNodeQuaternary(
// 	itok, childA, childB, childC, childD interface{}, nodeType TNodeType,
// ) (*ASTNode, error) {
// 	parent := NewASTNodeNestable(itok, nodeType)
// 	convertToQuaternary(parent, childA, childB, childC, childD)
// 	return parent, nil
// }
// 
// // Pass-through expressions in the grammar sometimes need to be turned from
// // (ASTNode) to (ASTNode, error)
// func Nestable(iparent interface{}) (*ASTNode, error) {
// 	return iparent.(*ASTNode), nil
// }
// 
// func convertToZary(iparent interface{}) {
// 	parent := iparent.(*ASTNode)
// 	children := make([]*ASTNode, 0)
// 	parent.Children = children
// }
// 
// // xxx inline this. can be a one-liner.
// func convertToUnary(iparent interface{}, childA interface{}) {
// 	parent := iparent.(*ASTNode)
// 	children := make([]*ASTNode, 1)
// 	children[0] = childA.(*ASTNode)
// 	parent.Children = children
// }
// 
// func convertToBinary(iparent interface{}, childA, childB interface{}) {
// 	parent := iparent.(*ASTNode)
// 	children := make([]*ASTNode, 2)
// 	children[0] = childA.(*ASTNode)
// 	children[1] = childB.(*ASTNode)
// 	parent.Children = children
// }
// 
// func convertToTernary(iparent interface{}, childA, childB, childC interface{}) {
// 	parent := iparent.(*ASTNode)
// 	children := make([]*ASTNode, 3)
// 	children[0] = childA.(*ASTNode)
// 	children[1] = childB.(*ASTNode)
// 	children[2] = childC.(*ASTNode)
// 	parent.Children = children
// }
// 
// func convertToQuaternary(iparent interface{}, childA, childB, childC, childD interface{}) {
// 	parent := iparent.(*ASTNode)
// 	children := make([]*ASTNode, 4)
// 	children[0] = childA.(*ASTNode)
// 	children[1] = childB.(*ASTNode)
// 	children[2] = childC.(*ASTNode)
// 	children[3] = childD.(*ASTNode)
// 	parent.Children = children
// }
// 
// func PrependChild(iparent interface{}, ichild interface{}) (*ASTNode, error) {
// 	parent := iparent.(*ASTNode)
// 	child := ichild.(*ASTNode)
// 	if parent.Children == nil {
// 		convertToUnary(iparent, ichild)
// 	} else {
// 		parent.Children = append([]*ASTNode{child}, parent.Children...)
// 	}
// 	return parent, nil
// }
// 
// func AppendChild(iparent interface{}, child interface{}) (*ASTNode, error) {
// 	parent := iparent.(*ASTNode)
// 	if parent.Children == nil {
// 		convertToUnary(iparent, child)
// 	} else {
// 		parent.Children = append(parent.Children, child.(*ASTNode))
// 	}
// 	return parent, nil
// }
// 
// func AdoptChildren(iparent interface{}, ichild interface{}) (*ASTNode, error) {
// 	parent := iparent.(*ASTNode)
// 	child := ichild.(*ASTNode)
// 	parent.Children = child.Children
// 	child.Children = nil
// 	return parent, nil
// }
// 
// // TODO: comment
// func Wrap(inode interface{}) (*ASTNode, error) {
// 	node := inode.(*ASTNode)
// 	return node, nil
// }
// 
// func (this *ASTNode) CheckArity(
// 	arity int,
// ) error {
// 	if len(this.Children) != arity {
// 		return errors.New(
// 			fmt.Sprintf(
// 				"AST node arity %d, expected %d",
// 				len(this.Children), arity,
// 			),
// 		)
// 	} else {
// 		return nil
// 	}
// }
// 
// // Tokens are produced by GOCC. However there is an exception: for the ternary
// // operator I want the AST to have a "?:" token, which GOCC doesn't produce
// // since nothing is actually spelled like that in the DSL.
// func NewASTToken(iliteral interface{}, iclonee interface{}) *token.Token {
// 	literal := iliteral.(string)
// 	clonee := iclonee.(*token.Token)
// 	return &token.Token{
// 		Type: clonee.Type,
// 		Lit:  []byte(literal),
// 		Pos:  clonee.Pos,
// 	}
// }
