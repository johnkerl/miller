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

const verbNameSeqgen = "seqgen"

var seqgenOptions = []OptionSpec{
	{Flag: "-f", Arg: "{name}", Type: "string", Desc: "Field name for counters. Default \"i\"."},
	{Flag: "--start", Arg: "{value}", Type: "float", Desc: "Inclusive start value. Default 1."},
	{Flag: "--step", Arg: "{value}", Type: "float", Desc: "Step value. Default 1. May be negative but not zero (unless start == stop)."},
	{Flag: "--stop", Arg: "{value}", Type: "float", Desc: "Inclusive stop value. Default 100."},
}

var SeqgenSetup = TransformerSetup{
	Verb:         verbNameSeqgen,
	UsageFunc:    transformerSeqgenUsage,
	ParseCLIFunc: transformerSeqgenParseCLI,
	IgnoresInput: true,
	Options:      seqgenOptions,
}

func transformerSeqgenUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameSeqgen)
	fmt.Fprintf(o, "Produces a sequence of counters.  Discards the input record stream. Produces\n")
	fmt.Fprintf(o, "output as specified by the options.\n")
	fmt.Fprintf(o, "\n")
	WriteVerbOptions(o, seqgenOptions)

	fmt.Fprintf(o, "Start, stop, and/or step may be floating-point. Output is integer if start,\n")
	fmt.Fprintf(o, "stop, and step are all integers. Step may be negative. It may not be zero\n")
	fmt.Fprintf(o, "unless start == stop.\n")
}

func transformerSeqgenParseCLI(
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

	fieldName := "i"
	startString := "1"
	stopString := "100"
	stepString := "1"

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
			transformerSeqgenUsage(os.Stdout)
			return nil, cli.ErrHelpRequested

		case "-f":
			fieldName, err = cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		case "--start":
			startString, err = cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		case "--stop":
			stopString, err = cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		case "--step":
			stepString, err = cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		default:
			return nil, cli.VerbErrorf(verb, "option \"%s\" not recognized", opt)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	transformer, err := NewTransformerSeqgen(
		fieldName,
		startString,
		stopString,
		stepString,
	)
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

type TransformerSeqgen struct {
	fieldName      string
	start          *mlrval.Mlrval
	stop           *mlrval.Mlrval
	step           *mlrval.Mlrval
	doneComparator bifs.BinaryFunc
	mdone          *mlrval.Mlrval
}

func NewTransformerSeqgen(
	fieldName string,
	startString string,
	stopString string,
	stepString string,
) (*TransformerSeqgen, error) {
	start := mlrval.FromInferredType(startString)
	stop := mlrval.FromInferredType(stopString)
	step := mlrval.FromInferredType(stepString)
	var doneComparator bifs.BinaryFunc

	fstart, startIsNumeric := start.GetNumericToFloatValue()
	if !startIsNumeric {
		return nil, fmt.Errorf("mlr seqgen: start value should be number; got \"%s\"", startString)
	}

	fstop, stopIsNumeric := stop.GetNumericToFloatValue()
	if !stopIsNumeric {
		return nil, fmt.Errorf("mlr seqgen: stop value should be number; got \"%s\"", stopString)
	}

	fstep, stepIsNumeric := step.GetNumericToFloatValue()
	if !stepIsNumeric {
		return nil, fmt.Errorf("mlr seqgen: step value should be number; got \"%s\"", stepString)
	}

	if fstep > 0 {
		doneComparator = bifs.BIF_greater_than
	} else if fstep < 0 {
		doneComparator = bifs.BIF_less_than
	} else {
		if fstart == fstop {
			doneComparator = bifs.BIF_equals
		} else {
			return nil, fmt.Errorf("mlr seqgen: step must not be zero unless start == stop")
		}
	}

	return &TransformerSeqgen{
		fieldName:      fieldName,
		start:          start,
		stop:           stop,
		step:           step,
		doneComparator: doneComparator,
		mdone:          mlrval.FALSE,
	}, nil
}

// Transform satisfies RecordTransformer. In normal operation seqgen is
// driven via ProduceStream (see StreamingProducer) instead, since
// ChainTransformer special-cases any transformer implementing that
// interface. This is kept as a correct, if non-streaming, fallback.
func (tr *TransformerSeqgen) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) error {
	if !inrecAndContext.EndOfStream {
		// Discard upstream records; generate output only when upstream is done.
		return nil
	}

	counter := tr.start
	context := types.NewNilContext()
	context.UpdateForStartOfFile("seqgen")

	for {
		select {
		case b := <-inputDownstreamDoneChannel:
			outputDownstreamDoneChannel <- b
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, types.NewEndOfStreamMarker(context))
			return nil
		default:
		}

		tr.mdone = tr.doneComparator(counter, tr.stop)
		done, _ := tr.mdone.GetBoolValue()
		if done {
			break
		}

		outrec := mlrval.NewMlrmapAsRecord()
		outrec.PutCopy(tr.fieldName, counter)

		context.UpdateForInputRecord()

		outrecAndContext := types.NewRecordAndContext(outrec, context)
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, outrecAndContext)

		counter = bifs.BIF_plus_binary(counter, tr.step)
	}

	*outputRecordsAndContexts = append(*outputRecordsAndContexts, types.NewEndOfStreamMarker(context))
	return nil
}

// ProduceStream implements StreamingProducer. Unlike Transform, which must
// return its entire output in one shot -- unusable for a NoInput verb like
// seqgen, since nothing downstream can run until that single return happens
// -- this writes bounded batches directly to outputRecordChannel, checking
// inputDownstreamDoneChannel between batches. This is what lets 'mlr seqgen
// --stop 1000000000 then head -n 10' finish quickly instead of generating
// (and OOMing on) the entire sequence before head ever sees a record.
// Mirrors pkg/input/pseudo_reader_gen.go, which solves the same problem for
// 'mlr --from gen'.
func (tr *TransformerSeqgen) ProduceStream(
	outputRecordChannel chan<- []*types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	counter := tr.start
	context := types.NewNilContext()
	context.UpdateForStartOfFile("seqgen")

	recordsPerBatch := int64(cli.DEFAULT_RECORDS_PER_BATCH)
	batch := make([]*types.RecordAndContext, 0, recordsPerBatch)

	for {
		tr.mdone = tr.doneComparator(counter, tr.stop)
		done, _ := tr.mdone.GetBoolValue()
		if done {
			break
		}

		outrec := mlrval.NewMlrmapAsRecord()
		outrec.PutCopy(tr.fieldName, counter)

		context.UpdateForInputRecord()
		batch = append(batch, types.NewRecordAndContext(outrec, context))

		counter = bifs.BIF_plus_binary(counter, tr.step)

		if int64(len(batch)) >= recordsPerBatch {
			outputRecordChannel <- batch
			batch = make([]*types.RecordAndContext, 0, recordsPerBatch)

			// See if downstream transformers will be ignoring further data
			// (e.g. mlr head -n 10). If so, stop generating. Checked only
			// between batches to avoid goroutine-scheduler thrash. Even
			// though downstream won't use any more real records, it still
			// needs the end-of-stream marker to know to stop reading and
			// drain cleanly -- so that's sent below rather than returning
			// directly from here.
			select {
			case b := <-inputDownstreamDoneChannel:
				outputDownstreamDoneChannel <- b
				outputRecordChannel <- []*types.RecordAndContext{types.NewEndOfStreamMarker(context)}
				return
			default:
			}
		}
	}

	batch = append(batch, types.NewEndOfStreamMarker(context))
	outputRecordChannel <- batch
}
