package transformers

import (
	"fmt"
	"os"
	"strings"

	"miller/src/cliutil"
	"miller/src/lib"
	"miller/src/transforming"
	"miller/src/types"
)

// ----------------------------------------------------------------
const verbNameHistogram = "histogram"

var HistogramSetup = transforming.TransformerSetup{
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
	argv0 := lib.MlrExeName()
	verb := verbNameHistogram
	fmt.Fprintf(o, "Just a histogram. Input values < lo or > hi are not counted.\n")
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprintf(o, "-f {a,b,c}    Value-field names for histogram counts\n")
	fmt.Fprintf(o, "--lo {lo}     Histogram low value\n")
	fmt.Fprintf(o, "--hi {hi}     Histogram high value\n")
	fmt.Fprintf(o, "--nbins {n}   Number of histogram bins\n")
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
	_ *cliutil.TReaderOptions,
	__ *cliutil.TWriterOptions,
) transforming.IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// Parse local flags
	var valueFieldNames []string = nil
	lo := 0.0
	nbins := 0
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
			valueFieldNames = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "--lo" {
			lo = cliutil.VerbGetFloatArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "--nbins" {
			nbins = cliutil.VerbGetIntArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "--hi" {
			hi = cliutil.VerbGetFloatArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "--auto" {
			doAuto = true

		} else if opt == "-o" {
			outputPrefix = cliutil.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerHistogramUsage(os.Stderr, true, 1)
		}
	}

	if valueFieldNames == nil {
		transformerHistogramUsage(os.Stderr, true, 1)
	}

	if nbins == 0 {
		transformerHistogramUsage(os.Stderr, true, 1)
	}

	if lo == hi && !doAuto {
		transformerHistogramUsage(os.Stderr, true, 1)
	}

	transformer, _ := NewTransformerHistogram(
		valueFieldNames,
		lo,
		nbins,
		hi,
		doAuto,
		outputPrefix,
	)

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

	recordTransformerFunc transforming.RecordTransformerFunc
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

	this := &TransformerHistogram{
		valueFieldNames: valueFieldNames,
		countsByField:   countsByField,
		outputPrefix:    outputPrefix,
		nbins:           nbins,
	}

	if !doAuto {
		this.recordTransformerFunc = this.transformNonAuto
		this.lo = lo
		this.hi = hi
		this.mul = float64(nbins) / (hi - lo)
	} else {
		this.vectorsByFieldName = make(map[string][]float64)
		for _, valueFieldName := range valueFieldNames {
			this.vectorsByFieldName[valueFieldName] = make([]float64, 0, histogramVectorInitialSize)
		}

		this.recordTransformerFunc = this.transformAuto
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerHistogram) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	this.recordTransformerFunc(inrecAndContext, outputChannel)
}

// ----------------------------------------------------------------
func (this *TransformerHistogram) transformNonAuto(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		this.ingestNonAuto(inrecAndContext)
	} else {
		this.emitNonAuto(&inrecAndContext.Context, outputChannel)
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}

func (this *TransformerHistogram) ingestNonAuto(
	inrecAndContext *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	for _, valueFieldName := range this.valueFieldNames {
		stringValue := inrec.Get(valueFieldName)
		if stringValue != nil {
			floatValue, ok := stringValue.GetNumericToFloatValue()
			if !ok {
				fmt.Fprintf(
					os.Stderr,
					"%s %s: cannot parse \"%s\" as float.\n",
					lib.MlrExeName(), verbNameHistogram, stringValue.String(),
				)
				os.Exit(1)
			}
			if (floatValue >= this.lo) && (floatValue < this.hi) {
				idx := int((floatValue - this.lo) * this.mul)
				this.countsByField[valueFieldName][idx]++
			} else if floatValue == this.hi {
				idx := this.nbins - 1
				this.countsByField[valueFieldName][idx]++
			}
		}
	}
}

func (this *TransformerHistogram) emitNonAuto(
	endOfStreamContext *types.Context,
	outputChannel chan<- *types.RecordAndContext,
) {
	countFieldNames := make(map[string]string)
	for _, valueFieldName := range this.valueFieldNames {
		countFieldNames[valueFieldName] = this.outputPrefix + valueFieldName + "_count"
	}
	for i := 0; i < this.nbins; i++ {
		outrec := types.NewMlrmapAsRecord()

		outrec.PutReference(
			this.outputPrefix+"bin_lo",
			types.MlrvalPointerFromFloat64((this.lo+float64(i))/this.mul),
		)
		outrec.PutReference(
			this.outputPrefix+"bin_hi",
			types.MlrvalPointerFromFloat64((this.lo+float64(i+1))/this.mul),
		)

		for _, valueFieldName := range this.valueFieldNames {
			outrec.PutReference(
				countFieldNames[valueFieldName],
				types.MlrvalPointerFromInt(this.countsByField[valueFieldName][i]),
			)
		}

		outputChannel <- types.NewRecordAndContext(outrec, endOfStreamContext)
	}
}

// ----------------------------------------------------------------
func (this *TransformerHistogram) transformAuto(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		this.ingestAuto(inrecAndContext)
	} else {
		this.emitAuto(&inrecAndContext.Context, outputChannel)
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}

func (this *TransformerHistogram) ingestAuto(
	inrecAndContext *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	for _, valueFieldName := range this.valueFieldNames {
		mvalue := inrec.Get(valueFieldName)
		if mvalue != nil {
			value := mvalue.GetNumericToFloatValueOrDie()
			this.vectorsByFieldName[valueFieldName] = append(this.vectorsByFieldName[valueFieldName], value)
		}
	}
}

func (this *TransformerHistogram) emitAuto(
	endOfStreamContext *types.Context,
	outputChannel chan<- *types.RecordAndContext,
) {
	haveLoHi := false
	lo := 0.0
	hi := 1.0
	nbins := this.nbins

	// Limits pass
	for _, valueFieldName := range this.valueFieldNames {
		vector := this.vectorsByFieldName[valueFieldName]
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
	for _, valueFieldName := range this.valueFieldNames {
		vector := this.vectorsByFieldName[valueFieldName]
		counts := this.countsByField[valueFieldName]
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
	for _, valueFieldName := range this.valueFieldNames {
		countFieldNames[valueFieldName] = this.outputPrefix + valueFieldName + "_count"
	}

	for i := 0; i < nbins; i++ {
		outrec := types.NewMlrmapAsRecord()

		outrec.PutReference(
			this.outputPrefix+"bin_lo",
			types.MlrvalPointerFromFloat64((lo+float64(i))/mul),
		)
		outrec.PutReference(
			this.outputPrefix+"bin_hi",
			types.MlrvalPointerFromFloat64((lo+float64(i+1))/mul),
		)

		for _, valueFieldName := range this.valueFieldNames {
			outrec.PutReference(
				countFieldNames[valueFieldName],
				types.MlrvalPointerFromInt(this.countsByField[valueFieldName][i]),
			)
		}

		outputChannel <- types.NewRecordAndContext(outrec, endOfStreamContext)
	}
}
