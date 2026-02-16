package transformers

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

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
	fmt.Fprintf(o, "-r        Treat field names as regular expressions. Matched fields are moved\n")
	fmt.Fprintf(o, "          to start or end in record order. Example: -r '^YYY,^XXX' puts all\n")
	fmt.Fprintf(o, "          YYY- and XXX-prefixed fields first (in record order), then the rest.\n")
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
	fmt.Fprintf(o, "%s %s -r '^YYY,^XXX' puts YYY- and XXX-prefixed fields first (record order), then rest.\n", argv0, verb)
}

func transformerReorderParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) (RecordTransformer, error) {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	var fieldNames []string = nil
	doRegexes := false
	putAfter := false
	centerFieldName := ""

	var err error
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
			return nil, cli.ErrHelpRequested

		} else if opt == "-f" {
			fieldNames, err = cli.VerbGetStringArrayArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			doRegexes = false

		} else if opt == "-r" {
			fieldNames, err = cli.VerbGetStringArrayArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			doRegexes = true

		} else if opt == "-b" {
			centerFieldName, err = cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			putAfter = false

		} else if opt == "-a" {
			centerFieldName, err = cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			putAfter = true

		} else if opt == "-e" {
			putAfter = true
			centerFieldName = ""

		} else {
			return nil, cli.VerbErrorf(verb, "option \"%s\" not recognized", opt)
		}
	}

	if fieldNames == nil {
		return nil, cli.VerbErrorf(verb, "-f field names required")
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	transformer, err := NewTransformerReorder(
		fieldNames,
		doRegexes,
		putAfter,
		centerFieldName,
	)
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

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
					"mlr", verbNameReorder, regexString,
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
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
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
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // end-of-stream marker
	}
}

func (tr *TransformerReorder) reorderToStartNoRegex(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
) {
	inrec := inrecAndContext.Record
	for _, fieldName := range tr.fieldNames {
		inrec.MoveToHead(fieldName)
	}
	*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
}

// reorderSplitByRegex splits record fields into matching (any regex) and rest, preserving record order.
func (tr *TransformerReorder) reorderSplitByRegex(inrec *mlrval.Mlrmap) (matching []*mlrval.MlrmapEntry, rest []*mlrval.MlrmapEntry) {
	for pe := inrec.Head; pe != nil; pe = pe.Next {
		found := false
		for _, regex := range tr.regexes {
			if regex.MatchString(pe.Key) {
				matching = append(matching, pe)
				found = true
				break
			}
		}
		if !found {
			rest = append(rest, pe)
		}
	}
	return matching, rest
}

func (tr *TransformerReorder) reorderToStartWithRegex(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
) {
	inrec := inrecAndContext.Record
	matching, rest := tr.reorderSplitByRegex(inrec)

	outrec := mlrval.NewMlrmapAsRecord()
	for _, pe := range matching {
		outrec.PutReference(pe.Key, pe.Value)
	}
	for _, pe := range rest {
		outrec.PutReference(pe.Key, pe.Value)
	}

	outrecAndContext := types.NewRecordAndContext(outrec, &inrecAndContext.Context)
	*outputRecordsAndContexts = append(*outputRecordsAndContexts, outrecAndContext)
}

func (tr *TransformerReorder) reorderToEndNoRegex(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
) {
	inrec := inrecAndContext.Record
	for _, fieldName := range tr.fieldNames {
		inrec.MoveToTail(fieldName)
	}
	*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)

}

func (tr *TransformerReorder) reorderToEndWithRegex(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
) {
	inrec := inrecAndContext.Record
	matching, rest := tr.reorderSplitByRegex(inrec)

	outrec := mlrval.NewMlrmapAsRecord()
	for _, pe := range rest {
		outrec.PutReference(pe.Key, pe.Value)
	}
	for _, pe := range matching {
		outrec.PutReference(pe.Key, pe.Value)
	}

	outrecAndContext := types.NewRecordAndContext(outrec, &inrecAndContext.Context)
	*outputRecordsAndContexts = append(*outputRecordsAndContexts, outrecAndContext)
}

func (tr *TransformerReorder) reorderBeforeOrAfterNoRegex(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
) {
	inrec := inrecAndContext.Record
	if inrec.Get(tr.centerFieldName) == nil {
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
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

	*outputRecordsAndContexts = append(*outputRecordsAndContexts, types.NewRecordAndContext(outrec, &inrecAndContext.Context))

}

func (tr *TransformerReorder) reorderBeforeOrAfterWithRegex(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
) {
	inrec := inrecAndContext.Record
	if inrec.Get(tr.centerFieldName) == nil {
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
		return
	}

	// Build matching set in record order (OrderedMap preserves insertion order)
	matchingFieldNamesSet := lib.NewOrderedMap[*mlrval.Mlrval]()
	for pe := inrec.Head; pe != nil; pe = pe.Next {
		for _, regex := range tr.regexes {
			if regex.MatchString(pe.Key) && pe.Key != tr.centerFieldName {
				matchingFieldNamesSet.Put(pe.Key, pe.Value)
				break
			}
		}
	}

	outrec := mlrval.NewMlrmapAsRecord()
	for pe := inrec.Head; pe != nil; pe = pe.Next {
		if pe.Key == tr.centerFieldName {
			if tr.putAfter {
				outrec.PutReference(pe.Key, pe.Value)
			}
			for pf := matchingFieldNamesSet.Head; pf != nil; pf = pf.Next {
				outrec.PutReference(pf.Key, pf.Value)
			}
			if !tr.putAfter {
				outrec.PutReference(pe.Key, pe.Value)
			}
		} else if !matchingFieldNamesSet.Has(pe.Key) {
			outrec.PutReference(pe.Key, pe.Value)
		}
	}

	*outputRecordsAndContexts = append(*outputRecordsAndContexts, types.NewRecordAndContext(outrec, &inrecAndContext.Context))
}
