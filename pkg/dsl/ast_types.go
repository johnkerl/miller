// ================================================================
// AST and ASTNode data structures for the Miller DSL parser
// ================================================================

package dsl

import (
	"github.com/johnkerl/miller/pkg/parsing/token"
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
	NodeTypeRegex                     TNodeType = "regular expression"                  // not in the BNF -- written during CST pre-build pass
	NodeTypeRegexCaseInsensitive      TNodeType = "case-insensitive regular expression" // E.g. "a.*b"i -- note the trailing 'i'
	NodeTypeIntLiteral                TNodeType = "int literal"
	NodeTypeFloatLiteral              TNodeType = "float literal"
	NodeTypeBoolLiteral               TNodeType = "bool literal"
	NodeTypeNullLiteral               TNodeType = "null literal"
	NodeTypeArrayLiteral              TNodeType = "array literal"
	NodeTypeMapLiteral                TNodeType = "map literal"
	NodeTypeMapLiteralKeyValuePair    TNodeType = "map-literal key-value pair"
	NodeTypeArrayOrMapIndexAccess     TNodeType = "array or map index access"
	NodeTypeArraySliceAccess          TNodeType = "array-slice access"
	NodeTypeArraySliceEmptyLowerIndex TNodeType = "array-slice empty lower index"
	NodeTypeArraySliceEmptyUpperIndex TNodeType = "array-slice empty upper index"

	NodeTypePositionalFieldName             TNodeType = "positionally-indexed field name"
	NodeTypePositionalFieldValue            TNodeType = "positionally-indexed field value"
	NodeTypeArrayOrMapPositionalNameAccess  TNodeType = "positionally-indexed map key"
	NodeTypeArrayOrMapPositionalValueAccess TNodeType = "positionally-indexed map value"

	NodeTypeContextVariable     TNodeType = "context variable"
	NodeTypeConstant            TNodeType = "mathematical constant"
	NodeTypeEnvironmentVariable TNodeType = "environment variable"

	NodeTypeDirectFieldValue    TNodeType = "direct field value"
	NodeTypeIndirectFieldValue  TNodeType = "indirect field value"
	NodeTypeFullSrec            TNodeType = "full record"
	NodeTypeDirectOosvarValue   TNodeType = "direct oosvar value"
	NodeTypeIndirectOosvarValue TNodeType = "indirect oosvar value"
	NodeTypeFullOosvar          TNodeType = "full oosvar"
	NodeTypeLocalVariable       TNodeType = "local variable"
	NodeTypeTypedecl            TNodeType = "type declaration"

	NodeTypeStatementBlock TNodeType = "statement block"
	NodeTypeAssignment     TNodeType = "assignment"
	NodeTypeUnset          TNodeType = "unset"

	NodeTypeBareBoolean     TNodeType = "bare boolean"
	NodeTypeFilterStatement TNodeType = "filter statement"

	NodeTypeTeeStatement     TNodeType = "tee statement"
	NodeTypeEmit1Statement   TNodeType = "emit1 statement"
	NodeTypeEmitStatement    TNodeType = "emit statement"
	NodeTypeEmitPStatement   TNodeType = "emitp statement"
	NodeTypeEmitFStatement   TNodeType = "emitf statement"
	NodeTypeEmittableList    TNodeType = "emittable list"
	NodeTypeEmitKeys         TNodeType = "emit keys"
	NodeTypeDumpStatement    TNodeType = "dump statement"
	NodeTypeEdumpStatement   TNodeType = "edump statement"
	NodeTypePrintStatement   TNodeType = "print statement"
	NodeTypeEprintStatement  TNodeType = "eprint statement"
	NodeTypePrintnStatement  TNodeType = "printn statement"
	NodeTypeEprintnStatement TNodeType = "eprintn statement"

	// For 'print > filename, "string"' et al.
	NodeTypeRedirectWrite        TNodeType = "redirect write"
	NodeTypeRedirectAppend       TNodeType = "redirect append"
	NodeTypeRedirectPipe         TNodeType = "redirect pipe"
	NodeTypeRedirectTargetStdout TNodeType = "stdout redirect target"
	NodeTypeRedirectTargetStderr TNodeType = "stderr redirect target"
	NodeTypeRedirectTarget       TNodeType = "redirect target"

	// This helps various emit-variant sub-ASTs have the same shape.  For
	// example, in 'emit > "foo.txt", @v' and 'emit @v', the latter has a no-op
	// for its redirect target.
	NodeTypeNoOp TNodeType = "no-op"

	// The dot operator is a little different from other operators since it's
	// type-dependent: for strings/int/bools etc it's just concatenation of
	// string representations, but if the left-hand side is a map, it's a
	// key-lookup with an unquoted literal on the right. E.g. mymap.foo is the
	// same as mymap["foo"].
	NodeTypeOperator           TNodeType = "operator"
	NodeTypeDotOperator        TNodeType = "dot operator"
	NodeTypeFunctionCallsite   TNodeType = "function callsite"
	NodeTypeSubroutineCallsite TNodeType = "subroutine callsite"

	NodeTypeBeginBlock           TNodeType = "begin block"
	NodeTypeEndBlock             TNodeType = "end block"
	NodeTypeIfChain              TNodeType = "if-chain"
	NodeTypeIfItem               TNodeType = "if-item"
	NodeTypeCondBlock            TNodeType = "cond block"
	NodeTypeWhileLoop            TNodeType = "while loop"
	NodeTypeDoWhileLoop          TNodeType = "do-while`loop"
	NodeTypeForLoopOneVariable   TNodeType = "single-variable for-loop"
	NodeTypeForLoopTwoVariable   TNodeType = "double-variable for-loop"
	NodeTypeForLoopMultivariable TNodeType = "multi-variable for-loop"
	NodeTypeTripleForLoop        TNodeType = "triple-for loop"
	NodeTypeBreak                TNodeType = "break"
	NodeTypeContinue             TNodeType = "continue"

	NodeTypeNamedFunctionDefinition   TNodeType = "function definition"
	NodeTypeUnnamedFunctionDefinition TNodeType = "function literal"
	NodeTypeSubroutineDefinition      TNodeType = "subroutine definition"
	NodeTypeParameterList             TNodeType = "parameter list"
	NodeTypeParameter                 TNodeType = "parameter"
	NodeTypeParameterName             TNodeType = "parameter name"
	NodeTypeReturn                    TNodeType = "return"

	// A special token which causes a panic when evaluated.  This is for
	// testing that AND/OR short-circuiting is implemented correctly: output =
	// input1 || panic should NOT panic the process when input1 is true.
	NodeTypePanic TNodeType = "panic token"
)
