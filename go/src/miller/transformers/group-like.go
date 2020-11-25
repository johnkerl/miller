package transformers

import (
	"container/list"
	"flag"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/lib"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
var GroupLikeSetup = transforming.TransformerSetup{
	Verb:         "group-like",
	ParseCLIFunc: transformerGroupLikeParseCLI,
	IgnoresInput: false,
}

func transformerGroupLikeParseCLI(
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

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		transformerGroupLikeUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	transformer, _ := NewTransformerGroupLike()

	*pargi = argi
	return transformer
}

func transformerGroupLikeUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprint(o,
		`Outputs records in batches having identical field names.
`)
	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v (default %v) %v\n", f.Name, f.Value, f.Usage) // f.Name, f.Value
	})
}

// ----------------------------------------------------------------
type TransformerGroupLike struct {
	// map from string to *list.List
	recordListsByGroup *lib.OrderedMap
}

func NewTransformerGroupLike() (*TransformerGroupLike, error) {

	this := &TransformerGroupLike{
		recordListsByGroup: lib.NewOrderedMap(),
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerGroupLike) Map(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream

		groupingKey := inrec.GetKeysJoined()

		recordListForGroup := this.recordListsByGroup.Get(groupingKey)
		if recordListForGroup == nil { // first time
			recordListForGroup = list.New()
			this.recordListsByGroup.Put(groupingKey, recordListForGroup)
		}

		recordListForGroup.(*list.List).PushBack(inrecAndContext)

	} else {
		for outer := this.recordListsByGroup.Head; outer != nil; outer = outer.Next {
			recordListForGroup := outer.Value.(*list.List)
			for inner := recordListForGroup.Front(); inner != nil; inner = inner.Next() {
				outputChannel <- inner.Value.(*types.RecordAndContext)
			}
		}
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
