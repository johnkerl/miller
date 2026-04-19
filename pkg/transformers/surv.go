package transformers

import (
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
	"github.com/kshedden/statmodel/duration"
	"github.com/kshedden/statmodel/statmodel"
)

const verbNameSurv = "surv"

// SurvSetup defines the surv verb: Kaplan-Meier survival curve.
var SurvSetup = TransformerSetup{
	Verb:         verbNameSurv,
	UsageFunc:    transformerSurvUsage,
	ParseCLIFunc: transformerSurvParseCLI,
	IgnoresInput: false,
}

func transformerSurvUsage(o *os.File) {
	fmt.Fprintf(o, "Usage: %s %s -d {duration-field} -s {status-field}\n", "mlr", verbNameSurv)
	fmt.Fprint(o, `
Estimate Kaplan-Meier survival curve (right-censored).
Options:
  -d {field}   Name of duration field (time-to-event or censoring).
  -s {field}   Name of status field (0=censored, 1=event).
  -h, --help   Show this message.
`)
}

func transformerSurvParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool,
) (RecordTransformer, error) {
	argi := *pargi
	verb := args[argi]
	argi++

	var durationField, statusField string

	for argi < argc {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break
		}
		if opt == "-h" || opt == "--help" {
			transformerSurvUsage(os.Stdout)
			return nil, cli.ErrHelpRequested
		} else if opt == "-d" {
			if argi+1 >= argc {
				return nil, cli.VerbErrorf(verb, "-d requires an argument")
			}
			argi++
			durationField = args[argi]
			argi++
		} else if opt == "-s" {
			if argi+1 >= argc {
				return nil, cli.VerbErrorf(verb, "-s requires an argument")
			}
			argi++
			statusField = args[argi]
			argi++
		} else {
			break
		}
	}
	*pargi = argi
	if !doConstruct {
		return nil, nil
	}
	if durationField == "" {
		return nil, fmt.Errorf("mlr %s: -d option is required", verb)
	}
	if statusField == "" {
		return nil, fmt.Errorf("mlr %s: -s option is required", verb)
	}
	return NewTransformerSurv(durationField, statusField)
}

// TransformerSurv holds fields for surv verb.
type TransformerSurv struct {
	durationField string
	statusField   string
	times         []float64
	events        []bool
}

// NewTransformerSurv constructs a new surv transformer.
func NewTransformerSurv(durationField, statusField string) (*TransformerSurv, error) {
	return &TransformerSurv{
		durationField: durationField,
		statusField:   statusField,
		times:         []float64{},
		events:        []bool{},
	}, nil
}

// Transform processes each record or emits results at end-of-stream.
func (tr *TransformerSurv) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		rec := inrecAndContext.Record
		mvDur := rec.Get(tr.durationField)
		if mvDur == nil {
			// Skip records missing the duration field
			return
		}
		duration := mvDur.GetNumericToFloatValueOrDie()
		mvStat := rec.Get(tr.statusField)
		if mvStat == nil {
			// Skip records missing the status field
			return
		}
		status := mvStat.GetNumericToFloatValueOrDie() != 0
		tr.times = append(tr.times, duration)
		tr.events = append(tr.events, status)
	} else {
		// Compute survival using kshedden/statmodel
		n := len(tr.times)
		if n == 0 {
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
			return
		}
		durations := tr.times
		statuses := make([]float64, n)
		for i, ev := range tr.events {
			if ev {
				statuses[i] = 1.0
			} else {
				statuses[i] = 0.0
			}
		}
		dataCols := [][]float64{durations, statuses}
		names := []string{tr.durationField, tr.statusField}
		ds := statmodel.NewDataset(dataCols, names)
		sf, err := duration.NewSurvfuncRight(ds, tr.durationField, tr.statusField, &duration.SurvfuncRightConfig{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "mlr surv: %v\n", err)
			os.Exit(1)
		}
		sf.Fit()
		times := sf.Time()
		survProbs := sf.SurvProb()
		for i, t := range times {
			newrec := mlrval.NewMlrmapAsRecord()
			newrec.PutCopy("time", mlrval.FromFloat(t))
			newrec.PutCopy("survival", mlrval.FromFloat(survProbs[i]))
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, types.NewRecordAndContext(newrec, &inrecAndContext.Context))
		}
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
	}
}
