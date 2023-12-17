// ================================================================
// Tracks everything needed for statement evaluation/assignment in the Miller
// DSL runtimne: current record/context (the latter being NF, NR, etc);
// out-of-stream variables; local-variable stack; etc.
// ================================================================

package runtime

import (
	"container/list"

	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/lib"
	"github.com/johnkerl/miller/pkg/mlrval"
	"github.com/johnkerl/miller/pkg/types"
)

type State struct {
	Inrec                    *mlrval.Mlrmap
	Context                  *types.Context
	Oosvars                  *mlrval.Mlrmap
	FilterExpression         *mlrval.Mlrval
	Stack                    *Stack
	OutputRecordsAndContexts *list.List // list of *types.RecordAndContext

	// For holding "\0".."\9" between where they are set via things like
	// '$x =~ "(..)_(...)"', and interpolated via things like '$y = "\2:\1"'.
	RegexCaptures []string
	Options       *cli.TOptions

	// StrictMode allows for runtime handling of absent-reads and untyped assignments.
	StrictMode bool
}

func NewEmptyState(options *cli.TOptions, strictMode bool) *State {
	oosvars := mlrval.NewMlrmap()
	return &State{
		Inrec:            nil,
		Context:          nil,
		Oosvars:          oosvars,
		FilterExpression: mlrval.TRUE,
		Stack:            NewStack(),

		// OutputRecordsAndContexts is assigned after construction

		// See lib.MakeEmptyCaptures for context.
		RegexCaptures: lib.MakeEmptyCaptures(),
		Options:       options,

		StrictMode: strictMode,
	}
}

func (state *State) Update(
	inrec *mlrval.Mlrmap,
	context *types.Context,
) {
	state.Inrec = inrec
	state.Context = context
	state.RegexCaptures = lib.MakeEmptyCaptures()
}
