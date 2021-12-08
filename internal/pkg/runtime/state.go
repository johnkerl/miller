// ================================================================
// Tracks everything needed for statement evaluation/assignment in the Miller
// DSL runtimne: current record/context (the latter being NF, NR, etc);
// out-of-stream variables; local-variable stack; etc.
// ================================================================

package runtime

import (
	"container/list"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/types"
)

type State struct {
	Inrec                    *types.Mlrmap
	Context                  *types.Context
	Oosvars                  *types.Mlrmap
	FilterExpression         *types.Mlrval
	Stack                    *Stack
	OutputRecordsAndContexts *list.List // list of *types.RecordAndContext
	// For holding "\0".."\9" between where they are set via things like
	// '$x =~ "(..)_(...)"', and interpolated via things like '$y = "\2:\1"'.
	RegexCaptures []string
	Options       *cli.TOptions
}

func NewEmptyState(options *cli.TOptions) *State {
	oosvars := types.NewMlrmap()
	return &State{
		Inrec:            nil,
		Context:          nil,
		Oosvars:          oosvars,
		FilterExpression: types.MLRVAL_TRUE,
		Stack:            NewStack(),

		// OutputRecordsAndContexts is assigned after construction

		// See lib.MakeEmptyRegexCaptures for context.
		RegexCaptures: lib.MakeEmptyRegexCaptures(),
		Options:       options,
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
