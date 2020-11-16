package dsl

import (
	"errors"
	"fmt"
	"strings"

	"miller/lib"
	"miller/parsing/token"
)

type TNodeType string

const (
	NodeTypeEmptyStatement TNodeType = "empty statement"

	NodeTypeStringLiteral             = "string literal"
	NodeTypeIntLiteral                = "int literal"
	NodeTypeFloatLiteral              = "float literal"
	NodeTypeBoolLiteral               = "bool literal"
	NodeTypeArrayLiteral              = "array literal"
	NodeTypeMapLiteral                = "map literal"
	NodeTypeMapLiteralKeyValuePair    = "map-literal key-value pair"
	NodeTypeArrayOrMapIndexAccess     = "array or map index access"
	NodeTypeArraySliceAccess          = "array-slice access"
	NodeTypeArraySliceEmptyLowerIndex = "array-slice empty lower index"
	NodeTypeArraySliceEmptyUpperIndex = "array-slice empty upper index"
	NodeTypeContextVariable           = "context variable"
	NodeTypeConstant                  = "mathematical constant"

	NodeTypeDirectFieldValue    = "direct field value"
	NodeTypeIndirectFieldValue  = "indirect field value"
	NodeTypeFullSrec            = "full record"
	NodeTypeDirectOosvarValue   = "direct oosvar value"
	NodeTypeIndirectOosvarValue = "indirect oosvar value"
	NodeTypeFullOosvar          = "full oosvar"
	NodeTypeLocalVariable       = "local variable"
	NodeTypeTypedecl            = "type declaration"

	NodeTypeStatementBlock = "statement block"
	NodeTypeAssignment     = "assignment"
	NodeTypeUnset          = "unset"

	NodeTypeBareBoolean      = "bare boolean"
	NodeTypeFilterStatement  = "filter statement"
	NodeTypeEmitStatement    = "emit statement"
	NodeTypeDumpStatement    = "dump statement"
	NodeTypeEdumpStatement   = "edump statement"
	NodeTypePrintStatement   = "print statement"
	NodeTypeEprintStatement  = "eprint statement"
	NodeTypePrintnStatement  = "printn statement"
	NodeTypeEprintnStatement = "eprintn statement"

	NodeTypeOperator         = "operator"
	NodeTypeFunctionCallsite = "function callsite"

	NodeTypeBeginBlock           = "begin block"
	NodeTypeEndBlock             = "end block"
	NodeTypeIfChain              = "if-chain"
	NodeTypeIfItem               = "if-item"
	NodeTypeCondBlock            = "cond block"
	NodeTypeWhileLoop            = "while loop"
	NodeTypeDoWhileLoop          = "do-while`loop"
	NodeTypeForLoopOneVariable   = "single-variable for-loop"
	NodeTypeForLoopTwoVariable   = "double-variable for-loop"
	NodeTypeForLoopMultivariable = "multi-variable for-loop"
	NodeTypeTripleForLoop        = "triple-for loop"
	NodeTypeBreak                = "break"
	NodeTypeContinue             = "continue"

	NodeTypeFunctionDefinition   = "function definition"
	NodeTypeSubroutineDefinition = "subroutine definition"
	NodeTypeParameterList        = "parameter list"
	NodeTypeParameter            = "parameter"
	NodeTypeParameterName        = "parameter name"
	NodeTypeReturn               = "return"

	// A special token which causes a panic when evaluated.  This is for
	// testing that AND/OR short-circuiting is implemented correctly: output =
	// input1 || panic should NOT panic the process when input1 is true.
	NodeTypePanic = "panic token"
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

func (this *AST) PrintParex() {
	this.RootNode.PrintParex()
}

func (this *AST) PrintParexOneLine() {
	this.RootNode.PrintParexOneLine()
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

// ----------------------------------------------------------------
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
		fmt.Printf(" \"%s\"", string(tok.Lit))
	}
	fmt.Println()
	if this.Children != nil {
		for _, child := range this.Children {
			child.PrintAux(depth + 1)
		}
	}
}

// ----------------------------------------------------------------
func (this *ASTNode) PrintParex() {
	this.PrintParexAux(0)
}

func (this *ASTNode) PrintParexAux(depth int) {
	if this.IsLeaf() {
		for i := 0; i < depth; i++ {
			fmt.Print("    ")
		}
		fmt.Print(this.Text())
		fmt.Println()

	} else if this.ChildrenAreAllLeaves() {
		for i := 0; i < depth; i++ {
			fmt.Print("    ")
		}
		fmt.Print("(")
		fmt.Print(this.Text())

		for _, child := range this.Children {
			fmt.Print(" ")
			fmt.Print(child.Text())
		}
		fmt.Println(")")

	} else {
		for i := 0; i < depth; i++ {
			fmt.Print("    ")
		}
		fmt.Print("(")
		fmt.Print(this.Text())
		fmt.Println()

		for _, child := range this.Children {
			child.PrintParexAux(depth + 1)
		}

		for i := 0; i < depth; i++ {
			fmt.Print("    ")
		}
		fmt.Println(")")
	}
}

// ----------------------------------------------------------------
func (this *ASTNode) PrintParexOneLine() {
	this.PrintParexOneLineAux()
	fmt.Println()
}

func (this *ASTNode) PrintParexOneLineAux() {
	if this.IsLeaf() {
		fmt.Print(this.Text())
	} else {
		fmt.Print("(")
		fmt.Print(this.Text())
		for _, child := range this.Children {
			fmt.Print(" ")
			child.PrintParexOneLineAux()
		}
		fmt.Print(")")
	}
}

// ----------------------------------------------------------------
func (this *ASTNode) IsLeaf() bool {
	return this.Children == nil || len(this.Children) == 0
}

func (this *ASTNode) ChildrenAreAllLeaves() bool {
	for _, child := range this.Children {
		if !child.IsLeaf() {
			return false
		}
	}
	return true
}

// ----------------------------------------------------------------
func NewASTNode(itok interface{}, nodeType TNodeType) (*ASTNode, error) {
	return NewASTNodeNestable(itok, nodeType), nil
}

// For handling empty expressions.
func NewASTNodeEmptyNestable(nodeType TNodeType) *ASTNode {
	return &ASTNode{
		Token:    nil,
		Type:     nodeType,
		Children: nil,
	}
}

// For handling empty expressions.
func NewASTNodeEmpty(nodeType TNodeType) (*ASTNode, error) {
	return NewASTNodeEmptyNestable(nodeType), nil
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
	return NewASTNodeNestable(newToken, nodeType), nil
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
	return NewASTNodeNestable(newToken, nodeType), nil
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
	contents = strings.Replace(contents, "\\\"", "\"", -1)
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
func NewASTNodeBinary(
	itok, childA, childB interface{}, nodeType TNodeType,
) (*ASTNode, error) {
	return NewASTNodeBinaryNestable(itok, childA, childB, nodeType), nil
}

func NewASTNodeTernary(itok, childA, childB, childC interface{}, nodeType TNodeType) (*ASTNode, error) {
	parent := NewASTNodeNestable(itok, nodeType)
	convertToTernary(parent, childA, childB, childC)
	return parent, nil
}

func NewASTNodeQuaternary(
	itok, childA, childB, childC, childD interface{}, nodeType TNodeType,
) (*ASTNode, error) {
	parent := NewASTNodeNestable(itok, nodeType)
	convertToQuaternary(parent, childA, childB, childC, childD)
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

func convertToQuaternary(iparent interface{}, childA, childB, childC, childD interface{}) {
	parent := iparent.(*ASTNode)
	children := make([]*ASTNode, 4)
	children[0] = childA.(*ASTNode)
	children[1] = childB.(*ASTNode)
	children[2] = childC.(*ASTNode)
	children[3] = childD.(*ASTNode)
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

func PrependTwoChildren(iparent interface{}, ichildA, ichildB interface{}) (*ASTNode, error) {
	parent := iparent.(*ASTNode)
	childA := ichildA.(*ASTNode)
	childB := ichildB.(*ASTNode)
	if parent.Children == nil {
		convertToBinary(iparent, ichildA, ichildB)
	} else {
		parent.Children = append([]*ASTNode{childA, childB}, parent.Children...)
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

// TODO: comment
func Wrap(inode interface{}) (*ASTNode, error) {
	node := inode.(*ASTNode)
	return node, nil
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

// ----------------------------------------------------------------
func (this *ASTNode) Text() string {
	tokenText := ""
	if this.Token != nil {
		tokenText = string(this.Token.Lit)
	}

	switch this.Type {

	case NodeTypeEmptyStatement:
		return "empty"
	case NodeTypeStringLiteral:
		return "\"" + strings.ReplaceAll(tokenText, "\"", "\\\"") + "\""
	case NodeTypeIntLiteral:
		return tokenText
	case NodeTypeFloatLiteral:
		return tokenText
	case NodeTypeBoolLiteral:
		return tokenText
	case NodeTypeArrayLiteral:
		return tokenText
	case NodeTypeMapLiteral:
		return tokenText
	case NodeTypeMapLiteralKeyValuePair:
		return tokenText

	case NodeTypeArrayOrMapIndexAccess:
		return "[]"
	case NodeTypeArraySliceAccess:
		return "[:]"
	case NodeTypeArraySliceEmptyLowerIndex:
		return "array-slice-empty-lower-index"
	case NodeTypeArraySliceEmptyUpperIndex:
		return "array-slice-empty-upper-index"
	case NodeTypeContextVariable:
		return tokenText
	case NodeTypeConstant:
		return tokenText

	case NodeTypeDirectFieldValue:
		return "$" + tokenText
	case NodeTypeIndirectFieldValue:
		return "$[" + tokenText + "]"
	case NodeTypeFullSrec:
		return tokenText
	case NodeTypeDirectOosvarValue:
		return "@" + tokenText
	case NodeTypeIndirectOosvarValue:
		return "@[" + tokenText + "]"
	case NodeTypeFullOosvar:
		return tokenText
	case NodeTypeLocalVariable:
		return tokenText
	case NodeTypeTypedecl:
		return tokenText

	case NodeTypeStatementBlock:
		return "statement-block"
	case NodeTypeAssignment:
		return tokenText
	case NodeTypeUnset:
		return tokenText

	case NodeTypeBareBoolean:
		return "bare-boolean"
	case NodeTypeFilterStatement:
		return tokenText
	case NodeTypeEmitStatement:
		return tokenText
	case NodeTypeDumpStatement:
		return tokenText
	case NodeTypeEdumpStatement:
		return tokenText
	case NodeTypePrintStatement:
		return tokenText
	case NodeTypeEprintStatement:
		return tokenText
		return tokenText
	case NodeTypePrintnStatement:
		return tokenText
	case NodeTypeEprintnStatement:
		return tokenText

	case NodeTypeOperator:
		return tokenText
	case NodeTypeFunctionCallsite:
		return tokenText

	case NodeTypeBeginBlock:
		return "begin"
	case NodeTypeEndBlock:
		return "end"
	case NodeTypeIfChain:
		return "if-chain"
	case NodeTypeIfItem:
		return tokenText
	case NodeTypeCondBlock:
		return "cond"
	case NodeTypeWhileLoop:
		return tokenText
	case NodeTypeDoWhileLoop:
		return tokenText
	case NodeTypeForLoopOneVariable:
		return tokenText
	case NodeTypeForLoopTwoVariable:
		return tokenText
	case NodeTypeForLoopMultivariable:
		return tokenText
	case NodeTypeTripleForLoop:
		return tokenText
	case NodeTypeBreak:
		return tokenText
	case NodeTypeContinue:
		return tokenText

	case NodeTypeFunctionDefinition:
		return "func"
	case NodeTypeSubroutineDefinition:
		return "subr"
	case NodeTypeParameterList:
		return "parameters"
	case NodeTypeParameter:
		return "parameter"
	case NodeTypeParameterName:
		return tokenText
	case NodeTypeReturn:
		return tokenText

	case NodeTypePanic:
		return tokenText

	}
	return "[ERROR]"
}
