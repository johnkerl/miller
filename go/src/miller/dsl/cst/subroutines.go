// ================================================================
// CST build/execute for subroutine nodes.
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
// Subroutine lookup:
//
// * Unlike for functions, There are no built-in subroutines -- the only ones
//   that exist are user-defined.
// * Try UDS lookup (i.e. the UDS has been defined before being called)
// * Absent a match there:
//   o Make a UDS-placeholder node with present signature but nil function-pointer
//   o Append that node to CST to-be-resolved list
//   o On a next pass, we will walk that list resolving against all encountered
//     UDS definitions. (It will be an error then if it's still unresolvable.)

func (this *RootNode) BuildSubroutineCallsiteNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(
		astNode.Type != dsl.NodeTypeSubroutineCallsite &&
			astNode.Type != dsl.NodeTypeOperator,
	)
	lib.InternalCodingErrorIf(astNode.Token == nil)
	lib.InternalCodingErrorIf(astNode.Children == nil)

	subroutineName := string(astNode.Token.Lit)

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Look for a user-defined subroutine with the given name.

	callsiteArity := len(astNode.Children)
	uds, err := this.udsManager.LookUp(subroutineName, callsiteArity)
	if err != nil {
		return nil, err
	}

	// AST snippet for '$z = f($x, $y)':
	// * Assignment "="
	//     * DirectFieldValue "z"
	//     * SubroutineCallsite "f"
	//         * DirectFieldValue "x"
	//         * DirectFieldValue "y"
	//
	// Here we need to make an array of our arguments at the callsite, to be
	// paired up with the parameters within he subroutine definition at runtime.
	argumentNodes := make([]IEvaluable, callsiteArity)
	for i, argumentASTNode := range astNode.Children {
		argumentNode, err := this.BuildEvaluableNode(argumentASTNode)
		if err != nil {
			return nil, err
		}
		argumentNodes[i] = argumentNode
	}

	if uds == nil {
		// Mark this as unresolved for an after-pass to see if a UDS with this
		// name/arity has been defined farther down in the DSL expression after
		// this callsite. This happens example when a subroutine is called before
		// it's defined.
		uds = NewUnresolvedUDS(subroutineName, callsiteArity)
		udsCallsiteNode := NewUDSCallsite(argumentNodes, uds)
		this.rememberUnresolvedSubroutineCallsite(udsCallsiteNode)
		return udsCallsiteNode, nil
	} else {
		udsCallsiteNode := NewUDSCallsite(argumentNodes, uds)
		return udsCallsiteNode, nil
	}
}
