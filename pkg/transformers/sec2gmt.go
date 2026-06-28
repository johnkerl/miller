package transformers

import (
	"fmt"
	"os"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

const verbNameSec2GMT = "sec2gmt"

var Sec2GMTSetup = TransformerSetup{
	Verb:         verbNameSec2GMT,
	UsageFunc:    transformerSec2GMTUsage,
	ParseCLIFunc: transformerSec2GMTParseCLI,
	IgnoresInput: false,
	Options: []OptionSpec{
		{Flag: "-1", Type: "bool", Desc: "Format seconds with 1 decimal place."},
		{Flag: "-2", Type: "bool", Desc: "Format seconds with 2 decimal places."},
		{Flag: "-3", Type: "bool", Desc: "Format seconds with 3 decimal places."},
		{Flag: "-4", Type: "bool", Desc: "Format seconds with 4 decimal places."},
		{Flag: "-5", Type: "bool", Desc: "Format seconds with 5 decimal places."},
		{Flag: "-6", Type: "bool", Desc: "Format seconds with 6 decimal places."},
		{Flag: "-7", Type: "bool", Desc: "Format seconds with 7 decimal places."},
		{Flag: "-8", Type: "bool", Desc: "Format seconds with 8 decimal places."},
		{Flag: "-9", Type: "bool", Desc: "Format seconds with 9 decimal places."},
		{Flag: "--millis", Type: "bool", Desc: "Input numbers are treated as milliseconds since the epoch."},
		{Flag: "--micros", Type: "bool", Desc: "Input numbers are treated as microseconds since the epoch."},
		{Flag: "--nanos", Type: "bool", Desc: "Input numbers are treated as nanoseconds since the epoch."},
	},
}

func transformerSec2GMTUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options] {comma-separated list of field names}\n", "mlr", verbNameSec2GMT)
	fmt.Fprintf(o, "Replaces a numeric field representing seconds since the epoch with the\n")
	fmt.Fprintf(o, "corresponding GMT timestamp; leaves non-numbers as-is. This is nothing\n")
	fmt.Fprintf(o, "more than a keystroke-saver for the sec2gmt function:\n")
	fmt.Fprintf(o, "  %s %s time1,time2\n", "mlr", verbNameSec2GMT)
	fmt.Fprintf(o, "is the same as\n")
	fmt.Fprintf(o, "  %s put '$time1 = sec2gmt($time1); $time2 = sec2gmt($time2)'\n", "mlr")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-1 through -9: format the seconds using 1..9 decimal places, respectively.\n")
	fmt.Fprintf(o, "--millis Input numbers are treated as milliseconds since the epoch.\n")
	fmt.Fprintf(o, "--micros Input numbers are treated as microseconds since the epoch.\n")
	fmt.Fprintf(o, "--nanos  Input numbers are treated as nanoseconds since the epoch.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
}

func transformerSec2GMTParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) (RecordTransformer, error) {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	argi++

	preDivide := 1.0
	numDecimalPlaces := 0

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if opt[0] != '-' {
			break // No more flag options to process
		}
		if args[argi] == "--" {
			break // All transformers must do this so main-flags can follow verb-flags
		}
		argi++

		switch opt {
		case "-h", "--help":
			transformerSec2GMTUsage(os.Stdout)
			return nil, cli.ErrHelpRequested

		case "-1":
			numDecimalPlaces = 1
		case "-2":
			numDecimalPlaces = 2
		case "-3":
			numDecimalPlaces = 3
		case "-4":
			numDecimalPlaces = 4
		case "-5":
			numDecimalPlaces = 5
		case "-6":
			numDecimalPlaces = 6
		case "-7":
			numDecimalPlaces = 7
		case "-8":
			numDecimalPlaces = 8
		case "-9":
			numDecimalPlaces = 9

		case "--millis":
			preDivide = 1.0e3
		case "--micros":
			preDivide = 1.0e6
		case "--nanos":
			preDivide = 1.0e9

		default:
			return nil, cli.VerbErrorf(verbNameSec2GMT, "option \"%s\" not recognized", opt)
		}
	}

	if argi >= argc {
		return nil, cli.VerbErrorf(verbNameSec2GMT, "field names required")
	}
	fieldNames := args[argi]
	argi++

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	transformer, err := NewTransformerSec2GMT(
		fieldNames,
		preDivide,
		numDecimalPlaces,
	)
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

type TransformerSec2GMT struct {
	fieldNameList    []string
	preDivide        float64
	numDecimalPlaces int
}

func NewTransformerSec2GMT(
	fieldNames string,
	preDivide float64,
	numDecimalPlaces int,
) (*TransformerSec2GMT, error) {
	tr := &TransformerSec2GMT{
		fieldNameList:    lib.SplitString(fieldNames, ","),
		preDivide:        preDivide,
		numDecimalPlaces: numDecimalPlaces,
	}
	return tr, nil
}

func (tr *TransformerSec2GMT) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		for _, fieldName := range tr.fieldNameList {
			value := inrec.Get(fieldName)
			if value != nil {
				floatval, ok := value.GetNumericToFloatValue()
				if ok {
					newValue := mlrval.FromString(lib.Sec2GMT(
						floatval/tr.preDivide,
						tr.numDecimalPlaces,
					))
					inrec.PutReference(fieldName, newValue)
				}
			}
		}
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)

	} else { // End of record stream
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // end-of-stream marker
	}
}
