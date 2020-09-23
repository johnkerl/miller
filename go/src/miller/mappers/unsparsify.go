package mappers

import (
	"container/list"
	"flag"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/lib"
	"miller/mapping"
	"miller/types"
)

// ----------------------------------------------------------------
var UnsparsifySetup = mapping.MapperSetup{
	Verb:         "unsparsify",
	ParseCLIFunc: mapperUnsparsifyParseCLI,
	IgnoresInput: false,
}

func mapperUnsparsifyParseCLI(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) mapping.IRecordMapper {

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

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		mapperUnsparsifyUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	mapper, _ := NewMapperUnsparsify(*pFillerString)

	*pargi = argi
	return mapper
}

func mapperUnsparsifyUsage(
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
type MapperUnsparsify struct {
	fillerString       string
	recordsAndContexts *list.List
	fieldNamesSeen     *lib.OrderedMap
}

func NewMapperUnsparsify(
	fillerString string,
) (*MapperUnsparsify, error) {

	this := &MapperUnsparsify{
		fillerString:       fillerString,
		recordsAndContexts: list.New(),
		fieldNamesSeen:     lib.NewOrderedMap(),
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *MapperUnsparsify) Map(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			key := *pe.Key
			if !this.fieldNamesSeen.Has(key) {
				this.fieldNamesSeen.Put(key, key)
			}
		}
		this.recordsAndContexts.PushBack(inrecAndContext)
	} else {
		fillerMlrval := types.MlrvalFromString(this.fillerString)
		for e := this.recordsAndContexts.Front(); e != nil; e = e.Next() {
			outrecAndContext := e.Value.(*types.RecordAndContext)
			outrec := outrecAndContext.Record

			newrec := types.NewMlrmapAsRecord()
			for pe := this.fieldNamesSeen.Head; pe != nil; pe = pe.Next {
				fieldName := pe.Key
				if !outrec.Has(&fieldName) {
					newrec.PutCopy(&fieldName, &fillerMlrval)
				} else {
					newrec.PutReference(&fieldName, outrec.Get(&fieldName))
				}
			}

			outputChannel <- types.NewRecordAndContext(newrec, &outrecAndContext.Context)
		}

		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
