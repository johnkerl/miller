// ================================================================
// CST build/execute for AST operator/function nodes, and subroutine nodes.
//
// Operators and functions are semantically the same thing -- they differ only
// syntactically. Binary operators are infix, like '1+2', while functions are
// prefix, like 'max(1,2)'. Both parse to the same AST shape.
//
// Subroutines can't be used as rvalues; their invocation must be the entire
// statement. Nonetheless, their name-resolution, argument/parameter binding,
// etc. are very similar to functions.
// ================================================================

package cst

import (
	"miller/dsl"
	"miller/lib"
)

// ----------------------------------------------------------------
// Function lookup:
//
// * Try builtins first
// * Absent a match there, try UDF lookup (UDF has been defined before being called)
// * Absent a match there:
//   o Make a UDF-placeholder node with present signature but nil function-pointer
//   o Append that node to CST to-be-resolved list
//   o On a next pass, we will walk that list resolving against all encountered
//     UDF definitions. (It will be an error then if it's still unresolvable.)

func (this *RootNode) BuildFunctionCallsiteNode(astNode *dsl.ASTNode) (IEvaluable, error) {
	lib.InternalCodingErrorIf(
		astNode.Type != dsl.NodeTypeFunctionCallsite &&
			astNode.Type != dsl.NodeTypeOperator,
	)
	lib.InternalCodingErrorIf(astNode.Token == nil)
	lib.InternalCodingErrorIf(astNode.Children == nil)

	functionName := string(astNode.Token.Lit)

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Look for a builtin function with the given name.

	builtinFunctionCallsiteNode, err := this.BuildBuiltinFunctionCallsiteNode(astNode)
	if err != nil {
		return nil, err
	}
	if builtinFunctionCallsiteNode != nil {
		return builtinFunctionCallsiteNode, nil
	}

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Look for a user-defined function with the given name.

	callsiteArity := len(astNode.Children)
	udf, err := this.udfManager.LookUp(functionName, callsiteArity)
	if err != nil {
		return nil, err
	}

	// AST snippet for '$z = f($x, $y)':
	// * Assignment "="
	//     * DirectFieldValue "z"
	//     * FunctionCallsite "f"
	//         * DirectFieldValue "x"
	//         * DirectFieldValue "y"
	//
	// Here we need to make an array of our arguments at the callsite, to be
	// paired up with the parameters within he function definition at runtime.
	argumentNodes := make([]IEvaluable, callsiteArity)
	for i, argumentASTNode := range astNode.Children {
		argumentNode, err := this.BuildEvaluableNode(argumentASTNode)
		if err != nil {
			return nil, err
		}
		argumentNodes[i] = argumentNode
	}

	if udf == nil {
		// Mark this as unresolved for an after-pass to see if a UDF with this
		// name/arity has been defined farther down in the DSL expression after
		// this callsite. This happens example when a function is called before
		// it's defined.
		udf = NewUnresolvedUDF(functionName, callsiteArity)
		udfCallsiteNode := NewUDFCallsite(argumentNodes, udf)
		this.rememberUnresolvedFunctionCallsite(udfCallsiteNode)
		return udfCallsiteNode, nil
	} else {
		udfCallsiteNode := NewUDFCallsite(argumentNodes, udf)
		return udfCallsiteNode, nil
	}
}

func (this *RootNode) BuildSubroutineCallsiteNode(astNode *dsl.ASTNode) (IExecutable, error) {
	return nil, nil // TODO
}
