package transformers

import (
	"fmt"
	"os"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameHead = "head"

var HeadSetup = TransformerSetup{
	Verb:         verbNameHead,
	UsageFunc:    transformerHeadUsage,
	ParseCLIFunc: transformerHeadParseCLI,
	IgnoresInput: false,
}

func transformerHeadUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameHead)
	fmt.Fprintf(o, "Passes through the first n records, optionally by category.\n")
	fmt.Fprintf(o, "Without -g, ceases consuming more input (i.e. is fast) when n records have been read.\n")

	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-g {a,b,c} Optional group-by-field names for head counts, e.g. a,b,c.\n")
	fmt.Fprintf(o, "-n {n} Head-count to print. Default 10.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerHeadParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	headCount := 10
	var groupByFieldNames []string = nil

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerHeadUsage(os.Stdout, true, 0)

		} else if opt == "-n" {
			headCount = cli.VerbGetIntArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-g" {
			groupByFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerHeadUsage(os.Stderr, true, 1)
		}
	}

	transformer, err := NewTransformerHead(
		headCount,
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
type TransformerHead struct {
	// input
	headCount         int
	groupByFieldNames []string

	// state
	recordTransformerFunc RecordTransformerFunc
	unkeyedRecordCount    int
	keyedRecordCounts     map[string]int

	// See ChainTransformer
	wroteDownstreamDone bool
}

func NewTransformerHead(
	headCount int,
	groupByFieldNames []string,
) (*TransformerHead, error) {

	tr := &TransformerHead{
		headCount:           headCount,
		groupByFieldNames:   groupByFieldNames,
		unkeyedRecordCount:  0,
		keyedRecordCounts:   make(map[string]int),
		wroteDownstreamDone: false,
	}

	if groupByFieldNames == nil {
		tr.recordTransformerFunc = tr.transformUnkeyed
	} else {
		tr.recordTransformerFunc = tr.transformKeyed
	}

	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerHead) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, inputDownstreamDoneChannel, outputDownstreamDoneChannel, outputChannel)
}

func (tr *TransformerHead) transformUnkeyed(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		tr.unkeyedRecordCount++
		if tr.unkeyedRecordCount <= tr.headCount {
			outputChannel <- inrecAndContext
		} else if !tr.wroteDownstreamDone {
			// Signify to data producers upstream that we'll ignore further
			// data, so as far as we're concerned they can stop sending it. See
			// ChainTransformer.
			outputDownstreamDoneChannel <- true
			tr.wroteDownstreamDone = true
		}
	} else {
		outputChannel <- inrecAndContext
	}
}

func (tr *TransformerHead) transformKeyed(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		groupingKey, ok := inrec.GetSelectedValuesJoined(tr.groupByFieldNames)
		if !ok {
			return
		}

		count, present := tr.keyedRecordCounts[groupingKey]
		if !present { // first time
			tr.keyedRecordCounts[groupingKey] = 1
			count = 1
		} else {
			tr.keyedRecordCounts[groupingKey] += 1
			count += 1
		}

		if count <= tr.headCount {
			outputChannel <- inrecAndContext
		}

	} else {
		outputChannel <- inrecAndContext
	}
}
