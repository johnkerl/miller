// ================================================================
// AST and ASTNode data structures for the Miller DSL parser
// ================================================================

package dsl

import (
	"mlr/internal/pkg/parsing/token"
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
	NodeTypeStringLiteral             TNodeType = "string literal"
	NodeTypeRegex                               = "regular expression"                  // not in the BNF -- written during CST pre-build pass
	NodeTypeRegexCaseInsensitive                = "case-insensitive regular expression" // E.g. "a.*b"i -- note the trailing 'i'
	NodeTypeIntLiteral                          = "int literal"
	NodeTypeFloatLiteral                        = "float literal"
	NodeTypeBoolLiteral                         = "bool literal"
	NodeTypeNullLiteral                         = "null literal"
	NodeTypeArrayLiteral                        = "array literal"
	NodeTypeMapLiteral                          = "map literal"
	NodeTypeMapLiteralKeyValuePair              = "map-literal key-value pair"
	NodeTypeArrayOrMapIndexAccess               = "array or map index access"
	NodeTypeArraySliceAccess                    = "array-slice access"
	NodeTypeArraySliceEmptyLowerIndex           = "array-slice empty lower index"
	NodeTypeArraySliceEmptyUpperIndex           = "array-slice empty upper index"

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

	NodeTypeBareBoolean     = "bare boolean"
	NodeTypeFilterStatement = "filter statement"

	NodeTypeTeeStatement     = "tee statement"
	NodeTypeEmit1Statement   = "emit1 statement"
	NodeTypeEmitStatement    = "emit statement"
	NodeTypeEmitPStatement   = "emitp statement"
	NodeTypeEmitFStatement   = "emitf statement"
	NodeTypeEmittableList    = "emittable list"
	NodeTypeEmitKeys         = "emit keys"
	NodeTypeDumpStatement    = "dump statement"
	NodeTypeEdumpStatement   = "edump statement"
	NodeTypePrintStatement   = "print statement"
	NodeTypeEprintStatement  = "eprint statement"
	NodeTypePrintnStatement  = "printn statement"
	NodeTypeEprintnStatement = "eprintn statement"

	// For 'print > filename, "string"' et al.
	NodeTypeRedirectWrite        = "redirect write"
	NodeTypeRedirectAppend       = "redirect append"
	NodeTypeRedirectPipe         = "redirect pipe"
	NodeTypeRedirectTargetStdout = "stdout redirect target"
	NodeTypeRedirectTargetStderr = "stderr redirect target"
	NodeTypeRedirectTarget       = "redirect target"

	// This helps various emit-variant sub-ASTs have the same shape.  For
	// example, in 'emit > "foo.txt", @v' and 'emit @v', the latter has a no-op
	// for its redirect target.
	NodeTypeNoOp = "no-op"

	// The dot operator is a little different from other operators since it's
	// type-dependent: for strings/int/bools etc it's just concatenation of
	// string representations, but if the left-hand side is a map, it's a
	// key-lookup with an unquoted literal on the right. E.g. mymap.foo is the
	// same as mymap["foo"].
	NodeTypeOperator           = "operator"
	NodeTypeDotOperator        = "dot operator"
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

	NodeTypeNamedFunctionDefinition   = "function definition"
	NodeTypeUnnamedFunctionDefinition = "function literal"
	NodeTypeSubroutineDefinition      = "subroutine definition"
	NodeTypeParameterList             = "parameter list"
	NodeTypeParameter                 = "parameter"
	NodeTypeParameterName             = "parameter name"
	NodeTypeReturn                    = "return"

	// A special token which causes a panic when evaluated.  This is for
	// testing that AND/OR short-circuiting is implemented correctly: output =
	// input1 || panic should NOT panic the process when input1 is true.
	NodeTypePanic = "panic token"
)
