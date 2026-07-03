package transformers

import (
	"os"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/types"
)

// RecordTransformer is the interface satisfied by all transformers, i.e.,
// Miller verbs. See stream.go for context on the channels, as well as for
// context on end-of-record-stream signaling.
type RecordTransformer interface {
	Transform(
		inrecAndContext *types.RecordAndContext,
		outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
		inputDownstreamDoneChannel <-chan bool,
		outputDownstreamDoneChannel chan<- bool,
	)
}

type RecordTransformerFunc func(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
)

// Used within some verbs
type RecordTransformerHelperFunc func(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
)

type TransformerUsageFunc func(
	ostream *os.File,
)

// TransformerParseCLIFunc parses verb options from the CLI. Returns (nil, nil) on
// success for pass one (doConstruct=false). Returns (transformer, nil) on success
// for pass two (doConstruct=true). Returns (nil, cli.ErrHelpRequested) when -h/--help
// was used (caller should exit 0; usage already printed). Returns (nil, err) on
// parse or constructor failure (caller should print and exit 1).
type TransformerParseCLIFunc func(
	pargi *int,
	argc int,
	args []string,
	mainOptions *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) (RecordTransformer, error)

// OptionSpec describes one verb option in a machine-readable form. It is the
// building block for the structured Tier-2 catalog (PR3 of the AI-friendly
// roadmap). Verbs migrate to this incrementally; verbs with Options == nil
// fall back to the prose UsageText in the JSON catalog.
//
// Type is one of: bool, string, int, float, csv-list, regex, filename,
// format, enum.  For Type=="enum", Values contains every valid choice.
// Arg is the placeholder shown in usage text, e.g. "{n}" or "{a,b,c}";
// it is empty for bool flags.  Repeatable marks flags that may appear
// multiple times on the command line.  Aliases lists alternate spellings
// of Flag (e.g. "--sorted-input" for "-s").
type OptionSpec struct {
	Flag       string   `json:"flag"`
	Aliases    []string `json:"aliases,omitempty"`
	Arg        string   `json:"arg,omitempty"`
	Type       string   `json:"type"`
	Desc       string   `json:"desc"`
	Repeatable bool     `json:"repeatable,omitempty"`
	Values     []string `json:"values,omitempty"`
}

type TransformerSetup struct {
	Verb         string
	UsageFunc    TransformerUsageFunc
	ParseCLIFunc TransformerParseCLIFunc

	// For seqgen only -- all other transformers process records sourced by the
	// record-reader.  The seqgen verb, by contrast, is a record-source of its
	// own. (The seqgen verb probably should have been designed as a zero-file
	// record "reader" object, rather than a verb, alas.)
	IgnoresInput bool

	// Options is the structured option list for this verb. nil means the verb
	// has not yet been migrated to Tier-2; agents fall back to UsageText.
	// An explicitly empty slice means the verb accepts no options.
	Options []OptionSpec
}

// HandleDefaultDownstreamDone is a utility function for most verbs other than
// head, tee, and seqgen to use for passing downstream-done flags back
// upstream.  See ChainTransformer for context.
func HandleDefaultDownstreamDone(
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	select {
	case b := <-inputDownstreamDoneChannel:
		outputDownstreamDoneChannel <- b
	default:
	}
}
