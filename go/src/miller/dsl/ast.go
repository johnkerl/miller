package dsl

import (
	"errors"
	"fmt"
	"strings"

	"miller/parsing/token"
)

type TNodeType string

const (
	NodeTypeEmptyStatement TNodeType = "EmptyStatement"

	NodeTypeStringLiteral             = "StringLiteral"
	NodeTypeIntLiteral                = "IntLiteral"
	NodeTypeFloatLiteral              = "FloatLiteral"
	NodeTypeBoolLiteral               = "BoolLiteral"
	NodeTypeArrayLiteral              = "ArrayLiteral"
	NodeTypeMapLiteral                = "MapLiteral"
	NodeTypeMapLiteralKeyValuePair    = "MapLiteralKeyValuePair"
	NodeTypeArrayOrMapIndexAccess     = "ArrayOrMapIndexAccess"
	NodeTypeArraySliceAccess          = "ArraySliceAccess"
	NodeTypeArraySliceEmptyLowerIndex = "ArraySliceEmptyLowerIndex"
	NodeTypeArraySliceEmptyUpperIndex = "ArraySliceEmptyUpperIndex"

	NodeTypeDirectFieldValue   = "DirectFieldValue"
	NodeTypeIndirectFieldValue = "IndirectFieldValue"
	NodeTypeFullSrec           = "FullSrec"
	NodeTypeIndexedLvalue      = "IndexedLvalue"

	NodeTypeStatementBlock  = "StatementBlock"
	NodeTypeAssignment      = "Assignment"
	NodeTypeOperator        = "Operator"
	NodeTypeContextVariable = "ContextVariable"

	// A special token which causes a panic when evaluated.  This is for
	// testing that AND/OR short-circuiting is implemented correctly: output =
	// input1 || panic should NOT panic the process when input1 is true.
	NodeTypePanic = "Panic"
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

// This is for the GOCC/BNF parser, which produces an AST
func NewAST(root interface{}) (*AST, error) {
	return &AST{
		Root: root.(*ASTNode),
	}, nil
}

func (this *AST) Print() {
	this.Root.Print(0)
}

// ----------------------------------------------------------------
type ASTNode struct {
	Token    *token.Token // Nil for tokenless/structural nodes
	Type     TNodeType
	Children []*ASTNode

	// xxx sketch:
	// * no longer have separate AST/CST as in the C version ?
	// * have a nullable evaluator function pointer attached to each node
	// * outrec := node.Evaluate(inrec, state) ?
	// * what about evaluator-state ?
	// * outrec := node.Evaluator.Evaluate(inrec, state) ?
	// * state:
	//   o string_array_t** ppregex_captures;
	//   o mlhmmv_root_t*   poosvars;
	//   o context_t*       pctx;
	//   o local_stack_t*   plocal_stack;
	//   o loop_stack_t*    ploop_stack;
	//   o return_state_t   return_state;
	//   o int              trace_execution; -- move out of state?
	//   o int              json_quote_int_keys; -- move out of state?
	//   o int              json_quote_non_string_values; -- move out of state?
	// * statements:
	// * node.Executor.Execute(inrec, state) ?
}

func (this *ASTNode) Print(depth int) {
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
			child.Print(depth + 1)
		}
	}
}

func NewASTNode(itok interface{}, nodeType TNodeType) (*ASTNode, error) {
	return NewASTNodeNestable(itok, nodeType), nil
}

// For handling empty expressions.
func NewASTNodeEmpty(nodeType TNodeType) (*ASTNode, error) {
	return &ASTNode{
		Token:    nil,
		Type:     nodeType,
		Children: nil,
	}, nil
}

// Strips the leading '$' from field names. Not done in the parser itself due
// to LR-1 conflicts.
func NewASTNodeStripDollarPlease(itok interface{}, nodeType TNodeType) (*ASTNode, error) {
	oldToken := itok.(*token.Token)
	newToken := &token.Token{
		Type: oldToken.Type,
		Lit:  oldToken.Lit[1:],
		Pos:  oldToken.Pos,
	}
	return NewASTNodeNestable(newToken, nodeType), nil
}

// Likewise for the leading/trailing double quotes on string literals.  Also,
// since string literals can have backslash-escaped double-quotes like
// "...\"...\"...", we also unbackslash here.
func NewASTNodeStripDoubleQuotePairPlease(
	itok interface{},
	nodeType TNodeType,
) (*ASTNode, error) {
	oldToken := itok.(*token.Token)
	n := len(oldToken.Lit)
	contents := string(oldToken.Lit[1 : n-1])
	contents = strings.ReplaceAll(contents, "\\\"", "\"")
	newToken := &token.Token{
		Type: oldToken.Type,
		Lit:  []byte(contents),
		Pos:  oldToken.Pos,
	}
	return NewASTNodeNestable(newToken, nodeType), nil
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

func NewASTNodeZary(itok interface{}, nodeType TNodeType) (*ASTNode, error) {
	parent := NewASTNodeNestable(itok, nodeType)
	convertToZary(parent)
	return parent, nil
}

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
func NewASTNodeBinary(itok, childA, childB interface{}, nodeType TNodeType) (*ASTNode, error) {
	return NewASTNodeBinaryNestable(itok, childA, childB, nodeType), nil
}

func NewASTNodeTernary(itok, childA, childB, childC interface{}, nodeType TNodeType) (*ASTNode, error) {
	parent := NewASTNodeNestable(itok, nodeType)
	convertToTernary(parent, childA, childB, childC)
	return parent, nil
}

// Pass-through expressions in the grammar sometimes need to be turned from
// (ASTNode) to (ASTNode, error)
func Nestable(iparent interface{}) (*ASTNode, error) {
	return iparent.(*ASTNode), nil
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

func convertToTernary(iparent interface{}, childA, childB, childC interface{}) {
	parent := iparent.(*ASTNode)
	children := make([]*ASTNode, 3)
	children[0] = childA.(*ASTNode)
	children[1] = childB.(*ASTNode)
	children[2] = childC.(*ASTNode)
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

func AppendChild(iparent interface{}, child interface{}) (*ASTNode, error) {
	parent := iparent.(*ASTNode)
	if parent.Children == nil {
		convertToUnary(iparent, child)
	} else {
		parent.Children = append(parent.Children, child.(*ASTNode))
	}
	return parent, nil
}

func AdoptChildren(iparent interface{}, ichild interface{}) (*ASTNode, error) {
	parent := iparent.(*ASTNode)
	child := ichild.(*ASTNode)
	parent.Children = child.Children
	child.Children = nil
	return parent, nil
}

func (this *ASTNode) CheckArity(
	arity int,
) error {
	if len(this.Children) != arity {
		return errors.New(
			fmt.Sprintf(
				"AST node arity %d, expected %d",
				len(this.Children), arity,
			),
		)
	} else {
		return nil
	}
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
