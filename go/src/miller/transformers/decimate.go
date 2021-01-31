package transformers

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"miller/clitypes"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
const verbNameDecimate = "decimate"

var DecimateSetup = transforming.TransformerSetup{
	Verb:         verbNameDecimate,
	ParseCLIFunc: transformerDecimateParseCLI,
	UsageFunc:    transformerDecimateUsage,
	IgnoresInput: false,
}

func transformerDecimateParseCLI(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) transforming.IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	decimateCount := 10
	atStart := false
	atEnd := false
	var groupByFieldNames []string = nil

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		if !strings.HasPrefix(args[argi], "-") {
			break // No more flag options to process

		} else if args[argi] == "-h" || args[argi] == "--help" {
			transformerDecimateUsage(os.Stdout, true, 0)
			return nil // help intentionally requested

		} else if args[argi] == "-n" {
			decimateCount = clitypes.VerbGetIntArgOrDie(verb, args, &argi, argc)
			if decimateCount <= 0 {
				transformerDecimateUsage(os.Stderr, true, 1)
			}

		} else if args[argi] == "-b" {
			atStart = true
			argi++

		} else if args[argi] == "-e" {
			atEnd = true
			argi++

		} else if args[argi] == "-g" {
			groupByFieldNames = clitypes.VerbGetStringArrayArgOrDie(verb, args, &argi, argc)

		} else {
			transformerDecimateUsage(os.Stderr, true, 1)
			os.Exit(1)
		}
	}

	transformer, _ := NewTransformerDecimate(
		decimateCount,
		atStart,
		atEnd,
		groupByFieldNames,
	)

	*pargi = argi
	return transformer
}

func transformerDecimateUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", os.Args[0], verbNameDecimate)
	fmt.Fprintf(o, "Passes through one of every n records, optionally by category.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, " -b Decimate by printing first of every n.\n")
	fmt.Fprintf(o, " -e Decimate by printing last of every n (default).\n")
	fmt.Fprintf(o, " -g {a,b,c} Optional group-by-field names for decimate counts, e.g. a,b,c.\n")
	fmt.Fprintf(o, " -n {n} Decimation factor (default 10).\n")

	if doExit {
		os.Exit(exitCode)
	}
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

	this := &TransformerDecimate{
		decimateCount:     decimateCount,
		remainderToKeep:   remainderToKeep,
		groupByFieldNames: groupByFieldNames,
		countsByGroup:     make(map[string]int),
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerDecimate) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		groupingKey, ok := inrec.GetSelectedValuesJoined(this.groupByFieldNames)
		if !ok {
			return // This particular record doesn't have the specified fields; ignore
		}

		countForGroup, ok := this.countsByGroup[groupingKey]
		if !ok {
			countForGroup = 0
			this.countsByGroup[groupingKey] = countForGroup
		}

		remainder := countForGroup % this.decimateCount
		if remainder == this.remainderToKeep {
			outputChannel <- inrecAndContext
		}

		countForGroup++
		this.countsByGroup[groupingKey] = countForGroup

	} else {
		outputChannel <- inrecAndContext // Emit the stream-terminating null record
	}
}
