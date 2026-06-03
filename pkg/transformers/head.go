package transformers

import (
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/types"
)

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
	fmt.Fprintf(o, "           A negative count, e.g. -n -2, passes through all but the last n records,\n")
	fmt.Fprintf(o, "           optionally by category.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
}

func transformerHeadParseCLI(
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
			return nil, cli.ErrHelpRequested

		} else if opt == "-n" {
			n, err := cli.VerbGetIntArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			headCount = n

		} else if opt == "-g" {
			names, err := cli.VerbGetStringArrayArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			groupByFieldNames = names

		} else {
			transformerHeadUsage(os.Stderr)
			return nil, fmt.Errorf("%s %s: option \"%s\" not recognized", "mlr", verb, opt)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	transformer, err := NewTransformerHead(
		headCount,
		groupByFieldNames,
	)
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

type TransformerHead struct {
	// input
	headCount         int64
	groupByFieldNames []string

	// state
	recordTransformerFunc RecordTransformerFunc
	unkeyedRecordCount    int64
	keyedRecordCounts     map[string]int64
	recordListsByGroup    map[string]*[]*types.RecordAndContext

	// See ChainTransformer
	wroteDownstreamDone bool
}

func NewTransformerHead(
	headCount int64,
	groupByFieldNames []string,
) (*TransformerHead, error) {

	allButLast := headCount < 0
	if allButLast {
		headCount = -headCount
	}

	tr := &TransformerHead{
		headCount:           headCount,
		groupByFieldNames:   groupByFieldNames,
		unkeyedRecordCount:  0,
		keyedRecordCounts:   make(map[string]int64),
		recordListsByGroup:  make(map[string]*[]*types.RecordAndContext),
		wroteDownstreamDone: false,
	}

	if allButLast {
		tr.recordTransformerFunc = tr.transformAllButLast
	} else if groupByFieldNames == nil {
		tr.recordTransformerFunc = tr.transformUnkeyed
	} else {
		tr.recordTransformerFunc = tr.transformKeyed
	}

	return tr, nil
}

func (tr *TransformerHead) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, outputRecordsAndContexts, inputDownstreamDoneChannel, outputDownstreamDoneChannel)
}

func (tr *TransformerHead) transformUnkeyed(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		tr.unkeyedRecordCount++
		if tr.unkeyedRecordCount <= tr.headCount {
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
		} else if !tr.wroteDownstreamDone {
			// Signify to data producers upstream that we'll ignore further
			// data, so as far as we're concerned they can stop sending it. See
			// ChainTransformer.
			//TODO: maybe remove: *outputRecordsAndContexts = append(*outputRecordsAndContexts, types.NewEndOfStreamMarker(&inrecAndContext.Context))
			outputDownstreamDoneChannel <- true
			tr.wroteDownstreamDone = true
		}
	} else {
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
	}
}

func (tr *TransformerHead) transformKeyed(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
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
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
		}

	} else {
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
	}
}

func (tr *TransformerHead) transformAllButLast(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		groupingKey, ok := inrec.GetSelectedValuesJoined(tr.groupByFieldNames)
		if !ok {
			return
		}

		recordListForGroup := tr.recordListsByGroup[groupingKey]
		if recordListForGroup == nil { // first time
			recordListForGroup = &[]*types.RecordAndContext{}
			tr.recordListsByGroup[groupingKey] = recordListForGroup
		}

		*recordListForGroup = append(*recordListForGroup, inrecAndContext)
		for int64(len(*recordListForGroup)) > tr.headCount {
			// Emit records that have fallen out of the window.
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, (*recordListForGroup)[0])
			(*recordListForGroup)[0] = nil // release the backing-array slot's reference
			*recordListForGroup = (*recordListForGroup)[1:]
		}

	} else {
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
	}
}
