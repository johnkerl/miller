package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
	"github.com/johnkerl/miller/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameLatin1ToUTF8 = "latin1-to-utf8"

var Latin1ToUTF8Setup = TransformerSetup{
	Verb:         verbNameLatin1ToUTF8,
	UsageFunc:    transformerLatin1ToUTF8Usage,
	ParseCLIFunc: transformerLatin1ToUTF8ParseCLI,
	IgnoresInput: false,
}

func transformerLatin1ToUTF8Usage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s, with no options.\n", "mlr", verbNameLatin1ToUTF8)
	fmt.Fprintf(o, "Recursively converts record strings from Latin-1 to UTF-8.\n")
	fmt.Fprintf(o, "For field-level control, please see the latin1_to_utf8 DSL function.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
}

func transformerLatin1ToUTF8ParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) IRecordTransformer {

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
			transformerLatin1ToUTF8Usage(os.Stdout)
			os.Exit(0)

		} else {
			transformerLatin1ToUTF8Usage(os.Stderr)
			os.Exit(1)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerLatin1ToUTF8()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerLatin1ToUTF8 struct {
}

func NewTransformerLatin1ToUTF8() (*TransformerLatin1ToUTF8, error) {
	tr := &TransformerLatin1ToUTF8{}
	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerLatin1ToUTF8) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		for pe := inrec.Head; pe != nil; pe = pe.Next {
			inval := pe.Value
			if inval.IsString() {
				output, err := lib.TryLatin1ToUTF8(pe.Value.String())
				if err == nil {
					pe.Value = mlrval.FromString(output)
				} else {
					pe.Value = mlrval.ERROR
				}
			}
		}

		outputRecordsAndContexts.PushBack(types.NewRecordAndContext(inrec, &inrecAndContext.Context))

	} else { // end of record stream
		outputRecordsAndContexts.PushBack(inrecAndContext)
	}
}
