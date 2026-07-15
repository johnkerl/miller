package transformers

import (
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/bifs"
	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

const verbNameHistogram = "histogram"
const histogramDefaultBinCount = int64(20)

var histogramOptions = []OptionSpec{
	{Flag: "-f", Arg: "{a,b,c}", Type: "csv-list", Desc: "Value-field names for histogram counts."},
	{Flag: "--lo", Arg: "{lo}", Type: "float", Desc: "Histogram low value."},
	{Flag: "--hi", Arg: "{hi}", Type: "float", Desc: "Histogram high value."},
	{Flag: "--nbins", Arg: "{n}", Type: "int", Desc: "Number of histogram bins. Defaults to 20."},
	{Flag: "--auto", Type: "bool", Desc: "Automatically computes limits, ignoring --lo and --hi. Holds all values in memory before producing any output."},
	{Flag: "-o", Arg: "{prefix}", Type: "string", Desc: "Prefix for output field name. Default: no prefix."},
	{Flag: "-s", Type: "bool", Desc: "Print a one-line Unicode sparkline per field instead of per-bin counts."},
}

var HistogramSetup = TransformerSetup{
	Verb:         verbNameHistogram,
	UsageFunc:    transformerHistogramUsage,
	ParseCLIFunc: transformerHistogramParseCLI,
	IgnoresInput: false,
	Options:      histogramOptions,
}

func transformerHistogramUsage(
	o *os.File,
) {
	argv0 := "mlr"
	verb := verbNameHistogram
	fmt.Fprintf(o, "Just a histogram. Input values < lo or > hi are not counted.\n")
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	WriteVerbOptions(o, histogramOptions)
	fmt.Fprintf(o, "With -s, output is one record per value-field, with a sparkline field\n")
	fmt.Fprintf(o, "instead of one record per bin.\n")
}

func transformerHistogramParseCLI(
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
	var valueFieldNames []string = nil
	lo := 0.0
	nbins := histogramDefaultBinCount
	hi := 0.0
	doAuto := false
	outputPrefix := ""
	doSparkline := false

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

		switch opt {
		case "-h", "--help":
			transformerHistogramUsage(os.Stdout)
			return nil, cli.ErrHelpRequested

		case "-f":
			valueFieldNames, err = cli.VerbGetStringArrayArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		case "--lo":
			lo, err = cli.VerbGetFloatArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		case "--nbins":
			nbins, err = cli.VerbGetIntArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		case "--hi":
			hi, err = cli.VerbGetFloatArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		case "--auto":
			doAuto = true

		case "-o":
			outputPrefix, err = cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		case "-s":
			doSparkline = true

		default:
			return nil, cli.VerbErrorf(verb, "option \"%s\" not recognized", opt)
		}
	}

	if valueFieldNames == nil {
		return nil, cli.VerbErrorf(verb, "-f field names required")
	}

	if nbins <= 0 {
		return nil, cli.VerbErrorf(verb, "number of bins must be positive")
	}

	if lo == hi && !doAuto {
		return nil, cli.VerbErrorf(verb, "lo and hi must differ, or use --auto")
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	transformer, err := NewTransformerHistogram(
		valueFieldNames,
		lo,
		nbins,
		hi,
		doAuto,
		outputPrefix,
		doSparkline,
	)
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

const histogramVectorInitialSize = 1024

type TransformerHistogram struct {
	valueFieldNames []string
	lo              float64
	nbins           int64
	hi              float64
	mul             float64

	countsByField      map[string][]int64
	vectorsByFieldName map[string][]float64 // For auto-mode
	outputPrefix       string
	doSparkline        bool

	recordTransformerFunc RecordTransformerFunc
}

func NewTransformerHistogram(
	valueFieldNames []string,
	lo float64,
	nbins int64,
	hi float64,
	doAuto bool,
	outputPrefix string,
	doSparkline bool,
) (*TransformerHistogram, error) {

	countsByField := make(map[string][]int64)
	for _, valueFieldName := range valueFieldNames {
		countsByField[valueFieldName] = make([]int64, nbins)
		for i := range nbins {
			countsByField[valueFieldName][i] = 0
		}
	}

	tr := &TransformerHistogram{
		valueFieldNames: valueFieldNames,
		countsByField:   countsByField,
		outputPrefix:    outputPrefix,
		doSparkline:     doSparkline,
		nbins:           nbins,
	}

	if !doAuto {
		tr.recordTransformerFunc = tr.transformNonAuto
		tr.lo = lo
		tr.hi = hi
		tr.mul = float64(nbins) / (hi - lo)
	} else {
		tr.vectorsByFieldName = make(map[string][]float64)
		for _, valueFieldName := range valueFieldNames {
			tr.vectorsByFieldName[valueFieldName] = make([]float64, 0, histogramVectorInitialSize)
		}

		tr.recordTransformerFunc = tr.transformAuto
	}

	return tr, nil
}

func (tr *TransformerHistogram) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) error {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	return tr.recordTransformerFunc(inrecAndContext, outputRecordsAndContexts, inputDownstreamDoneChannel, outputDownstreamDoneChannel)
}

func (tr *TransformerHistogram) transformNonAuto(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) error {
	if !inrecAndContext.EndOfStream {
		if err := tr.ingestNonAuto(inrecAndContext); err != nil {
			return err
		}
	} else {
		tr.emitNonAuto(&inrecAndContext.Context, outputRecordsAndContexts)
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // end-of-stream marker
	}
	return nil
}

func (tr *TransformerHistogram) ingestNonAuto(
	inrecAndContext *types.RecordAndContext,
) error {
	inrec := inrecAndContext.Record
	for _, valueFieldName := range tr.valueFieldNames {
		stringValue := inrec.Get(valueFieldName)
		if stringValue != nil {
			floatValue, ok := stringValue.GetNumericToFloatValue()
			if !ok {
				return cli.VerbErrorf(verbNameHistogram,
					"cannot parse \"%s\" as float.", stringValue.String())
			}
			if (floatValue >= tr.lo) && (floatValue < tr.hi) {
				idx := int((floatValue - tr.lo) * tr.mul)
				tr.countsByField[valueFieldName][idx]++
			} else if floatValue == tr.hi {
				idx := tr.nbins - 1
				tr.countsByField[valueFieldName][idx]++
			}
		}
	}
	return nil
}

// sparklineRecord summarizes a field's binned counts as a single record with
// a Unicode-block sparkline field, rather than one record per bin.
func (tr *TransformerHistogram) sparklineRecord(
	valueFieldName string,
	lo float64,
	hi float64,
) *mlrval.Mlrmap {
	counts := tr.countsByField[valueFieldName]
	countMlrvals := make([]*mlrval.Mlrval, len(counts))
	for i, count := range counts {
		countMlrvals[i] = mlrval.FromInt(count)
	}
	sparkline := bifs.BIF_sparkline(mlrval.FromArray(countMlrvals))

	outrec := mlrval.NewMlrmapAsRecord()
	outrec.PutReference(tr.outputPrefix+"field", mlrval.FromString(valueFieldName))
	outrec.PutReference(tr.outputPrefix+"lo", mlrval.FromFloat(lo))
	outrec.PutReference(tr.outputPrefix+"hi", mlrval.FromFloat(hi))
	outrec.PutReference(tr.outputPrefix+"sparkline", sparkline)
	return outrec
}

func (tr *TransformerHistogram) emitNonAuto(
	endOfStreamContext *types.Context,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
) {
	if tr.doSparkline {
		for _, valueFieldName := range tr.valueFieldNames {
			outrec := tr.sparklineRecord(valueFieldName, tr.lo, tr.hi)
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, types.NewRecordAndContext(outrec, endOfStreamContext))
		}
		return
	}

	countFieldNames := make(map[string]string)
	for _, valueFieldName := range tr.valueFieldNames {
		countFieldNames[valueFieldName] = tr.outputPrefix + valueFieldName + "_count"
	}
	for i := int64(0); i < tr.nbins; i++ {
		outrec := mlrval.NewMlrmapAsRecord()

		outrec.PutReference(
			tr.outputPrefix+"bin_lo",
			mlrval.FromFloat(tr.lo+float64(i)/tr.mul),
		)
		outrec.PutReference(
			tr.outputPrefix+"bin_hi",
			mlrval.FromFloat(tr.lo+float64(i+1)/tr.mul),
		)

		for _, valueFieldName := range tr.valueFieldNames {
			outrec.PutReference(
				countFieldNames[valueFieldName],
				mlrval.FromInt(tr.countsByField[valueFieldName][i]),
			)
		}

		*outputRecordsAndContexts = append(*outputRecordsAndContexts, types.NewRecordAndContext(outrec, endOfStreamContext))
	}
}

func (tr *TransformerHistogram) transformAuto(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) error {
	if !inrecAndContext.EndOfStream {
		tr.ingestAuto(inrecAndContext)
	} else {
		tr.emitAuto(&inrecAndContext.Context, outputRecordsAndContexts)
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // end-of-stream marker
	}
	return nil
}

func (tr *TransformerHistogram) ingestAuto(
	inrecAndContext *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	for _, valueFieldName := range tr.valueFieldNames {
		mvalue := inrec.Get(valueFieldName)
		if mvalue != nil {
			value := mvalue.GetNumericToFloatValueOrDie()
			tr.vectorsByFieldName[valueFieldName] = append(tr.vectorsByFieldName[valueFieldName], value)
		}
	}
}

func (tr *TransformerHistogram) emitAuto(
	endOfStreamContext *types.Context,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
) {
	haveLoHi := false
	lo := 0.0
	hi := 1.0
	nbins := tr.nbins

	// Limits pass
	for _, valueFieldName := range tr.valueFieldNames {
		vector := tr.vectorsByFieldName[valueFieldName]
		n := len(vector)
		for i := range n {
			value := vector[i]
			if haveLoHi {
				if lo > value {
					lo = value
				}
				if hi < value {
					hi = value
				}
			} else {
				lo = value
				hi = value
				haveLoHi = true
			}
		}
	}

	// Binning pass
	mul := float64(nbins) / (hi - lo)
	for _, valueFieldName := range tr.valueFieldNames {
		vector := tr.vectorsByFieldName[valueFieldName]
		counts := tr.countsByField[valueFieldName]
		lib.InternalCodingErrorIf(counts == nil)
		n := len(vector)
		for i := range n {
			value := vector[i]
			if (value >= lo) && (value < hi) {
				idx := int(((value - lo) * mul))
				counts[idx]++
			} else if value == hi {
				idx := nbins - 1
				counts[idx]++
			}
		}
	}

	if tr.doSparkline {
		for _, valueFieldName := range tr.valueFieldNames {
			outrec := tr.sparklineRecord(valueFieldName, lo, hi)
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, types.NewRecordAndContext(outrec, endOfStreamContext))
		}
		return
	}

	// Emission pass
	countFieldNames := make(map[string]string)
	for _, valueFieldName := range tr.valueFieldNames {
		countFieldNames[valueFieldName] = tr.outputPrefix + valueFieldName + "_count"
	}

	for i := range nbins {
		outrec := mlrval.NewMlrmapAsRecord()

		outrec.PutReference(
			tr.outputPrefix+"bin_lo",
			mlrval.FromFloat(lo+(float64(i)/mul)),
		)
		outrec.PutReference(
			tr.outputPrefix+"bin_hi",
			mlrval.FromFloat(lo+(float64(i+1)/mul)),
		)

		for _, valueFieldName := range tr.valueFieldNames {
			outrec.PutReference(
				countFieldNames[valueFieldName],
				mlrval.FromInt(tr.countsByField[valueFieldName][i]),
			)
		}

		*outputRecordsAndContexts = append(*outputRecordsAndContexts, types.NewRecordAndContext(outrec, endOfStreamContext))
	}
}
