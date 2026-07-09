package transformers

import (
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/bifs"
	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

const verbNameSparkline = "sparkline"

var sparklineOptions = []OptionSpec{
	{Flag: "-f", Arg: "{a,b,c}", Type: "csv-list", Desc: "Field names to sparkline."},
}

var SparklineSetup = TransformerSetup{
	Verb:         verbNameSparkline,
	UsageFunc:    transformerSparklineUsage,
	ParseCLIFunc: transformerSparklineParseCLI,
	IgnoresInput: false,
	Options:      sparklineOptions,
}

func transformerSparklineUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameSparkline)
	fmt.Fprintf(o, "Reduces numeric field(s), across all records in input order, to a compact\n")
	fmt.Fprintf(o, "Unicode sparkline -- one block character per record -- for visualizing\n")
	fmt.Fprintf(o, "trends. Emits one output record per field. Holds all records in memory\n")
	fmt.Fprintf(o, "before producing any output.\n")
	WriteVerbOptions(o, sparklineOptions)
}

func transformerSparklineParseCLI(
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

	// Parse local flags
	var err error
	var fieldNames []string = nil

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		if args[argi] == "--" {
			break // All transformers must do this so main-flags can follow verb-flags
		}
		argi++

		switch opt {
		case "-h", "--help":
			transformerSparklineUsage(os.Stdout)
			return nil, cli.ErrHelpRequested

		case "-f":
			fieldNames, err = cli.VerbGetStringArrayArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		default:
			return nil, cli.VerbErrorf(verb, "option \"%s\" not recognized", opt)
		}
	}

	if fieldNames == nil {
		return nil, cli.VerbErrorf(verb, "-f field names required")
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	transformer, err := NewTransformerSparkline(fieldNames)
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

type TransformerSparkline struct {
	fieldNames    []string
	valuesByField map[string][]*mlrval.Mlrval
}

func NewTransformerSparkline(
	fieldNames []string,
) (*TransformerSparkline, error) {
	valuesByField := make(map[string][]*mlrval.Mlrval)
	for _, fieldName := range fieldNames {
		valuesByField[fieldName] = make([]*mlrval.Mlrval, 0)
	}
	return &TransformerSparkline{
		fieldNames:    fieldNames,
		valuesByField: valuesByField,
	}, nil
}

func (tr *TransformerSparkline) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)

	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		for _, fieldName := range tr.fieldNames {
			mvalue := inrec.Get(fieldName)
			if mvalue != nil {
				tr.valuesByField[fieldName] = append(tr.valuesByField[fieldName], mvalue.Copy())
			}
		}
		return
	}

	// Else, end of stream: emit one summary record per field.
	for _, fieldName := range tr.fieldNames {
		values := tr.valuesByField[fieldName]

		outrec := mlrval.NewMlrmapAsRecord()
		outrec.PutReference("field", mlrval.FromString(fieldName))
		outrec.PutReference("n", mlrval.FromInt(int64(len(values))))

		sparkline := bifs.BIF_sparkline(mlrval.FromArray(values))
		if !sparkline.IsError() {
			lo, hi, haveLoHi := floatRangeOf(values)
			if haveLoHi {
				outrec.PutReference("lo", mlrval.FromFloat(lo))
				outrec.PutReference("hi", mlrval.FromFloat(hi))
			}
		}
		outrec.PutReference("sparkline", sparkline)

		*outputRecordsAndContexts = append(*outputRecordsAndContexts, types.NewRecordAndContext(outrec, &inrecAndContext.Context))
	}

	*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // Emit the end-of-stream marker
}

// floatRangeOf returns the min and max of the numeric values, and whether
// any were found.
func floatRangeOf(values []*mlrval.Mlrval) (lo float64, hi float64, haveLoHi bool) {
	for _, value := range values {
		floatValue, isFloat := value.GetNumericToFloatValue()
		if !isFloat {
			continue
		}
		if !haveLoHi {
			lo = floatValue
			hi = floatValue
			haveLoHi = true
		} else {
			if floatValue < lo {
				lo = floatValue
			}
			if floatValue > hi {
				hi = floatValue
			}
		}
	}
	return lo, hi, haveLoHi
}
