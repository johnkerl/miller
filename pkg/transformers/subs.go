package transformers

import (
	"container/list"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/johnkerl/miller/pkg/bifs"
	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/lib"
	"github.com/johnkerl/miller/pkg/mlrval"
	"github.com/johnkerl/miller/pkg/types"
)

// ----------------------------------------------------------------
const verbNameSub = "sub"
const verbNameGsub = "gsub"
const verbNameSsub = "ssub"

var SubSetup = TransformerSetup{
	Verb:         verbNameSub,
	UsageFunc:    transformerSubUsage,
	ParseCLIFunc: transformerSubParseCLI,
	IgnoresInput: false,
}

var GsubSetup = TransformerSetup{
	Verb:         verbNameGsub,
	UsageFunc:    transformerGsubUsage,
	ParseCLIFunc: transformerGsubParseCLI,
	IgnoresInput: false,
}

var SsubSetup = TransformerSetup{
	Verb:         verbNameSsub,
	UsageFunc:    transformerSsubUsage,
	ParseCLIFunc: transformerSsubParseCLI,
	IgnoresInput: false,
}

func transformerSubUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameSub)
	fmt.Fprintf(o, "Replaces old string with new string in specified field(s), with regex support\n")
	fmt.Fprintf(o, "for the old string and not handling multiple matches, like the `sub` DSL function.\n")
	fmt.Fprintf(o, "See also the `gsub` and `ssub` verbs.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-f {a,b,c}  Field names to convert.\n")
	fmt.Fprintf(o, "-h|--help   Show this message.\n")
}

func transformerGsubUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameGsub)
	fmt.Fprintf(o, "Replaces old string with new string in specified field(s), with regex support\n")
	fmt.Fprintf(o, "for the old string and handling multiple matches, like the `gsub` DSL function.\n")
	fmt.Fprintf(o, "See also the `sub` and `ssub` verbs.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-f {a,b,c}  Field names to convert.\n")
	fmt.Fprintf(o, "-h|--help   Show this message.\n")
}

func transformerSsubUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameSsub)
	fmt.Fprintf(o, "Replaces old string with new string in specified field(s), without regex support for\n")
	fmt.Fprintf(o, "the old string, like the `ssub` DSL function. See also the `gsub` and `sub` verbs.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-f {a,b,c}  Field names to convert.\n")
	fmt.Fprintf(o, "-h|--help   Show this message.\n")
}

type subConstructorFunc func(
	fieldNames []string,
	doAllFieldNames bool,
	doRegexes bool,
	oldText string,
	newText string,
) (IRecordTransformer, error)

type fieldAcceptorFunc func(
	fieldName string,
) bool

func transformerSubParseCLI(
	pargi *int,
	argc int,
	args []string,
	opts *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) IRecordTransformer {
	return transformerSubsParseCLI(pargi, argc, args, opts, doConstruct, transformerSubUsage, NewTransformerSub)
}

func transformerGsubParseCLI(
	pargi *int,
	argc int,
	args []string,
	opts *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) IRecordTransformer {
	return transformerSubsParseCLI(pargi, argc, args, opts, doConstruct, transformerGsubUsage, NewTransformerGsub)
}

func transformerSsubParseCLI(
	pargi *int,
	argc int,
	args []string,
	opts *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) IRecordTransformer {
	return transformerSubsParseCLI(pargi, argc, args, opts, doConstruct, transformerSsubUsage, NewTransformerSsub)
}

// transformerSubsParseCLI is a shared CLI-parser for the sub, gsub, and ssub verbs.
func transformerSubsParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
	usageFunc TransformerUsageFunc,
	constructorFunc subConstructorFunc,
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// Parse local flags
	var fieldNames []string = nil
	doAllFieldNames := false
	doRegexes := false
	var oldText string
	var newText string

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
			usageFunc(os.Stdout)
			os.Exit(0)

		} else if opt == "-a" {
			doAllFieldNames = true
			doRegexes = false
			fieldNames = nil

		} else if opt == "-r" {
			doRegexes = true

		} else if opt == "-f" {
			fieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)
			doAllFieldNames = false
		} else {
			usageFunc(os.Stderr)
			os.Exit(1)
		}
	}

	if fieldNames == nil && !doAllFieldNames {
		usageFunc(os.Stderr)
		os.Exit(1)
	}

	// Get the old and new text from the command line
	if (argc - argi) < 2 {
		usageFunc(os.Stderr)
		os.Exit(1)
	}
	oldText = args[argi]
	newText = args[argi+1]

	argi += 2

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := constructorFunc(
		fieldNames,
		doAllFieldNames,
		doRegexes,
		oldText,
		newText,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

type TransformerSubs struct {
	fieldNamesSet map[string]bool  // for -f
	regexes       []*regexp.Regexp // for -r
	oldText       *mlrval.Mlrval
	newText       *mlrval.Mlrval
	fieldAcceptor fieldAcceptorFunc // for -f, -r, -a
	subber        bifs.TernaryFunc  // for sub, gsub, ssub
}

func NewTransformerSub(
	fieldNames []string,
	doAllFieldNames bool,
	doRegexes bool,
	oldText string,
	newText string,
) (IRecordTransformer, error) {
	return NewTransformerSubs(fieldNames, doAllFieldNames, doRegexes, oldText, newText, safe_sub)
}

func NewTransformerGsub(
	fieldNames []string,
	doAllFieldNames bool,
	doRegexes bool,
	oldText string,
	newText string,
) (IRecordTransformer, error) {
	return NewTransformerSubs(fieldNames, doAllFieldNames, doRegexes, oldText, newText, safe_gsub)
}

func NewTransformerSsub(
	fieldNames []string,
	doAllFieldNames bool,
	doRegexes bool,
	oldText string,
	newText string,
) (IRecordTransformer, error) {
	return NewTransformerSubs(fieldNames, doAllFieldNames, doRegexes, oldText, newText, safe_ssub)
}

func NewTransformerSubs(
	fieldNames []string,
	doAllFieldNames bool,
	doRegexes bool,
	oldText string,
	newText string,
	subber bifs.TernaryFunc,
) (IRecordTransformer, error) {
	tr := &TransformerSubs{
		fieldNamesSet: lib.StringListToSet(fieldNames),
		oldText:       mlrval.FromString(oldText),
		newText:       mlrval.FromString(newText),
		subber:        subber,
	}
	if doAllFieldNames {
		tr.fieldAcceptor = tr.fieldAcceptorAll
	} else if doRegexes {
		tr.fieldAcceptor = tr.fieldAcceptorByRegexes

		tr.regexes = make([]*regexp.Regexp, len(fieldNames))
		for i, regexString := range fieldNames {
			// Handles "a.*b"i Miller case-insensitive-regex specification
			regex, err := lib.CompileMillerRegex(regexString)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s %s: cannot compile regex [%s]\n", "mlr", verbNameCut, regexString)
				os.Exit(1)
			}
			tr.regexes[i] = regex
		}
	} else {
		tr.fieldAcceptor = tr.fieldAcceptorByNames
	}
	return tr, nil
}

func (tr *TransformerSubs) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)

	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		// Run sub, gsub, or ssub on the user-specified field names
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if tr.fieldAcceptor(pe.Key) {
				pe.Value = tr.subber(pe.Value, tr.oldText, tr.newText)
			}
		}
	}
	// Including emit of end-of-stream marker
	outputRecordsAndContexts.PushBack(inrecAndContext)
}

// fieldAcceptorByNames implements -f
func (tr *TransformerSubs) fieldAcceptorByNames(
	fieldName string,
) bool {
	return tr.fieldNamesSet[fieldName]
}

// fieldAcceptorByNames implements -r
func (tr *TransformerSubs) fieldAcceptorByRegexes(
	fieldName string,
) bool {
	for _, regex := range tr.regexes {
		if regex.MatchString(fieldName) {
			return true
		}
	}
	return false
}

// fieldAcceptorByNames implements -a
func (tr *TransformerSubs) fieldAcceptorAll(
	fieldName string,
) bool {
	return true
}

// safe_sub implements sub, but doesn't produce error-type on non-string input.
func safe_sub(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsString() {
		return bifs.BIF_sub(input1, input2, input3)
	} else {
		return input1
	}
}

// safe_gsub implements gsub, but doesn't produce error-type on non-string input.
func safe_gsub(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsString() {
		return bifs.BIF_gsub(input1, input2, input3)
	} else {
		return input1
	}
}

// safe_ssub implements ssub, but doesn't produce error-type on non-string input.
func safe_ssub(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsString() {
		return bifs.BIF_ssub(input1, input2, input3)
	} else {
		return input1
	}
}
