package transformers

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameAltkv = "altkv"

var AltkvSetup = TransformerSetup{
	Verb:         verbNameAltkv,
	UsageFunc:    transformerAltkvUsage,
	ParseCLIFunc: transformerAltkvParseCLI,
	IgnoresInput: false,
}

func transformerAltkvUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameAltkv)
	fmt.Fprintf(o, "Given fields with values of the form a,b,c,d,e,f emits a=b,c=d,e=f pairs.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
	if doExit {
		os.Exit(exitCode)
	}
}

func transformerAltkvParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	argi++

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerAltkvUsage(os.Stdout, true, 0)

		} else {
			transformerAltkvUsage(os.Stderr, true, 1)
		}
	}

	transformer, err := NewTransformerAltkv()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerAltkv struct {
}

func NewTransformerAltkv() (*TransformerAltkv, error) {
	tr := &TransformerAltkv{}
	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerAltkv) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
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
