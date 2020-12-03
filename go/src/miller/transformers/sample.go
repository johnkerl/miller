package transformers

import (
	"flag"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/lib"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
var SampleSetup = transforming.TransformerSetup{
	Verb:         "sample",
	ParseCLIFunc: transformerSampleParseCLI,
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

	// Get the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// Parse local flags
	flagSet := flag.NewFlagSet(verb, errorHandling)

	// TODO: Needs to be 64-bit friendly
	pSampleCount := flagSet.Int64(
		"k",
		-1,
		`Required: number of records to output in total, or by group if using -g.`,
	)

	pGroupByFieldNames := flagSet.String(
		"g",
		"",
		"Optional: group-by-field names for samples, e.g. a,b,c",
	)

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		transformerSampleUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	if *pSampleCount < 0 {
		transformerSampleUsage(os.Stderr, args[0], verb, flagSet)
		os.Exit(1)
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	transformer, _ := NewTransformerSample(
		*pSampleCount,
		*pGroupByFieldNames,
	)

	*pargi = argi
	return transformer
}

func transformerSampleUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprintf(o,
		`Reservoir sampling (subsampling without replacement), optionally by category.
See also %s bootstrap and %s shuffle.
`, argv0, argv0)
	fmt.Fprintf(o, "Options:\n")

	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v %v\n", f.Name, f.Usage)
	})
}

// ----------------------------------------------------------------
type sampleBucketType struct {
	nalloc             int64
	nused              int64
	recordsAndContexts []*types.RecordAndContext
}

type TransformerSample struct {
	groupByFieldNameList []string
	sampleCount          int64
	bucketsByGroup       *lib.OrderedMap
}

func NewTransformerSample(
	sampleCount int64,
	groupByFieldNames string,
) (*TransformerSample, error) {
	this := &TransformerSample{
		sampleCount:          sampleCount,
		groupByFieldNameList: lib.SplitString(groupByFieldNames, ","),
		bucketsByGroup:       lib.NewOrderedMap(),
	}
	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerSample) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	// Not end of input stream: retain the record, and emit nothing until end of stream.
	if inrec != nil {
		groupingKey, ok := inrec.GetSelectedValuesJoined(this.groupByFieldNameList)
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
			var i int64 = 0
			for i = 0; i < sampleBucket.nused; i++ {
				outputChannel <- sampleBucket.recordsAndContexts[i]

			}
		}

		// Emit the stream-terminating null record
		outputChannel <- inrecAndContext
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
func (this *sampleBucketType) handleRecord(
	inrecAndContext *types.RecordAndContext,
	recordNumber int64,
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
		r := lib.RandInt63() % recordNumber
		if r < this.nalloc {
			this.recordsAndContexts[r] = inrecAndContext.Copy()
		}
	}
}
