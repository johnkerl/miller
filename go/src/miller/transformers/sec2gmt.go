package transformers

import (
	"flag"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/lib"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
const verbNameSec2GMT = "sec2gmt"

var Sec2GMTSetup = transforming.TransformerSetup{
	Verb:         verbNameSec2GMT,
	ParseCLIFunc: transformerSec2GMTParseCLI,
	UsageFunc:    transformerSec2GMTUsage,
	IgnoresInput: false,
}

func transformerSec2GMTParseCLI(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) transforming.IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	argi++

	numDecimalPlaces := 0

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		if args[argi][0] != '-' {
			break // No more flag options to process
		} else if args[argi] == "-h" || args[argi] == "--help" {
			transformerSec2GMTUsage(os.Stdout, true, 0)
			argi += 1

		} else if args[argi] == "-1" {
			numDecimalPlaces = 1
			argi++
		} else if args[argi] == "-2" {
			numDecimalPlaces = 2
			argi++
		} else if args[argi] == "-3" {
			numDecimalPlaces = 3
			argi++
		} else if args[argi] == "-4" {
			numDecimalPlaces = 4
			argi++
		} else if args[argi] == "-5" {
			numDecimalPlaces = 5
			argi++
		} else if args[argi] == "-6" {
			numDecimalPlaces = 6
			argi++
		} else if args[argi] == "-7" {
			numDecimalPlaces = 7
			argi++
		} else if args[argi] == "-8" {
			numDecimalPlaces = 8
			argi++
		} else if args[argi] == "-9" {
			numDecimalPlaces = 9
			argi++

		} else {
			transformerSec2GMTUsage(os.Stderr, true, 1)
		}
	}

	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	if argi >= argc {
		transformerSec2GMTUsage(os.Stderr, true, 1)
	}
	fieldNames := args[argi]
	argi++

	transformer, _ := NewTransformerSec2GMT(
		fieldNames,
		numDecimalPlaces,
	)

	*pargi = argi
	return transformer
}

func transformerSec2GMTUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options] {comma-separated list of field names}\n", os.Args[0], verbNameSec2GMT)
	fmt.Fprintf(o, "Replaces a numeric field representing seconds since the epoch with the\n")
	fmt.Fprintf(o, "corresponding GMT timestamp; leaves non-numbers as-is. This is nothing\n")
	fmt.Fprintf(o, "more than a keystroke-saver for the sec2gmt function:\n")
	fmt.Fprintf(o, "  %s %s time1,time2\n", os.Args[0], verbNameSec2GMT)
	fmt.Fprintf(o, "is the same as\n")
	fmt.Fprintf(o, "  %s put '$time1 = sec2gmt($time1); $time2 = sec2gmt($time2)'\n", os.Args[0])
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-1 through -9: format the seconds using 1..9 decimal places, respectively.\n")

	if doExit {
		os.Exit(exitCode)
	}
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
