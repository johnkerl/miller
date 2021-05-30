// ================================================================
// Tracks everything needed for statement evaluation/assignment in the Miller
// DSL runtimne: current record/context (the latter being NF, NR, etc);
// out-of-stream variables; local-variable stack; etc.
// ================================================================
package runtime

import (
	"miller/src/types"
)

type State struct {
	Inrec            *types.Mlrmap
	Context          *types.Context
	Oosvars          *types.Mlrmap
	FilterExpression *types.Mlrval
	Stack            *Stack
	OutputChannel    chan<- *types.RecordAndContext
}

func NewEmptyState() *State {
	oosvars := types.NewMlrmap()
	return &State{
		Inrec:            nil,
		Context:          nil,
		Oosvars:          oosvars,
		FilterExpression: types.MLRVAL_TRUE,
		Stack:            NewStack(),
		// OutputChannel is assigned after construction
	}
}

func (state *State) Update(
	inrec *types.Mlrmap,
	context *types.Context,
) {
	state.Inrec = inrec
	state.Context = context
}
