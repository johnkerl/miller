// ================================================================
// This handles emitf statements. This produces new records (in addition to $*)
// into the output record stream.
// ================================================================

package cst

import (
	"errors"
	"fmt"
	"os"

	"miller/dsl"
	"miller/lib"
	"miller/types"
)

// ================================================================
// Examples:
//   emitf @a
//   emitf @a, @b
//
// Each argument must be a non-indexed oosvar/localvar/fieldname, so we can use
// their names as keys in the emitted record.  These restrictions are enforced
// in the CST logic, to keep this parser/AST logic simpler.

type EmitFStatementNode struct {
	emitfNames      []string
	emitfEvaluables []IEvaluable
	// xxx redirect
}

// ----------------------------------------------------------------
// $ mlr -n put -v 'emitf a,$b,@c'
// DSL EXPRESSION:
// emitf a,$b,@c
// RAW AST:
// * statement block
//     * dump statement "emitf"
//         * local variable "a"
//         * direct field value "b"
//         * direct oosvar value "c"

func (this *RootNode) BuildEmitFStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeEmitFStatement)
	lib.InternalCodingErrorIf(len(astNode.Children) != 2)
	expressionsNode := astNode.Children[0]

	n := len(expressionsNode.Children)
	lib.InternalCodingErrorIf(n < 1)

	emitfNames := make([]string, n)
	emitfEvaluables := make([]IEvaluable, n)
	for i, childNode := range expressionsNode.Children {
		emitfName, err := getNameFromNamedNode(childNode, "emitf")
		if err != nil {
			return nil, err
		}
		emitfEvaluable, err := this.BuildEvaluableNode(childNode)
		if err != nil {
			return nil, err
		}
		emitfNames[i] = emitfName
		emitfEvaluables[i] = emitfEvaluable
	}
	return &EmitFStatementNode{
		emitfNames:      emitfNames,
		emitfEvaluables: emitfEvaluables,
	}, nil
}

func (this *EmitFStatementNode) Execute(state *State) (*BlockExitPayload, error) {
	newrec := types.NewMlrmapAsRecord()
	for i, emitfEvaluable := range this.emitfEvaluables {
		emitfName := this.emitfNames[i]
		emitfValue := emitfEvaluable.Evaluate(state)

		if !emitfValue.IsAbsent() {
			newrec.PutCopy(&emitfName, &emitfValue)
		}
	}
	state.OutputChannel <- types.NewRecordAndContext(
		newrec,
		state.Context.Copy(),
	)

	return nil, nil
}

// Gets the name of a non-indexed oosvar, localvar, or field name; otherwise,
// returns error.
//
// TODO: support indirects like 'emitf @[x."_sum"]'

func getNameFromNamedNode(astNode *dsl.ASTNode, description string) (string, error) {
	if astNode.Type == dsl.NodeTypeDirectOosvarValue {
		return string(astNode.Token.Lit), nil
	} else if astNode.Type == dsl.NodeTypeLocalVariable {
		return string(astNode.Token.Lit), nil
	} else if astNode.Type == dsl.NodeTypeDirectFieldValue {
		return string(astNode.Token.Lit), nil
	}
	return "", errors.New(
		fmt.Sprintf(
			"%s: can't get name of node type \"%s\" for %s.",
			os.Args[0], string(astNode.Type), description,
		),
	)
}
