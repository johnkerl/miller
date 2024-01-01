package transformers

import (
	"container/list"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/lib"
	"github.com/johnkerl/miller/pkg/mlrval"
	"github.com/johnkerl/miller/pkg/types"
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
}

func transformerReorderParseCLI(
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

	var fieldNames []string = nil
	doRegexes := false
	putAtEnd := false
	beforeFieldName := ""
	afterFieldName := ""

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
			transformerReorderUsage(os.Stdout)
			os.Exit(0)

		} else if opt == "-f" {
			fieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-r" {
			fieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)
			doRegexes = true

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
			transformerReorderUsage(os.Stderr)
			os.Exit(1)
		}
	}

	if fieldNames == nil {
		transformerReorderUsage(os.Stderr)
		os.Exit(1)
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerReorder(
		fieldNames,
		doRegexes,
		putAtEnd,
		beforeFieldName,
		afterFieldName,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerReorder struct {
	// input
	fieldNames      []string
	fieldNamesSet   map[string]bool
	regexes         []*regexp.Regexp
	beforeFieldName string
	afterFieldName  string

	// state
	recordTransformerFunc RecordTransformerFunc
}

func NewTransformerReorder(
	fieldNames []string,
	doRegexes bool,
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

	if doRegexes {
		tr.regexes = make([]*regexp.Regexp, len(fieldNames))
		for i, regexString := range fieldNames {
			// Handles "a.*b"i Miller case-insensitive-regex specification
			regex, err := lib.CompileMillerRegex(regexString)
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"%s %s: cannot compile regex [%s]\n",
					"mlr", verbNameCut, regexString,
				)
				os.Exit(1)
			}
			tr.regexes[i] = regex
		}
	}

	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerReorder) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, outputRecordsAndContexts, inputDownstreamDoneChannel, outputDownstreamDoneChannel)
}

// ----------------------------------------------------------------
func (tr *TransformerReorder) reorderToStart(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		if tr.regexes == nil {
			for _, fieldName := range tr.fieldNames {
				inrec.MoveToHead(fieldName)
			}
			outputRecordsAndContexts.PushBack(inrecAndContext)

		} else {
			outrec := mlrval.NewMlrmapAsRecord()
			atEnds := list.New()
			for pe := inrec.Head; pe != nil; pe = pe.Next {
				found := false
				for _, regex := range tr.regexes {
					if regex.MatchString(pe.Key) {
						outrec.PutReference(pe.Key, pe.Value)
						found = true
						break
					}
				}
				if !found {
					atEnds.PushBack(pe)
				}
			}

			for atEnd := atEnds.Front(); atEnd != nil; atEnd = atEnd.Next() {
				// Ownership transfer; no copy needed
				pe := atEnd.Value.(*mlrval.MlrmapEntry)
				outrec.PutReference(pe.Key, pe.Value)
			}

			outrecAndContext := types.NewRecordAndContext(outrec, &inrecAndContext.Context)
			outputRecordsAndContexts.PushBack(outrecAndContext)
		}

	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext) // end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (tr *TransformerReorder) reorderToEnd(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		if tr.regexes == nil {
			for _, fieldName := range tr.fieldNames {
				inrec.MoveToTail(fieldName)
			}
			outputRecordsAndContexts.PushBack(inrecAndContext)

		} else {
			outrec := mlrval.NewMlrmapAsRecord()
			atEnds := list.New()
			for pe := inrec.Head; pe != nil; pe = pe.Next {
				found := false
				for _, regex := range tr.regexes {
					if regex.MatchString(pe.Key) {
						atEnds.PushBack(pe)
						found = true
						break
					}
				}
				if !found {
					outrec.PutReference(pe.Key, pe.Value)
				}
			}

			for atEnd := atEnds.Front(); atEnd != nil; atEnd = atEnd.Next() {
				// Ownership transfer; no copy needed
				pe := atEnd.Value.(*mlrval.MlrmapEntry)
				outrec.PutReference(pe.Key, pe.Value)
			}

			outrecAndContext := types.NewRecordAndContext(outrec, &inrecAndContext.Context)
			outputRecordsAndContexts.PushBack(outrecAndContext)
		}
	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext) // end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (tr *TransformerReorder) reorderBefore(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		if inrec.Get(tr.beforeFieldName) == nil {
			outputRecordsAndContexts.PushBack(inrecAndContext)
			return
		}

		outrec := mlrval.NewMlrmapAsRecord()
		pe := inrec.Head

		// * inrec will be GC'ed
		// * We will use outrec.PutReference not output.PutCopy since inrec will be GC'ed

		for ; pe != nil; pe = pe.Next {
			if pe.Key == tr.beforeFieldName {
				break
			}
			if tr.regexes == nil {
				if !tr.fieldNamesSet[pe.Key] {
					outrec.PutReference(pe.Key, pe.Value)
				}
			} else {
				// XXX TO DO
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
			if tr.regexes == nil {
				if pe.Key != tr.beforeFieldName && !tr.fieldNamesSet[pe.Key] {
					outrec.PutReference(pe.Key, pe.Value)
				}
			} else {
				// XXX TO DO
			}
		}

		outputRecordsAndContexts.PushBack(types.NewRecordAndContext(outrec, &inrecAndContext.Context))

	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext) // end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (tr *TransformerReorder) reorderAfter(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		if inrec.Get(tr.afterFieldName) == nil {
			outputRecordsAndContexts.PushBack(inrecAndContext)
			return
		}

		outrec := mlrval.NewMlrmapAsRecord()
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

		outputRecordsAndContexts.PushBack(types.NewRecordAndContext(outrec, &inrecAndContext.Context))

	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext) // end-of-stream marker
	}
}
