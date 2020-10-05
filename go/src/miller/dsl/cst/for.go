package cst

import (
	"errors"

	"miller/dsl"
	"miller/lib"
	"miller/types"
)

// ================================================================
// This is for various flavors of for-loop
// ================================================================

// ----------------------------------------------------------------
// Sample AST:

// mlr -n put -v 'for (k in $*) { emit { k : k } }'
// DSL EXPRESSION:
// for (k in $*) { emit { k : k} }
// RAW AST:
// * StatementBlock
//     * ForLoopKeyOnly "for"
//         * LocalVariable "k"
//         * FullSrec "$*"
//         * StatementBlock
//             * EmitStatement "emit"
//                 * MapLiteral "{}"
//                     * MapLiteralKeyValuePair ":"
//                         * LocalVariable "k"
//                         * LocalVariable "k"

func BuildForLoopKeyOnlyNode(astNode *dsl.ASTNode) (*ForLoopKeyValueNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeForLoopKeyOnly)
	lib.InternalCodingErrorIf(len(astNode.Children) != 3)

	keyVariableASTNode := astNode.Children[0]
	indexableASTNode := astNode.Children[1]
	blockASTNode := astNode.Children[2]

	lib.InternalCodingErrorIf(keyVariableASTNode.Type != dsl.NodeTypeLocalVariable)
	lib.InternalCodingErrorIf(keyVariableASTNode.Token == nil)
	keyVariableName := string(keyVariableASTNode.Token.Lit)

	// TODO: error if loop-over node isn't Mappable (inasmuch as can be
	// detected at CST-build time)
	// TODO: support arrays too.
	indexableNode, err := BuildEvaluableNode(indexableASTNode)
	if err != nil {
		return nil, err
	}

	statementBlockNode, err := BuildStatementBlockNode(blockASTNode)
	if err != nil {
		return nil, err
	}

	return NewForLoopKeyValueNode(
		keyVariableName,
		"", // No value-variable here, but re-use the code
		indexableNode,
		statementBlockNode,
	), nil
}

// ================================================================
type ForLoopKeyValueNode struct {
	keyVariableName    string
	valueVariableName  string
	indexableNode      IEvaluable
	statementBlockNode *StatementBlockNode
}

func NewForLoopKeyValueNode(
	keyVariableName string,
	valueVariableName string,
	indexableNode IEvaluable,
	statementBlockNode *StatementBlockNode,
) *ForLoopKeyValueNode {
	return &ForLoopKeyValueNode{
		keyVariableName,
		valueVariableName,
		indexableNode,
		statementBlockNode,
	}
}

// ----------------------------------------------------------------
// Sample AST:

// mlr -n put -v 'for (k, v in $*) { emit { k : v } }'
// DSL EXPRESSION:
// for (k, v in $*) { emit { k : v } }
// RAW AST:
// * StatementBlock
//     * ForLoopKeyValue "for"
//         * LocalVariable "k"
//         * LocalVariable "v"
//         * FullSrec "$*"
//         * StatementBlock
//             * EmitStatement "emit"
//                 * MapLiteral "{}"
//                     * MapLiteralKeyValuePair ":"
//                         * LocalVariable "k"
//                         * LocalVariable "v"

func BuildForLoopKeyValueNode(astNode *dsl.ASTNode) (*ForLoopKeyValueNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeForLoopKeyValue)
	lib.InternalCodingErrorIf(len(astNode.Children) != 4)

	keyVariableASTNode := astNode.Children[0]
	valueVariableASTNode := astNode.Children[1]
	indexableASTNode := astNode.Children[2]
	blockASTNode := astNode.Children[3]

	lib.InternalCodingErrorIf(keyVariableASTNode.Type != dsl.NodeTypeLocalVariable)
	lib.InternalCodingErrorIf(keyVariableASTNode.Token == nil)
	keyVariableName := string(keyVariableASTNode.Token.Lit)

	lib.InternalCodingErrorIf(valueVariableASTNode.Type != dsl.NodeTypeLocalVariable)
	lib.InternalCodingErrorIf(valueVariableASTNode.Token == nil)
	valueVariableName := string(valueVariableASTNode.Token.Lit)

	// TODO: error if loop-over node isn't Mappable (inasmuch as can be
	// detected at CST-build time)
	indexableNode, err := BuildEvaluableNode(indexableASTNode)
	if err != nil {
		return nil, err
	}

	statementBlockNode, err := BuildStatementBlockNode(blockASTNode)
	if err != nil {
		return nil, err
	}

	return NewForLoopKeyValueNode(
		keyVariableName,
		valueVariableName,
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

func (this *ForLoopKeyValueNode) Execute(state *State) error {
	mlrval := this.indexableNode.Evaluate(state)

	if mlrval.IsMap() {

		mapval := mlrval.GetMap()

		state.stack.PushStackFrame()
		defer state.stack.PopStackFrame()
		for pe := mapval.Head; pe != nil; pe = pe.Next {
			mapkey := types.MlrvalFromString(*pe.Key)

			state.stack.BindVariable(this.keyVariableName, &mapkey)
			if this.valueVariableName != "" { // 'for (k in ...)' not 'for (k,v in ...)'
				state.stack.BindVariable(this.valueVariableName, pe.Value)
			}
			//state.stack.Dump()
			err := this.statementBlockNode.Execute(state)
			if err != nil {
				return err
			}
		}

	} else if mlrval.IsArray() {

		arrayval := mlrval.GetArray()

		// Miller user-space array indices ("mindex") are 1-up. Internal Go
		// storage ("zindex") is 0-up.
		state.stack.PushStackFrame()
		defer state.stack.PopStackFrame()
		for zindex, element := range arrayval {
			mindex := types.MlrvalFromInt64(int64(zindex + 1))

			state.stack.BindVariable(this.keyVariableName, &mindex)
			if this.valueVariableName != "" { // 'for (k in ...)' not 'for (k,v in ...)'
				state.stack.BindVariable(this.valueVariableName, &element)
			}
			// state.stack.Dump()
			err := this.statementBlockNode.Execute(state)
			if err != nil {
				return err
			}
		}

	} else {
		// TODO: give more context
		return errors.New("Miller: looped-over item is not a map.")
	}

	return nil
}

// ================================================================
type TripleForLoopNode struct {
	startBlockNode             *StatementBlockNode
	continuationExpressionNode IEvaluable
	updateBlockNode            *StatementBlockNode
	bodyBlockNode              *StatementBlockNode
}

func NewTripleForLoopNode(
	startBlockNode *StatementBlockNode,
	continuationExpressionNode IEvaluable,
	updateBlockNode *StatementBlockNode,
	bodyBlockNode *StatementBlockNode,
) *TripleForLoopNode {
	return &TripleForLoopNode{
		startBlockNode,
		continuationExpressionNode,
		updateBlockNode,
		bodyBlockNode,
	}
}

// ----------------------------------------------------------------
// Sample ASTs:

// DSL EXPRESSION:
// for (;;) {}
// RAW AST:
// * StatementBlock
//     * TripleForLoop "for"
//         * StatementBlock
//         * StatementBlock
//         * StatementBlock
//         * StatementBlock

// mlr --from u/s.dkvp put -v for (i = 0; i < NR; i += 1) { $i += i }
// DSL EXPRESSION:
// for (i = 0; i < NR; i += 1) { $i += i }
// RAW AST:
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

func BuildTripleForLoopNode(astNode *dsl.ASTNode) (*TripleForLoopNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeTripleForLoop)
	lib.InternalCodingErrorIf(len(astNode.Children) != 4)

	startBlockASTNode := astNode.Children[0]
	continuationExpressionASTNode := astNode.Children[1]
	updateBlockASTNode := astNode.Children[2]
	bodyBlockASTNode := astNode.Children[3]

	lib.InternalCodingErrorIf(startBlockASTNode.Type != dsl.NodeTypeStatementBlock)
	lib.InternalCodingErrorIf(continuationExpressionASTNode.Type != dsl.NodeTypeStatementBlock)
	lib.InternalCodingErrorIf(len(continuationExpressionASTNode.Children) > 1)
	lib.InternalCodingErrorIf(updateBlockASTNode.Type != dsl.NodeTypeStatementBlock)
	lib.InternalCodingErrorIf(bodyBlockASTNode.Type != dsl.NodeTypeStatementBlock)

	startBlockNode, err := BuildStatementBlockNode(startBlockASTNode)
	if err != nil {
		return nil, err
	}

	var continuationExpressionNode IEvaluable = nil
	// empty is true
	if len(continuationExpressionASTNode.Children) == 1 {
		bareBooleanASTNode := continuationExpressionASTNode.Children[0]
		lib.InternalCodingErrorIf(bareBooleanASTNode.Type != dsl.NodeTypeBareBoolean)
		lib.InternalCodingErrorIf(len(bareBooleanASTNode.Children) != 1)

		continuationExpressionNode, err = BuildEvaluableNode(bareBooleanASTNode.Children[0])
		if err != nil {
			return nil, err
		}
	}

	updateBlockNode, err := BuildStatementBlockNode(updateBlockASTNode)
	if err != nil {
		return nil, err
	}

	bodyBlockNode, err := BuildStatementBlockNode(bodyBlockASTNode)
	if err != nil {
		return nil, err
	}

	return NewTripleForLoopNode(
		startBlockNode,
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
//     x = 1;           <--- frame #1 main
//     for (k in $*) {  <--- frame #2 for for-loop bindvars (right here)
//       x = 2          <--- frame #3 for for-loop locals
//     }
//     x = 3;           <--- back in frame #1 main
//   '
//

func (this *TripleForLoopNode) Execute(state *State) error {

	// TODO: new stack frames

	err := this.startBlockNode.Execute(state)
	if err != nil {
		return err
	}

	for {
		// empty is true
		if this.continuationExpressionNode != nil {
			continuationValue := this.continuationExpressionNode.Evaluate(state)
			boolValue, isBool := continuationValue.GetBoolValue()
			if !isBool {
				// TODO: propagate line-number context
				return errors.New("Miller: for-loop continuation did not evaluate to boolean.")
			}
			if boolValue == false {
				break
			}
		}

		err = this.bodyBlockNode.Execute(state)
		if err != nil {
			return err
		}

		err = this.updateBlockNode.Execute(state)
		if err != nil {
			return err
		}
	}

	return nil
}
