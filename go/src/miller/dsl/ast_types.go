// ================================================================
// AST and ASTNode data structures for the Miller DSL parser
// ================================================================

package dsl

import (
	"miller/parsing/token"
)

// ----------------------------------------------------------------
type AST struct {
	RootNode *ASTNode
}

// ----------------------------------------------------------------
type ASTNode struct {
	Token    *token.Token // Nil for tokenless/structural nodes
	Type     TNodeType
	Children []*ASTNode
}

// ----------------------------------------------------------------
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

	NodeTypePositionalFieldName             = "positionally-indexed field name"
	NodeTypePositionalFieldValue            = "positionally-indexed field value"
	NodeTypeArrayOrMapPositionalNameAccess  = "positionally-indexed map key"
	NodeTypeArrayOrMapPositionalValueAccess = "positionally-indexed map value"

	NodeTypeContextVariable     = "context variable"
	NodeTypeConstant            = "mathematical constant"
	NodeTypeEnvironmentVariable = "environment variable"

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
	NodeTypeEmitPStatement   = "emitp statement"
	NodeTypeEmitFStatement   = "emitf statement"
	NodeTypeDumpStatement    = "dump statement"
	NodeTypeEdumpStatement   = "edump statement"
	NodeTypePrintStatement   = "print statement"
	NodeTypeEprintStatement  = "eprint statement"
	NodeTypePrintnStatement  = "printn statement"
	NodeTypeEprintnStatement = "eprintn statement"

	NodeTypeOperator           = "operator"
	NodeTypeFunctionCallsite   = "function callsite"
	NodeTypeSubroutineCallsite = "subroutine callsite"

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
