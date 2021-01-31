package transformers

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"miller/clitypes"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
const verbNameAltkv = "altkv"

var AltkvSetup = transforming.TransformerSetup{
	Verb:         verbNameAltkv,
	ParseCLIFunc: transformerAltkvParseCLI,
	UsageFunc:    transformerAltkvUsage,
	IgnoresInput: false,
}

func transformerAltkvParseCLI(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) transforming.IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	argi++

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		if !strings.HasPrefix(args[argi], "-") {
			break // No more flag options to process

		} else if args[argi] == "-h" || args[argi] == "--help" {
			transformerAltkvUsage(os.Stdout, true, 0)
			return nil // help intentionally requested

		} else {
			transformerAltkvUsage(os.Stderr, true, 1)
			os.Exit(1)
		}
	}

	transformer, _ := NewTransformerAltkv()

	*pargi = argi
	return transformer
}

func transformerAltkvUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s {no options}\n", os.Args[0], verbNameAltkv)
	fmt.Fprintf(o, "Given fields with values of the form a,b,c,d,e,f emits a=b,c=d,e=f pairs.\n")
}

// ----------------------------------------------------------------
type TransformerAltkv struct {
}

func NewTransformerAltkv() (*TransformerAltkv, error) {
	this := &TransformerAltkv{}
	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerAltkv) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		newrec := types.NewMlrmapAsRecord()
		outputFieldNumber := 1

		for pe := inrec.Head; pe != nil; /* increment in loop body */ {
			if pe.Next != nil { // Not at end of record with odd-numbered field count
				key := pe.Value.String()
				value := pe.Next.Value
				// Transferring ownership from old record to new record; no copy needed
				newrec.PutReference(key, value)
			} else { // At end of record with odd-numbered field count
				key := strconv.Itoa(outputFieldNumber)
				value := pe.Value
				// Transferring ownership from old record to new record; no copy needed
				newrec.PutReference(key, value)
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
