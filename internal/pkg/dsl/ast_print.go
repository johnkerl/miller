// ================================================================
// Print routines for AST and ASTNode
// ================================================================

package dsl

import (
	"fmt"
	"strings"
)

// ================================================================
// Print is indent-style multiline print.
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

func (ast *AST) Print() {
	ast.RootNode.Print()
}

// PrintParex is parenthesized-expression print.
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

func (ast *AST) PrintParex() {
	ast.RootNode.PrintParex()
}

// PrintParexOneLine is parenthesized-expression print, all on one line.
// Example, given parse of '$y = 2 * $x + 1':
//
// (statement-block (= $y (+ (* 2 $x) 1)))

func (ast *AST) PrintParexOneLine() {
	ast.RootNode.PrintParexOneLine()
}

// ================================================================

// Print is indent-style multiline print.
func (node *ASTNode) Print() {
	node.printAux(0)
}

// printAux is a recursion-helper for Print.
func (node *ASTNode) printAux(depth int) {
	// Indent
	for i := 0; i < depth; i++ {
		fmt.Print("    ")
	}

	// Token text (if non-nil) and token type
	tok := node.Token
	fmt.Print("* " + node.Type)
	if tok != nil {
		fmt.Printf(" \"%s\"", string(tok.Lit))
	}
	fmt.Println()

	// Children, indented one level further
	if node.Children != nil {
		for _, child := range node.Children {
			child.printAux(depth + 1)
		}
	}
}

// ----------------------------------------------------------------

// PrintParex is parenthesized-expression print.
func (node *ASTNode) PrintParex() {
	node.printParexAux(0)
}

// printParexAux is a recursion-helper for PrintParex.
func (node *ASTNode) printParexAux(depth int) {
	if node.IsLeaf() {
		for i := 0; i < depth; i++ {
			fmt.Print("    ")
		}
		fmt.Println(node.Text())

	} else if node.ChildrenAreAllLeaves() {
		// E.g. (= sum 0) or (+ 1 2)
		for i := 0; i < depth; i++ {
			fmt.Print("    ")
		}
		fmt.Print("(")
		fmt.Print(node.Text())

		for _, child := range node.Children {
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
		fmt.Println(node.Text())

		// Children on their own lines
		for _, child := range node.Children {
			child.printParexAux(depth + 1)
		}

		// Closing parenthesis on last line
		for i := 0; i < depth; i++ {
			fmt.Print("    ")
		}
		fmt.Println(")")
	}
}

// ----------------------------------------------------------------

// PrintParexOneLine is parenthesized-expression print, all on one line.
func (node *ASTNode) PrintParexOneLine() {
	node.printParexOneLineAux()
	fmt.Println()
}

// printParexOneLineAux is a recursion-helper for PrintParexOneLine.
func (node *ASTNode) printParexOneLineAux() {
	if node.IsLeaf() {
		fmt.Print(node.Text())
	} else {
		fmt.Print("(")
		fmt.Print(node.Text())
		for _, child := range node.Children {
			fmt.Print(" ")
			child.printParexOneLineAux()
		}
		fmt.Print(")")
	}
}

// ----------------------------------------------------------------

// IsLeaf determines if an AST node is a leaf node.
func (node *ASTNode) IsLeaf() bool {
	return node.Children == nil || len(node.Children) == 0
}

// ChildrenAreAllLeaves determines if an AST node's children are all leaf nodes.
func (node *ASTNode) ChildrenAreAllLeaves() bool {
	for _, child := range node.Children {
		if !child.IsLeaf() {
			return false
		}
	}
	return true
}

// ----------------------------------------------------------------

// Text makes a human-readable, whitespace-free name for an AST node. Some
// nodes have non-nil tokens; other, nil. And token-types can have spaces in
// them. In this method we use custom mappings to always get a whitespace-free
// representation of the content of a single AST node.
func (node *ASTNode) Text() string {
	tokenText := ""
	if node.Token != nil {
		tokenText = string(node.Token.Lit)
	}

	switch node.Type {

	case NodeTypeStringLiteral:
		return "\"" + strings.ReplaceAll(tokenText, "\"", "\\\"") + "\""
	case NodeTypeIntLiteral:
		return tokenText
	case NodeTypeFloatLiteral:
		return tokenText
	case NodeTypeBoolLiteral:
		return tokenText
	case NodeTypeNullLiteral:
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
	case NodeTypeEmit1Statement:
		return tokenText
	case NodeTypeEmitStatement:
		return tokenText
	case NodeTypeEmitPStatement:
		return tokenText
	case NodeTypeEmitFStatement:
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

	case NodeTypeNoOp:
		return "no-op"

	case NodeTypeOperator:
		return tokenText
	case NodeTypeDotOperator:
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

	case NodeTypeNamedFunctionDefinition:
		return "named-udf"
	case NodeTypeUnnamedFunctionDefinition:
		return "unnamed-udf"
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
