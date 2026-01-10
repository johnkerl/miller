// ================================================================
// Tracks everything needed for statement evaluation/assignment in the Miller
// DSL runtimne: current record/context (the latter being NF, NR, etc);
// out-of-stream variables; local-variable stack; etc.
// ================================================================

package runtime

import (
	"container/list"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
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
	//
	// Each top-level block and user-defined function has its own captures.
	//
	// For example, in function `f()`, one can do `somevar =~ someregex`, then
	// call some function `g()` which also uses `=~`, and then when `g()` returns,
	// `f()` will have its "\1", "\2", etc intact.
	//
	// This is necessary for the stateful semantics of `=~` and "\1", "\2", etc.
	// Those are avoided when the user calls `matchx`, which is newer, and
	// stateless. However, `=~` exists in the Miller DSL and we must support it.
	regexCapturesByFrame [][]string

	Options *cli.TOptions

	// StrictMode allows for runtime handling of absent-reads and untyped assignments.
	StrictMode bool
}

func NewEmptyState(options *cli.TOptions, strictMode bool) *State {

	// See lib.MakeEmptyCaptures for context.
	regexCapturesByFrame := make([][]string, 1)
	regexCapturesByFrame[0] = lib.MakeEmptyCaptures()

	oosvars := mlrval.NewMlrmap()
	return &State{
		Inrec:                nil,
		Context:              nil,
		Oosvars:              oosvars,
		// XXX
		//FilterExpression:     mlrval.TRUE,
		FilterExpression:     mlrval.NULL,
		Stack:                NewStack(),
		regexCapturesByFrame: regexCapturesByFrame,

		// OutputRecordsAndContexts is assigned after construction

		Options: options,

		StrictMode: strictMode,
	}
}

func (state *State) Update(
	inrec *mlrval.Mlrmap,
	context *types.Context,
) {
	state.Inrec = inrec
	state.Context = context
	state.regexCapturesByFrame[0] = lib.MakeEmptyCaptures()
}

func (state *State) SetRegexCaptures(
	captures []string,
) {
	state.regexCapturesByFrame[0] = lib.CopyStringArray(captures)
}

func (state *State) GetRegexCaptures() []string {
	regexCaptures := state.regexCapturesByFrame[0]
	return lib.CopyStringArray(regexCaptures)
}

func (state *State) PushRegexCapturesFrame() {
	state.regexCapturesByFrame = append([][]string{lib.MakeEmptyCaptures()}, state.regexCapturesByFrame...)
}

func (state *State) PopRegexCapturesFrame() {
	state.regexCapturesByFrame = state.regexCapturesByFrame[1:]
}
