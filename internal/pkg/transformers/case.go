package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
	"github.com/johnkerl/miller/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameCase = "case"

var CaseSetup = TransformerSetup{
	Verb:         verbNameCase,
	UsageFunc:    transformerCaseUsage,
	ParseCLIFunc: transformerCaseParseCLI,
	IgnoresInput: false,
}

const (
	e_UNSPECIFIED_CASE = iota
	e_UPPER_CASE
	e_LOWER_CASE
	e_SENTENCE_CASE
	e_TITLE_CASE
)

func transformerCaseUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameCase)
	fmt.Fprintf(o, "Uppercases strings in record keys and/or values.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-k  Case only keys, not keys and values.\n")
	fmt.Fprintf(o, "-v  Case only values, not keys and values.\n")
	fmt.Fprintf(o, "-f  {a,b,c} Specify which field names to case (default: all)\n")
	fmt.Fprintf(o, "-u  Convert to uppercase\n")
	fmt.Fprintf(o, "-l  Convert to lowercase\n")
	fmt.Fprintf(o, "-s  Convert to sentence case (capitalize first letter)\n")
	fmt.Fprintf(o, "-t  Convert to title case (capitalize words)\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
}

func transformerCaseParseCLI(
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

	which := "keys_and_values"
	style := e_UNSPECIFIED_CASE
	var fieldNames []string = nil

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
			transformerCaseUsage(os.Stdout)
			os.Exit(0)

		} else if opt == "-k" {
			which = "keys_only"

		} else if opt == "-v" {
			which = "values_only"

		} else if opt == "-f" {
			fieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-u" {
			style = e_UPPER_CASE
		} else if opt == "-l" {
			style = e_LOWER_CASE
		} else if opt == "-s" {
			style = e_SENTENCE_CASE
		} else if opt == "-t" {
			style = e_TITLE_CASE

		} else {
			transformerCaseUsage(os.Stderr)
			os.Exit(1)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerCase(which, fieldNames, style)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type caserFuncT func(input string) string

type TransformerCase struct {
	recordTransformerFunc RecordTransformerFunc
	fieldNameSet          map[string]bool
	caserFunc             caserFuncT
}

func NewTransformerCase(
	which string,
	fieldNames []string,
	style int,
) (*TransformerCase, error) {
	tr := &TransformerCase{}

	if which == "keys_only" {
		tr.recordTransformerFunc = tr.transformKeysOnly
	} else if which == "values_only" {
		tr.recordTransformerFunc = tr.transformValuesOnly
	} else {
		tr.recordTransformerFunc = tr.transformKeysAndValues
	}

	if fieldNames != nil {
		tr.fieldNameSet = lib.StringListToSet(fieldNames)
	}

	switch style {
	case e_UPPER_CASE:
		tr.caserFunc = cases.Upper(language.Und).String
	case e_LOWER_CASE:
		tr.caserFunc = cases.Lower(language.Und).String
	case e_SENTENCE_CASE:
		tr.caserFunc = caseSentenceFunc
	case e_TITLE_CASE:
		tr.caserFunc = cases.Title(language.Und).String
	default:
		return nil, fmt.Errorf(
			"mlr %s: case option must be specified using one of -u, -l, -s, -t.",
			verbNameCase,
		)
	}

	return tr, nil
}

func caseSentenceFunc(input string) string {
	runes := []rune(input)
	if len(runes) == 0 {
		return input
	}
	first := string(runes[0])
	rest := string(runes[1:])
	return strings.ToUpper(first) + strings.ToLower(rest)
}

func (tr *TransformerCase) Transform(
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
			inputDownstreamDoneChannel,
			outputDownstreamDoneChannel,
		)
	} else { // end of record stream
		outputRecordsAndContexts.PushBack(inrecAndContext)
	}
}

func (tr *TransformerCase) transformKeysOnly(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	_ <-chan bool,
	__ chan<- bool,
) {
	inrec := inrecAndContext.Record
	newrec := mlrval.NewMlrmapAsRecord()
	for pe := inrec.Head; pe != nil; pe = pe.Next {
		if tr.fieldNameSet == nil || tr.fieldNameSet[pe.Key] {
			newkey := tr.caserFunc(pe.Key)
			// Reference not copy since this is ownership transfer of the value from the now-abandoned inrec
			newrec.PutReference(newkey, pe.Value)
		} else {
			newrec.PutReference(pe.Key, pe.Value)
		}
	}
	outputRecordsAndContexts.PushBack(types.NewRecordAndContext(newrec, &inrecAndContext.Context))
}

func (tr *TransformerCase) transformValuesOnly(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	_ <-chan bool,
	__ chan<- bool,
) {
	inrec := inrecAndContext.Record
	for pe := inrec.Head; pe != nil; pe = pe.Next {
		if tr.fieldNameSet == nil || tr.fieldNameSet[pe.Key] {
			stringval, ok := pe.Value.GetStringValue()
			if ok {
				pe.Value = mlrval.FromString(tr.caserFunc(stringval))
			}
		}
	}
	outputRecordsAndContexts.PushBack(types.NewRecordAndContext(inrec, &inrecAndContext.Context))
}

func (tr *TransformerCase) transformKeysAndValues(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	_ <-chan bool,
	__ chan<- bool,
) {
	inrec := inrecAndContext.Record
	newrec := mlrval.NewMlrmapAsRecord()
	for pe := inrec.Head; pe != nil; pe = pe.Next {
		if tr.fieldNameSet == nil || tr.fieldNameSet[pe.Key] {
			newkey := tr.caserFunc(pe.Key)
			stringval, ok := pe.Value.GetStringValue()
			if ok {
				stringval = tr.caserFunc(stringval)
				newrec.PutReference(newkey, mlrval.FromString(stringval))
			} else {
				newrec.PutReference(newkey, pe.Value)
			}
		} else {
			newrec.PutReference(pe.Key, pe.Value)
		}
	}
	outputRecordsAndContexts.PushBack(types.NewRecordAndContext(newrec, &inrecAndContext.Context))
}
