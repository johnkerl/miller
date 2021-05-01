package transformers

import (
	"fmt"
	"os"
	"strings"

	"miller/src/cliutil"
	"miller/src/lib"
	"miller/src/transforming"
	"miller/src/types"
)

// ----------------------------------------------------------------
const verbNameFillDown = "fill-down"

var FillDownSetup = transforming.TransformerSetup{
	Verb:         verbNameFillDown,
	UsageFunc:    transformerFillDownUsage,
	ParseCLIFunc: transformerFillDownParseCLI,
	IgnoresInput: false,
}

func transformerFillDownUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", lib.MlrExeName(), verbNameFillDown)
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

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerFillDownParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cliutil.TReaderOptions,
	__ *cliutil.TWriterOptions,
) transforming.IRecordTransformer {

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
		argi++

		if opt == "-h" || opt == "--help" {
			transformerFillDownUsage(os.Stdout, true, 0)

		} else if opt == "-f" {
			fillDownFieldNames = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "--all" {
			doAll = true

		} else if opt == "-a" {
			onlyIfAbsent = true

		} else if opt == "--only-if-absent" {
			onlyIfAbsent = true

		} else {
			transformerFillDownUsage(os.Stderr, true, 1)
		}
	}

	if fillDownFieldNames == nil && !doAll {
		transformerFillDownUsage(os.Stderr, true, 1)
	}

	transformer, _ := NewTransformerFillDown(
		fillDownFieldNames,
		doAll,
		onlyIfAbsent,
	)

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerFillDown struct {
	// input
	fillDownFieldNames []string
	doAll              bool
	onlyIfAbsent       bool

	// state
	lastNonNullValues map[string]*types.Mlrval

	recordTransformerFunc transforming.RecordTransformerFunc
}

func NewTransformerFillDown(
	fillDownFieldNames []string,
	doAll bool,
	onlyIfAbsent bool,
) (*TransformerFillDown, error) {
	this := &TransformerFillDown{
		fillDownFieldNames: fillDownFieldNames,
		onlyIfAbsent:       onlyIfAbsent,
		lastNonNullValues:  make(map[string]*types.Mlrval),
	}

	if doAll {
		this.recordTransformerFunc = this.transformAll
	} else {
		this.recordTransformerFunc = this.transformSpecified
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerFillDown) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	this.recordTransformerFunc(inrecAndContext, outputChannel)
}

// ----------------------------------------------------------------
func (this *TransformerFillDown) transformSpecified(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		for _, fillDownFieldName := range this.fillDownFieldNames {
			present := false
			value := inrec.Get(fillDownFieldName)
			if this.onlyIfAbsent {
				present = value != nil
			} else {
				present = value != nil && !value.IsEmpty()
			}

			if present {
				// Remember it for a subsequent record lacking this field
				this.lastNonNullValues[fillDownFieldName] = value.Copy()
			} else {
				// Reuse previously seen value, if any
				prev, ok := this.lastNonNullValues[fillDownFieldName]
				if ok {
					inrec.PutCopy(fillDownFieldName, prev)
				}
			}
		}

		outputChannel <- inrecAndContext

	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (this *TransformerFillDown) transformAll(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		for pe := inrec.Head; pe != nil; pe = pe.Next {
			fillDownFieldName := pe.Key
			present := false
			value := inrec.Get(fillDownFieldName)
			if this.onlyIfAbsent {
				present = value != nil
			} else {
				present = value != nil && !value.IsEmpty()
			}

			if present {
				// Remember it for a subsequent record lacking this field
				this.lastNonNullValues[fillDownFieldName] = value.Copy()
			} else {
				// Reuse previously seen value, if any
				prev, ok := this.lastNonNullValues[fillDownFieldName]
				if ok {
					inrec.PutCopy(fillDownFieldName, prev)
				}
			}
		}

		outputChannel <- inrecAndContext

	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
