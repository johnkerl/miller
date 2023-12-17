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
	regexCapturesByFrame *list.List // list of []string

	Options *cli.TOptions

	// StrictMode allows for runtime handling of absent-reads and untyped assignments.
	StrictMode bool
}

func NewEmptyState(options *cli.TOptions, strictMode bool) *State {

	// See lib.MakeEmptyRegexCaptures for context.
	regexCapturesByFrame := list.New()
	regexCapturesByFrame.PushFront(lib.MakeEmptyRegexCaptures())

	oosvars := mlrval.NewMlrmap()
	return &State{
		Inrec:                nil,
		Context:              nil,
		Oosvars:              oosvars,
		FilterExpression:     mlrval.TRUE,
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
	state.regexCapturesByFrame.Front().Value = lib.MakeEmptyRegexCaptures()
}

func (state *State) SetRegexCaptures(
	captures []string,
) {
	state.regexCapturesByFrame.Front().Value = lib.CopyStringArray(captures)
}

func (state *State) GetRegexCaptures() []string {
	regexCaptures := state.regexCapturesByFrame.Front().Value.([]string)
	return lib.CopyStringArray(regexCaptures)
}

func (state *State) PushRegexCapturesFrame() {
	state.regexCapturesByFrame.PushFront(lib.MakeEmptyRegexCaptures())
}

func (state *State) PopRegexCapturesFrame() {
	// There is no PopFront
	state.regexCapturesByFrame.Remove(state.regexCapturesByFrame.Front())
}
