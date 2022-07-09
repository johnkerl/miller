package transformers

import (
	"container/list"
	"fmt"
	"os"

	"github.com/johnkerl/miller/internal/pkg/bifs"
	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/types"
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
) {
	fmt.Fprintf(o, "Usage: ../c/mlr sec2gmtdate {comma-separated list of field names}\n")
	fmt.Fprintf(o, "Replaces a numeric field representing seconds since the epoch with the\n")
	fmt.Fprintf(o, "corresponding GMT year-month-day timestamp; leaves non-numbers as-is.\n")
	fmt.Fprintf(o, "This is nothing more than a keystroke-saver for the sec2gmtdate function:\n")
	fmt.Fprintf(o, "  ../c/mlr sec2gmtdate time1,time2\n")
	fmt.Fprintf(o, "is the same as\n")
	fmt.Fprintf(o, "  ../c/mlr put '$time1=sec2gmtdate($time1);$time2=sec2gmtdate($time2)'\n")
}

func transformerSec2GMTDateParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	argi++

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if opt[0] != '-' {
			break // No more flag options to process
		}
		if args[argi] == "--" {
			break // All transformers must do this so main-flags can follow verb-flags
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerSec2GMTDateUsage(os.Stdout)
			os.Exit(0)

		} else {
			transformerSec2GMTDateUsage(os.Stderr)
			os.Exit(1)
		}
	}

	if argi >= argc {
		transformerSec2GMTDateUsage(os.Stderr)
		os.Exit(1)
	}
	fieldNames := args[argi]
	argi++

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerSec2GMTDate(
		fieldNames,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

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
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		for _, fieldName := range tr.fieldNameList {
			value := inrec.Get(fieldName)
			if value != nil {
				inrec.PutReference(fieldName, bifs.BIF_sec2gmtdate(value))
			}
		}
		outputRecordsAndContexts.PushBack(inrecAndContext)

	} else { // End of record stream
		outputRecordsAndContexts.PushBack(inrecAndContext) // end-of-stream marker
	}
}
