// ================================================================
// Tracks everything needed for statement evaluation/assignment in the Miller
// DSL runtimne: current record/context (the latter being NF, NR, etc);
// out-of-stream variables; local-variable stack; etc.
// ================================================================

package cst

import (
	"miller/types"
)

func NewEmptyState() *State {
	oosvars := types.NewMlrmap()
	return &State{
		Inrec:        nil,
		Context:      nil,
		Oosvars:      oosvars,
		FilterResult: true,
		// OutputChannel is assigned after construction
		stack: NewStack(),
	}
}

func (this *State) Update(
	inrec *types.Mlrmap,
	context *types.Context,
) {
	this.Inrec = inrec
	this.Context = context
}
