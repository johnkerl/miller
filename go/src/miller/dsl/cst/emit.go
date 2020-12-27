// ================================================================
// This handles emit statements.
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
type EmitStatementNode struct {
	emitEvaluable IEvaluable

	// emitEvaluables []IEvaluable
	// keyEvaluables []IEvaluable
}

// ----------------------------------------------------------------
func (this *RootNode) BuildEmitStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeEmitStatement)
	nchild := len(astNode.Children)
	lib.InternalCodingErrorIf(nchild != 1 && nchild != 2)

	emitEvaluable, err := this.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	return &EmitStatementNode{
		emitEvaluable: emitEvaluable,
	}, nil
}

func (this *EmitStatementNode) Execute(state *State) (*BlockExitPayload, error) {
	emitResult := this.emitEvaluable.Evaluate(state)

	if emitResult.IsAbsent() {
		return nil, nil
	}

	if emitResult.IsMap() {
		state.OutputChannel <- types.NewRecordAndContext(
			emitResult.Copy().GetMap(),
			state.Context, // xxx clone ?
		)
	}

	return nil, nil
}

// cases:
// * 'emit (@count, @sum)' -- convert to mlrmap "count=1,sum=2"
// * 'emit (@count, @sum), "a"' -- convert to mlrmap "a=foo,count=2,sum=3.4'
// ?? maybe alter from mlr-c syntax to require a map here -- ?
// * 'emit {"a": @a, "b": @b}' -- ?
// * 'for k in @u { emit {"a": k, "u": @u[k], "v": @v[k] }' -- ?

// possibles:
// * maps -- as-is
//   o what about nameless bases such as @* and $*?
// * srecs -- key-value pairs into a new map
// * oosvars -- key-value pairs into a new map
// * localvars -- key-value pairs into a new map
// * otherwise error

// * Given @count = 2 and @sum = 3.4:
//   o 'emit (@sum, @count)' => [{ "sum": 2, "count": 3.4 }]

// * Given @count = {"pan": 2, "eks": 3} and @sum = {"pan" 3.4, "eks": 5.6 }:
//   o 'emit (@sum, @count)' => [{
//       "count": {"pan": 2, "eks": 3},
//       "sum": {"pan" 3.4, "eks": 5.6 }
//     }]

// * Given @count = {"pan": 2, "eks": 3} and @sum = {"pan" 3.4, "eks": 5.6 }:
//   o 'emit (@sum, @count), $a' =>
//     [
//       {
//         "a": "pan",
//         "count": 2,
//         "sum": 3.4
//       },
//       {
//         "a": "eks",
//         "count": 3,
//         "sum": 5.6
//       }
//     ]

// ================================================================
type EmitPStatementNode struct {
	emitpEvaluable IEvaluable
	// xxx to do:
	// * required array of evaluables
	// * optional array of indexing keys
}

// ----------------------------------------------------------------
func (this *RootNode) BuildEmitPStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeEmitPStatement)
	lib.InternalCodingErrorIf(len(astNode.Children) < 1)

	emitpEvaluable, err := this.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	return &EmitPStatementNode{
		emitpEvaluable: emitpEvaluable,
	}, nil
}

func (this *EmitPStatementNode) Execute(state *State) (*BlockExitPayload, error) {
	emitpResult := this.emitpEvaluable.Evaluate(state)

	if emitpResult.IsAbsent() {
		return nil, nil
	}

	if emitpResult.IsMap() {
		state.OutputChannel <- types.NewRecordAndContext(
			emitpResult.Copy().GetMap(),
			state.Context, // xxx clone ?
		)
	}

	// xxx WIP
	// xxx need to reshape rvalue mlrvals -> mlrmaps; publish w/ contexts; method for that

	//	outputChannel <- types.NewRecordAndContext(
	//		mlrmap goes here,
	//		&context,
	//	)

	return nil, nil
}

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
	n := len(astNode.Children)
	lib.InternalCodingErrorIf(n < 1)

	emitfNames := make([]string, n)
	emitfEvaluables := make([]IEvaluable, n)
	for i, childNode := range astNode.Children {
		emitfName, err := getNameFromNamedNode(childNode)
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

// ================================================================
// HELPER FUNCTIONS

// For emitf: gets the name of a non-indexed oosvar, localvar, or field name;
// otherwise, returns error.
//
// TODO: support indirects like 'emitf @[x."_sum"]'

func getNameFromNamedNode(astNode *dsl.ASTNode) (string, error) {
	if astNode.Type == dsl.NodeTypeDirectOosvarValue {
		return string(astNode.Token.Lit), nil
	} else if astNode.Type == dsl.NodeTypeLocalVariable {
		return string(astNode.Token.Lit), nil
	} else if astNode.Type == dsl.NodeTypeDirectFieldValue {
		return string(astNode.Token.Lit), nil
	}
	return "", errors.New(
		fmt.Sprintf(
			"%s: can't get name of node type \"%s\" for emitf.",
			os.Args[0], string(astNode.Type),
		),
	)
}
