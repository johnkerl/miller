package transformers

import (
	"fmt"
	"os"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameSec2GMT = "sec2gmt"

var Sec2GMTSetup = TransformerSetup{
	Verb:         verbNameSec2GMT,
	UsageFunc:    transformerSec2GMTUsage,
	ParseCLIFunc: transformerSec2GMTParseCLI,
	IgnoresInput: false,
}

func transformerSec2GMTUsage(
	o *os.File,
	doExit bool,
	exitCode int,
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

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerSec2GMTParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
) IRecordTransformer {

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
		argi++

		if opt == "-h" || opt == "--help" {
			transformerSec2GMTUsage(os.Stdout, true, 0)

		} else if opt == "-1" {
			numDecimalPlaces = 1
		} else if opt == "-2" {
			numDecimalPlaces = 2
		} else if opt == "-3" {
			numDecimalPlaces = 3
		} else if opt == "-4" {
			numDecimalPlaces = 4
		} else if opt == "-5" {
			numDecimalPlaces = 5
		} else if opt == "-6" {
			numDecimalPlaces = 6
		} else if opt == "-7" {
			numDecimalPlaces = 7
		} else if opt == "-8" {
			numDecimalPlaces = 8
		} else if opt == "-9" {
			numDecimalPlaces = 9

		} else if opt == "--millis" {
			preDivide = 1.0e3
		} else if opt == "--micros" {
			preDivide = 1.0e6
		} else if opt == "--nanos" {
			preDivide = 1.0e9

		} else {
			transformerSec2GMTUsage(os.Stderr, true, 1)
		}
	}

	if argi >= argc {
		transformerSec2GMTUsage(os.Stderr, true, 1)
	}
	fieldNames := args[argi]
	argi++

	transformer, err := NewTransformerSec2GMT(
		fieldNames,
		preDivide,
		numDecimalPlaces,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
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

// ----------------------------------------------------------------

func (tr *TransformerSec2GMT) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		for _, fieldName := range tr.fieldNameList {
			value := inrec.Get(fieldName)
			if value != nil {
				floatval, ok := value.GetNumericToFloatValue()
				if ok {
					newValue := types.MlrvalFromString(lib.Sec2GMT(
						floatval/tr.preDivide,
						tr.numDecimalPlaces,
					))
					inrec.PutReference(fieldName, newValue)
				}
			}
		}
		outputChannel <- inrecAndContext

	} else { // End of record stream
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
