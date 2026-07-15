package transformers

import (
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/types"
)

const verbNameDecimate = "decimate"

var decimateOptions = []OptionSpec{
	{Flag: "-b", Type: "bool", Desc: "Decimate by printing first of every n."},
	{Flag: "-e", Type: "bool", Desc: "Decimate by printing last of every n (default)."},
	{Flag: "-g", Arg: "{a,b,c}", Type: "csv-list", Desc: "Optional group-by-field names for decimate counts, e.g. a,b,c."},
	{Flag: "-n", Arg: "{n}", Type: "int", Desc: "Decimation factor (default 10)."},
}

var DecimateSetup = TransformerSetup{
	Verb:         verbNameDecimate,
	UsageFunc:    transformerDecimateUsage,
	ParseCLIFunc: transformerDecimateParseCLI,
	IgnoresInput: false,
	Options:      decimateOptions,
}

func transformerDecimateUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameDecimate)
	fmt.Fprintf(o, "Passes through one of every n records, optionally by category.\n")
	WriteVerbOptions(o, decimateOptions)
}

func transformerDecimateParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) (RecordTransformer, error) {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	var err error
	decimateCount := int64(10)
	atStart := false
	atEnd := false
	var groupByFieldNames []string = nil

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		if args[argi] == "--" {
			break // All transformers must do this so main-flags can follow verb-flags
		}
		argi++

		switch opt {
		case "-h", "--help":
			transformerDecimateUsage(os.Stdout)
			return nil, cli.ErrHelpRequested

		case "-n":
			decimateCount, err = cli.VerbGetIntArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			if decimateCount <= 0 {
				return nil, cli.VerbErrorf(verb, "-n must be positive")
			}

		case "-b":
			atStart = true

		case "-e":
			atEnd = true

		case "-g":
			groupByFieldNames, err = cli.VerbGetStringArrayArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		default:
			return nil, cli.VerbErrorf(verb, "option \"%s\" not recognized", opt)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	transformer, err := NewTransformerDecimate(
		decimateCount,
		atStart,
		atEnd,
		groupByFieldNames,
	)
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

type TransformerDecimate struct {
	decimateCount     int64
	remainderToKeep   int64
	groupByFieldNames []string

	countsByGroup map[string]int64
}

func NewTransformerDecimate(
	decimateCount int64,
	atStart bool,
	atEnd bool,
	groupByFieldNames []string,
) (*TransformerDecimate, error) {

	remainderToKeep := decimateCount - 1
	if atStart && !atEnd {
		remainderToKeep = 0
	}

	tr := &TransformerDecimate{
		decimateCount:     decimateCount,
		remainderToKeep:   remainderToKeep,
		groupByFieldNames: groupByFieldNames,
		countsByGroup:     make(map[string]int64),
	}

	return tr, nil
}

func (tr *TransformerDecimate) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) error {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		groupingKey, ok := inrec.GetSelectedValuesJoined(tr.groupByFieldNames)
		if !ok {
			return nil // This particular record doesn't have the specified fields; ignore
		}

		countForGroup, ok := tr.countsByGroup[groupingKey]
		if !ok {
			countForGroup = 0
			tr.countsByGroup[groupingKey] = countForGroup
		}

		remainder := countForGroup % tr.decimateCount
		if remainder == tr.remainderToKeep {
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
		}

		countForGroup++
		tr.countsByGroup[groupingKey] = countForGroup

	} else {
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // Emit the stream-terminating null record
	}
	return nil
}
