// ================================================================
// Support for user-defined functions
// ================================================================

package cst

import (
	"errors"
	"fmt"
	"os"

	"miller/src/dsl"
	"miller/src/lib"
	"miller/src/runtime"
	"miller/src/types"
)

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
	signature := NewSignature(functionName, callsiteArity, nil, nil)
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

func (this *UDFCallsite) Evaluate(state *runtime.State) types.Mlrval {
	lib.InternalCodingErrorIf(this.argumentNodes == nil)
	lib.InternalCodingErrorIf(this.udf == nil)
	lib.InternalCodingErrorIf(this.udf.functionBody == nil)

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
	// we push a new frameset and BindVariable using the callee's frameset.

	// Evaluate the arguments
	numArguments := len(this.udf.signature.typeGatedParameterNames)
	arguments := make([]types.Mlrval, numArguments)

	for i, typeGatedParameterName := range this.udf.signature.typeGatedParameterNames {
		argument := this.argumentNodes[i].Evaluate(state)

		err := typeGatedParameterName.Check(&argument)
		if err != nil {
			// TODO: put error-return in the Evaluate API
			fmt.Fprint(
				os.Stderr,
				err,
			)
			os.Exit(1)
		}

		arguments[i] = argument
	}

	// Bind the arguments to the parameters
	state.Stack.PushStackFrameSet()
	defer state.Stack.PopStackFrameSet()

	for i, argument := range arguments {
		state.Stack.BindVariable(this.udf.signature.typeGatedParameterNames[i].Name, &argument)
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

	err = this.udf.signature.typeGatedReturnValue.Check(blockExitPayload.blockReturnValue)
	if err != nil {
		// TODO: put error-return in the Evaluate API
		fmt.Fprint(
			os.Stderr,
			err,
		)
		os.Exit(1)
	}

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
	this.functions[udf.signature.funcOrSubrName] = udf
}

func (this *UDFManager) ExistsByName(name string) bool {
	_, ok := this.functions[name]
	return ok
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

	if BuiltinFunctionManagerInstance.LookUp(functionName) != nil {
		return errors.New(
			fmt.Sprintf(
				"Miller: function named \"%s\" must not override a built-in function of the same name.",
				functionName,
			),
		)
	}

	if !this.allowUDFUDSRedefinitions {
		if this.udfManager.ExistsByName(functionName) {
			return errors.New(
				fmt.Sprintf(
					"Miller: function named \"%s\" has already been defined.",
					functionName,
				),
			)
		}
	}

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
			return err
		}

		typeGatedParameterNames[i] = typeGatedParameterName
	}

	signature := NewSignature(functionName, arity, typeGatedParameterNames, typeGatedReturnValue)

	functionBody, err := this.BuildStatementBlockNode(functionBodyASTNode)
	if err != nil {
		return err
	}

	udf := NewUDF(signature, functionBody)

	this.udfManager.Install(udf)

	return nil
}
