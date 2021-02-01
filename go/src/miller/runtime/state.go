// ================================================================
// Tracks everything needed for statement evaluation/assignment in the Miller
// DSL runtimne: current record/context (the latter being NF, NR, etc);
// out-of-stream variables; local-variable stack; etc.
// ================================================================
package runtime

import (
	"miller/types"
)

type State struct {
	Inrec         *types.Mlrmap
	Context       *types.Context
	Oosvars       *types.Mlrmap
	FilterResult  bool
	OutputChannel chan<- *types.RecordAndContext
	Stack         *Stack
}

func NewEmptyState() *State {
	oosvars := types.NewMlrmap()
	return &State{
		Inrec:        nil,
		Context:      nil,
		Oosvars:      oosvars,
		FilterResult: true,
		// OutputChannel is assigned after construction
		Stack: NewStack(),
	}
}

func (this *State) Update(
	inrec *types.Mlrmap,
	context *types.Context,
) {
	this.Inrec = inrec
	this.Context = context
}
