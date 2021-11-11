package transformers

import (
	"fmt"
	"os"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameSec2GMTDate = "sec2gmtdate"

var Sec2GMTDateSetup = TransformerSetup{
	Verb:         verbNameSec2GMTDate,
	UsageFunc:    transformerSec2GMTDateUsage,
	ParseCLIFunc: transformerSec2GMTDateParseCLI,
	IgnoresInput: false,
}

func transformerSec2GMTDateUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: ../c/mlr sec2gmtdate {comma-separated list of field names}\n")
	fmt.Fprintf(o, "Replaces a numeric field representing seconds since the epoch with the\n")
	fmt.Fprintf(o, "corresponding GMT year-month-day timestamp; leaves non-numbers as-is.\n")
	fmt.Fprintf(o, "This is nothing more than a keystroke-saver for the sec2gmtdate function:\n")
	fmt.Fprintf(o, "  ../c/mlr sec2gmtdate time1,time2\n")
	fmt.Fprintf(o, "is the same as\n")
	fmt.Fprintf(o, "  ../c/mlr put '$time1=sec2gmtdate($time1);$time2=sec2gmtdate($time2)'\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerSec2GMTDateParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	argi++

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if opt[0] != '-' {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerSec2GMTDateUsage(os.Stdout, true, 0)

		} else {
			transformerSec2GMTDateUsage(os.Stderr, true, 1)
		}
	}

	if argi >= argc {
		transformerSec2GMTDateUsage(os.Stderr, true, 1)
	}
	fieldNames := args[argi]
	argi++

	transformer, err := NewTransformerSec2GMTDate(
		fieldNames,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerSec2GMTDate struct {
	fieldNameList []string
}

func NewTransformerSec2GMTDate(
	fieldNames string,
) (*TransformerSec2GMTDate, error) {
	tr := &TransformerSec2GMTDate{
		fieldNameList: lib.SplitString(fieldNames, ","),
	}
	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerSec2GMTDate) Transform(
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
				inrec.PutReference(fieldName, types.BIF_sec2gmtdate(value))
			}
		}
		outputChannel <- inrecAndContext

	} else { // End of record stream
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
