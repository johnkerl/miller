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
