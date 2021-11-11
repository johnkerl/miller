// ================================================================
// Support for user-defined subroutines
// ================================================================

package cst

import (
	"errors"
	"fmt"

	"mlr/internal/pkg/dsl"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/runtime"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
type UDS struct {
	signature      *Signature
	subroutineBody *StatementBlockNode
}

func NewUDS(
	signature *Signature,
	subroutineBody *StatementBlockNode,
) *UDS {
	return &UDS{
		signature:      signature,
		subroutineBody: subroutineBody,
	}
}

// For when a subroutine is called before being defined. This gives us something
// to go back and fill in later once we've encountered the subroutine definition.
func NewUnresolvedUDS(
	subroutineName string,
	callsiteArity int,
) *UDS {
	signature := NewSignature(subroutineName, callsiteArity, nil, nil)
	uds := NewUDS(signature, nil)
	return uds
}

// ----------------------------------------------------------------
type UDSCallsite struct {
	argumentNodes []IEvaluable
	uds           *UDS
}

func NewUDSCallsite(
	argumentNodes []IEvaluable,
	uds *UDS,
) *UDSCallsite {
	return &UDSCallsite{
		argumentNodes: argumentNodes,
		uds:           uds,
	}
}

func (site *UDSCallsite) Execute(state *runtime.State) (*BlockExitPayload, error) {
	lib.InternalCodingErrorIf(site.argumentNodes == nil)
	lib.InternalCodingErrorIf(site.uds == nil)
	lib.InternalCodingErrorIf(site.uds.subroutineBody == nil)

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
	// 'var': the subroutine's locals are fenced off from the caller's locals.
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
	arguments := make([]*types.Mlrval, len(site.uds.signature.typeGatedParameterNames))

	for i, typeGatedParameterName := range site.uds.signature.typeGatedParameterNames {
		arguments[i] = site.argumentNodes[i].Evaluate(state)

		err := typeGatedParameterName.Check(arguments[i])
		if err != nil {
			return nil, err
		}
	}

	// Bind the arguments to the parameters
	state.Stack.PushStackFrameSet()
	defer state.Stack.PopStackFrameSet()

	for i := range arguments {
		err := state.Stack.DefineTypedAtScope(
			runtime.NewStackVariable(site.uds.signature.typeGatedParameterNames[i].Name),
			site.uds.signature.typeGatedParameterNames[i].TypeName,
			arguments[i],
		)
		if err != nil {
			return nil, err
		}
	}

	// Execute the subroutine body.
	blockExitPayload, err := site.uds.subroutineBody.Execute(state)

	if err != nil {
		return nil, err
	}

	// Fell off end of subroutine with no return
	if blockExitPayload == nil {
		return nil, nil
	}

	// TODO: should be an internal coding error. This would be break or
	// continue not in a loop, or return-void, both of which should have been
	// reported as syntax errors during the parsing pass.
	lib.InternalCodingErrorIf(blockExitPayload.blockExitStatus != BLOCK_EXIT_RETURN_VOID)

	// Subroutines can't return values: 'return' not 'return x'. This should
	// have been caught in the AST validator.
	lib.InternalCodingErrorIf(blockExitPayload.blockReturnValue != nil)

	return blockExitPayload, nil
}

// ----------------------------------------------------------------
type UDSManager struct {
	subroutines map[string]*UDS
}

func NewUDSManager() *UDSManager {
	return &UDSManager{
		subroutines: make(map[string]*UDS),
	}
}

func (manager *UDSManager) LookUp(subroutineName string, callsiteArity int) (*UDS, error) {
	uds := manager.subroutines[subroutineName]
	if uds == nil {
		return nil, nil
	}
	if uds.signature.arity != callsiteArity {
		return nil, errors.New(
			fmt.Sprintf(
				"mlr: subroutine %s invoked with %d argument%s; expected %d",
				subroutineName,
				callsiteArity,
				lib.Plural(callsiteArity),
				uds.signature.arity,
			),
		)
	}
	return uds, nil
}

func (manager *UDSManager) Install(uds *UDS) {
	manager.subroutines[uds.signature.funcOrSubrName] = uds
}

func (manager *UDSManager) ExistsByName(name string) bool {
	_, ok := manager.subroutines[name]
	return ok
}

// ----------------------------------------------------------------
// Example AST for UDS definition and callsite:

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
//     * SubroutineDefinition "f"
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
//         * SubroutineCallsite "f"
//             * DirectFieldValue "x"

func (root *RootNode) BuildAndInstallUDS(astNode *dsl.ASTNode) error {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeSubroutineDefinition)
	lib.InternalCodingErrorIf(astNode.Children == nil)
	lib.InternalCodingErrorIf(len(astNode.Children) != 2 && len(astNode.Children) != 3)

	subroutineName := string(astNode.Token.Lit)

	if !root.allowUDFUDSRedefinitions {
		if root.udsManager.ExistsByName(subroutineName) {
			return errors.New(
				fmt.Sprintf(
					"mlr: subroutine named \"%s\" has already been defined.",
					subroutineName,
				),
			)
		}
	}

	parameterListASTNode := astNode.Children[0]
	subroutineBodyASTNode := astNode.Children[1]

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

	signature := NewSignature(subroutineName, arity, typeGatedParameterNames, nil)

	subroutineBody, err := root.BuildStatementBlockNode(subroutineBodyASTNode)
	if err != nil {
		return err
	}

	uds := NewUDS(signature, subroutineBody)

	root.udsManager.Install(uds)

	return nil
}
