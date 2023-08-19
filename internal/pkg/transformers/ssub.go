package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/bifs"
	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
	"github.com/johnkerl/miller/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameSsub = "ssub"

var SsubSetup = TransformerSetup{
	Verb:         verbNameSsub,
	UsageFunc:    transformerSsubUsage,
	ParseCLIFunc: transformerSsubParseCLI,
	IgnoresInput: false,
}

func transformerSsubUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameSsub)
	fmt.Fprintf(o, "Replaces old string with new string in specified field(s), without regex support for\n")
	fmt.Fprintf(o, "the old string, like the `ssub` DSL function. See also the `gsub` and `sub` verbs.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-f {a,b,c}  Field names to convert.\n")
	fmt.Fprintf(o, "-h|--help   Show this message.\n")
}

func transformerSsubParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// Parse local flags
	var fieldNames []string = nil
	var oldText string
	var newText string

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
			transformerSsubUsage(os.Stdout)
			os.Exit(0)

		} else if opt == "-f" {
			fieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)
		} else {
			transformerSsubUsage(os.Stderr)
			os.Exit(1)
		}
	}

	if fieldNames == nil {
		transformerSsubUsage(os.Stderr)
		os.Exit(1)
	}

	// Get the old and new text from the command line
	if (argc - argi) < 2 {
		transformerSsubUsage(os.Stderr)
		os.Exit(1)
	}
	oldText = args[argi]
	newText = args[argi+1]

	argi += 2

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerSsub(
		fieldNames,
		oldText,
		newText,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerSsub struct {
	fieldNames []string
	oldText    *mlrval.Mlrval
	newText    *mlrval.Mlrval
}

// ----------------------------------------------------------------
func NewTransformerSsub(
	fieldNames []string,
	oldText string,
	newText string,
) (*TransformerSsub, error) {
	tr := &TransformerSsub{
		fieldNames: fieldNames,
		oldText:    mlrval.FromString(oldText),
		newText:    mlrval.FromString(newText),
	}
	return tr, nil
}

func (tr *TransformerSsub) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)

	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		for _, fieldName := range tr.fieldNames {
			oldValue := inrec.Get(fieldName)
			if oldValue == nil {
				continue
			}

			newValue := bifs.BIF_ssub(oldValue, tr.oldText, tr.newText)

			inrec.PutReference(fieldName, newValue)
		}

		outputRecordsAndContexts.PushBack(inrecAndContext)
	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext) // emit end-of-stream marker
	}
}
