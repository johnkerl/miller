package transformers

import (
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

const verbNameTemplate = "template"

var TemplateSetup = TransformerSetup{
	Verb:         verbNameTemplate,
	UsageFunc:    transformerTemplateUsage,
	ParseCLIFunc: transformerTemplateParseCLI,
	IgnoresInput: false,
}

func transformerTemplateUsage(
	o *os.File,
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
}

func transformerTemplateParseCLI(
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
	fillWith := ""

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
			transformerTemplateUsage(os.Stdout)
			return nil, cli.ErrHelpRequested

		} else if opt == "-f" {
			fieldNames, err = cli.VerbGetStringArrayArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		} else if opt == "-t" {
			templateFileName, err := cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			temp, err := lib.ReadCSVHeader(templateFileName)
			if err != nil {
				return nil, fmt.Errorf("mlr %s: cannot read template file: %w", verb, err)
			}
			fieldNames = temp

		} else if opt == "--fill-with" {
			fillWith, err = cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		} else {
			return nil, cli.VerbErrorf(verb, "option \"%s\" not recognized", opt)
		}
	}

	if fieldNames == nil {
		return nil, cli.VerbErrorf(verb, "-f or -t is required")
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	transformer, err := NewTransformerTemplate(
		fieldNames,
		fillWith,
	)
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

type TransformerTemplate struct {
	fieldNameList []string
	fieldNameSet  map[string]bool
	fillWith      *mlrval.Mlrval
}

func NewTransformerTemplate(
	fieldNames []string,
	fillWith string,
) (*TransformerTemplate, error) {
	return &TransformerTemplate{
		fieldNameList: fieldNames,
		fieldNameSet:  lib.StringListToSet(fieldNames),
		fillWith:      mlrval.FromString(fillWith),
	}, nil
}

func (tr *TransformerTemplate) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		outrec := mlrval.NewMlrmap()
		for _, fieldName := range tr.fieldNameList {
			value := inrec.Get(fieldName)
			if value != nil {
				outrec.PutReference(fieldName, value) // inrec will be GC'ed
			} else {
				outrec.PutCopy(fieldName, tr.fillWith)
			}
		}
		outrecAndContext := types.NewRecordAndContext(outrec, &inrecAndContext.Context)
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, outrecAndContext)
	} else {
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
	}
}
