package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/types"
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
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameHead)
	fmt.Fprintf(o, "Passes through the first n records, optionally by category.\n")
	fmt.Fprintf(o, "Without -g, ceases consuming more input (i.e. is fast) when n records have been read.\n")

	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-g {a,b,c} Optional group-by-field names for head counts, e.g. a,b,c.\n")
	fmt.Fprintf(o, "-n {n} Head-count to print. Default 10.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
}

func transformerHeadParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	headCount := int64(10)
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

		if opt == "-h" || opt == "--help" {
			transformerHeadUsage(os.Stdout)
			os.Exit(0)

		} else if opt == "-n" {
			headCount = cli.VerbGetIntArgOrDie(verb, opt, args, &argi, argc)

			// This is a bit of a hack. In our Getoptify routine we preprocess
			// the command line sending '-xyz' to '-x -y -z', but leaving
			// '--xyz' as-is. Also, Unix-like tools often support 'head -n4'
			// and 'tail -n4' in addition to 'head -n 4' and 'tail -n 4'.  Our
			// getoptify paradigm, combined with syntax familiar to users,
			// means we get '-n -4' here. So, take the absolute value to handle this.
			if headCount < 0 {
				headCount = -headCount
			}

		} else if opt == "-g" {
			groupByFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerHeadUsage(os.Stderr)
			os.Exit(1)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerHead(
		headCount,
		groupByFieldNames,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerHead struct {
	// input
	headCount         int64
	groupByFieldNames []string

	// state
	recordTransformerFunc RecordTransformerFunc
	unkeyedRecordCount    int64
	keyedRecordCounts     map[string]int64

	// See ChainTransformer
	wroteDownstreamDone bool
}

func NewTransformerHead(
	headCount int64,
	groupByFieldNames []string,
) (*TransformerHead, error) {

	tr := &TransformerHead{
		headCount:           headCount,
		groupByFieldNames:   groupByFieldNames,
		unkeyedRecordCount:  0,
		keyedRecordCounts:   make(map[string]int64),
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
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, outputRecordsAndContexts, inputDownstreamDoneChannel, outputDownstreamDoneChannel)
}

func (tr *TransformerHead) transformUnkeyed(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		tr.unkeyedRecordCount++
		if tr.unkeyedRecordCount <= tr.headCount {
			outputRecordsAndContexts.PushBack(inrecAndContext)
		} else if !tr.wroteDownstreamDone {
			// Signify to data producers upstream that we'll ignore further
			// data, so as far as we're concerned they can stop sending it. See
			// ChainTransformer.
			//TODO: maybe remove: outputRecordsAndContexts.PushBack(types.NewEndOfStreamMarker(&inrecAndContext.Context))
			outputDownstreamDoneChannel <- true
			tr.wroteDownstreamDone = true
		}
	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext)
	}
}

func (tr *TransformerHead) transformKeyed(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
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
			outputRecordsAndContexts.PushBack(inrecAndContext)
		}

	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext)
	}
}
