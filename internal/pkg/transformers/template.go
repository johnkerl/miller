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
const verbNameTemplate = "template"

var TemplateSetup = TransformerSetup{
	Verb:         verbNameTemplate,
	UsageFunc:    transformerTemplateUsage,
	ParseCLIFunc: transformerTemplateParseCLI,
	IgnoresInput: false,
}

func transformerTemplateUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameTemplate)
	fmt.Fprintf(o, "Places input-record fields in the order specified by list of column names.\n")
	fmt.Fprintf(o, "If the input record is missing a specified field, it will be filled with the fill-with.\n")
	fmt.Fprintf(o, "If the input record possesses an unspecified field, it will be discarded.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, " -f {a,b,c} Comma-separated field names for template, e.g. a,b,c.\n")
	fmt.Fprintf(o, " -t {filename} CSV file whose header line will be used for template.\n")
	fmt.Fprintf(o, "--fill-with {filler string}  What to fill absent fields with. Defaults to the empty string.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
	fmt.Fprintf(o, "Example:\n")
	fmt.Fprintf(o, "* Specified fields are a,b,c.\n")
	fmt.Fprintf(o, "* Input record is c=3,a=1,f=6.\n")
	fmt.Fprintf(o, "* Output record is a=1,b=,c=3.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerTemplateParseCLI(
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
	fillWith := ""

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerTemplateUsage(os.Stdout, true, 0)

		} else if opt == "-f" {
			fieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-t" {
			templateFileName := cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)
			temp, err := lib.ReadCSVHeader(templateFileName)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fieldNames = temp

		} else if opt == "--fill-with" {
			fillWith = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerTemplateUsage(os.Stderr, true, 1)
		}
	}

	if fieldNames == nil {
		transformerTemplateUsage(os.Stderr, true, 1)
	}

	transformer, err := NewTransformerTemplate(
		fieldNames,
		fillWith,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerTemplate struct {
	fieldNameList []string
	fieldNameSet  map[string]bool
	fillWith      *types.Mlrval
}

func NewTransformerTemplate(
	fieldNames []string,
	fillWith string,
) (*TransformerTemplate, error) {
	return &TransformerTemplate{
		fieldNameList: fieldNames,
		fieldNameSet:  lib.StringListToSet(fieldNames),
		fillWith:      types.MlrvalFromString(fillWith),
	}, nil
}

// ----------------------------------------------------------------

func (tr *TransformerTemplate) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		outrec := types.NewMlrmap()
		for _, fieldName := range tr.fieldNameList {
			value := inrec.Get(fieldName)
			if value != nil {
				outrec.PutReference(fieldName, value) // inrec will be GC'ed
			} else {
				outrec.PutCopy(fieldName, tr.fillWith)
			}
		}
		outrecAndContext := types.NewRecordAndContext(outrec, &inrecAndContext.Context)
		outputChannel <- outrecAndContext
	} else {
		outputChannel <- inrecAndContext
	}
}
