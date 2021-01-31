package transformers

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"miller/clitypes"
	"miller/lib"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
const verbNameSample = "sample"

var SampleSetup = transforming.TransformerSetup{
	Verb:         verbNameSample,
	ParseCLIFunc: transformerSampleParseCLI,
	UsageFunc:    transformerSampleUsage,
	IgnoresInput: false,
}

func transformerSampleParseCLI(
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

	sampleCount := -1
	var groupByFieldNames []string = nil

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		if !strings.HasPrefix(args[argi], "-") {
			break // No more flag options to process

		} else if args[argi] == "-h" || args[argi] == "--help" {
			transformerSampleUsage(os.Stdout, true, 0)
			return nil // help intentionally requested

		} else if args[argi] == "-k" {
			sampleCount = clitypes.VerbGetIntArgOrDie(verb, args, &argi, argc)

		} else if args[argi] == "-g" {
			groupByFieldNames = clitypes.VerbGetStringArrayArgOrDie(verb, args, &argi, argc)

		} else {
			transformerSampleUsage(os.Stderr, true, 1)
			os.Exit(1)
		}
	}

	if sampleCount < 0 {
		transformerSampleUsage(os.Stderr, true, 1)
	}

	transformer, _ := NewTransformerSample(
		sampleCount,
		groupByFieldNames,
	)

	*pargi = argi
	return transformer
}

func transformerSampleUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", os.Args[0], verbNameSample)
	fmt.Fprintf(o,
		`Reservoir sampling (subsampling without replacement), optionally by category.
See also %s bootstrap and %s shuffle.
`, os.Args[0], os.Args[0])
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-g {a,b,c} Optional: group-by-field names for samples, e.g. a,b,c.\n")
	fmt.Fprintf(o, "-k {k} Required: number of records to output in total, or by group if using -g.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

// ----------------------------------------------------------------
type sampleBucketType struct {
	nalloc             int
	nused              int
	recordsAndContexts []*types.RecordAndContext
}

type TransformerSample struct {
	groupByFieldNames []string
	sampleCount       int
	bucketsByGroup    *lib.OrderedMap
}

func NewTransformerSample(
	sampleCount int,
	groupByFieldNames []string,
) (*TransformerSample, error) {
	this := &TransformerSample{
		sampleCount:       sampleCount,
		groupByFieldNames: groupByFieldNames,
		bucketsByGroup:    lib.NewOrderedMap(),
	}
	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerSample) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	// Not end of input stream: retain the record, and emit nothing until end of stream.
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		groupingKey, ok := inrec.GetSelectedValuesJoined(this.groupByFieldNames)
		if ok {
			sampleBucket := this.bucketsByGroup.Get(groupingKey)
			if sampleBucket == nil {
				sampleBucket = newSampleBucket(this.sampleCount)
				this.bucketsByGroup.Put(groupingKey, sampleBucket)
			}
			sampleBucket.(*sampleBucketType).handleRecord(inrecAndContext, inrecAndContext.Context.NR)
		} // else, specified keys aren't present in this record, so ignore it

	} else { // end of record stream

		for pe := this.bucketsByGroup.Head; pe != nil; pe = pe.Next {
			sampleBucket := pe.Value.(*sampleBucketType)
			for i := 0; i < sampleBucket.nused; i++ {
				outputChannel <- sampleBucket.recordsAndContexts[i]

			}
		}

		// Emit the stream-terminating null record
		outputChannel <- inrecAndContext
	}
}

// ----------------------------------------------------------------
func newSampleBucket(sampleCount int) *sampleBucketType {
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
func (this *sampleBucketType) handleRecord(
	inrecAndContext *types.RecordAndContext,
	recordNumber int,
) {
	if this.nused < this.nalloc {
		// Always accept new entries until the bucket is full.
		//
		// Note: we need to copy the record since Go is concurrent and all
		// Miller transformers execute in their own goroutine -- if we just keep a
		// pointer, a downstream transformer mutate the pointed-to record between
		// our saving it and our re-using it.
		this.recordsAndContexts[this.nused] = inrecAndContext.Copy()
		this.nused++
	} else {
		r := int(lib.RandInt63()) % recordNumber
		if r < this.nalloc {
			this.recordsAndContexts[r] = inrecAndContext.Copy()
		}
	}
}
