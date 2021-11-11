// ================================================================
// Support for user-defined functions
// ================================================================

package cst

import (
	"errors"
	"fmt"
	"os"

	"mlr/internal/pkg/dsl"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/runtime"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
type UDF struct {
	signature    *Signature
	functionBody *StatementBlockNode
	// Function literals can access locals in their enclosing scope; named
	// functions cannot.
	isFunctionLiteral bool
}

func NewUDF(
	signature *Signature,
	functionBody *StatementBlockNode,
	isFunctionLiteral bool,
) *UDF {
	return &UDF{
		signature:         signature,
		functionBody:      functionBody,
		isFunctionLiteral: isFunctionLiteral,
	}
}

// For when a function is called before being defined. This gives us something
// to go back and fill in later once we've encountered the function definition.
func NewUnresolvedUDF(
	functionName string,
	callsiteArity int,
) *UDF {
	signature := NewSignature(functionName, callsiteArity, nil, nil)
	udf := NewUDF(signature, nil, false)
	return udf
}

// ----------------------------------------------------------------
type UDFCallsite struct {
	argumentNodes []IEvaluable

	// Non-nil if name was resolved at CST-build time, including named UDFs
	// mutually-recursively calling each other. Nil if the function is in
	// a local variable, like 'f = func(a,b) { return a*b }; z = f(x,y)'.
	udf *UDF

	// Used if the function is in a local variable.
	stackVariable *runtime.StackVariable
	functionName  string
	arity         int
}

// NewUDFCallsite is for the normal UDF callsites outside of sortaf/sortmf,
// e.g. $z = f($a+$b, $c/2). The argument nodes are evaluables since they need
// to be computed, e.g. binding the field names a,b,c, evaluating the
// arithmetic operators, etc.
func NewUDFCallsite(
	argumentNodes []IEvaluable,
	udf *UDF,
) *UDFCallsite {
	functionName := udf.signature.funcOrSubrName
	arity := udf.signature.arity
	return &UDFCallsite{
		argumentNodes: argumentNodes,
		udf:           udf,
		stackVariable: runtime.NewStackVariable(functionName),
		functionName:  functionName,
		arity:         arity,
	}
}

// NewUDFCallsiteForHigherOrderFunction is for UDF callsites such as
// sortaf/sortmf.  Here, the array/map to be sorted has already been evaluated
// and is an array of *types.Mlrval.  The UDF needs to be invoked on pairs of
// array elements.
func NewUDFCallsiteForHigherOrderFunction(
	udf *UDF,
	arity int,
) *UDFCallsite {
	return &UDFCallsite{
		udf:   udf,
		arity: arity,
	}
}

func (site *UDFCallsite) findUDF(state *runtime.State) *UDF {
	if site.udf != nil {
		// Name already resolved at CST-build time
		return site.udf
	}

	// Try stack variable, e.g. the "f" in '$z = f($x, $y)', and supposing
	// there was 'f = func(a, b) { return a*b }' in scope.
	v := state.Stack.Get(site.stackVariable)
	if v == nil { // Nothing in scope on the stack with that name
		// StackVariable
		return nil
	}

	iudf := v.GetFunction()
	if iudf == nil { // Something in scope on the stack with that name, but it's not a function
		return nil
	}

	// func-type mlrvals have only interface{} as value, to avoid what would
	// otherwise be a cyclic package dependency. Here, we deference it.
	return iudf.(*UDF)
}

// Evaluate is for the normal UDF callsites outside of sortaf/sortmf.
// See comments above NewUDFCallsite.
func (site *UDFCallsite) Evaluate(
	state *runtime.State,
) *types.Mlrval {

	udf := site.findUDF(state)
	if udf == nil {
		fmt.Fprintln(os.Stderr, "mlr: function name not found: "+site.functionName)
		os.Exit(1)
	}
	lib.InternalCodingErrorIf(udf.functionBody == nil)
	lib.InternalCodingErrorIf(site.argumentNodes == nil)

	// Evaluate and pair up the callsite arguments with our parameters,
	// positionally.
	//
	// This needs to be a two-step process, for the following reason.
	//
	// The Miller-DSL stack has 'framesets' and 'frames'. For example:
	//
	//   x = 1;                        | Frameset 1
	//   y = 2;                        | Frame 1a: x=1, y=2
	//   if (NR > 10) {                  | Frameset 1b:
	//     x = 3;                        | updates 1a's x; new y=4
	//     var y = 4;                    |
	//   }                             |
	//   func f() {                        | Frameset 2
	//                                     | Frame 2a
	//     x = 5;                          | x = 5, doesn't affect caller's frames
	//     if (some condition) {           |
	//       x = 6;                          | Frame 2b: updates x from from 2a
	//     }                               |
	//   }                                 |
	//
	// We allow scope-walk within a frameset -- so the 1b reference to x
	// updates 1a's x, while 1b's reference to y binds its own y (due to
	// 'var'). But we don't allow scope-walks across framesets with or without
	// 'var': the function's locals are fenced off from the caller's locals.
	//
	// All well and good. What affects us here is callsites of the form
	//
	//   x = 1;
	//   y = f(x);
	//   func f(n) {
	//     return n**2;
	//   }
	//
	// The code in this method implements the line 'y = f(x)', setting up for
	// the call to f(n). Due to the fencing mentioned above, we need to
	// evaluate the argument 'x' using the caller's frameset, but bind it to
	// the callee's parameter 'n' using the callee's frameset.
	//
	// That's why we have two loops here: the first evaluates the arguments
	// using the caller's frameset, stashing them in the arguments array.  Then
	// we push a new frameset and DefineTypedAtScope using the callee's frameset.

	// Evaluate the arguments
	numArguments := len(site.argumentNodes)
	numParameters := len(udf.signature.typeGatedParameterNames)

	if numArguments != numParameters {
		fmt.Fprintf(
			os.Stderr,
			"mlr: function \"%s\" invoked with argument count %d; expected %d.\n",
			udf.signature.funcOrSubrName, numArguments, numParameters)
		os.Exit(1)
	}

	arguments := make([]*types.Mlrval, numArguments)

	for i := range udf.signature.typeGatedParameterNames {
		arguments[i] = site.argumentNodes[i].Evaluate(state)

		err := udf.signature.typeGatedParameterNames[i].Check(arguments[i])
		if err != nil {
			// TODO: put error-return in the Evaluate API
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	}

	return site.EvaluateWithArguments(state, udf, arguments)
}

// EvaluateWithArguments is for UDF callsites in sortaf/sortmf, where the
// arguments are already evaluated. Or, for normal UDF callsites, as a helper
// function for Evaluate.
func (site *UDFCallsite) EvaluateWithArguments(
	state *runtime.State,
	udf *UDF,
	arguments []*types.Mlrval,
) *types.Mlrval {

	// Bind the arguments to the parameters.  Function literals can access
	// locals in their enclosing scope; named functions cannot. Hence stack
	// frame (scope-walkable) vs stack frame set (not scope-walkable).
	if udf.isFunctionLiteral {
		state.Stack.PushStackFrame()
		defer state.Stack.PopStackFrame()
	} else {
		state.Stack.PushStackFrameSet()
		defer state.Stack.PopStackFrameSet()
	}

	cacheable := !udf.isFunctionLiteral

	for i := range arguments {
		// TODO: comment
		err := state.Stack.DefineTypedAtScope(
			runtime.NewStackVariableAux(
				udf.signature.typeGatedParameterNames[i].Name,
				cacheable,
			),
			udf.signature.typeGatedParameterNames[i].TypeName,
			arguments[i],
		)
		// TODO: put error-return in the Evaluate API
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	}

	// Execute the function body.
	blockExitPayload, err := udf.functionBody.Execute(state)

	// TODO: rethink error-propagation here: blockExitPayload.blockReturnValue
	// being MT_ERROR should be mapped to MT_ERROR here (nominally,
	// data-dependent). But error-return could be something not data-dependent.
	if err != nil {
		err = udf.signature.typeGatedReturnValue.Check(types.MLRVAL_ERROR)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
		return types.MLRVAL_ERROR
	}

	// Fell off end of function with no return
	if blockExitPayload == nil {
		err = udf.signature.typeGatedReturnValue.Check(types.MLRVAL_ABSENT)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
		return types.MLRVAL_ABSENT
	}

	// TODO: should be an internal coding error. This would be break or
	// continue not in a loop, or return-void, both of which should have been
	// reported as syntax errors during the parsing pass.
	if blockExitPayload.blockExitStatus != BLOCK_EXIT_RETURN_VALUE {
		err = udf.signature.typeGatedReturnValue.Check(types.MLRVAL_ABSENT)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
		return types.MLRVAL_ABSENT
	}

	// Definitely a Miller internal coding error if the user put 'return x' in
	// their UDF but we lost the return value.
	lib.InternalCodingErrorIf(blockExitPayload.blockReturnValue == nil)

	err = udf.signature.typeGatedReturnValue.Check(blockExitPayload.blockReturnValue)
	if err != nil {
		// TODO: put error-return in the Evaluate API
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	return blockExitPayload.blockReturnValue.Copy()
}

// ----------------------------------------------------------------

// UDFManager tracks named UDFs like 'func f(a, b) { return b - a }'
type UDFManager struct {
	functions map[string]*UDF
}

// NewUDFManager creates an empty UDFManager.
func NewUDFManager() *UDFManager {
	return &UDFManager{
		functions: make(map[string]*UDF),
	}
}

func (manager *UDFManager) Install(udf *UDF) {
	manager.functions[udf.signature.funcOrSubrName] = udf
}

func (manager *UDFManager) ExistsByName(name string) bool {
	_, ok := manager.functions[name]
	return ok
}

// LookUp is for callsites invoking UDFs whose names are known at CST-build time.
func (manager *UDFManager) LookUp(functionName string, callsiteArity int) (*UDF, error) {
	udf := manager.functions[functionName]
	if udf == nil {
		return nil, nil
	}
	if udf.signature.arity != callsiteArity {
		return nil, errors.New(
			fmt.Sprintf(
				"mlr: function %s invoked with %d argument%s; expected %d",
				functionName,
				callsiteArity,
				lib.Plural(callsiteArity),
				udf.signature.arity,
			),
		)
	}
	return udf, nil
}

// LookUpDisregardingArity is used for evaluating right-hand sides of 'f = udf'
// where f will be a local variable of type funct and udf is an existing UDF.
func (manager *UDFManager) LookUpDisregardingArity(functionName string) *UDF {
	return manager.functions[functionName] // nil if not found
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
// AST:
// * StatementBlock
//     * NamedFunctionDefinition "f"
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

// BuildAndInstallUDF is for named UDFs, like `func f(a, b) { return b - a}'.
func (root *RootNode) BuildAndInstallUDF(astNode *dsl.ASTNode) error {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeNamedFunctionDefinition)
	lib.InternalCodingErrorIf(astNode.Children == nil)
	lib.InternalCodingErrorIf(len(astNode.Children) != 2 && len(astNode.Children) != 3)

	functionName := string(astNode.Token.Lit)

	if BuiltinFunctionManagerInstance.LookUp(functionName) != nil {
		return errors.New(
			fmt.Sprintf(
				"mlr: function named \"%s\" must not override a built-in function of the same name.",
				functionName,
			),
		)
	}

	if !root.allowUDFUDSRedefinitions {
		if root.udfManager.ExistsByName(functionName) {
			return errors.New(
				fmt.Sprintf(
					"mlr: function named \"%s\" has already been defined.",
					functionName,
				),
			)
		}
	}

	udf, err := root.BuildUDF(astNode, functionName, false)
	if err != nil {
		return err
	}

	root.udfManager.Install(udf)

	return nil
}

// ================================================================

var namelessFunctionCounter int = 0

// genFunctionLiteralName provides a UUID for function-literal nodes like `func (a, b) { return b - a }'.
// Even nameless function literals need some sort of name for caching purposes.
func genFunctionLiteralName() string {
	namelessFunctionCounter++
	return fmt.Sprintf("function-literal-%06d", namelessFunctionCounter)
}

// ----------------------------------------------------------------

// UnnamedUDFNode holds function literals like 'func (a, b) { return b - a }'.
type UnnamedUDFNode struct {
	udfAsMlrval *types.Mlrval
}

func (root *RootNode) BuildUnnamedUDFNode(astNode *dsl.ASTNode) (IEvaluable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeUnnamedFunctionDefinition)

	name := genFunctionLiteralName()

	udf, err := root.BuildUDF(astNode, name, true)
	if err != nil {
		return nil, err
	}

	udfAsMlrval := types.MlrvalFromFunction(udf, name)

	return &UnnamedUDFNode{
		udfAsMlrval: udfAsMlrval,
	}, nil
}

func (node *UnnamedUDFNode) Evaluate(state *runtime.State) *types.Mlrval {
	return node.udfAsMlrval
}

// ================================================================

// BuildUDF is for named UDFs, like `func f(a, b) { return b - a}', or,
// unnamed UDFs like `func (a, b) { return b - a }'.
func (root *RootNode) BuildUDF(
	astNode *dsl.ASTNode,
	functionName string,
	isFunctionLiteral bool,
) (*UDF, error) {
	lib.InternalCodingErrorIf(
		(astNode.Type != dsl.NodeTypeNamedFunctionDefinition) &&
			(astNode.Type != dsl.NodeTypeUnnamedFunctionDefinition))

	lib.InternalCodingErrorIf(astNode.Children == nil)
	lib.InternalCodingErrorIf(len(astNode.Children) != 2 && len(astNode.Children) != 3)

	parameterListASTNode := astNode.Children[0]
	functionBodyASTNode := astNode.Children[1]

	returnValueTypeName := "any"
	if len(astNode.Children) == 3 {
		typeNode := astNode.Children[2]
		lib.InternalCodingErrorIf(typeNode.Type != dsl.NodeTypeTypedecl)
		returnValueTypeName = string(typeNode.Token.Lit)
	}
	typeGatedReturnValue, err := types.NewTypeGatedMlrvalName(
		"function return value",
		returnValueTypeName,
	)

	lib.InternalCodingErrorIf(parameterListASTNode.Type != dsl.NodeTypeParameterList)
	lib.InternalCodingErrorIf(parameterListASTNode.Children == nil)
	arity := len(parameterListASTNode.Children)
	typeGatedParameterNames := make([]*types.TypeGatedMlrvalName, arity)
	for i, parameterASTNode := range parameterListASTNode.Children {
		lib.InternalCodingErrorIf(parameterASTNode.Type != dsl.NodeTypeParameter)
		lib.InternalCodingErrorIf(parameterASTNode.Children == nil)
		lib.InternalCodingErrorIf(len(parameterASTNode.Children) != 1)
		typeGatedParameterNameASTNode := parameterASTNode.Children[0]

		lib.InternalCodingErrorIf(typeGatedParameterNameASTNode.Type != dsl.NodeTypeParameterName)
		variableName := string(typeGatedParameterNameASTNode.Token.Lit)
		typeName := "any"
		if typeGatedParameterNameASTNode.Children != nil { // typed parameter like 'num x'
			lib.InternalCodingErrorIf(len(typeGatedParameterNameASTNode.Children) != 1)
			typeNode := typeGatedParameterNameASTNode.Children[0]
			lib.InternalCodingErrorIf(typeNode.Type != dsl.NodeTypeTypedecl)
			typeName = string(typeNode.Token.Lit)
		}
		typeGatedParameterName, err := types.NewTypeGatedMlrvalName(
			variableName,
			typeName,
		)
		if err != nil {
			return nil, err
		}

		typeGatedParameterNames[i] = typeGatedParameterName
	}

	signature := NewSignature(functionName, arity, typeGatedParameterNames, typeGatedReturnValue)

	functionBody, err := root.BuildStatementBlockNode(functionBodyASTNode)
	if err != nil {
		return nil, err
	}

	return NewUDF(signature, functionBody, isFunctionLiteral), nil
}
