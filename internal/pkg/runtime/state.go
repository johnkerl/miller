// ================================================================
// Tracks everything needed for statement evaluation/assignment in the Miller
// DSL runtimne: current record/context (the latter being NF, NR, etc);
// out-of-stream variables; local-variable stack; etc.
// ================================================================

package runtime

import (
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/types"
)

type State struct {
	Inrec            *types.Mlrmap
	Context          *types.Context
	Oosvars          *types.Mlrmap
	FilterExpression *types.Mlrval
	Stack            *Stack
	OutputChannel    chan<- *types.RecordAndContext
	// For holding "\0".."\9" between where they are set via things like
	// '$x =~ "(..)_(...)"', and interpolated via things like '$y = "\2:\1"'.
	RegexCaptures []string
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

		// See lib.MakeEmptyRegexCaptures for context.
		RegexCaptures: lib.MakeEmptyRegexCaptures(),
	}
}

func (state *State) Update(
	inrec *types.Mlrmap,
	context *types.Context,
) {
	state.Inrec = inrec
	state.Context = context
	state.RegexCaptures = lib.MakeEmptyRegexCaptures()
}
