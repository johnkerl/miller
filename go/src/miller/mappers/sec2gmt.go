package mappers

import (
	"flag"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/lib"
	"miller/mapping"
	"miller/types"
)

// ----------------------------------------------------------------
var Sec2GMTSetup = mapping.MapperSetup{
	Verb:         "sec2gmt",
	ParseCLIFunc: mapperSec2GMTParseCLI,
	IgnoresInput: false,
}

func mapperSec2GMTParseCLI(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) mapping.IRecordMapper {

	// Get the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// We don't use flagSet here since we want -1..-9 to all point to the same
	// numDecimalPlaces. That's super-easy if we roll our own.

	numDecimalPlaces := 0

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		if args[argi][0] != '-' {
			break // No more flag options to process
		} else if args[argi] == "-h" || args[argi] == "--help" {
			mapperSec2GMTUsage(os.Stdout, 0, errorHandling, args[0], verb)
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
			mapperSec2GMTUsage(os.Stderr, 1, flag.ExitOnError, args[0], verb)
			os.Exit(1)
		}
	}

	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	if argi >= argc {
		mapperSec2GMTUsage(os.Stderr, 1, flag.ExitOnError, args[0], verb)
		os.Exit(1)
	}
	fieldNames := args[argi]
	argi++

	mapper, _ := NewMapperSec2GMT(
		fieldNames,
		numDecimalPlaces,
	)

	*pargi = argi
	return mapper
}

func mapperSec2GMTUsage(
	o *os.File,
	exitCode int,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	argv0 string,
	verb string,
) {
	fmt.Fprintf(o, "Usage: %s %s [options] {comma-separated list of field names}\n", argv0, verb)
	fmt.Fprintf(o, "Replaces a numeric field representing seconds since the epoch with the\n")
	fmt.Fprintf(o, "corresponding GMT timestamp; leaves non-numbers as-is. This is nothing\n")
	fmt.Fprintf(o, "more than a keystroke-saver for the sec2gmt function:\n")
	fmt.Fprintf(o, "  %s %s time1,time2\n", argv0, verb)
	fmt.Fprintf(o, "is the same as\n")
	fmt.Fprintf(o, "  %s put '$time1=sec2gmt($time1);$time2=sec2gmt($time2)'\n", argv0)
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-1 through -9: format the seconds using 1..9 decimal places, respectively.\n")
	if errorHandling == flag.ExitOnError {
		os.Exit(exitCode)
	}
}

// ----------------------------------------------------------------
type MapperSec2GMT struct {
	fieldNameList    []string
	numDecimalPlaces int
}

func NewMapperSec2GMT(
	fieldNames string,
	numDecimalPlaces int,
) (*MapperSec2GMT, error) {
	this := &MapperSec2GMT{
		fieldNameList:    lib.SplitString(fieldNames, ","),
		numDecimalPlaces: numDecimalPlaces,
	}
	return this, nil
}

// ----------------------------------------------------------------
func (this *MapperSec2GMT) Map(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	if inrec != nil { // Not end of record stream
		for _, fieldName := range this.fieldNameList {
			value := inrec.Get(&fieldName)
			if value != nil {
				floatval, ok := value.GetNumericToFloatValue()
				if ok {
					newValue := types.MlrvalFromString(lib.Sec2GMT(floatval, this.numDecimalPlaces))
					inrec.PutReference(&fieldName, &newValue)
				}
			}
		}
		outputChannel <- inrecAndContext // end-of-stream marker

	} else { // End of record stream
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
