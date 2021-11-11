package transformers

import (
	"fmt"
	"os"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameReorder = "reorder"

var ReorderSetup = TransformerSetup{
	Verb:         verbNameReorder,
	UsageFunc:    transformerReorderUsage,
	ParseCLIFunc: transformerReorderParseCLI,
	IgnoresInput: false,
}

func transformerReorderUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	argv0 := "mlr"
	verb := verbNameReorder
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprint(o,
		`Moves specified names to start of record, or end of record.
`)
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-e Put specified field names at record end: default is to put them at record start.\n")
	fmt.Fprintf(o, "-f {a,b,c} Field names to reorder.\n")
	fmt.Fprintf(o, "-b {x}     Put field names specified with -f before field name specified by {x},\n")
	fmt.Fprintf(o, "           if any. If {x} isn't present in a given record, the specified fields\n")
	fmt.Fprintf(o, "           will not be moved.\n")
	fmt.Fprintf(o, "-a {x}     Put field names specified with -f after field name specified by {x},\n")
	fmt.Fprintf(o, "           if any. If {x} isn't present in a given record, the specified fields\n")
	fmt.Fprintf(o, "           will not be moved.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Examples:\n")
	fmt.Fprintf(o, "%s %s    -f a,b sends input record \"d=4,b=2,a=1,c=3\" to \"a=1,b=2,d=4,c=3\".\n", argv0, verb)
	fmt.Fprintf(o, "%s %s -e -f a,b sends input record \"d=4,b=2,a=1,c=3\" to \"d=4,c=3,a=1,b=2\".\n", argv0, verb)

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerReorderParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	var fieldNames []string = nil
	putAtEnd := false
	beforeFieldName := ""
	afterFieldName := ""

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerReorderUsage(os.Stdout, true, 0)

		} else if opt == "-f" {
			fieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-b" {
			beforeFieldName = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)
			afterFieldName = ""
			putAtEnd = false

		} else if opt == "-a" {
			afterFieldName = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)
			beforeFieldName = ""
			putAtEnd = false

		} else if opt == "-e" {
			putAtEnd = true
			beforeFieldName = ""
			afterFieldName = ""

		} else {
			transformerReorderUsage(os.Stderr, true, 1)
		}
	}

	if fieldNames == nil {
		transformerReorderUsage(os.Stderr, true, 1)
	}

	transformer, err := NewTransformerReorder(
		fieldNames,
		putAtEnd,
		beforeFieldName,
		afterFieldName,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerReorder struct {
	// input
	fieldNames      []string
	fieldNamesSet   map[string]bool
	beforeFieldName string
	afterFieldName  string

	// state
	recordTransformerFunc RecordTransformerFunc
}

func NewTransformerReorder(
	fieldNames []string,
	putAtEnd bool,
	beforeFieldName string,
	afterFieldName string,
) (*TransformerReorder, error) {

	tr := &TransformerReorder{
		fieldNames:      fieldNames,
		fieldNamesSet:   lib.StringListToSet(fieldNames),
		beforeFieldName: beforeFieldName,
		afterFieldName:  afterFieldName,
	}

	if putAtEnd {
		tr.recordTransformerFunc = tr.reorderToEnd
	} else if beforeFieldName != "" {
		tr.recordTransformerFunc = tr.reorderBefore
	} else if afterFieldName != "" {
		tr.recordTransformerFunc = tr.reorderAfter
	} else {
		tr.recordTransformerFunc = tr.reorderToStart
		lib.ReverseStringList(tr.fieldNames)
	}

	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerReorder) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, inputDownstreamDoneChannel, outputDownstreamDoneChannel, outputChannel)
}

// ----------------------------------------------------------------
func (tr *TransformerReorder) reorderToStart(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		for _, fieldName := range tr.fieldNames {
			inrec.MoveToHead(fieldName)
		}
		outputChannel <- inrecAndContext

	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (tr *TransformerReorder) reorderToEnd(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		for _, fieldName := range tr.fieldNames {
			inrec.MoveToTail(fieldName)
		}
		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (tr *TransformerReorder) reorderBefore(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		if inrec.Get(tr.beforeFieldName) == nil {
			outputChannel <- inrecAndContext
			return
		}

		outrec := types.NewMlrmapAsRecord()
		pe := inrec.Head

		// * inrec will be GC'ed
		// * We will use outrec.PutReference not output.PutCopy since inrec will be GC'ed

		for ; pe != nil; pe = pe.Next {
			if pe.Key == tr.beforeFieldName {
				break
			}
			if !tr.fieldNamesSet[pe.Key] {
				outrec.PutReference(pe.Key, pe.Value)
			}
		}

		for _, fieldName := range tr.fieldNames {
			value := inrec.Get(fieldName)
			if value != nil {
				outrec.PutReference(fieldName, value)
			}
		}

		value := inrec.Get(tr.beforeFieldName)
		if value != nil {
			outrec.PutReference(tr.beforeFieldName, value)
		}

		for ; pe != nil; pe = pe.Next {
			if pe.Key != tr.beforeFieldName && !tr.fieldNamesSet[pe.Key] {
				outrec.PutReference(pe.Key, pe.Value)
			}
		}

		for _, fieldName := range tr.fieldNames {
			inrec.MoveToHead(fieldName)
		}
		outputChannel <- types.NewRecordAndContext(outrec, &inrecAndContext.Context)

	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (tr *TransformerReorder) reorderAfter(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		if inrec.Get(tr.afterFieldName) == nil {
			outputChannel <- inrecAndContext
			return
		}

		outrec := types.NewMlrmapAsRecord()
		pe := inrec.Head

		// * inrec will be GC'ed
		// * We will use outrec.PutReference not output.PutCopy since inrec will be GC'ed

		for ; pe != nil; pe = pe.Next {
			if pe.Key == tr.afterFieldName {
				break
			}
			if !tr.fieldNamesSet[pe.Key] {
				outrec.PutReference(pe.Key, pe.Value)
			}
		}

		value := inrec.Get(tr.afterFieldName)
		if value != nil {
			outrec.PutReference(tr.afterFieldName, value)
		}

		for _, fieldName := range tr.fieldNames {
			value := inrec.Get(fieldName)
			if value != nil {
				outrec.PutReference(fieldName, value)
			}
		}

		for ; pe != nil; pe = pe.Next {
			if pe.Key != tr.afterFieldName && !tr.fieldNamesSet[pe.Key] {
				outrec.PutReference(pe.Key, pe.Value)
			}
		}

		for _, fieldName := range tr.fieldNames {
			inrec.MoveToHead(fieldName)
		}
		outputChannel <- types.NewRecordAndContext(outrec, &inrecAndContext.Context)

	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
