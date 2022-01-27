package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameSample = "sample"

var SampleSetup = TransformerSetup{
	Verb:         verbNameSample,
	UsageFunc:    transformerSampleUsage,
	ParseCLIFunc: transformerSampleParseCLI,
	IgnoresInput: false,
}

func transformerSampleUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameSample)
	fmt.Fprintf(o,
		`Reservoir sampling (subsampling without replacement), optionally by category.
See also %s bootstrap and %s shuffle.
`, "mlr", "mlr")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-g {a,b,c} Optional: group-by-field names for samples, e.g. a,b,c.\n")
	fmt.Fprintf(o, "-k {k} Required: number of records to output in total, or by group if using -g.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerSampleParseCLI(
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

	sampleCount := int64(-1)
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
			transformerSampleUsage(os.Stdout, true, 0)

		} else if opt == "-k" {
			sampleCount = cli.VerbGetIntArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-g" {
			groupByFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerSampleUsage(os.Stderr, true, 1)
		}
	}

	if sampleCount < 0 {
		transformerSampleUsage(os.Stderr, true, 1)
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerSample(
		sampleCount,
		groupByFieldNames,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type sampleBucketType struct {
	nalloc             int64
	nused              int64
	recordsAndContexts []*types.RecordAndContext
}

type TransformerSample struct {
	groupByFieldNames []string
	sampleCount       int64
	bucketsByGroup    *lib.OrderedMap
}

func NewTransformerSample(
	sampleCount int64,
	groupByFieldNames []string,
) (*TransformerSample, error) {
	tr := &TransformerSample{
		sampleCount:       sampleCount,
		groupByFieldNames: groupByFieldNames,
		bucketsByGroup:    lib.NewOrderedMap(),
	}
	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerSample) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	// Not end of input stream: retain the record, and emit nothing until end of stream.
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		groupingKey, ok := inrec.GetSelectedValuesJoined(tr.groupByFieldNames)
		if ok {
			sampleBucket := tr.bucketsByGroup.Get(groupingKey)
			if sampleBucket == nil {
				sampleBucket = newSampleBucket(tr.sampleCount)
				tr.bucketsByGroup.Put(groupingKey, sampleBucket)
			}
			sampleBucket.(*sampleBucketType).handleRecord(inrecAndContext, inrecAndContext.Context.NR)
		} // else, specified keys aren't present in this record, so ignore it

	} else { // end of record stream

		for pe := tr.bucketsByGroup.Head; pe != nil; pe = pe.Next {
			sampleBucket := pe.Value.(*sampleBucketType)
			for i := int64(0); i < sampleBucket.nused; i++ {
				outputRecordsAndContexts.PushBack(sampleBucket.recordsAndContexts[i])

			}
		}

		// Emit the stream-terminating null record
		outputRecordsAndContexts.PushBack(inrecAndContext)
	}
}

// ----------------------------------------------------------------
func newSampleBucket(sampleCount int64) *sampleBucketType {
	return &sampleBucketType{
		nalloc:             sampleCount,
		nused:              0,
		recordsAndContexts: make([]*types.RecordAndContext, sampleCount),
	}
}

// ----------------------------------------------------------------
// This is the reservoir-sampling algorithm.  Here we retain an input record
// (if retained in the sample) or let it be GC'ed (if not retained in the
// sample).
func (t *sampleBucketType) handleRecord(
	inrecAndContext *types.RecordAndContext,
	recordNumber int64,
) {
	if t.nused < t.nalloc {
		// Always accept new entries until the bucket is full.
		//
		// Note: we need to copy the record since Go is concurrent and all
		// Miller transformers execute in their own goroutine -- if we just keep a
		// pointer, a downstream transformer mutate the pointed-to record between
		// our saving it and our re-using it.
		t.recordsAndContexts[t.nused] = inrecAndContext.Copy()
		t.nused++
	} else {
		r := int64(lib.RandInt63()) % recordNumber
		if r < t.nalloc {
			t.recordsAndContexts[r] = inrecAndContext.Copy()
		}
	}
}
