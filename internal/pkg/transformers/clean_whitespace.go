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
const verbNameCleanWhitespace = "clean-whitespace"

var CleanWhitespaceSetup = TransformerSetup{
	Verb:         verbNameCleanWhitespace,
	UsageFunc:    transformerCleanWhitespaceUsage,
	ParseCLIFunc: transformerCleanWhitespaceParseCLI,
	IgnoresInput: false,
}

func transformerCleanWhitespaceUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameCleanWhitespace)
	fmt.Fprintf(o, "For each record, for each field in the record, whitespace-cleans the keys and/or\n")
	fmt.Fprintf(o, "values. Whitespace-cleaning entails stripping leading and trailing whitespace,\n")
	fmt.Fprintf(o, "and replacing multiple whitespace with singles. For finer-grained control,\n")
	fmt.Fprintf(o, "please see the DSL functions lstrip, rstrip, strip, collapse_whitespace,\n")
	fmt.Fprintf(o, "and clean_whitespace.\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-k|--keys-only    Do not touch values.\n")
	fmt.Fprintf(o, "-v|--values-only  Do not touch keys.\n")
	fmt.Fprintf(o, "It is an error to specify -k as well as -v -- to clean keys and values,\n")
	fmt.Fprintf(o, "leave off -k as well as -v.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerCleanWhitespaceParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) IRecordTransformer {

	doKeys := true
	doValues := true

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
			transformerCleanWhitespaceUsage(os.Stdout, true, 0)

		} else if opt == "-k" || opt == "--keys-only" {
			doKeys = true
			doValues = false
		} else if opt == "-v" || opt == "--values-only" {
			doKeys = false
			doValues = true

		} else {
			transformerCleanWhitespaceUsage(os.Stderr, true, 1)
		}
	}

	if !doKeys && !doValues {
		transformerCleanWhitespaceUsage(os.Stderr, true, 1)
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerCleanWhitespace(
		doKeys,
		doValues,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerCleanWhitespace struct {
	recordTransformerFunc RecordTransformerFunc
}

// ----------------------------------------------------------------
func NewTransformerCleanWhitespace(
	doKeys bool,
	doValues bool,
) (*TransformerCleanWhitespace, error) {

	tr := &TransformerCleanWhitespace{}

	if doKeys && doValues {
		tr.recordTransformerFunc = tr.cleanWhitespaceInKeysAndValues
	} else if doKeys {
		tr.recordTransformerFunc = tr.cleanWhitespaceInKeys
	} else {
		tr.recordTransformerFunc = tr.cleanWhitespaceInValues
	}

	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerCleanWhitespace) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, outputRecordsAndContexts, inputDownstreamDoneChannel, outputDownstreamDoneChannel)
}

// ----------------------------------------------------------------
func (tr *TransformerCleanWhitespace) cleanWhitespaceInKeysAndValues(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		newrec := mlrval.NewMlrmapAsRecord()

		for pe := inrecAndContext.Record.Head; pe != nil; pe = pe.Next {
			oldKey := mlrval.FromString(pe.Key)
			// xxx temp
			newKey := bifs.BIF_clean_whitespace(oldKey)
			newValue := bifs.BIF_clean_whitespace(pe.Value)
			// Transferring ownership from old record to new record; no copy needed
			newrec.PutReference(newKey.String(), newValue)
		}

		outputRecordsAndContexts.PushBack(types.NewRecordAndContext(newrec, &inrecAndContext.Context))
	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext)
	}
}

// ----------------------------------------------------------------
func (tr *TransformerCleanWhitespace) cleanWhitespaceInKeys(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		newrec := mlrval.NewMlrmapAsRecord()

		for pe := inrecAndContext.Record.Head; pe != nil; pe = pe.Next {
			oldKey := mlrval.FromString(pe.Key)
			newKey := bifs.BIF_clean_whitespace(oldKey)
			// Transferring ownership from old record to new record; no copy needed
			newrec.PutReference(newKey.String(), pe.Value)
		}

		outputRecordsAndContexts.PushBack(types.NewRecordAndContext(newrec, &inrecAndContext.Context))
	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext)
	}
}

// ----------------------------------------------------------------
func (tr *TransformerCleanWhitespace) cleanWhitespaceInValues(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		for pe := inrecAndContext.Record.Head; pe != nil; pe = pe.Next {
			pe.Value = bifs.BIF_clean_whitespace(pe.Value)
		}
		outputRecordsAndContexts.PushBack(inrecAndContext)
	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext)
	}
}
