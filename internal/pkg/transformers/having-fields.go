package transformers

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/types"
)

type tHavingFieldsCriterion int

const (
	havingFieldsCriterionUnspecified tHavingFieldsCriterion = iota
	havingFieldsAtLeast
	havingFieldsWhichAre
	havingFieldsAtMost
	havingAllFieldsMatching
	havingAnyFieldsMatching
	havingNoFieldsMatching
)

// ----------------------------------------------------------------
const verbNameHavingFields = "having-fields"

var HavingFieldsSetup = TransformerSetup{
	Verb:         verbNameHavingFields,
	UsageFunc:    transformerHavingFieldsUsage,
	ParseCLIFunc: transformerHavingFieldsParseCLI,

	IgnoresInput: false,
}

func transformerHavingFieldsUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	exeName := "mlr"
	verb := verbNameHavingFields
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameHavingFields)

	fmt.Fprintf(o, "Conditionally passes through records depending on each record's field names.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "  --at-least      {comma-separated names}\n")
	fmt.Fprintf(o, "  --which-are     {comma-separated names}\n")
	fmt.Fprintf(o, "  --at-most       {comma-separated names}\n")
	fmt.Fprintf(o, "  --all-matching  {regular expression}\n")
	fmt.Fprintf(o, "  --any-matching  {regular expression}\n")
	fmt.Fprintf(o, "  --none-matching {regular expression}\n")
	fmt.Fprintf(o, "Examples:\n")
	fmt.Fprintf(o, "  %s %s --which-are amount,status,owner\n", exeName, verb)
	fmt.Fprintf(o, "  %s %s --any-matching 'sda[0-9]'\n", exeName, verb)
	fmt.Fprintf(o, "  %s %s --any-matching '\"sda[0-9]\"'\n", exeName, verb)
	fmt.Fprintf(o, "  %s %s --any-matching '\"sda[0-9]\"i' (this is case-insensitive)\n", exeName, verb)

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerHavingFieldsParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
) IRecordTransformer {

	havingFieldsCriterion := havingFieldsCriterionUnspecified
	var fieldNames []string = nil
	regexString := ""

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerHavingFieldsUsage(os.Stdout, true, 0)

		} else if opt == "--at-least" {
			havingFieldsCriterion = havingFieldsAtLeast
			fieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)
			regexString = ""

		} else if opt == "--which-are" {
			havingFieldsCriterion = havingFieldsWhichAre
			fieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)
			regexString = ""

		} else if opt == "--at-most" {
			havingFieldsCriterion = havingFieldsAtMost
			fieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)
			regexString = ""

		} else if opt == "--all-matching" {
			havingFieldsCriterion = havingAllFieldsMatching
			regexString = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)
			fieldNames = nil

		} else if opt == "--any-matching" {
			havingFieldsCriterion = havingAnyFieldsMatching
			regexString = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)
			fieldNames = nil

		} else if opt == "--none-matching" {
			havingFieldsCriterion = havingNoFieldsMatching
			regexString = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)
			fieldNames = nil

		} else {
			transformerHavingFieldsUsage(os.Stderr, true, 1)
		}
	}

	if havingFieldsCriterion == havingFieldsCriterionUnspecified {
		transformerHavingFieldsUsage(os.Stderr, true, 1)
	}
	if fieldNames == nil && regexString == "" {
		transformerHavingFieldsUsage(os.Stderr, true, 1)
	}

	transformer, err := NewTransformerHavingFields(
		havingFieldsCriterion,
		fieldNames,
		regexString,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerHavingFields struct {
	fieldNames    []string
	numFieldNames int
	fieldNameSet  map[string]bool

	regex *regexp.Regexp

	recordTransformerFunc RecordTransformerFunc
}

// ----------------------------------------------------------------
func NewTransformerHavingFields(
	havingFieldsCriterion tHavingFieldsCriterion,
	fieldNames []string,
	regexString string,
) (*TransformerHavingFields, error) {

	tr := &TransformerHavingFields{}

	if fieldNames != nil {
		tr.fieldNames = fieldNames
		tr.numFieldNames = len(fieldNames)
		tr.fieldNameSet = lib.StringListToSet(fieldNames)

		if havingFieldsCriterion == havingFieldsAtLeast {
			tr.recordTransformerFunc = tr.transformHavingFieldsAtLeast
		} else if havingFieldsCriterion == havingFieldsWhichAre {
			tr.recordTransformerFunc = tr.transformHavingFieldsWhichAre
		} else if havingFieldsCriterion == havingFieldsAtMost {
			tr.recordTransformerFunc = tr.transformHavingFieldsAtMost
		} else {
			lib.InternalCodingErrorIf(true)
		}

	} else {
		// Let them type in a.*b if they want, or "a.*b", or "a.*b"i.
		// Strip off the leading " and trailing " or "i.
		regex, err := lib.CompileMillerRegex(regexString)
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"%s %s: cannot compile regex \"%s\"\n",
				"mlr",
				verbNameHavingFields,
				regexString,
			)
			os.Exit(1)
			// return nil, err
		}
		tr.regex = regex

		if havingFieldsCriterion == havingAllFieldsMatching {
			tr.recordTransformerFunc = tr.transformHavingAllFieldsMatching
		} else if havingFieldsCriterion == havingAnyFieldsMatching {
			tr.recordTransformerFunc = tr.transformHavingAnyFieldsMatching
		} else if havingFieldsCriterion == havingNoFieldsMatching {
			tr.recordTransformerFunc = tr.transformHavingNoFieldsMatching
		} else {
			lib.InternalCodingErrorIf(true)
		}
	}

	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerHavingFields) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, inputDownstreamDoneChannel, outputDownstreamDoneChannel, outputChannel)
}

// ----------------------------------------------------------------
func (tr *TransformerHavingFields) transformHavingFieldsAtLeast(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		numFound := 0
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if tr.fieldNameSet[pe.Key] {
				numFound++
				if numFound == tr.numFieldNames {
					outputChannel <- inrecAndContext
					return
				}
			}
		}

	} else {
		outputChannel <- inrecAndContext
	}
}

func (tr *TransformerHavingFields) transformHavingFieldsWhichAre(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		if inrec.FieldCount != tr.numFieldNames {
			return
		}
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if !tr.fieldNameSet[pe.Key] {
				return
			}
		}
		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext
	}
}

func (tr *TransformerHavingFields) transformHavingFieldsAtMost(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if !tr.fieldNameSet[pe.Key] {
				return
			}
		}
		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext
	}
}

// ----------------------------------------------------------------
func (tr *TransformerHavingFields) transformHavingAllFieldsMatching(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if !tr.regex.MatchString(pe.Key) {
				return
			}
		}
		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext
	}
}

func (tr *TransformerHavingFields) transformHavingAnyFieldsMatching(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if tr.regex.MatchString(pe.Key) {
				outputChannel <- inrecAndContext
				return
			}
		}
	} else {
		outputChannel <- inrecAndContext
	}
}

func (tr *TransformerHavingFields) transformHavingNoFieldsMatching(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if tr.regex.MatchString(pe.Key) {
				return
			}
		}
		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext
	}
}
