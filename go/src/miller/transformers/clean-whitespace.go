package transformers

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"miller/clitypes"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
var CleanWhitespaceSetup = transforming.TransformerSetup{
	Verb:         "clean-whitespace",
	ParseCLIFunc: transformerCleanWhitespaceParseCLI,
	IgnoresInput: false,
}

func transformerCleanWhitespaceParseCLI(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) transforming.IRecordTransformer {

	doKeys := true
	doValues := true

	// Get the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		if !strings.HasPrefix(args[argi], "-") {
			break // No more flag options to process

		} else if args[argi] == "-h" || args[argi] == "--help" {
			transformerPutUsage(os.Stdout, 0, errorHandling, args[0], verb)
			return nil // help intentionally requested

		} else if args[argi] == "-k" || args[argi] == "--keys-only" {
			doKeys = true
			doValues = false
			argi++
		} else if args[argi] == "-v" || args[argi] == "--values-only" {
			doKeys = false
			doValues = true
			argi++

		} else {
			transformerPutUsage(os.Stderr, 1, flag.ExitOnError, args[0], verb)
			os.Exit(1)
		}
	}

	if !doKeys && !doValues {
		transformerCleanWhitespaceUsage(os.Stderr, 1, flag.ExitOnError, args[0], verb)
		os.Exit(1)
	}

	transformer, _ := NewTransformerCleanWhitespace(
		doKeys,
		doValues,
	)

	*pargi = argi
	return transformer
}

func transformerCleanWhitespaceUsage(
	o *os.File,
	exitCode int,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	argv0 string,
	verb string,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
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
}

// ----------------------------------------------------------------
type TransformerCleanWhitespace struct {
	recordTransformerFunc transforming.RecordTransformerFunc
}

// ----------------------------------------------------------------
func NewTransformerCleanWhitespace(
	doKeys bool,
	doValues bool,
) (*TransformerCleanWhitespace, error) {

	this := &TransformerCleanWhitespace{}

	if doKeys && doValues {
		this.recordTransformerFunc = this.cleanWhitespaceInKeysAndValues
	} else if doKeys {
		this.recordTransformerFunc = this.cleanWhitespaceInKeys
	} else {
		this.recordTransformerFunc = this.cleanWhitespaceInValues
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerCleanWhitespace) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	this.recordTransformerFunc(inrecAndContext, outputChannel)
}

// ----------------------------------------------------------------
func (this *TransformerCleanWhitespace) cleanWhitespaceInKeysAndValues(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		newrec := types.NewMlrmapAsRecord()

		for pe := inrecAndContext.Record.Head; pe != nil; pe = pe.Next {
			oldKey := types.MlrvalFromString(pe.Key)
			newKey := types.MlrvalCleanWhitespace(&oldKey)
			newValue := types.MlrvalCleanWhitespace(pe.Value)
			// Transferring ownership from old record to new record; no copy needed
			newrec.PutReference(newKey.String(), &newValue)
		}

		outputChannel <- types.NewRecordAndContext(newrec, &inrecAndContext.Context)
	} else {
		outputChannel <- inrecAndContext
	}
}

// ----------------------------------------------------------------
func (this *TransformerCleanWhitespace) cleanWhitespaceInKeys(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		newrec := types.NewMlrmapAsRecord()

		for pe := inrecAndContext.Record.Head; pe != nil; pe = pe.Next {
			oldKey := types.MlrvalFromString(pe.Key)
			newKey := types.MlrvalCleanWhitespace(&oldKey)
			// Transferring ownership from old record to new record; no copy needed
			newrec.PutReference(newKey.String(), pe.Value)
		}

		outputChannel <- types.NewRecordAndContext(newrec, &inrecAndContext.Context)
	} else {
		outputChannel <- inrecAndContext
	}
}

// ----------------------------------------------------------------
func (this *TransformerCleanWhitespace) cleanWhitespaceInValues(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		for pe := inrecAndContext.Record.Head; pe != nil; pe = pe.Next {
			newValue := types.MlrvalCleanWhitespace(pe.Value)
			pe.Value = &newValue
		}
		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext
	}
}
