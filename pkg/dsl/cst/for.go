// This is for various flavors of for-loop

package cst

import (
	"fmt"

	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/runtime"
	"github.com/johnkerl/pgpg/go/lib/pkg/asts"
	"github.com/johnkerl/pgpg/go/lib/pkg/tokens"
)

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

func (root *RootNode) BuildForLoopOneVariableNode(astNode *asts.ASTNode) (*ForLoopOneVariableNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != asts.NodeType(NodeTypeForLoopOneVariable))
	lib.InternalCodingErrorIf(len(astNode.Children) != 3)

	variableASTNode := astNode.Children[0]
	indexableASTNode := astNode.Children[1]
	blockASTNode := astNode.Children[2]

	lib.InternalCodingErrorIf(variableASTNode.Type != asts.NodeType(NodeTypeLocalVariable))
	lib.InternalCodingErrorIf(variableASTNode.Token == nil)
	variableName := tokenLit(variableASTNode)

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
			mapkey := mlrval.FromString(pe.Key)

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
			err := state.Stack.SetAtScope(node.indexVariable, element)
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
	//		return nil, fmt.Errorf(
	//			"mlr: looped-over item is not a map or array; got %s",
	//			indexMlrval.GetTypeName(),
	//		)

	return nil, nil
}

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

func (root *RootNode) BuildForLoopTwoVariableNode(astNode *asts.ASTNode) (*ForLoopTwoVariableNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != asts.NodeType(NodeTypeForLoopTwoVariable))
	lib.InternalCodingErrorIf(len(astNode.Children) != 4)

	keyVariableASTNode := astNode.Children[0]
	valueVariableASTNode := astNode.Children[1]
	indexableASTNode := astNode.Children[2]
	blockASTNode := astNode.Children[3]

	lib.InternalCodingErrorIf(keyVariableASTNode.Type != asts.NodeType(NodeTypeLocalVariable))
	lib.InternalCodingErrorIf(keyVariableASTNode.Token == nil)
	keyVariableName := tokenLit(keyVariableASTNode)
	keyIndexVariable := runtime.NewStackVariable(keyVariableName)

	lib.InternalCodingErrorIf(valueVariableASTNode.Type != asts.NodeType(NodeTypeLocalVariable))
	lib.InternalCodingErrorIf(valueVariableASTNode.Token == nil)
	valueVariableName := tokenLit(valueVariableASTNode)
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
			mapkey := mlrval.FromString(pe.Key)

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
			mindex := mlrval.FromInt(int64(zindex + 1))

			err := state.Stack.SetAtScope(node.keyIndexVariable, mindex)
			if err != nil {
				return nil, err
			}
			err = state.Stack.SetAtScope(node.valueIndexVariable, element)
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
	//		return nil, fmt.Errorf(
	//			"mlr: looped-over item is not a map or array; got %s",
	//			indexMlrval.GetTypeName(),
	//		)

	return nil, nil
}

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
	astNode *asts.ASTNode,
) (*ForLoopMultivariableNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != asts.NodeType(NodeTypeForLoopMultivariable))
	lib.InternalCodingErrorIf(len(astNode.Children) != 4)

	keyVariablesASTNode := astNode.Children[0]
	valueVariableASTNode := astNode.Children[1]
	indexableASTNode := astNode.Children[2]
	blockASTNode := astNode.Children[3]

	// PGPG produces MultiIndex; legacy produced ParameterList. Both have LocalVariable children.
	lib.InternalCodingErrorIf(keyVariablesASTNode.Type != asts.NodeType(NodeTypeParameterList) &&
		keyVariablesASTNode.Type != asts.NodeType(NodeTypeMultiIndex))
	lib.InternalCodingErrorIf(keyVariablesASTNode.Children == nil)
	keyIndexVariables := make([]*runtime.StackVariable, len(keyVariablesASTNode.Children))
	for i, keyVariableASTNode := range keyVariablesASTNode.Children {
		lib.InternalCodingErrorIf(keyVariableASTNode.Token == nil)
		keyIndexVariableName := tokenLit(keyVariableASTNode)
		keyIndexVariables[i] = runtime.NewStackVariable(keyIndexVariableName)
	}

	lib.InternalCodingErrorIf(valueVariableASTNode.Type != asts.NodeType(NodeTypeLocalVariable))
	lib.InternalCodingErrorIf(valueVariableASTNode.Token == nil)
	valueVariableName := tokenLit(valueVariableASTNode)
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
	}
	if blockExitPayload.blockExitStatus == BLOCK_EXIT_BREAK {
		return nil, err
	} else {
		return blockExitPayload, err
	}
}

func (node *ForLoopMultivariableNode) executeOuter(
	mv *mlrval.Mlrval,
	keyIndexVariables []*runtime.StackVariable,
	state *runtime.State,
) (*BlockExitPayload, error) {
	if len(keyIndexVariables) == 1 {
		return node.executeInner(mv, keyIndexVariables[0], state)
	}
	// else, recurse

	if mv.IsMap() {
		mapval := mv.GetMap()

		for pe := mapval.Head; pe != nil; pe = pe.Next {
			mapkey := mlrval.FromString(pe.Key)

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

	} else if mv.IsArray() {
		arrayval := mv.GetArray()

		// Note: Miller user-space array indices ("mindex") are 1-up. Internal
		// Go storage ("zindex") is 0-up.

		for zindex, element := range arrayval {
			mindex := mlrval.FromInt(int64(zindex + 1))

			err := state.Stack.SetAtScope(keyIndexVariables[0], mindex)
			if err != nil {
				return nil, err
			}

			blockExitPayload, err := node.executeOuter(element, keyIndexVariables[1:], state)
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

	} else if mv.IsAbsent() {
		// Data-heterogeneity no-op
	}

	// TODO: backwards compatibility with the C port means we treat this as
	// silent zero-pass. But maybe we should surface it as an error. Maybe
	// with a "mlr put --errors" flag or something.
	//	} else {
	//		return nil, fmt.Errorf(
	//			"mlr: looped-over item is not a map or array; got %s",
	//			mv.GetTypeName(),
	//		)

	return nil, nil
}

func (node *ForLoopMultivariableNode) executeInner(
	mv *mlrval.Mlrval,
	keyIndexVariable *runtime.StackVariable,
	state *runtime.State,
) (*BlockExitPayload, error) {
	if mv.IsMap() {
		mapval := mv.GetMap()

		for pe := mapval.Head; pe != nil; pe = pe.Next {
			mapkey := mlrval.FromString(pe.Key)

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

	} else if mv.IsArray() {
		arrayval := mv.GetArray()

		// Note: Miller user-space array indices ("mindex") are 1-up. Internal
		// Go storage ("zindex") is 0-up.

		for zindex, element := range arrayval {
			mindex := mlrval.FromInt(int64(zindex + 1))

			err := state.Stack.SetAtScope(keyIndexVariable, mindex)
			if err != nil {
				return nil, err
			}
			err = state.Stack.SetAtScope(node.valueIndexVariable, element)
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

	} else if mv.IsAbsent() {
		// Data-heterogeneity no-op
	}

	// TODO: backwards compatibility with the C port means we treat this as
	// silent zero-pass. But maybe we should surface it as an error. Maybe
	// with a "mlr put --errors" flag or something.
	//	} else {
	//		return nil, fmt.Errorf(
	//			"mlr: looped-over item is not a map or array; got %s",
	//			mv.GetTypeName(),
	//		)

	return nil, nil
}

type TripleForLoopNode struct {
	startBlockNode              *StatementBlockNode
	precontinuationAssignments  []IExecutable
	continuationExpressionNode  IEvaluable
	continuationExpressionToken *tokens.Token
	updateBlockNode             *StatementBlockNode
	bodyBlockNode               *StatementBlockNode
}

func NewTripleForLoopNode(
	startBlockNode *StatementBlockNode,
	precontinuationAssignments []IExecutable,
	continuationExpressionNode IEvaluable,
	continuationExpressionToken *tokens.Token,
	updateBlockNode *StatementBlockNode,
	bodyBlockNode *StatementBlockNode,
) *TripleForLoopNode {
	return &TripleForLoopNode{
		startBlockNode,
		precontinuationAssignments,
		continuationExpressionNode,
		continuationExpressionToken,
		updateBlockNode,
		bodyBlockNode,
	}
}

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

func (root *RootNode) BuildTripleForLoopNode(astNode *asts.ASTNode) (*TripleForLoopNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != asts.NodeType(NodeTypeTripleForLoop))
	lib.InternalCodingErrorIf(len(astNode.Children) != 4)

	startBlockASTNode := astNode.Children[0]
	continuationExpressionASTNode := astNode.Children[1]
	updateBlockASTNode := astNode.Children[2]
	bodyBlockASTNode := astNode.Children[3]

	// PGPG: body is StatementBlockInBraces; unwrap to get StatementBlock
	if bodyBlockASTNode.Type == asts.NodeType(NodeTypeStatementBlockInBraces) &&
		len(bodyBlockASTNode.Children) == 1 {
		bodyBlockASTNode = bodyBlockASTNode.Children[0]
	}

	lib.InternalCodingErrorIf(startBlockASTNode.Type != asts.NodeType(NodeTypeStatementBlock))
	lib.InternalCodingErrorIf(continuationExpressionASTNode.Type != asts.NodeType(NodeTypeStatementBlock))
	lib.InternalCodingErrorIf(updateBlockASTNode.Type != asts.NodeType(NodeTypeStatementBlock))
	lib.InternalCodingErrorIf(bodyBlockASTNode.Type != asts.NodeType(NodeTypeStatementBlock))

	startBlockNode, err := root.BuildStatementBlockNode(startBlockASTNode)
	if err != nil {
		return nil, err
	}

	// Enforced here, not in the grammar: the last must be a bare boolean; the ones
	// before must be assignments. Example:
	// for (int i = 0; c += 1, i < 10; i += 1) { ... }
	var precontinuationAssignments []IExecutable = nil
	var continuationExpressionNode IEvaluable = nil
	var continuationExpressionToken *tokens.Token = nil
	if len(continuationExpressionASTNode.Children) > 0 { // empty is true
		n := len(continuationExpressionASTNode.Children)
		if n > 1 {
			precontinuationAssignments = make([]IExecutable, n-1)
			for i := 0; i < n-1; i++ {
				childType := continuationExpressionASTNode.Children[i].Type
				if childType != asts.NodeType(NodeTypeAssignment) &&
					childType != asts.NodeType(NodeTypeCompoundAssignment) {
					return nil, fmt.Errorf(
						"the non-final triple-for continuation statements must be assignments",
					)
				}
				var precontinuationAssignment IExecutable
				var err error
				if childType == asts.NodeType(NodeTypeCompoundAssignment) {
					precontinuationAssignment, err = root.BuildCompoundAssignmentNode(continuationExpressionASTNode.Children[i])
				} else {
					precontinuationAssignment, err = root.BuildAssignmentNode(continuationExpressionASTNode.Children[i])
				}
				if err != nil {
					return nil, err
				}
				precontinuationAssignments[i] = precontinuationAssignment
			}
		}

		bareBooleanASTNode := continuationExpressionASTNode.Children[n-1]
		if bareBooleanASTNode.Type != asts.NodeType(NodeTypeBareBoolean) {
			if n == 1 {
				return nil, fmt.Errorf(
					"the triple-for continuation statement must be a bare boolean",
				)
			} else {
				return nil, fmt.Errorf(
					"the final triple-for continuation statement must be a bare boolean",
				)
			}
		}
		lib.InternalCodingErrorIf(len(bareBooleanASTNode.Children) != 1)
		continuationExpressionNode, err = root.BuildEvaluableNode(bareBooleanASTNode.Children[0])
		continuationExpressionToken = bareBooleanASTNode.Children[0].Token
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
		continuationExpressionToken,
		updateBlockNode,
		bodyBlockNode,
	), nil
}

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
		for _, precontinuationAssignment := range node.precontinuationAssignments {
			_, err := precontinuationAssignment.Execute(state)
			if err != nil {
				return nil, err
			}
		}
		if node.continuationExpressionNode != nil { // empty is true
			continuationValue := node.continuationExpressionNode.Evaluate(state)
			boolValue, isBool := continuationValue.GetBoolValue()
			if !isBool {
				return nil, fmt.Errorf(
					"for-loop continuation did not evaluate to boolean%s",
					pgpgTokenToLocationInfo(node.continuationExpressionToken),
				)
			}
			if !boolValue {
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
