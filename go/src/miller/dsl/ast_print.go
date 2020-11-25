// ================================================================
// Print routines for AST and ASTNode
// ================================================================

package dsl

import (
	"fmt"
	"strings"
)

// ================================================================
// Indent-style multiline print.
// Example, given parse of '$y = 2 * $x + 1':
//
// * statement block
//     * assignment "="
//         * direct field value "y"
//         * operator "+"
//             * operator "*"
//                 * int literal "2"
//                 * direct field value "x"
//             * int literal "1"

func (this *AST) Print() {
	this.RootNode.Print()
}

// Parenthesized-expression print.
// Example, given parse of '$y = 2 * $x + 1':
//
// (statement-block
//     (=
//         $y
//         (+
//             (* 2 $x)
//             1
//         )
//     )
// )

func (this *AST) PrintParex() {
	this.RootNode.PrintParex()
}

// Parenthesized-expression print, all on one line.
// Example, given parse of '$y = 2 * $x + 1':
//
// (statement-block (= $y (+ (* 2 $x) 1)))

func (this *AST) PrintParexOneLine() {
	this.RootNode.PrintParexOneLine()
}

// ================================================================
// Indent-style multiline print.
func (this *ASTNode) Print() {
	this.PrintAux(0)
}

func (this *ASTNode) PrintAux(depth int) {
	// Indent
	for i := 0; i < depth; i++ {
		fmt.Print("    ")
	}

	// Token text (if non-nil) and token type
	tok := this.Token
	fmt.Print("* " + this.Type)
	if tok != nil {
		fmt.Printf(" \"%s\"", string(tok.Lit))
	}
	fmt.Println()

	// Children, indented one level further
	if this.Children != nil {
		for _, child := range this.Children {
			child.PrintAux(depth + 1)
		}
	}
}

// ----------------------------------------------------------------
// Parenthesized-expression print.

func (this *ASTNode) PrintParex() {
	this.PrintParexAux(0)
}

func (this *ASTNode) PrintParexAux(depth int) {
	if this.IsLeaf() {
		for i := 0; i < depth; i++ {
			fmt.Print("    ")
		}
		fmt.Println(this.Text())

	} else if this.ChildrenAreAllLeaves() {
		// E.g. (= sum 0) or (+ 1 2)
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
		// Parent and opening parenthesis on first line
		for i := 0; i < depth; i++ {
			fmt.Print("    ")
		}
		fmt.Print("(")
		fmt.Println(this.Text())

		// Children on their own lines
		for _, child := range this.Children {
			child.PrintParexAux(depth + 1)
		}

		// Closing parenthesis on last line
		for i := 0; i < depth; i++ {
			fmt.Print("    ")
		}
		fmt.Println(")")
	}
}

// ----------------------------------------------------------------
// Parenthesized-expression print, all on one line.

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
// Some nodes have non-nil tokens; other, nil. And token-types can have spaces
// in them. In this method we use custom mappings to always get a
// whitespace-free representation of the content of a single AST node.

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
	case NodeTypeEnvironmentVariable:
		return "ENV[\"" + tokenText + "\"]"

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
