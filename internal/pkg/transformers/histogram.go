package transformers

import (
	"fmt"
	"os"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameHistogram = "histogram"
const histogramDefaultBinCount = 20

var HistogramSetup = TransformerSetup{
	Verb:         verbNameHistogram,
	UsageFunc:    transformerHistogramUsage,
	ParseCLIFunc: transformerHistogramParseCLI,
	IgnoresInput: false,
}

func transformerHistogramUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	argv0 := "mlr"
	verb := verbNameHistogram
	fmt.Fprintf(o, "Just a histogram. Input values < lo or > hi are not counted.\n")
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprintf(o, "-f {a,b,c}    Value-field names for histogram counts\n")
	fmt.Fprintf(o, "--lo {lo}     Histogram low value\n")
	fmt.Fprintf(o, "--hi {hi}     Histogram high value\n")
	fmt.Fprintf(o, "--nbins {n}   Number of histogram bins. Defaults to %d.\n", histogramDefaultBinCount)
	fmt.Fprintf(o, "--auto        Automatically computes limits, ignoring --lo and --hi.\n")
	fmt.Fprintf(o, "              Holds all values in memory before producing any output.\n")
	fmt.Fprintf(o, "-o {prefix}   Prefix for output field name. Default: no prefix.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerHistogramParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
) IRecordTransformer {

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

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerHistogramUsage(os.Stdout, true, 0)

		} else if opt == "-f" {
			valueFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "--lo" {
			lo = cli.VerbGetFloatArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "--nbins" {
			nbins = cli.VerbGetIntArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "--hi" {
			hi = cli.VerbGetFloatArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "--auto" {
			doAuto = true

		} else if opt == "-o" {
			outputPrefix = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerHistogramUsage(os.Stderr, true, 1)
		}
	}

	if valueFieldNames == nil {
		transformerHistogramUsage(os.Stderr, true, 1)
	}

	if nbins <= 0 {
		transformerHistogramUsage(os.Stderr, true, 1)
	}

	if lo == hi && !doAuto {
		transformerHistogramUsage(os.Stderr, true, 1)
	}

	transformer, err := NewTransformerHistogram(
		valueFieldNames,
		lo,
		nbins,
		hi,
		doAuto,
		outputPrefix,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
const histogramVectorInitialSize = 1024

type TransformerHistogram struct {
	valueFieldNames []string
	lo              float64
	nbins           int
	hi              float64
	mul             float64

	countsByField      map[string][]int
	vectorsByFieldName map[string][]float64 // For auto-mode
	outputPrefix       string

	recordTransformerFunc RecordTransformerFunc
}

// ----------------------------------------------------------------
func NewTransformerHistogram(
	valueFieldNames []string,
	lo float64,
	nbins int,
	hi float64,
	doAuto bool,
	outputPrefix string,
) (*TransformerHistogram, error) {

	countsByField := make(map[string][]int)
	for _, valueFieldName := range valueFieldNames {
		countsByField[valueFieldName] = make([]int, nbins)
		for i := 0; i < nbins; i++ {
			countsByField[valueFieldName][i] = 0
		}
	}

	tr := &TransformerHistogram{
		valueFieldNames: valueFieldNames,
		countsByField:   countsByField,
		outputPrefix:    outputPrefix,
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

// ----------------------------------------------------------------

func (tr *TransformerHistogram) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, inputDownstreamDoneChannel, outputDownstreamDoneChannel, outputChannel)
}

// ----------------------------------------------------------------
func (tr *TransformerHistogram) transformNonAuto(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		tr.ingestNonAuto(inrecAndContext)
	} else {
		tr.emitNonAuto(&inrecAndContext.Context, outputChannel)
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}

func (tr *TransformerHistogram) ingestNonAuto(
	inrecAndContext *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	for _, valueFieldName := range tr.valueFieldNames {
		stringValue := inrec.Get(valueFieldName)
		if stringValue != nil {
			floatValue, ok := stringValue.GetNumericToFloatValue()
			if !ok {
				fmt.Fprintf(
					os.Stderr,
					"%s %s: cannot parse \"%s\" as float.\n",
					"mlr", verbNameHistogram, stringValue.String(),
				)
				os.Exit(1)
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
}

func (tr *TransformerHistogram) emitNonAuto(
	endOfStreamContext *types.Context,
	outputChannel chan<- *types.RecordAndContext,
) {
	countFieldNames := make(map[string]string)
	for _, valueFieldName := range tr.valueFieldNames {
		countFieldNames[valueFieldName] = tr.outputPrefix + valueFieldName + "_count"
	}
	for i := 0; i < tr.nbins; i++ {
		outrec := types.NewMlrmapAsRecord()

		outrec.PutReference(
			tr.outputPrefix+"bin_lo",
			types.MlrvalFromFloat64((tr.lo+float64(i))/tr.mul),
		)
		outrec.PutReference(
			tr.outputPrefix+"bin_hi",
			types.MlrvalFromFloat64((tr.lo+float64(i+1))/tr.mul),
		)

		for _, valueFieldName := range tr.valueFieldNames {
			outrec.PutReference(
				countFieldNames[valueFieldName],
				types.MlrvalFromInt(tr.countsByField[valueFieldName][i]),
			)
		}

		outputChannel <- types.NewRecordAndContext(outrec, endOfStreamContext)
	}
}

// ----------------------------------------------------------------
func (tr *TransformerHistogram) transformAuto(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		tr.ingestAuto(inrecAndContext)
	} else {
		tr.emitAuto(&inrecAndContext.Context, outputChannel)
		outputChannel <- inrecAndContext // end-of-stream marker
	}
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
	outputChannel chan<- *types.RecordAndContext,
) {
	haveLoHi := false
	lo := 0.0
	hi := 1.0
	nbins := tr.nbins

	// Limits pass
	for _, valueFieldName := range tr.valueFieldNames {
		vector := tr.vectorsByFieldName[valueFieldName]
		n := len(vector)
		for i := 0; i < n; i++ {
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
		for i := 0; i < n; i++ {
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

	// Emission pass
	countFieldNames := make(map[string]string)
	for _, valueFieldName := range tr.valueFieldNames {
		countFieldNames[valueFieldName] = tr.outputPrefix + valueFieldName + "_count"
	}

	for i := 0; i < nbins; i++ {
		outrec := types.NewMlrmapAsRecord()

		outrec.PutReference(
			tr.outputPrefix+"bin_lo",
			types.MlrvalFromFloat64((lo+float64(i))/mul),
		)
		outrec.PutReference(
			tr.outputPrefix+"bin_hi",
			types.MlrvalFromFloat64((lo+float64(i+1))/mul),
		)

		for _, valueFieldName := range tr.valueFieldNames {
			outrec.PutReference(
				countFieldNames[valueFieldName],
				types.MlrvalFromInt(tr.countsByField[valueFieldName][i]),
			)
		}

		outputChannel <- types.NewRecordAndContext(outrec, endOfStreamContext)
	}
}
