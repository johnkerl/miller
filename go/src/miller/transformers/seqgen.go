package transformers

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
const verbNameSeqgen = "seqgen"

var SeqgenSetup = transforming.TransformerSetup{
	Verb:         verbNameSeqgen,
	ParseCLIFunc: transformerSeqgenParseCLI,
	IgnoresInput: true,
}

func transformerSeqgenParseCLI(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) transforming.IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// Parse local flags
	flagSet := flag.NewFlagSet(verb, errorHandling)

	pFieldName := flagSet.String(
		"f",
		"i",
		"Field name for counters.",
	)

	pStartString := flagSet.String(
		"start",
		"1",
		"Inclusive start value.",
	)

	pStopString := flagSet.String(
		"stop",
		"100",
		"Inclusive stop value.",
	)

	pStepString := flagSet.String(
		"step",
		"1",
		"Step value",
	)

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		transformerSeqgenUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	transformer, err := NewTransformerSeqgen(
		*pFieldName,
		*pStartString,
		*pStopString,
		*pStepString,
	)
	// TODO: put error return into this API
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

func transformerSeqgenUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprintf(o, "Passes input records directly to output. Most useful for format conversion.\n")
	fmt.Fprintf(o, "Produces a sequence of counters.  Discards the input record stream. Produces\n")
	fmt.Fprintf(o, "output as specified by the following options:\n")
	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v (default %v) %v\n", f.Name, f.Value, f.Usage) // f.Name, f.Value
	})
	fmt.Fprintf(o, "Start, stop, and/or step may be floating-point. Output is integer if start,\n")
	fmt.Fprintf(o, "stop, and step are all integers. Step may be negative. It may not be zero\n")
	fmt.Fprintf(o, "unless start == stop.\n")
}

// ----------------------------------------------------------------
type TransformerSeqgen struct {
	fieldName      string
	start          types.Mlrval
	stop           types.Mlrval
	step           types.Mlrval
	doneComparator types.BinaryFunc
}

// ----------------------------------------------------------------
func NewTransformerSeqgen(
	fieldName string,
	startString string,
	stopString string,
	stepString string,
) (*TransformerSeqgen, error) {
	start := types.MlrvalFromInferredType(startString)
	stop := types.MlrvalFromInferredType(stopString)
	step := types.MlrvalFromInferredType(stepString)
	var doneComparator types.BinaryFunc = nil

	fstart, startIsNumeric := start.GetNumericToFloatValue()
	if !startIsNumeric {
		return nil, errors.New(
			fmt.Sprintf(
				"mlr seqgen: start value should be number; got \"%s\"",
				startString,
			),
		)
	}

	fstop, stopIsNumeric := stop.GetNumericToFloatValue()
	if !stopIsNumeric {
		return nil, errors.New(
			fmt.Sprintf(
				"mlr seqgen: stop value should be number; got \"%s\"",
				stopString,
			),
		)
	}

	fstep, stepIsNumeric := step.GetNumericToFloatValue()
	if !stepIsNumeric {
		return nil, errors.New(
			fmt.Sprintf(
				"mlr seqgen: step value should be number; got \"%s\"",
				stepString,
			),
		)
	}

	if fstep > 0 {
		doneComparator = types.MlrvalGreaterThan
	} else if fstep < 0 {
		doneComparator = types.MlrvalLessThan
	} else {
		if fstart == fstop {
			doneComparator = types.MlrvalEquals
		} else {
			return nil, errors.New(
				"mlr seqgen: step must not be zero unless start == stop.",
			)
		}
	}

	return &TransformerSeqgen{
		fieldName:      fieldName,
		start:          start,
		stop:           stop,
		step:           step,
		doneComparator: doneComparator,
	}, nil
}

// ----------------------------------------------------------------
func (this *TransformerSeqgen) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	counter := this.start
	context := types.NewContext(nil)
	context.UpdateForStartOfFile("seqgen")

	for {
		mdone := this.doneComparator(&counter, &this.stop)
		done, _ := mdone.GetBoolValue()
		if done {
			break
		}

		outrec := types.NewMlrmapAsRecord()
		outrec.PutCopy(this.fieldName, &counter)

		context.UpdateForInputRecord()

		outrecAndContext := types.NewRecordAndContext(outrec, context)
		outputChannel <- outrecAndContext

		counter = types.MlrvalBinaryPlus(&counter, &this.step)
	}

	outputChannel <- types.NewEndOfStreamMarker(context)
}
