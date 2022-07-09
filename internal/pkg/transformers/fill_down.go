package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
	"github.com/johnkerl/miller/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameFillDown = "fill-down"

var FillDownSetup = TransformerSetup{
	Verb:         verbNameFillDown,
	UsageFunc:    transformerFillDownUsage,
	ParseCLIFunc: transformerFillDownParseCLI,
	IgnoresInput: false,
}

func transformerFillDownUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameFillDown)
	fmt.Fprintln(o, "If a given record has a missing value for a given field, fill that from")
	fmt.Fprintln(o, "the corresponding value from a previous record, if any.")
	fmt.Fprintln(o, "By default, a 'missing' field either is absent, or has the empty-string value.")
	fmt.Fprintln(o, "With -a, a field is 'missing' only if it is absent.")
	fmt.Fprintln(o, "")
	fmt.Fprintln(o, "Options:")
	fmt.Fprintln(o, " --all Operate on all fields in the input.")
	fmt.Fprintln(o, " -a|--only-if-absent If a given record has a missing value for a given field,")
	fmt.Fprintln(o, "     fill that from the corresponding value from a previous record, if any.")
	fmt.Fprintln(o, "     By default, a 'missing' field either is absent, or has the empty-string value.")
	fmt.Fprintln(o, "     With -a, a field is 'missing' only if it is absent.")
	fmt.Fprintln(o, " -f  Field names for fill-down.")
	fmt.Fprintln(o, " -h|--help Show this message.")
}

func transformerFillDownParseCLI(
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

	var fillDownFieldNames []string = nil
	doAll := false
	onlyIfAbsent := false

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
			transformerFillDownUsage(os.Stdout)
			os.Exit(0)

		} else if opt == "-f" {
			fillDownFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "--all" {
			doAll = true

		} else if opt == "-a" {
			onlyIfAbsent = true

		} else if opt == "--only-if-absent" {
			onlyIfAbsent = true

		} else {
			transformerFillDownUsage(os.Stderr)
			os.Exit(1)
		}
	}

	if fillDownFieldNames == nil && !doAll {
		transformerFillDownUsage(os.Stderr)
		os.Exit(1)
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerFillDown(
		fillDownFieldNames,
		doAll,
		onlyIfAbsent,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerFillDown struct {
	// input
	fillDownFieldNames []string
	doAll              bool
	onlyIfAbsent       bool

	// state
	lastNonNullValues map[string]*mlrval.Mlrval

	recordTransformerFunc RecordTransformerFunc
}

func NewTransformerFillDown(
	fillDownFieldNames []string,
	doAll bool,
	onlyIfAbsent bool,
) (*TransformerFillDown, error) {
	tr := &TransformerFillDown{
		fillDownFieldNames: fillDownFieldNames,
		onlyIfAbsent:       onlyIfAbsent,
		lastNonNullValues:  make(map[string]*mlrval.Mlrval),
	}

	if doAll {
		tr.recordTransformerFunc = tr.transformAll
	} else {
		tr.recordTransformerFunc = tr.transformSpecified
	}

	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerFillDown) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, outputRecordsAndContexts, inputDownstreamDoneChannel, outputDownstreamDoneChannel)
}

// ----------------------------------------------------------------
func (tr *TransformerFillDown) transformSpecified(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		for _, fillDownFieldName := range tr.fillDownFieldNames {
			present := false
			value := inrec.Get(fillDownFieldName)
			if tr.onlyIfAbsent {
				present = value != nil
			} else {
				present = value != nil && !value.IsVoid()
			}

			if present {
				// Remember it for a subsequent record lacking this field
				tr.lastNonNullValues[fillDownFieldName] = value.Copy()
			} else {
				// Reuse previously seen value, if any
				prev, ok := tr.lastNonNullValues[fillDownFieldName]
				if ok {
					inrec.PutCopy(fillDownFieldName, prev)
				}
			}
		}

		outputRecordsAndContexts.PushBack(inrecAndContext)

	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext) // end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (tr *TransformerFillDown) transformAll(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		for pe := inrec.Head; pe != nil; pe = pe.Next {
			fillDownFieldName := pe.Key
			present := false
			value := inrec.Get(fillDownFieldName)
			if tr.onlyIfAbsent {
				present = value != nil
			} else {
				present = value != nil && !value.IsVoid()
			}

			if present {
				// Remember it for a subsequent record lacking this field
				tr.lastNonNullValues[fillDownFieldName] = value.Copy()
			} else {
				// Reuse previously seen value, if any
				prev, ok := tr.lastNonNullValues[fillDownFieldName]
				if ok {
					inrec.PutCopy(fillDownFieldName, prev)
				}
			}
		}

		outputRecordsAndContexts.PushBack(inrecAndContext)

	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext) // end-of-stream marker
	}
}
