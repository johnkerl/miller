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

	// * Try already-found UDFs first
	// * Try builtins second
	// * Absent either of those, make a UDF-placeholder with present signature but nil function-pointer
	//   o Append node to CST to-be-resolved list
	// * Next pass: we will walk that list resolving against all encountered UDF definitions
	//   o Error then if unresolvable

	// callsiteArity := len(astNode.Children)
	// udfInfo, err := this.udfManager.LookUp(functionName, callsiteArity)
	// if err != nil {
	// 	return nil, err
	// }

	builtinFunctionCallsiteNode, err := this.BuildBuiltinFunctionCallsiteNode(astNode)
	if err != nil {
		return nil, err
	}
	if builtinFunctionCallsiteNode != nil {
		return builtinFunctionCallsiteNode, nil
	}

	// retval := NewUDFCallsitePlaceholder(name, arity)
	// this.RememberUDFCallsitePlaceholder(retval)
	// return retval, nil

	return nil, errors.New(
		"CST BuildFunctionCallsiteNode: function name not found: " +
			functionName,
	)
}
