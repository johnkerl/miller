package mappers

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"miller/clitypes"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
var AltkvSetup = transforming.TransformerSetup{
	Verb:         "altkv",
	ParseCLIFunc: mapperAltkvParseCLI,
	IgnoresInput: false,
}

func mapperAltkvParseCLI(
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
		mapperAltkvUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	transformer, _ := NewTransformerAltkv()

	*pargi = argi
	return transformer
}

func mapperAltkvUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s {no options}\n", argv0, verb)
	fmt.Fprintf(o, "Given fields with values of the form a,b,c,d,e,f emits a=b,c=d,e=f pairs.\n")

	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v (default %v) %v\n", f.Name, f.Value, f.Usage) // f.Name, f.Value
	})
}

// ----------------------------------------------------------------
type TransformerAltkv struct {
}

func NewTransformerAltkv() (*TransformerAltkv, error) {
	this := &TransformerAltkv{}
	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerAltkv) Map(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream
		newrec := types.NewMlrmapAsRecord()
		outputFieldNumber := 1

		for pe := inrec.Head; pe != nil; /* increment in loop body */ {
			if pe.Next != nil { // Not at end of record with odd-numbered field count
				key := pe.Value.String()
				value := pe.Next.Value
				// Transferring ownership from old record to new record; no copy needed
				newrec.PutReference(&key, value)
			} else { // At end of record with odd-numbered field count
				key := strconv.Itoa(outputFieldNumber)
				value := pe.Value
				// Transferring ownership from old record to new record; no copy needed
				newrec.PutReference(&key, value)
			}
			outputFieldNumber++

			pe = pe.Next
			if pe == nil {
				break
			}
			pe = pe.Next
		}

		outputChannel <- types.NewRecordAndContext(newrec, &inrecAndContext.Context)

	} else { // end of record stream
		outputChannel <- inrecAndContext
	}
}
