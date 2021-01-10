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
var UnsparsifySetup = transforming.TransformerSetup{
	Verb:         "unsparsify",
	ParseCLIFunc: transformerUnsparsifyParseCLI,
	IgnoresInput: false,
}

func transformerUnsparsifyParseCLI(
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

	pFillerString := flagSet.String(
		"fill-with",
		"",
		"Prepend field {name} to each record with record-counter starting at 1",
	)

	pSpecifiedFieldNames := flagSet.String(
		"f",
		"",
		`Specify field names to be operated on. Any other fields won't be
modified, and operation will be streaming.`,
	)

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		transformerUnsparsifyUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	transformer, _ := NewTransformerUnsparsify(
		*pFillerString,
		*pSpecifiedFieldNames,
	)

	*pargi = argi
	return transformer
}

func transformerUnsparsifyUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprint(o,
		`Prints records with the union of field names over all input records.
For field names absent in a given record but present in others, fills in
a value. This verb retains all input before producing any output.

Example: if the input is two records, one being 'a=1,b=2' and the other
being 'b=3,c=4', then the output is the two records 'a=1,b=2,c=' and
'a=,b=3,c=4'.
`)
	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v (default %v) %v\n", f.Name, f.Value, f.Usage) // f.Name, f.Value
	})
}

// ----------------------------------------------------------------
type TransformerUnsparsify struct {
	fillerMlrval          types.Mlrval
	recordsAndContexts    *list.List
	fieldNamesSeen        *lib.OrderedMap
	recordTransformerFunc transforming.RecordTransformerFunc
}

func NewTransformerUnsparsify(
	fillerString string,
	specifiedFieldNames string,
) (*TransformerUnsparsify, error) {

	specifiedFieldNameList := lib.SplitString(specifiedFieldNames, ",")
	fieldNamesSeen := lib.NewOrderedMap()
	for _, specifiedFieldName := range specifiedFieldNameList {
		fieldNamesSeen.Put(specifiedFieldName, specifiedFieldName)
	}

	this := &TransformerUnsparsify{
		fillerMlrval:       types.MlrvalFromString(fillerString),
		recordsAndContexts: list.New(),
		fieldNamesSeen:     fieldNamesSeen,
	}

	if len(specifiedFieldNameList) == 0 {
		this.recordTransformerFunc = this.mapNonStreaming
	} else {
		this.recordTransformerFunc = this.mapStreaming
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerUnsparsify) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	this.recordTransformerFunc(inrecAndContext, outputChannel)
}

// ----------------------------------------------------------------
func (this *TransformerUnsparsify) mapNonStreaming(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			key := *pe.Key
			if !this.fieldNamesSeen.Has(key) {
				this.fieldNamesSeen.Put(key, key)
			}
		}
		this.recordsAndContexts.PushBack(inrecAndContext)
	} else {
		for e := this.recordsAndContexts.Front(); e != nil; e = e.Next() {
			outrecAndContext := e.Value.(*types.RecordAndContext)
			outrec := outrecAndContext.Record

			newrec := types.NewMlrmapAsRecord()
			for pe := this.fieldNamesSeen.Head; pe != nil; pe = pe.Next {
				fieldName := pe.Key
				if !outrec.Has(&fieldName) {
					newrec.PutCopy(&fieldName, &this.fillerMlrval)
				} else {
					newrec.PutReference(&fieldName, outrec.Get(&fieldName))
				}
			}

			outputChannel <- types.NewRecordAndContext(newrec, &outrecAndContext.Context)
		}

		outputChannel <- inrecAndContext // end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (this *TransformerUnsparsify) mapStreaming(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		for pe := this.fieldNamesSeen.Head; pe != nil; pe = pe.Next {
			if !inrec.Has(&pe.Key) {
				inrec.PutCopy(&pe.Key, &this.fillerMlrval)
			}
		}

		outputChannel <- inrecAndContext

	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
