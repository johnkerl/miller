package transformers

import (
	"fmt"
	"os"

	"miller/src/cliutil"
	"miller/src/lib"
	"miller/src/transforming"
	"miller/src/types"
)

// ----------------------------------------------------------------
const verbNameSec2GMT = "sec2gmt"

var Sec2GMTSetup = transforming.TransformerSetup{
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
	fmt.Fprintf(o, "Usage: %s %s [options] {comma-separated list of field names}\n", lib.MlrExeName(), verbNameSec2GMT)
	fmt.Fprintf(o, "Replaces a numeric field representing seconds since the epoch with the\n")
	fmt.Fprintf(o, "corresponding GMT timestamp; leaves non-numbers as-is. This is nothing\n")
	fmt.Fprintf(o, "more than a keystroke-saver for the sec2gmt function:\n")
	fmt.Fprintf(o, "  %s %s time1,time2\n", lib.MlrExeName(), verbNameSec2GMT)
	fmt.Fprintf(o, "is the same as\n")
	fmt.Fprintf(o, "  %s put '$time1 = sec2gmt($time1); $time2 = sec2gmt($time2)'\n", lib.MlrExeName())
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-1 through -9: format the seconds using 1..9 decimal places, respectively.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerSec2GMTParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cliutil.TReaderOptions,
	__ *cliutil.TWriterOptions,
) transforming.IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	argi++

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
	numDecimalPlaces int
}

func NewTransformerSec2GMT(
	fieldNames string,
	numDecimalPlaces int,
) (*TransformerSec2GMT, error) {
	this := &TransformerSec2GMT{
		fieldNameList:    lib.SplitString(fieldNames, ","),
		numDecimalPlaces: numDecimalPlaces,
	}
	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerSec2GMT) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		for _, fieldName := range this.fieldNameList {
			value := inrec.Get(fieldName)
			if value != nil {
				floatval, ok := value.GetNumericToFloatValue()
				if ok {
					newValue := types.MlrvalFromString(lib.Sec2GMT(floatval, this.numDecimalPlaces))
					inrec.PutReference(fieldName, &newValue)
				}
			}
		}
		outputChannel <- inrecAndContext

	} else { // End of record stream
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
