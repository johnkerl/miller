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
	putAfter := false
	centerFieldName := ""

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
			doRegexes = false

		} else if opt == "-r" {
			fieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)
			doRegexes = true

		} else if opt == "-b" {
			centerFieldName = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)
			putAfter = false

		} else if opt == "-a" {
			centerFieldName = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)
			putAfter = true

		} else if opt == "-e" {
			putAfter = true
			centerFieldName = ""

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
		putAfter,
		centerFieldName,
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
	centerFieldName string
	putAfter        bool

	// state
	recordTransformerFunc RecordTransformerHelperFunc
}

func NewTransformerReorder(
	fieldNames []string,
	doRegexes bool,
	putAfter bool,
	centerFieldName string,
) (*TransformerReorder, error) {

	tr := &TransformerReorder{
		fieldNames:      fieldNames,
		fieldNamesSet:   lib.StringListToSet(fieldNames),
		centerFieldName: centerFieldName,
		putAfter:        putAfter,
	}

	if centerFieldName == "" {
		if putAfter {
			if doRegexes {
				tr.recordTransformerFunc = tr.reorderToEndWithRegex
			} else {
				tr.recordTransformerFunc = tr.reorderToEndNoRegex
			}
		} else {
			if doRegexes {
				tr.recordTransformerFunc = tr.reorderToStartWithRegex
			} else {
				tr.recordTransformerFunc = tr.reorderToStartNoRegex
				lib.ReverseStringList(tr.fieldNames)
			}
		}
	} else {
		if doRegexes {
			tr.recordTransformerFunc = tr.reorderBeforeOrAfterWithRegex
		} else {
			tr.recordTransformerFunc = tr.reorderBeforeOrAfterNoRegex
		}
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

func (tr *TransformerReorder) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		tr.recordTransformerFunc(
			inrecAndContext,
			outputRecordsAndContexts,
		)
	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext) // end-of-stream marker
	}
}

func (tr *TransformerReorder) reorderToStartNoRegex(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
) {
	inrec := inrecAndContext.Record
	for _, fieldName := range tr.fieldNames {
		inrec.MoveToHead(fieldName)
	}
	outputRecordsAndContexts.PushBack(inrecAndContext)
}

func (tr *TransformerReorder) reorderToStartWithRegex(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
) {
	inrec := inrecAndContext.Record

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

func (tr *TransformerReorder) reorderToEndNoRegex(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
) {
	inrec := inrecAndContext.Record
	for _, fieldName := range tr.fieldNames {
		inrec.MoveToTail(fieldName)
	}
	outputRecordsAndContexts.PushBack(inrecAndContext)

}

func (tr *TransformerReorder) reorderToEndWithRegex(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
) {
	inrec := inrecAndContext.Record
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

func (tr *TransformerReorder) reorderBeforeOrAfterNoRegex(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
) {
	inrec := inrecAndContext.Record
	if inrec.Get(tr.centerFieldName) == nil {
		outputRecordsAndContexts.PushBack(inrecAndContext)
		return
	}

	outrec := mlrval.NewMlrmapAsRecord()
	pe := inrec.Head

	// We use outrec.PutReference not output.PutCopy since inrec will be GC'ed

	for ; pe != nil; pe = pe.Next {
		if pe.Key == tr.centerFieldName {
			break
		}
		if !tr.fieldNamesSet[pe.Key] {
			outrec.PutReference(pe.Key, pe.Value)
		}
	}

	if !tr.putAfter {
		for _, fieldName := range tr.fieldNames {
			value := inrec.Get(fieldName)
			if value != nil {
				outrec.PutReference(fieldName, value)
			}
		}
	}

	value := inrec.Get(tr.centerFieldName)
	if value != nil {
		outrec.PutReference(tr.centerFieldName, value)
	}

	if tr.putAfter {
		for _, fieldName := range tr.fieldNames {
			value := inrec.Get(fieldName)
			if value != nil {
				outrec.PutReference(fieldName, value)
			}
		}
	}

	for ; pe != nil; pe = pe.Next {
		if pe.Key != tr.centerFieldName && !tr.fieldNamesSet[pe.Key] {
			outrec.PutReference(pe.Key, pe.Value)
		}
	}

	outputRecordsAndContexts.PushBack(types.NewRecordAndContext(outrec, &inrecAndContext.Context))

}

func (tr *TransformerReorder) reorderBeforeOrAfterWithRegex(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
) {
	inrec := inrecAndContext.Record
	if inrec.Get(tr.centerFieldName) == nil {
		outputRecordsAndContexts.PushBack(inrecAndContext)
		return
	}

	matchingFieldNamesSet := lib.NewOrderedMap()
	for pe := inrec.Head; pe != nil; pe = pe.Next {
		for _, regex := range tr.regexes {
			if regex.MatchString(pe.Key) {
				if pe.Key != tr.centerFieldName {
					matchingFieldNamesSet.Put(pe.Key, pe.Value)
					break
				}
			}
		}
	}

	// We use outrec.PutReference not output.PutCopy since inrec will be GC'ed
	outrec := mlrval.NewMlrmapAsRecord()
	for pe := inrec.Head; pe != nil; pe = pe.Next {
		if pe.Key == tr.centerFieldName {
			if tr.putAfter {
				outrec.PutReference(pe.Key, pe.Value)
			}
			for pf := matchingFieldNamesSet.Head; pf != nil; pf = pf.Next {
				outrec.PutReference(pf.Key, pf.Value.(*mlrval.Mlrval))
			}
			if !tr.putAfter {
				outrec.PutReference(pe.Key, pe.Value)
			}
		} else if !matchingFieldNamesSet.Has(pe.Key) {
			outrec.PutReference(pe.Key, pe.Value)
		}
	}

	outputRecordsAndContexts.PushBack(types.NewRecordAndContext(outrec, &inrecAndContext.Context))
}
