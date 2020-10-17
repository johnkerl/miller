package cst

import (
	"errors"

	"miller/dsl"
	"miller/lib"
)

// ================================================================
// CST build/execute for AST operator/function nodes
//
// Operators and functions are semantically the same thing -- they differ only
// syntactically. Binary operators are infix, like '1+2', while functions are
// prefix, like 'max(1,2)'. Both parse to the same AST shape.
// ================================================================

func (this *RootNode) BuildFunctionCallsiteNode(astNode *dsl.ASTNode) (IEvaluable, error) {
	lib.InternalCodingErrorIf(
		astNode.Type != dsl.NodeTypeFunctionCallsite &&
			astNode.Type != dsl.NodeTypeOperator,
	)
	lib.InternalCodingErrorIf(astNode.Token == nil)
	lib.InternalCodingErrorIf(astNode.Children == nil)

	functionName := string(astNode.Token.Lit)

	// TODO: try UDFs first

	builtinFunctionInfo := BuiltinFunctionManagerInstance.LookUp(functionName)
	if builtinFunctionInfo != nil {
		if builtinFunctionInfo.hasMultipleArities { // E.g. "+" and "-"
			return this.BuildMultipleArityFunctionCallsiteNode(astNode, builtinFunctionInfo)
		} else if builtinFunctionInfo.zaryFunc != nil {
			return this.BuildZaryFunctionCallsiteNode(astNode, builtinFunctionInfo)
		} else if builtinFunctionInfo.unaryFunc != nil {
			return this.BuildUnaryFunctionCallsiteNode(astNode, builtinFunctionInfo)
		} else if builtinFunctionInfo.binaryFunc != nil {
			return this.BuildBinaryFunctionCallsiteNode(astNode, builtinFunctionInfo)
		} else if builtinFunctionInfo.ternaryFunc != nil {
			return this.BuildTernaryFunctionCallsiteNode(astNode, builtinFunctionInfo)
		} else if builtinFunctionInfo.variadicFunc != nil {
			return this.BuildVariadicFunctionCallsiteNode(astNode, builtinFunctionInfo)
		} else {
			return nil, errors.New(
				"CST BuildFunctionCallsiteNode: function not implemented yet: " +
					functionName,
			)
		}
	} else {
		return nil, errors.New(
			"CST BuildFunctionCallsiteNode: function name not found: " +
				functionName,
		)
	}
}
