// ================================================================
// This handles emit statements.
// ================================================================

package cst

import (
	"miller/dsl"
	"miller/lib"
	"miller/types"
)

// ----------------------------------------------------------------
type EmitStatementNode struct {
	emitEvaluable IEvaluable
	// xxx to do:
	// * required array of evaluables
	// * optional array of indexing keys
}

// ----------------------------------------------------------------
func (this *RootNode) BuildEmitStatementNode(astNode *dsl.ASTNode) (IExecutable, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeEmitStatement)
	lib.InternalCodingErrorIf(len(astNode.Children) < 1)

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

	// xxx WIP
	// xxx need to reshape rvalue mlrvals -> mlrmaps; publish w/ contexts; method for that

	//	outputChannel <- types.NewRecordAndContext(
	//		mlrmap goes here,
	//		&context,
	//	)

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
