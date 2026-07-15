package transformers

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

const verbNameReorder = "reorder"

var reorderOptions = []OptionSpec{
	{Flag: "-e", Type: "bool", Desc: "Put specified field names at record end: default is to put them at record start."},
	{Flag: "-f", Arg: "{a,b,c}", Type: "csv-list", Desc: "Field names to reorder."},
	{Flag: "-r", Arg: "{a,b,c}", Type: "csv-list", Desc: "Treat field names as regular expressions. Matched fields are moved to start or end, grouped by the order the regexes are given; within each group, fields keep their record order. Example: -r '^YYY,^XXX' puts all YYY-prefixed fields first, then all XXX-prefixed fields, then the rest."},
	{Flag: "-b", Arg: "{x}", Type: "string", Desc: "Put field names specified with -f before field name specified by {x}, if any. If {x} isn't present in a given record, the specified fields will not be moved."},
	{Flag: "-a", Arg: "{x}", Type: "string", Desc: "Put field names specified with -f after field name specified by {x}, if any. If {x} isn't present in a given record, the specified fields will not be moved."},
}

var ReorderSetup = TransformerSetup{
	Verb:         verbNameReorder,
	UsageFunc:    transformerReorderUsage,
	ParseCLIFunc: transformerReorderParseCLI,
	IgnoresInput: false,
	Options:      reorderOptions,
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
	WriteVerbOptions(o, reorderOptions)
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Examples:\n")
	fmt.Fprintf(o, "%s %s    -f a,b sends input record \"d=4,b=2,a=1,c=3\" to \"a=1,b=2,d=4,c=3\".\n", argv0, verb)
	fmt.Fprintf(o, "%s %s -e -f a,b sends input record \"d=4,b=2,a=1,c=3\" to \"d=4,c=3,a=1,b=2\".\n", argv0, verb)
	fmt.Fprintf(o, "%s %s -r '^YYY,^XXX' puts YYY-prefixed fields first, then XXX-prefixed fields, then rest.\n", argv0, verb)
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

		switch opt {
		case "-h", "--help":
			transformerReorderUsage(os.Stdout)
			return nil, cli.ErrHelpRequested

		case "-f":
			fieldNames, err = cli.VerbGetStringArrayArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			doRegexes = false

		case "-r":
			fieldNames, err = cli.VerbGetStringArrayArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			doRegexes = true

		case "-b":
			centerFieldName, err = cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			putAfter = false

		case "-a":
			centerFieldName, err = cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			putAfter = true

		case "-e":
			putAfter = true
			centerFieldName = ""

		default:
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
				slices.Reverse(tr.fieldNames)
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
				return nil, cli.VerbErrorf(verbNameReorder, "cannot compile regex [%s]", regexString)
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
) error {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		tr.recordTransformerFunc(
			inrecAndContext,
			outputRecordsAndContexts,
		)
	} else {
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // end-of-stream marker
	}
	return nil
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

// reorderSplitByRegex splits record fields into matching (any regex) and rest. Matching fields
// are grouped by the order the regexes were given on the command line; within each group, and
// within rest, fields keep their record order. A field matching multiple regexes is claimed by
// the first regex that matches it.
func (tr *TransformerReorder) reorderSplitByRegex(inrec *mlrval.Mlrmap) (matching []*mlrval.MlrmapEntry, rest []*mlrval.MlrmapEntry) {
	claimed := make(map[string]bool)
	for _, regex := range tr.regexes {
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if !claimed[pe.Key] && regex.MatchString(pe.Key) {
				matching = append(matching, pe)
				claimed[pe.Key] = true
			}
		}
	}
	for pe := inrec.Head; pe != nil; pe = pe.Next {
		if !claimed[pe.Key] {
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

	// Build matching set grouped by regex order (OrderedMap preserves insertion order);
	// within each group, fields keep their record order.
	matchingFieldNamesSet := lib.NewOrderedMap[*mlrval.Mlrval]()
	for _, regex := range tr.regexes {
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if pe.Key != tr.centerFieldName && !matchingFieldNamesSet.Has(pe.Key) && regex.MatchString(pe.Key) {
				matchingFieldNamesSet.Put(pe.Key, pe.Value)
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
