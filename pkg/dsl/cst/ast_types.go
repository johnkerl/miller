// PGPG AST node type constants. These match the "type" field values
// produced by the mlr.bnf grammar. The CST uses these when
// comparing asts.ASTNode.Type (asts.NodeType is string).

package cst

import (
	"fmt"

	"github.com/johnkerl/pgpg/go/lib/pkg/asts"
	"github.com/johnkerl/pgpg/go/lib/pkg/tokens"
)

// PGPG grammar node types (camelCase from grammar hints)
const (
	NodeTypeStatementBlock            = "StatementBlock"
	NodeTypeStatementBlockInBraces    = "StatementBlockInBraces"
	NodeTypeAssignment                = "Assignment"
	NodeTypeCompoundAssignment        = "CompoundAssignment"
	NodeTypeUnset                     = "Unset"
	NodeTypeBareBoolean               = "BareBoolean"
	NodeTypeFilterStatement           = "FilterStatement"
	NodeTypeRedirectWrite             = "RedirectWrite"
	NodeTypeRedirectAppend            = "RedirectAppend"
	NodeTypeRedirectPipe              = "RedirectPipe"
	NodeTypeRedirectTargetStdout      = "RedirectTargetStdout"
	NodeTypeRedirectTargetStderr      = "RedirectTargetStderr"
	NodeTypeRedirectTargetRvalue      = "RedirectTargetRvalue"
	NodeTypePrintStatement            = "PrintStatement"
	NodeTypePrintnStatement           = "PrintnStatement"
	NodeTypeEprintStatement           = "EprintStatement"
	NodeTypeEprintnStatement          = "EprintnStatement"
	NodeTypeDumpStatement             = "DumpStatement"
	NodeTypeEdumpStatement            = "EdumpStatement"
	NodeTypeTeeStatement              = "TeeStatement"
	NodeTypeEmit1Statement            = "Emit1Statement"
	NodeTypeEmitStatement             = "EmitStatement"
	NodeTypeEmitPStatement            = "EmitPStatement"
	NodeTypeEmitFStatement            = "EmitFStatement"
	NodeTypeBeginBlock                = "BeginBlock"
	NodeTypeEndBlock                  = "EndBlock"
	NodeTypeCondBlock                 = "CondBlock"
	NodeTypeIfChain                   = "IfChain"
	NodeTypeIfItem                    = "IfItem"
	NodeTypeWhileLoop                 = "WhileLoop"
	NodeTypeDoWhileLoop               = "DoWhileLoop"
	NodeTypeForLoopOneVariable        = "ForLoopOneVariable"
	NodeTypeForLoopTwoVariable        = "ForLoopTwoVariable"
	NodeTypeForLoopMultivariable      = "ForLoopMultivariable"
	NodeTypeTripleForLoop             = "TripleForLoop"
	NodeTypeBreakStatement            = "BreakStatement"
	NodeTypeContinueStatement         = "ContinueStatement"
	NodeTypeReturnStatement           = "ReturnStatement"
	NodeTypeSubroutineCallsite        = "SubroutineCallsite"
	NodeTypeNamedFunctionDefinition   = "NamedFunctionDefinition"
	NodeTypeSubroutineDefinition      = "SubroutineDefinition"
	NodeTypeUnnamedFunctionDefinition = "UnnamedFunctionDefinition"
	NodeTypeParameterList             = "ParameterList"
	NodeTypeParameter                 = "Parameter"
	NodeTypeDirectFieldValue          = "DirectFieldValue"
	NodeTypeIndirectFieldValue        = "IndirectFieldValue"
	NodeTypeBracedFieldValue          = "BracedFieldValue"
	NodeTypeFullSrec                  = "FullSrec"
	NodeTypeDirectOosvarValue         = "DirectOosvarValue"
	NodeTypeIndirectOosvarValue       = "IndirectOosvarValue"
	NodeTypeBracedOosvarValue         = "BracedOosvarValue"
	NodeTypeFullOosvar                = "FullOosvar"
	NodeTypeLocalVariable             = "LocalVariable"
	NodeTypeOperator                  = "Operator"
	NodeTypeDotOperator               = "DotOperator"
	NodeTypeFunctionCallsite          = "FunctionCallsite"
	NodeTypeArrayLiteral              = "ArrayLiteral"
	NodeTypeMapLiteral                = "MapLiteral"
	NodeTypeMapLiteralKeyValuePair    = "MapLiteralKeyValuePair"
	NodeTypeArrayOrMapIndexAccess     = "ArrayOrMapIndexAccess"
	NodeTypeArraySliceLoHi            = "ArraySliceLoHi"
	NodeTypeArraySliceHiOnly          = "ArraySliceHiOnly"
	NodeTypeArraySliceLoOnly          = "ArraySliceLoOnly"
	NodeTypeArraySliceFull            = "ArraySliceFull"
	NodeTypeParenthesized             = "Parenthesized"
	NodeTypeFcnArgs                   = "FcnArgs"
	NodeTypeIntLiteral                = "int_literal"
	NodeTypeFloatLiteral              = "float_literal"
	NodeTypeStringLiteral             = "string_literal"
	NodeTypeBoolLiteral               = "bool_literal"
	NodeTypeNullLiteral               = "null_literal"
	NodeTypeMultiIndex                = "MultiIndex"
	NodeTypeTypedecl                  = "Typedecl"
	// Synthetic nodes for $[[n]] / $[[[n]]] (positional indexing - not in PGPG grammar)
	NodeTypePositionalFieldName  = "PositionalFieldName"
	NodeTypePositionalFieldValue = "PositionalFieldValue"
	// For array slice empty bounds: [lo:], [:hi], [:]
	NodeTypeArraySliceEmptyLowerIndex = "ArraySliceEmptyLowerIndex"
	NodeTypeArraySliceEmptyUpperIndex = "ArraySliceEmptyUpperIndex"
)

// Injected during regexProtectPrePass (not from grammar)
const NodeTypeRegex = "Regex"

// NoOp: used when optional redirect/expressions are absent (PGPG may not produce these)
const NodeTypeNoOp = "NoOp"

// NodeTypePanic: special token for testing short-circuit (not in grammar)
const NodeTypePanic = "Panic"

// Types not in mlr.bnf (positional indexing, ENV, context vars excluded).
// Added so CST compiles; extend grammar to support.
const (
	NodeTypeArraySliceAccess                = "ArraySliceAccess" // alias for any of ArraySliceLoHi etc
	NodeTypeArrayOrMapPositionalNameAccess  = "ArrayOrMapPositionalNameAccess"
	NodeTypeArrayOrMapPositionalValueAccess = "ArrayOrMapPositionalValueAccess"
	NodeTypeEmitKeys                        = "EmitKeys"
	NodeTypeEnvironmentVariable             = "EnvironmentVariable"
	NodeTypeRegexCaseInsensitive            = "RegexCaseInsensitive"
	NodeTypeContextVariable                 = "ContextVariable"
	NodeTypeConstant                        = "Constant"
)

// tokenLit returns the lexeme text for an AST node's token, or "" if nil.
func tokenLit(node *asts.ASTNode) string {
	if node == nil || node.Token == nil {
		return ""
	}
	return node.Token.LexemeText()
}

// tokenLitStripDollarOrAt returns the lexeme with leading $ or @ stripped.
// Used for DirectFieldValue ($n -> n) and DirectOosvarValue (@x -> x).
func tokenLitStripDollarOrAt(node *asts.ASTNode) string {
	s := tokenLit(node)
	if len(s) >= 1 && (s[0] == '$' || s[0] == '@') {
		return s[1:]
	}
	return s
}

// tokenLitStripBraced returns the lexeme with ${ / @{ and } stripped.
// Used for BracedFieldValue (${foo} -> foo) and BracedOosvarValue (@{bar} -> bar).
func tokenLitStripBraced(node *asts.ASTNode) string {
	s := tokenLit(node)
	if len(s) >= 4 && (s[0] == '$' || s[0] == '@') && s[1] == '{' && s[len(s)-1] == '}' {
		return s[2 : len(s)-1]
	}
	return s
}

// pgpgTokenToLocationInfo formats location info for error messages.
func pgpgTokenToLocationInfo(tok *tokens.Token) string {
	if tok == nil {
		return ""
	}
	return fmt.Sprintf(" at DSL expression line %d column %d", tok.Location.LineNumber, tok.Location.ColumnNumber)
}
