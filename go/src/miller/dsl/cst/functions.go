package cst

import (
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

// ----------------------------------------------------------------
// Function lookup:
//
// * Try builtins first
// * Absent a match there, try UDF lookup (UDF has been defined before being called)
// * Absent a match there:
//   o Make a UDF-placeholder node with present signature but nil function-pointer
//   o Append that node to CST to-be-resolved list
//   o On a next pass, we will walk that list resolving against all encountered
//     UDF definitions
//     - Error then if still unresolvable

func (this *RootNode) BuildFunctionCallsiteNode(astNode *dsl.ASTNode) (IEvaluable, error) {
	lib.InternalCodingErrorIf(
		astNode.Type != dsl.NodeTypeFunctionCallsite &&
			astNode.Type != dsl.NodeTypeOperator,
	)
	lib.InternalCodingErrorIf(astNode.Token == nil)
	lib.InternalCodingErrorIf(astNode.Children == nil)

	functionName := string(astNode.Token.Lit)

	builtinFunctionCallsiteNode, err := this.BuildBuiltinFunctionCallsiteNode(astNode)
	if err != nil {
		return nil, err
	}
	if builtinFunctionCallsiteNode != nil {
		return builtinFunctionCallsiteNode, nil
	}

	callsiteArity := len(astNode.Children)
	udf, err := this.udfManager.LookUp(functionName, callsiteArity)
	if err != nil {
		return nil, err
	}

    // AST snippet for '$z = f($x, $y)':
    //     * Assignment "="
    //         * DirectFieldValue "z"
    //         * FunctionCallsite "f"
    //             * DirectFieldValue "x"
    //             * DirectFieldValue "y"
	argumentNodes := make([]IEvaluable, callsiteArity)
	for i, argumentASTNode := range(astNode.Children) {
		argumentNode, err := this.BuildEvaluableNode(argumentASTNode)
		if err != nil {
			return nil, err
		}
		argumentNodes[i] = argumentNode
	}

	udfCallsiteNode := NewUDFCallsite(argumentNodes, udf)
	return udfCallsiteNode, nil

	// TODO:
	// retval := NewUDFPlaceholder(name, arity)
	// this.RememberUDFCallsitePlaceholder(retval)
	// return retval, nil

	// TODO: move to builder
	//	return nil, errors.New(
	//		"CST BuildFunctionCallsiteNode: function name not found: " +
	//			functionName,
	//	)
}

// ----------------------------------------------------------------
// TODO: resolver function
