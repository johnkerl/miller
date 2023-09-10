package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/lib"
	"github.com/johnkerl/miller/pkg/mlrval"
	"github.com/johnkerl/miller/pkg/types"
)

// ----------------------------------------------------------------
const verbNameUTF8ToLatin1 = "utf8-to-latin1"

var UTF8ToLatin1Setup = TransformerSetup{
	Verb:         verbNameUTF8ToLatin1,
	UsageFunc:    transformerUTF8ToLatin1Usage,
	ParseCLIFunc: transformerUTF8ToLatin1ParseCLI,
	IgnoresInput: false,
}

func transformerUTF8ToLatin1Usage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s, with no options.\n", "mlr", verbNameUTF8ToLatin1)
	fmt.Fprintf(o, "Recursively converts record strings from Latin-1 to UTF-8.\n")
	fmt.Fprintf(o, "For field-level control, please see the utf8_to_latin1 DSL function.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
}

func transformerUTF8ToLatin1ParseCLI(
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
			transformerUTF8ToLatin1Usage(os.Stdout)
			os.Exit(0)

		} else {
			transformerUTF8ToLatin1Usage(os.Stderr)
			os.Exit(1)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerUTF8ToLatin1()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerUTF8ToLatin1 struct {
}

func NewTransformerUTF8ToLatin1() (*TransformerUTF8ToLatin1, error) {
	tr := &TransformerUTF8ToLatin1{}
	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerUTF8ToLatin1) Transform(
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
				output, err := lib.TryUTF8ToLatin1(pe.Value.String())
				if err == nil {
					pe.Value = mlrval.FromString(output)
				} else {
					pe.Value = mlrval.FromError(err)
				}
			}
		}

		outputRecordsAndContexts.PushBack(types.NewRecordAndContext(inrec, &inrecAndContext.Context))

	} else { // end of record stream
		outputRecordsAndContexts.PushBack(inrecAndContext)
	}
}
