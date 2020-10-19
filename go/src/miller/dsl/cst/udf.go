package cst

import (
	"errors"
	"fmt"

	"miller/dsl"
	"miller/lib"
	"miller/types"
)

// ================================================================
// Support for user-defined functions and subroutines
// ================================================================

// ----------------------------------------------------------------
type Signature struct {
	functionName   string
	arity          int // Computable from len(parameterNames) at callee, not at caller
	parameterNames []string

	// TODO: parameter typedecls
	// TODO: return-value typedecls
}

func NewSignature(
	functionName string,
	arity int,
	parameterNames []string,
) *Signature {
	return &Signature{
		functionName:   functionName,
		arity:          arity,
		parameterNames: parameterNames,
	}
}

// ----------------------------------------------------------------
type UDF struct {
	signature    *Signature
	functionBody *StatementBlockNode
}

func NewUDF(
	signature *Signature,
	functionBody *StatementBlockNode,
) *UDF {
	return &UDF{
		signature:    signature,
		functionBody: functionBody,
	}
}

// For when a function is called before being defined. This gives us something
// to go back and fill in later once we've encountered the function definition.
func NewUnresolvedUDF(
	functionName string,
	callsiteArity int,
) *UDF {
	signature := NewSignature(functionName, callsiteArity, nil)
	udf := NewUDF(signature, nil)
	return udf
}

// ----------------------------------------------------------------
type UDFCallsite struct {
	argumentNodes []IEvaluable
	udf           *UDF
}

func NewUDFCallsite(
	argumentNodes []IEvaluable,
	udf *UDF,
) *UDFCallsite {
	return &UDFCallsite{
		argumentNodes: argumentNodes,
		udf:           udf,
	}
}

func (this *UDFCallsite) Evaluate(state *State) types.Mlrval {
	lib.InternalCodingErrorIf(this.argumentNodes == nil)
	lib.InternalCodingErrorIf(this.udf == nil)
	lib.InternalCodingErrorIf(this.udf.functionBody == nil)

	// Evaluate and pair up the callsite arguments with our parameters,
	// positionally.

	state.stack.PushStackFrame()
	defer state.stack.PopStackFrame()

	for i, parameterName := range this.udf.signature.parameterNames {
		argument := this.argumentNodes[i].Evaluate(state)
		state.stack.BindVariable(parameterName, &argument)
	}

	// Execute the function body.
	blockExitPayload, err := this.udf.functionBody.Execute(state)

	// TODO: rethink error-propagation here: blockExitPayload.blockReturnValue
	// being MT_ERROR should be mapped to MT_ERROR here (nominally,
	// data-dependent). But error-return could be something not data-dependent.
	if err != nil {
		return types.MlrvalFromError()
	}

	// Fell off end of function with no return
	if blockExitPayload == nil {
		return types.MlrvalFromAbsent()
	}

	// TODO: should be an internal coding error. This would be break or
	// continue not in a loop, or return-void, both of which should have been
	// reported as syntax errors during the parsing pass.
	if blockExitPayload.blockExitStatus != BLOCK_EXIT_RETURN_VALUE {
		return types.MlrvalFromAbsent()
	}

	// Definitely a Miller internal coding error if the user put 'return x' in
	// their UDF but we lost the return value.
	lib.InternalCodingErrorIf(blockExitPayload.blockReturnValue == nil)

	return *blockExitPayload.blockReturnValue
}

// ----------------------------------------------------------------
type UDFManager struct {
	functions map[string]*UDF
}

func NewUDFManager() *UDFManager {
	return &UDFManager{
		functions: make(map[string]*UDF),
	}
}

func (this *UDFManager) LookUp(functionName string, callsiteArity int) (*UDF, error) {
	udf := this.functions[functionName]
	if udf == nil {
		return nil, nil
	}
	if udf.signature.arity != callsiteArity {
		return nil, errors.New(
			fmt.Sprintf(
				"Miller: function %s invoked with %d argument%s; expected %d",
				functionName,
				callsiteArity,
				lib.Plural(callsiteArity),
				udf.signature.arity,
			),
		)
	}
	return udf, nil
}

func (this *UDFManager) Install(udf *UDF) {
	this.functions[udf.signature.functionName] = udf
}

// ----------------------------------------------------------------
// Example AST for UDF definition and callsite:

// DSL EXPRESSION:
// func f(x) {
//   if (x >= 0) {
//     return x
//   } else {
//     return -x
//   }
// }
//
// $y = f($x)
//
// RAW AST:
// * StatementBlock
//     * FunctionDefinition "f"
//         * ParameterList
//             * Parameter
//                 * ParameterName "x"
//         * StatementBlock
//             * IfChain
//                 * IfItem "if"
//                     * Operator ">="
//                         * LocalVariable "x"
//                         * IntLiteral "0"
//                     * StatementBlock
//                         * Return "return"
//                             * LocalVariable "x"
//                 * IfItem "else"
//                     * StatementBlock
//                         * Return "return"
//                             * Operator "-"
//                                 * LocalVariable "x"
//     * Assignment "="
//         * DirectFieldValue "y"
//         * FunctionCallsite "f"
//             * DirectFieldValue "x"

func (this *RootNode) BuildAndInstallUDF(astNode *dsl.ASTNode) error {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeFunctionDefinition)
	lib.InternalCodingErrorIf(astNode.Children == nil)
	lib.InternalCodingErrorIf(len(astNode.Children) != 2 && len(astNode.Children) != 3)

	functionName := string(astNode.Token.Lit)
	parameterListASTNode := astNode.Children[0]
	functionBodyASTNode := astNode.Children[1]
	// TODO: optional typedecl 3rd arg

	lib.InternalCodingErrorIf(parameterListASTNode.Type != dsl.NodeTypeParameterList)
	lib.InternalCodingErrorIf(parameterListASTNode.Children == nil)
	arity := len(parameterListASTNode.Children)
	parameterNames := make([]string, arity)
	for i, parameterASTNode := range parameterListASTNode.Children {
		lib.InternalCodingErrorIf(parameterASTNode.Type != dsl.NodeTypeParameter)
		lib.InternalCodingErrorIf(parameterASTNode.Children == nil)
		lib.InternalCodingErrorIf(len(parameterASTNode.Children) != 1)
		parameterNameASTNode := parameterASTNode.Children[0]

		lib.InternalCodingErrorIf(parameterNameASTNode.Type != dsl.NodeTypeParameterName)
		if parameterNameASTNode.Children == nil {
			// untyped parameter
		} else {
			// typed parameter
			lib.InternalCodingErrorIf(len(parameterNameASTNode.Children) != 1)
		}

		parameterNames[i] = string(parameterNameASTNode.Token.Lit)
	}

	signature := NewSignature(functionName, arity, parameterNames)

	functionBody, err := this.BuildStatementBlockNode(functionBodyASTNode)
	if err != nil {
		return err
	}

	udf := NewUDF(signature, functionBody)

	this.udfManager.Install(udf)

	return nil
}
