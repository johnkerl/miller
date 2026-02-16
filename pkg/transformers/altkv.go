package transformers

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

const verbNameAltkv = "altkv"

var AltkvSetup = TransformerSetup{
	Verb:         verbNameAltkv,
	UsageFunc:    transformerAltkvUsage,
	ParseCLIFunc: transformerAltkvParseCLI,
	IgnoresInput: false,
}

func transformerAltkvUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameAltkv)
	fmt.Fprintf(o, "Given fields with values of the form a,b,c,d,e,f emits a=b,c=d,e=f pairs.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
}

func transformerAltkvParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) (RecordTransformer, error) {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	argi++

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
			transformerAltkvUsage(os.Stdout)
			return nil, cli.ErrHelpRequested
		}
		transformerAltkvUsage(os.Stderr)
		return nil, fmt.Errorf("%s %s: option \"%s\" not recognized", "mlr", verbNameAltkv, opt)
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	transformer, err := NewTransformerAltkv()
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

type TransformerAltkv struct {
}

func NewTransformerAltkv() (*TransformerAltkv, error) {
	tr := &TransformerAltkv{}
	return tr, nil
}

func (tr *TransformerAltkv) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		newrec := mlrval.NewMlrmapAsRecord()
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

		*outputRecordsAndContexts = append(*outputRecordsAndContexts, types.NewRecordAndContext(newrec, &inrecAndContext.Context))

	} else { // end of record stream
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
	}
}
