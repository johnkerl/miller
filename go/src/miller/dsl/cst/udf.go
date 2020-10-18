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

	state.stack.PushStackFrame()
	defer state.stack.PopStackFrame()

	// TODO: argument-parameter bindings
	for i, parameterName := range this.udf.signature.parameterNames {
		argument := this.argumentNodes[i].Evaluate(state)
		state.stack.BindVariable(parameterName, &argument)
		//fmt.Println("BIND", parameterName, argument.String())
	}

	blockExitPayload, err := this.udf.functionBody.Execute(state)
	if err != nil {
		// TODO: rethink error-propagation here
		return types.MlrvalFromError()
	}
	if blockExitPayload == nil { // Fell off end of function with no return
		return types.MlrvalFromAbsent()
	}
	if blockExitPayload.blockExitStatus != BLOCK_EXIT_RETURN_VALUE {
		// TODO: rethink error-propagation here
		return types.MlrvalFromAbsent()
	}
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
	lib.InternalCodingErrorIf(len(astNode.Children) != 2)

	functionName := string(astNode.Token.Lit)
	parameterListASTNode := astNode.Children[0]
	functionBodyASTNode := astNode.Children[1]

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
		lib.InternalCodingErrorIf(parameterNameASTNode.Children != nil)

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
