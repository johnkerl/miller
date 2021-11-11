// ================================================================
// This is for various flavors of for-loop
// ================================================================

package cst

import (
	"errors"

	"mlr/internal/pkg/dsl"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/runtime"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
// Sample AST:

// mlr -n put -v 'for (k in $*) { emit { k : k } }'
// DSL EXPRESSION:
// for (k in $*) { emit { k : k} }
// AST:
// * StatementBlock
//     * ForLoopOneVariable "for"
//         * LocalVariable "k"
//         * FullSrec "$*"
//         * StatementBlock
//             * EmitStatement "emit"
//                 * MapLiteral "{}"
//                     * MapLiteralTwoVariablePair ":"
//                         * LocalVariable "k"
//                         * LocalVariable "k"

// ================================================================
type ForLoopOneVariableNode struct {
	indexVariable      *runtime.StackVariable
	indexableNode      IEvaluable
	statementBlockNode *StatementBlockNode
}

func NewForLoopOneVariableNode(
	variableName string,
	indexableNode IEvaluable,
	statementBlockNode *StatementBlockNode,
) *ForLoopOneVariableNode {
	return &ForLoopOneVariableNode{
		runtime.NewStackVariable(variableName),
		indexableNode,
		statementBlockNode,
	}
}

// ----------------------------------------------------------------
// Sample AST:

// mlr -n put -v 'for (k in $*) { emit { k : k } }'
// DSL EXPRESSION:
// for (k, v in $*) { emit { k : v } }
// AST:
// * StatementBlock
//     * ForLoopOneVariable "for"
//         * LocalVariable "k"
//         * LocalVariable "v"
//         * FullSrec "$*"
//         * StatementBlock
//             * EmitStatement "emit"
//                 * MapLiteral "{}"
//                     * MapLiteralOneVariablePair ":"
//                         * LocalVariable "k"
//                         * LocalVariable "v"

func (root *RootNode) BuildForLoopOneVariableNode(astNode *dsl.ASTNode) (*ForLoopOneVariableNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeForLoopOneVariable)
	lib.InternalCodingErrorIf(len(astNode.Children) != 3)

	variableASTNode := astNode.Children[0]
	indexableASTNode := astNode.Children[1]
	blockASTNode := astNode.Children[2]

	lib.InternalCodingErrorIf(variableASTNode.Type != dsl.NodeTypeLocalVariable)
	lib.InternalCodingErrorIf(variableASTNode.Token == nil)
	variableName := string(variableASTNode.Token.Lit)

	// TODO: error if loop-over node isn't map/array (inasmuch as can be
	// detected at CST-build time)
	indexableNode, err := root.BuildEvaluableNode(indexableASTNode)
	if err != nil {
		return nil, err
	}

	statementBlockNode, err := root.BuildStatementBlockNode(blockASTNode)
	if err != nil {
		return nil, err
	}

	return NewForLoopOneVariableNode(
		variableName,
		indexableNode,
		statementBlockNode,
	), nil
}

// ----------------------------------------------------------------
// Note: The statement-block has its own push/pop for its localvars.
// Meanwhile we need to restrict scope of the bindvar to the for-loop.
//
// So we have:
//
//   mlr put '
//     x = 1;           <--- frame #1 main
//     for (k in $*) {  <--- frame #2 for for-loop bindvars (right here)
//       x = 2          <--- frame #3 for for-loop locals
//     }
//     x = 3;           <--- back in frame #1 main
//   '
//

func (node *ForLoopOneVariableNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	indexMlrval := node.indexableNode.Evaluate(state)

	if indexMlrval.IsMap() {

		mapval := indexMlrval.GetMap()

		// Make a frame for the loop variable(s)
		state.Stack.PushStackFrame()
		defer state.Stack.PopStackFrame()
		for pe := mapval.Head; pe != nil; pe = pe.Next {
			mapkey := types.MlrvalFromString(pe.Key)

			err := state.Stack.SetAtScope(node.indexVariable, mapkey)
			if err != nil {
				return nil, err
			}
			// The loop body will push its own frame
			blockExitPayload, err := node.statementBlockNode.Execute(state)
			if err != nil {
				return nil, err
			}
			if blockExitPayload != nil {
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_BREAK {
					break
				}
				// If BLOCK_EXIT_CONTINUE, keep going -- this means the body was exited
				// early but we keep going at this level
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_RETURN_VOID {
					return blockExitPayload, nil
				}
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_RETURN_VALUE {
					return blockExitPayload, nil
				}
			}
		}

	} else if indexMlrval.IsArray() {

		arrayval := indexMlrval.GetArray()

		// Note: Miller user-space array indices ("mindex") are 1-up. Internal
		// Go storage ("zindex") is 0-up.

		// Make a frame for the loop variable(s)
		state.Stack.PushStackFrame()
		defer state.Stack.PopStackFrame()
		for _, element := range arrayval {
			err := state.Stack.SetAtScope(node.indexVariable, &element)
			if err != nil {
				return nil, err
			}
			// The loop body will push its own frame
			blockExitPayload, err := node.statementBlockNode.Execute(state)
			if err != nil {
				return nil, err
			}
			if blockExitPayload != nil {
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_BREAK {
					break
				}
				// If BLOCK_EXIT_CONTINUE, keep going -- this means the body was exited
				// early but we keep going at this level
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_RETURN_VOID {
					return blockExitPayload, nil
				}
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_RETURN_VALUE {
					return blockExitPayload, nil
				}
			}
		}

	} else if indexMlrval.IsAbsent() {
		// Data-heterogeneity no-op
	}

	// TODO: backwards compatibility with the C port means we treat this as
	// silent zero-pass. But maybe we should surface it as an error. Maybe
	// with a "mlr put --errors" flag or something.
	//	} else {
	//		return nil, errors.New(
	//			fmt.Sprintf(
	//				"mlr: looped-over item is not a map or array; got %s",
	//				indexMlrval.GetTypeName(),
	//			),
	//		)

	return nil, nil
}

// ================================================================
type ForLoopTwoVariableNode struct {
	keyIndexVariable   *runtime.StackVariable
	valueIndexVariable *runtime.StackVariable
	indexableNode      IEvaluable
	statementBlockNode *StatementBlockNode
}

func NewForLoopTwoVariableNode(
	keyIndexVariable *runtime.StackVariable,
	valueIndexVariable *runtime.StackVariable,
	indexableNode IEvaluable,
	statementBlockNode *StatementBlockNode,
) *ForLoopTwoVariableNode {
	return &ForLoopTwoVariableNode{
		keyIndexVariable,
		valueIndexVariable,
		indexableNode,
		statementBlockNode,
	}
}

// ----------------------------------------------------------------
// Sample AST:

// mlr -n put -v 'for (k, v in $*) { emit { k : v } }'
// DSL EXPRESSION:
// for (k, v in $*) { emit { k : v } }
// AST:
// * StatementBlock
//     * ForLoopTwoVariable "for"
//         * LocalVariable "k"
//         * LocalVariable "v"
//         * FullSrec "$*"
//         * StatementBlock
//             * EmitStatement "emit"
//                 * MapLiteral "{}"
//                     * MapLiteralTwoVariablePair ":"
//                         * LocalVariable "k"
//                         * LocalVariable "v"

func (root *RootNode) BuildForLoopTwoVariableNode(astNode *dsl.ASTNode) (*ForLoopTwoVariableNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeForLoopTwoVariable)
	lib.InternalCodingErrorIf(len(astNode.Children) != 4)

	keyVariableASTNode := astNode.Children[0]
	valueVariableASTNode := astNode.Children[1]
	indexableASTNode := astNode.Children[2]
	blockASTNode := astNode.Children[3]

	lib.InternalCodingErrorIf(keyVariableASTNode.Type != dsl.NodeTypeLocalVariable)
	lib.InternalCodingErrorIf(keyVariableASTNode.Token == nil)
	keyVariableName := string(keyVariableASTNode.Token.Lit)
	keyIndexVariable := runtime.NewStackVariable(keyVariableName)

	lib.InternalCodingErrorIf(valueVariableASTNode.Type != dsl.NodeTypeLocalVariable)
	lib.InternalCodingErrorIf(valueVariableASTNode.Token == nil)
	valueVariableName := string(valueVariableASTNode.Token.Lit)
	valueIndexVariable := runtime.NewStackVariable(valueVariableName)

	// TODO: error if loop-over node isn't map/array (inasmuch as can be
	// detected at CST-build time)
	indexableNode, err := root.BuildEvaluableNode(indexableASTNode)
	if err != nil {
		return nil, err
	}

	statementBlockNode, err := root.BuildStatementBlockNode(blockASTNode)
	if err != nil {
		return nil, err
	}

	return NewForLoopTwoVariableNode(
		keyIndexVariable,
		valueIndexVariable,
		indexableNode,
		statementBlockNode,
	), nil
}

// ----------------------------------------------------------------
// Note: The statement-block has its own push/pop for its localvars.
// Meanwhile we need to restrict scope of the bindvar to the for-loop.
//
// So we have:
//
//   mlr put '
//     x = 1;             <--- frame #1 main
//     for (k,v in $*) {  <--- frame #2 for for-loop bindvars (right here)
//       x = 2            <--- frame #3 for for-loop locals
//     }
//     x = 3;             <--- back in frame #1 main
//   '
//

func (node *ForLoopTwoVariableNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	indexMlrval := node.indexableNode.Evaluate(state)

	if indexMlrval.IsMap() {

		mapval := indexMlrval.GetMap()

		// Make a frame for the loop variable(s)
		state.Stack.PushStackFrame()
		defer state.Stack.PopStackFrame()
		for pe := mapval.Head; pe != nil; pe = pe.Next {
			mapkey := types.MlrvalFromString(pe.Key)

			err := state.Stack.SetAtScope(node.keyIndexVariable, mapkey)
			if err != nil {
				return nil, err
			}
			err = state.Stack.SetAtScope(node.valueIndexVariable, pe.Value)
			if err != nil {
				return nil, err
			}
			// The loop body will push its own frame
			blockExitPayload, err := node.statementBlockNode.Execute(state)
			if err != nil {
				return nil, err
			}
			if blockExitPayload != nil {
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_BREAK {
					break
				}
				// If BLOCK_EXIT_CONTINUE, keep going -- this means the body was exited
				// early but we keep going at this level
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_RETURN_VOID {
					return blockExitPayload, nil
				}
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_RETURN_VALUE {
					return blockExitPayload, nil
				}
			}
		}

	} else if indexMlrval.IsArray() {

		arrayval := indexMlrval.GetArray()

		// Note: Miller user-space array indices ("mindex") are 1-up. Internal
		// Go storage ("zindex") is 0-up.

		// Make a frame for the loop variable(s)
		state.Stack.PushStackFrame()
		defer state.Stack.PopStackFrame()
		for zindex, element := range arrayval {
			mindex := types.MlrvalFromInt(int(zindex + 1))

			err := state.Stack.SetAtScope(node.keyIndexVariable, mindex)
			if err != nil {
				return nil, err
			}
			err = state.Stack.SetAtScope(node.valueIndexVariable, &element)
			if err != nil {
				return nil, err
			}
			// The loop body will push its own frame
			blockExitPayload, err := node.statementBlockNode.Execute(state)
			if err != nil {
				return nil, err
			}
			if blockExitPayload != nil {
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_BREAK {
					break
				}
				// If BLOCK_EXIT_CONTINUE, keep going -- this means the body was exited
				// early but we keep going at this level
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_RETURN_VOID {
					return blockExitPayload, nil
				}
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_RETURN_VALUE {
					return blockExitPayload, nil
				}
			}
		}

	} else if indexMlrval.IsAbsent() {
		// Data-heterogeneity no-op
	}

	// TODO: backwards compatibility with the C port means we treat this as
	// silent zero-pass. But maybe we should surface it as an error. Maybe
	// with a "mlr put --errors" flag or something.
	//	} else {
	//		return nil, errors.New(
	//			fmt.Sprintf(
	//				"mlr: looped-over item is not a map or array; got %s",
	//				indexMlrval.GetTypeName(),
	//			),
	//		)

	return nil, nil
}

// ================================================================
type ForLoopMultivariableNode struct {
	keyIndexVariables  []*runtime.StackVariable
	valueIndexVariable *runtime.StackVariable
	indexableNode      IEvaluable
	statementBlockNode *StatementBlockNode
}

func NewForLoopMultivariableNode(
	keyIndexVariables []*runtime.StackVariable,
	valueIndexVariable *runtime.StackVariable,
	indexableNode IEvaluable,
	statementBlockNode *StatementBlockNode,
) *ForLoopMultivariableNode {
	return &ForLoopMultivariableNode{
		keyIndexVariables,
		valueIndexVariable,
		indexableNode,
		statementBlockNode,
	}
}

// ----------------------------------------------------------------
// Sample AST:

// mlr -n put -v 'for ((k1, k2), v in $*) { }'
// DSL EXPRESSION:
// for ((k1, k2), v in $*) { }
// AST:
// * statement block
//     * multi-variable for-loop "for"
//         * parameter list
//             * local variable "k1"
//             * local variable "k2"
//         * local variable "v"
//         * full record "$*"
//         * statement block

func (root *RootNode) BuildForLoopMultivariableNode(
	astNode *dsl.ASTNode,
) (*ForLoopMultivariableNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeForLoopMultivariable)
	lib.InternalCodingErrorIf(len(astNode.Children) != 4)

	keyVariablesASTNode := astNode.Children[0]
	valueVariableASTNode := astNode.Children[1]
	indexableASTNode := astNode.Children[2]
	blockASTNode := astNode.Children[3]

	lib.InternalCodingErrorIf(keyVariablesASTNode.Type != dsl.NodeTypeParameterList)
	lib.InternalCodingErrorIf(keyVariablesASTNode.Children == nil)
	keyIndexVariables := make([]*runtime.StackVariable, len(keyVariablesASTNode.Children))
	for i, keyVariableASTNode := range keyVariablesASTNode.Children {
		lib.InternalCodingErrorIf(keyVariableASTNode.Token == nil)
		keyIndexVariableName := string(keyVariableASTNode.Token.Lit)
		keyIndexVariables[i] = runtime.NewStackVariable(keyIndexVariableName)
	}

	lib.InternalCodingErrorIf(valueVariableASTNode.Type != dsl.NodeTypeLocalVariable)
	lib.InternalCodingErrorIf(valueVariableASTNode.Token == nil)
	valueVariableName := string(valueVariableASTNode.Token.Lit)
	valueIndexVariable := runtime.NewStackVariable(valueVariableName)

	// TODO: error if loop-over node isn't map/array (inasmuch as can be
	// detected at CST-build time)
	indexableNode, err := root.BuildEvaluableNode(indexableASTNode)
	if err != nil {
		return nil, err
	}

	statementBlockNode, err := root.BuildStatementBlockNode(blockASTNode)
	if err != nil {
		return nil, err
	}

	return NewForLoopMultivariableNode(
		keyIndexVariables,
		valueIndexVariable,
		indexableNode,
		statementBlockNode,
	), nil
}

// ----------------------------------------------------------------
// Note: The statement-block has its own push/pop for its localvars.
// Meanwhile we need to restrict scope of the bindvar to the for-loop.
//
// So we have:
//
//   mlr put '
//     x = 1;                   <--- frame #1 main
//     for ((k1,k2),v in $*) {  <--- frame #2 for for-loop bindvars (right here)
//       x = 2                  <--- frame #3 for for-loop locals
//     }
//     x = 3;                   <--- back in frame #1 main
//   '
//

func (node *ForLoopMultivariableNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	indexMlrval := node.indexableNode.Evaluate(state)

	// Make a frame for the loop variables
	state.Stack.PushStackFrame()
	defer state.Stack.PopStackFrame()

	// Miller's multi-variable loops, in the Miller DSL, have a single {...}
	// but are implemented in Go via multiple, recursive functions.  A break
	// from any of the latter is a break from all.  However, at this point, the
	// break has been "broken" and should not be returned to the caller.
	// Return-statements should, though.
	blockExitPayload, err := node.executeOuter(indexMlrval, node.keyIndexVariables, state)
	if blockExitPayload == nil {
		return nil, err
	} else {
		if blockExitPayload.blockExitStatus == BLOCK_EXIT_BREAK {
			return nil, err
		} else {
			return blockExitPayload, err
		}
	}
}

// ----------------------------------------------------------------
func (node *ForLoopMultivariableNode) executeOuter(
	mlrval *types.Mlrval,
	keyIndexVariables []*runtime.StackVariable,
	state *runtime.State,
) (*BlockExitPayload, error) {
	if len(keyIndexVariables) == 1 {
		return node.executeInner(mlrval, keyIndexVariables[0], state)
	}
	// else, recurse

	if mlrval.IsMap() {
		mapval := mlrval.GetMap()

		for pe := mapval.Head; pe != nil; pe = pe.Next {
			mapkey := types.MlrvalFromString(pe.Key)

			err := state.Stack.SetAtScope(keyIndexVariables[0], mapkey)
			if err != nil {
				return nil, err
			}

			blockExitPayload, err := node.executeOuter(pe.Value, keyIndexVariables[1:], state)
			if err != nil {
				return nil, err
			}
			if blockExitPayload != nil {
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_BREAK {
					return blockExitPayload, nil
				}
				// If BLOCK_EXIT_CONTINUE, keep going -- this means the body was exited
				// early but we keep going at this level
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_RETURN_VOID {
					return blockExitPayload, nil
				}
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_RETURN_VALUE {
					return blockExitPayload, nil
				}
			}
		}

	} else if mlrval.IsArray() {
		arrayval := mlrval.GetArray()

		// Note: Miller user-space array indices ("mindex") are 1-up. Internal
		// Go storage ("zindex") is 0-up.

		for zindex, element := range arrayval {
			mindex := types.MlrvalFromInt(int(zindex + 1))

			err := state.Stack.SetAtScope(keyIndexVariables[0], mindex)
			if err != nil {
				return nil, err
			}

			blockExitPayload, err := node.executeOuter(&element, keyIndexVariables[1:], state)
			if err != nil {
				return nil, err
			}
			if blockExitPayload != nil {
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_BREAK {
					return blockExitPayload, nil
					break
				}
				// If BLOCK_EXIT_CONTINUE, keep going -- this means the body was exited
				// early but we keep going at this level
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_RETURN_VOID {
					return blockExitPayload, nil
				}
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_RETURN_VALUE {
					return blockExitPayload, nil
				}
			}
		}

	} else if mlrval.IsAbsent() {
		// Data-heterogeneity no-op
	}

	// TODO: backwards compatibility with the C port means we treat this as
	// silent zero-pass. But maybe we should surface it as an error. Maybe
	// with a "mlr put --errors" flag or something.
	//	} else {
	//		return nil, errors.New(
	//			fmt.Sprintf(
	//				"mlr: looped-over item is not a map or array; got %s",
	//				mlrval.GetTypeName(),
	//			),
	//		)

	return nil, nil
}

// ----------------------------------------------------------------
func (node *ForLoopMultivariableNode) executeInner(
	mlrval *types.Mlrval,
	keyIndexVariable *runtime.StackVariable,
	state *runtime.State,
) (*BlockExitPayload, error) {
	if mlrval.IsMap() {
		mapval := mlrval.GetMap()

		for pe := mapval.Head; pe != nil; pe = pe.Next {
			mapkey := types.MlrvalFromString(pe.Key)

			err := state.Stack.SetAtScope(keyIndexVariable, mapkey)
			if err != nil {
				return nil, err
			}
			err = state.Stack.SetAtScope(node.valueIndexVariable, pe.Value)
			if err != nil {
				return nil, err
			}

			// The loop body will push its own frame
			blockExitPayload, err := node.statementBlockNode.Execute(state)
			if err != nil {
				return nil, err
			}
			if blockExitPayload != nil {
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_BREAK {
					return blockExitPayload, nil
				}
				// If BLOCK_EXIT_CONTINUE, keep going -- this means the body was exited
				// early but we keep going at this level
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_RETURN_VOID {
					return blockExitPayload, nil
				}
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_RETURN_VALUE {
					return blockExitPayload, nil
				}
			}
		}

	} else if mlrval.IsArray() {
		arrayval := mlrval.GetArray()

		// Note: Miller user-space array indices ("mindex") are 1-up. Internal
		// Go storage ("zindex") is 0-up.

		for zindex, element := range arrayval {
			mindex := types.MlrvalFromInt(int(zindex + 1))

			err := state.Stack.SetAtScope(keyIndexVariable, mindex)
			if err != nil {
				return nil, err
			}
			err = state.Stack.SetAtScope(node.valueIndexVariable, &element)
			if err != nil {
				return nil, err
			}

			// The loop body will push its own frame
			blockExitPayload, err := node.statementBlockNode.Execute(state)
			if err != nil {
				return nil, err
			}
			if blockExitPayload != nil {
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_BREAK {
					return blockExitPayload, nil
				}
				// If BLOCK_EXIT_CONTINUE, keep going -- this means the body was exited
				// early but we keep going at this level
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_RETURN_VOID {
					return blockExitPayload, nil
				}
				if blockExitPayload.blockExitStatus == BLOCK_EXIT_RETURN_VALUE {
					return blockExitPayload, nil
				}
			}
		}

	} else if mlrval.IsAbsent() {
		// Data-heterogeneity no-op
	}

	// TODO: backwards compatibility with the C port means we treat this as
	// silent zero-pass. But maybe we should surface it as an error. Maybe
	// with a "mlr put --errors" flag or something.
	//	} else {
	//		return nil, errors.New(
	//			fmt.Sprintf(
	//				"mlr: looped-over item is not a map or array; got %s",
	//				mlrval.GetTypeName(),
	//			),
	//		)

	return nil, nil
}

// ================================================================
type TripleForLoopNode struct {
	startBlockNode             *StatementBlockNode
	precontinuationAssignments []IExecutable
	continuationExpressionNode IEvaluable
	updateBlockNode            *StatementBlockNode
	bodyBlockNode              *StatementBlockNode
}

func NewTripleForLoopNode(
	startBlockNode *StatementBlockNode,
	precontinuationAssignments []IExecutable,
	continuationExpressionNode IEvaluable,
	updateBlockNode *StatementBlockNode,
	bodyBlockNode *StatementBlockNode,
) *TripleForLoopNode {
	return &TripleForLoopNode{
		startBlockNode,
		precontinuationAssignments,
		continuationExpressionNode,
		updateBlockNode,
		bodyBlockNode,
	}
}

// ----------------------------------------------------------------
// Sample ASTs:

// DSL EXPRESSION:
// for (;;) {}
// AST:
// * StatementBlock
//     * TripleForLoop "for"
//         * StatementBlock
//         * StatementBlock
//         * StatementBlock
//         * StatementBlock

// mlr --from u/s.dkvp put -v for (i = 0; i < NR; i += 1) { $i += i }
// DSL EXPRESSION:
// for (i = 0; i < NR; i += 1) { $i += i }
// AST:
// * StatementBlock
//     * TripleForLoop "for"
//         * StatementBlock
//             * Assignment "="
//                 * LocalVariable "i"
//                 * IntLiteral "0"
//         * StatementBlock
//             * BareBoolean
//                 * Operator "<"
//                     * LocalVariable "i"
//                     * ContextVariable "NR"
//         * StatementBlock
//             * Assignment "="
//                 * LocalVariable "i"
//                 * Operator "+"
//                     * LocalVariable "i"
//                     * IntLiteral "1"
//         * StatementBlock
//             * Assignment "="
//                 * DirectFieldValue "i"
//                 * Operator "+"
//                     * DirectFieldValue "i"
//                     * LocalVariable "i"

func (root *RootNode) BuildTripleForLoopNode(astNode *dsl.ASTNode) (*TripleForLoopNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeTripleForLoop)
	lib.InternalCodingErrorIf(len(astNode.Children) != 4)

	startBlockASTNode := astNode.Children[0]
	continuationExpressionASTNode := astNode.Children[1]
	updateBlockASTNode := astNode.Children[2]
	bodyBlockASTNode := astNode.Children[3]

	lib.InternalCodingErrorIf(startBlockASTNode.Type != dsl.NodeTypeStatementBlock)
	lib.InternalCodingErrorIf(continuationExpressionASTNode.Type != dsl.NodeTypeStatementBlock)
	lib.InternalCodingErrorIf(updateBlockASTNode.Type != dsl.NodeTypeStatementBlock)
	lib.InternalCodingErrorIf(bodyBlockASTNode.Type != dsl.NodeTypeStatementBlock)

	startBlockNode, err := root.BuildStatementBlockNode(startBlockASTNode)
	if err != nil {
		return nil, err
	}

	// Enforced here, not in the grammar: the last must be a bare boolean; the ones
	// before must be assignments. Example:
	// for (int i = 0; c += 1, i < 10; i += 1) { ... }
	var precontinuationAssignments []IExecutable = nil
	var continuationExpressionNode IEvaluable = nil
	if len(continuationExpressionASTNode.Children) > 0 { // empty is true
		n := len(continuationExpressionASTNode.Children)
		if n > 1 {
			precontinuationAssignments = make([]IExecutable, n-1)
			for i := 0; i < n-1; i++ {
				if continuationExpressionASTNode.Children[i].Type != dsl.NodeTypeAssignment {
					return nil, errors.New(
						"mlr: the non-final triple-for continutation statements must be assignments.",
					)
				}
				precontinuationAssignment, err := root.BuildAssignmentNode(
					continuationExpressionASTNode.Children[i],
				)
				if err != nil {
					return nil, err
				}
				precontinuationAssignments[i] = precontinuationAssignment
			}
		}

		bareBooleanASTNode := continuationExpressionASTNode.Children[n-1]
		if bareBooleanASTNode.Type != dsl.NodeTypeBareBoolean {
			if n == 1 {
				return nil, errors.New(
					"mlr: the triple-for continutation statement must be a bare boolean.",
				)
			} else {
				return nil, errors.New(
					"mlr: the final triple-for continutation statement must be a bare boolean.",
				)
			}
		}
		lib.InternalCodingErrorIf(len(bareBooleanASTNode.Children) != 1)
		continuationExpressionNode, err = root.BuildEvaluableNode(bareBooleanASTNode.Children[0])
		if err != nil {
			return nil, err
		}
	}

	updateBlockNode, err := root.BuildStatementBlockNode(updateBlockASTNode)
	if err != nil {
		return nil, err
	}

	bodyBlockNode, err := root.BuildStatementBlockNode(bodyBlockASTNode)
	if err != nil {
		return nil, err
	}

	return NewTripleForLoopNode(
		startBlockNode,
		precontinuationAssignments,
		continuationExpressionNode,
		updateBlockNode,
		bodyBlockNode,
	), nil
}

// ----------------------------------------------------------------
// Note: The statement-block has its own push/pop for its localvars.
// Meanwhile we need to restrict scope of the bindvar to the for-loop.
//
// So we have:
//
//   mlr put '
//     x = 1;                             <--- frame #1 main
//     for (int i = 0; i < 10; i += 1) {  <--- frame #2 for for-loop bindvars (right here)
//       x = 2                            <--- frame #3 for for-loop locals
//     }
//     x = 3;                             <--- back in frame #1 main
//   '
//

func (node *TripleForLoopNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	// Make a frame for the loop variables.
	state.Stack.PushStackFrame()
	defer state.Stack.PopStackFrame()

	// Use ExecuteFrameless here, otherwise the start-statements would be
	// within an ephemeral, isolated frame and not accessible to the remaining
	// parts of the for-loop.
	_, err := node.startBlockNode.ExecuteFrameless(state)
	if err != nil {
		return nil, err
	}

	for {
		if node.precontinuationAssignments != nil {
			for _, precontinuationAssignment := range node.precontinuationAssignments {
				_, err := precontinuationAssignment.Execute(state)
				if err != nil {
					return nil, err
				}
			}
		}
		if node.continuationExpressionNode != nil { // empty is true
			continuationValue := node.continuationExpressionNode.Evaluate(state)
			boolValue, isBool := continuationValue.GetBoolValue()
			if !isBool {
				// TODO: propagate line-number context
				return nil, errors.New("mlr: for-loop continuation did not evaluate to boolean.")
			}
			if boolValue == false {
				break
			}
		}

		blockExitPayload, err := node.bodyBlockNode.Execute(state)
		if err != nil {
			return nil, err
		}
		if blockExitPayload != nil {
			if blockExitPayload.blockExitStatus == BLOCK_EXIT_BREAK {
				break
			}
			// If BLOCK_EXIT_CONTINUE, keep going -- this means the body was exited
			// early but we keep going at this level. In particular we still
			// need to execute the update-block.
			if blockExitPayload.blockExitStatus == BLOCK_EXIT_RETURN_VOID {
				return blockExitPayload, nil
			}
			if blockExitPayload.blockExitStatus == BLOCK_EXIT_RETURN_VALUE {
				return blockExitPayload, nil
			}
		}

		// The loop body will push its own frame.
		state.Stack.PushStackFrame()
		_, err = node.updateBlockNode.ExecuteFrameless(state)
		if err != nil {
			state.Stack.PopStackFrame()
			return nil, err
		}
		state.Stack.PopStackFrame()
	}

	return nil, nil
}
