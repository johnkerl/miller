package transformers

import (
	"fmt"
	"os"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameDecimate = "decimate"

var DecimateSetup = TransformerSetup{
	Verb:         verbNameDecimate,
	UsageFunc:    transformerDecimateUsage,
	ParseCLIFunc: transformerDecimateParseCLI,
	IgnoresInput: false,
}

func transformerDecimateUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameDecimate)
	fmt.Fprintf(o, "Passes through one of every n records, optionally by category.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, " -b Decimate by printing first of every n.\n")
	fmt.Fprintf(o, " -e Decimate by printing last of every n (default).\n")
	fmt.Fprintf(o, " -g {a,b,c} Optional group-by-field names for decimate counts, e.g. a,b,c.\n")
	fmt.Fprintf(o, " -n {n} Decimation factor (default 10).\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerDecimateParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	decimateCount := 10
	atStart := false
	atEnd := false
	var groupByFieldNames []string = nil

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerDecimateUsage(os.Stdout, true, 0)

		} else if opt == "-n" {
			decimateCount = cli.VerbGetIntArgOrDie(verb, opt, args, &argi, argc)
			if decimateCount <= 0 {
				transformerDecimateUsage(os.Stderr, true, 1)
			}

		} else if opt == "-b" {
			atStart = true

		} else if opt == "-e" {
			atEnd = true

		} else if opt == "-g" {
			groupByFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerDecimateUsage(os.Stderr, true, 1)
		}
	}

	transformer, err := NewTransformerDecimate(
		decimateCount,
		atStart,
		atEnd,
		groupByFieldNames,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerDecimate struct {
	decimateCount     int
	remainderToKeep   int
	groupByFieldNames []string

	countsByGroup map[string]int
}

// ----------------------------------------------------------------
func NewTransformerDecimate(
	decimateCount int,
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
		countsByGroup:     make(map[string]int),
	}

	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerDecimate) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		groupingKey, ok := inrec.GetSelectedValuesJoined(tr.groupByFieldNames)
		if !ok {
			return // This particular record doesn't have the specified fields; ignore
		}

		countForGroup, ok := tr.countsByGroup[groupingKey]
		if !ok {
			countForGroup = 0
			tr.countsByGroup[groupingKey] = countForGroup
		}

		remainder := countForGroup % tr.decimateCount
		if remainder == tr.remainderToKeep {
			outputChannel <- inrecAndContext
		}

		countForGroup++
		tr.countsByGroup[groupingKey] = countForGroup

	} else {
		outputChannel <- inrecAndContext // Emit the stream-terminating null record
	}
}
