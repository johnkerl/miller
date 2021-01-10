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
var HeadSetup = transforming.TransformerSetup{
	Verb:         "head",
	ParseCLIFunc: transformerHeadParseCLI,
	IgnoresInput: false,
}

func transformerHeadParseCLI(
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

	//Usage: mlr head [options]
	//-n {count}    Head count to print; default 10
	//-g {a,b,c}    Optional group-by-field names for head counts

	pHeadCount := flagSet.Uint64(
		"n",
		10,
		`Head count to print`,
	)

	pGroupByFieldNames := flagSet.String(
		"g",
		"",
		"Optional group-by-field names for head counts, e.g. a,b,c",
	)

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		transformerHeadUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	transformer, _ := NewTransformerHead(
		*pHeadCount,
		*pGroupByFieldNames,
	)

	*pargi = argi
	return transformer
}

func transformerHeadUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprint(o,
		`Passes through the first n records, optionally by category.
`)
	// TODO: work on this, keeping in mind https://github.com/johnkerl/miller/issues/291
	//	fmt.Fprint(o,
	//		`Without -g, ceases consuming more input (i.e. is fast) when n records
	//have been read.
	//`)

	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v (default %v) %v\n", f.Name, f.Value, f.Usage) // f.Name, f.Value
	})
}

// ----------------------------------------------------------------
type TransformerHead struct {
	// input
	headCount            uint64
	groupByFieldNameList []string

	// state
	recordTransformerFunc transforming.RecordTransformerFunc
	unkeyedRecordCount    uint64
	keyedRecordCounts     map[string]uint64
}

func NewTransformerHead(
	headCount uint64,
	groupByFieldNames string,
) (*TransformerHead, error) {

	groupByFieldNameList := lib.SplitString(groupByFieldNames, ",")

	this := &TransformerHead{
		headCount:            headCount,
		groupByFieldNameList: groupByFieldNameList,

		unkeyedRecordCount: 0,
		keyedRecordCounts:  make(map[string]uint64),
	}

	if len(groupByFieldNameList) == 0 {
		this.recordTransformerFunc = this.mapUnkeyed
	} else {
		this.recordTransformerFunc = this.mapKeyed
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerHead) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	this.recordTransformerFunc(inrecAndContext, outputChannel)
}

func (this *TransformerHead) mapUnkeyed(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		this.unkeyedRecordCount++
		if this.unkeyedRecordCount <= this.headCount {
			outputChannel <- inrecAndContext
		}
	} else {
		outputChannel <- inrecAndContext
	}
}

func (this *TransformerHead) mapKeyed(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		groupingKey, ok := inrec.GetSelectedValuesJoined(this.groupByFieldNameList)
		if !ok {
			return
		}

		count, present := this.keyedRecordCounts[groupingKey]
		if !present { // first time
			this.keyedRecordCounts[groupingKey] = 1
			count = 1
		} else {
			this.keyedRecordCounts[groupingKey] += 1
			count += 1
		}

		if count <= this.headCount {
			outputChannel <- inrecAndContext
		}

	} else {
		outputChannel <- inrecAndContext
	}
}
