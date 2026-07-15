package transformers

import (
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/transformers/utils"
	"github.com/johnkerl/miller/v6/pkg/types"
)

const verbNameBootstrapCI = "bootstrap-ci"

const bootstrapCIDefaultNumResamples = int64(1000)
const bootstrapCIDefaultConfidenceLevel = 0.95

var bootstrapCIOptions = []OptionSpec{
	{Flag: "-a", Arg: "{mean,...}", Type: "enum", Desc: "Names of statistics to bootstrap: one or more of the listed values, as in mlr stats1 -a. Also accepts median (same as p50) and percentiles p{n} for n in 0..100. Defaults to mean.", Values: []string{"count", "null_count", "distinct_count", "mode", "antimode", "sum", "mean", "mad", "var", "stddev", "meaneb", "skewness", "kurtosis", "min", "max", "minlen", "maxlen"}},
	{Flag: "-f", Arg: "{a,b,c}", Type: "csv-list", Desc: "Value-field names on which to compute statistics. Required."},
	{Flag: "-g", Arg: "{d,e,f}", Type: "csv-list", Desc: "Optional group-by-field names."},
	{Flag: "-n", Arg: "{n}", Type: "int", Desc: "Number of bootstrap resamples. Must be positive. Defaults to 1000."},
	{Flag: "-c", Arg: "{level}", Type: "float", Desc: "Confidence level, strictly between 0 and 1. Defaults to 0.95."},
	{Flag: "-i", Type: "bool", Desc: "Use interpolated percentiles, like R's type=7, for percentile statistics as well as for the confidence-interval endpoints; default like type=1."},
}

var BootstrapCISetup = TransformerSetup{
	Verb:         verbNameBootstrapCI,
	UsageFunc:    transformerBootstrapCIUsage,
	ParseCLIFunc: transformerBootstrapCIParseCLI,
	IgnoresInput: false,
	Options:      bootstrapCIOptions,
}

func transformerBootstrapCIUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameBootstrapCI)
	fmt.Fprintf(o,
		`Computes bootstrap confidence intervals for statistics of given fields,
accumulated across the input record stream: values are resampled with
replacement many times, the statistic is computed on each resample, and the
confidence interval is taken from percentiles of the resampled statistics.
For each value field and statistic, outputs the full-data statistic in
{field}_{stat}, along with confidence-interval endpoints in {field}_{stat}_lo
and {field}_{stat}_hi. Use mlr --seed for reproducible results.
See also %s bootstrap and %s stats1.
`, "mlr", "mlr")
	WriteVerbOptions(o, bootstrapCIOptions)
	fmt.Fprintln(o,
		"Example: mlr --seed 12345 bootstrap-ci -f x,y")
	fmt.Fprintln(o,
		"Example: mlr --seed 12345 bootstrap-ci -a mean,median -f x -g shape -n 5000 -c 0.99")
}

func transformerBootstrapCIParseCLI(
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

	accumulatorNameList := []string{"mean"}
	valueFieldNameList := []string{}
	groupByFieldNameList := []string{}
	numResamples := bootstrapCIDefaultNumResamples
	confidenceLevel := bootstrapCIDefaultConfidenceLevel
	doInterpolatedPercentiles := false

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
			transformerBootstrapCIUsage(os.Stdout)
			return nil, cli.ErrHelpRequested

		case "-a":
			accumulatorNameList, err = cli.VerbGetStringArrayArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		case "-f":
			valueFieldNameList, err = cli.VerbGetStringArrayArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		case "-g":
			groupByFieldNameList, err = cli.VerbGetStringArrayArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		case "-n":
			numResamples, err = cli.VerbGetIntArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		case "-c":
			confidenceLevel, err = cli.VerbGetFloatArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		case "-i":
			doInterpolatedPercentiles = true

		default:
			return nil, cli.VerbErrorf(verb, "option \"%s\" not recognized", opt)
		}
	}

	if len(valueFieldNameList) == 0 {
		return nil, cli.VerbErrorf(verb, "-f option is required")
	}
	if numResamples <= 0 {
		return nil, cli.VerbErrorf(verb, "-n argument must be positive; got %d", numResamples)
	}
	if confidenceLevel <= 0.0 || confidenceLevel >= 1.0 {
		return nil, cli.VerbErrorf(verb, "-c argument must be strictly between 0 and 1; got %v", confidenceLevel)
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	transformer, err := NewTransformerBootstrapCI(
		accumulatorNameList,
		valueFieldNameList,
		groupByFieldNameList,
		numResamples,
		confidenceLevel,
		doInterpolatedPercentiles,
	)
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

type TransformerBootstrapCI struct {
	// Input:
	accumulatorNameList       []string
	valueFieldNameList        []string
	groupByFieldNameList      []string
	numResamples              int64
	confidenceLevel           float64
	doInterpolatedPercentiles bool

	// State:
	accumulatorFactory *utils.Stats1AccumulatorFactory

	// Retained field values are indexed by
	//   groupingKey -> valueFieldName -> array of values
	// using maps that preserve insertion order.
	valuesByGroup *lib.OrderedMap[*lib.OrderedMap[[]*mlrval.Mlrval]]

	// map[groupingKey]OrderedMap[groupByFieldName]*mlrval.Mlrval
	groupingKeysToGroupByFieldValues map[string]*lib.OrderedMap[*mlrval.Mlrval]
}

func NewTransformerBootstrapCI(
	accumulatorNameList []string,
	valueFieldNameList []string,
	groupByFieldNameList []string,
	numResamples int64,
	confidenceLevel float64,
	doInterpolatedPercentiles bool,
) (*TransformerBootstrapCI, error) {
	for _, name := range accumulatorNameList {
		if !utils.ValidateStats1AccumulatorName(name) {
			return nil, fmt.Errorf(`mlr %s: accumulator "%s" not found`, verbNameBootstrapCI, name)
		}
	}

	tr := &TransformerBootstrapCI{
		accumulatorNameList:       accumulatorNameList,
		valueFieldNameList:        valueFieldNameList,
		groupByFieldNameList:      groupByFieldNameList,
		numResamples:              numResamples,
		confidenceLevel:           confidenceLevel,
		doInterpolatedPercentiles: doInterpolatedPercentiles,

		accumulatorFactory:               utils.NewStats1AccumulatorFactory(),
		valuesByGroup:                    lib.NewOrderedMap[*lib.OrderedMap[[]*mlrval.Mlrval]](),
		groupingKeysToGroupByFieldValues: make(map[string]*lib.OrderedMap[*mlrval.Mlrval]),
	}

	return tr, nil
}

// Transform is the function executed for every input record, as well as for
// the end-of-stream marker.
func (tr *TransformerBootstrapCI) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) error {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		tr.handleInputRecord(inrecAndContext)
	} else {
		tr.handleEndOfRecordStream(inrecAndContext, outputRecordsAndContexts)
	}
	return nil
}

func (tr *TransformerBootstrapCI) handleInputRecord(
	inrecAndContext *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record

	// E.g. if grouping by "a" and "b", and the current record has a=circle,
	// b=blue, then groupingKey is the string "circle,blue".
	groupingKey, ok := inrec.GetSelectedValuesJoined(tr.groupByFieldNameList)
	if !ok {
		return
	}

	valuesByFieldName := tr.valuesByGroup.Get(groupingKey)
	if valuesByFieldName == nil {
		valuesByFieldName = lib.NewOrderedMap[[]*mlrval.Mlrval]()
		tr.valuesByGroup.Put(groupingKey, valuesByFieldName)

		// E.g. if grouping by "color" and "shape", and the current record has
		// color=blue, shape=circle, then groupByFieldValues is the map
		// {"color": "blue", "shape": "circle"}.
		groupByFieldValuesArray, ok := inrec.GetSelectedValues(tr.groupByFieldNameList)
		if !ok {
			return
		}
		groupByFieldValues := lib.NewOrderedMap[*mlrval.Mlrval]()
		for i, groupByFieldValue := range groupByFieldValuesArray {
			groupByFieldValues.Put(tr.groupByFieldNameList[i], groupByFieldValue.Copy())
		}
		tr.groupingKeysToGroupByFieldValues[groupingKey] = groupByFieldValues
	}

	for _, valueFieldName := range tr.valueFieldNameList {
		valueFieldValue := inrec.Get(valueFieldName)
		if valueFieldValue == nil || valueFieldValue.IsVoid() {
			continue
		}
		values := valuesByFieldName.Get(valueFieldName)
		valuesByFieldName.Put(valueFieldName, append(values, valueFieldValue.Copy()))
	}
}

func (tr *TransformerBootstrapCI) handleEndOfRecordStream(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
) {
	for pa := tr.valuesByGroup.Head; pa != nil; pa = pa.Next {
		groupingKey := pa.Key
		valuesByFieldName := pa.Value

		newrec := mlrval.NewMlrmapAsRecord()

		groupByFieldValues := tr.groupingKeysToGroupByFieldValues[groupingKey]
		for pb := groupByFieldValues.Head; pb != nil; pb = pb.Next {
			newrec.PutCopy(pb.Key, pb.Value)
		}

		for pc := valuesByFieldName.Head; pc != nil; pc = pc.Next {
			valueFieldName := pc.Key
			values := pc.Value
			if len(values) == 0 {
				continue
			}
			for _, accumulatorName := range tr.accumulatorNameList {
				tr.emitConfidenceInterval(groupingKey, valueFieldName, accumulatorName, values, newrec)
			}
		}

		*outputRecordsAndContexts = append(*outputRecordsAndContexts,
			types.NewRecordAndContext(newrec, &inrecAndContext.Context))
	}

	*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // end-of-stream marker
}

// emitConfidenceInterval computes, for one group's values of one field, the
// full-data statistic along with its bootstrap confidence interval, placing
// the results into outrec.
func (tr *TransformerBootstrapCI) emitConfidenceInterval(
	groupingKey string,
	valueFieldName string,
	accumulatorName string,
	values []*mlrval.Mlrval,
	outrec *mlrval.Mlrmap,
) {
	// Point estimate: the statistic computed over the full data.
	pointEstimate := tr.computeStatistic(groupingKey, valueFieldName, accumulatorName, values, false)

	// Compute the statistic over each of the bootstrap resamples, retaining
	// the resampled statistics in a percentile-keeper.
	replicateStats := utils.NewPercentileKeeper(tr.doInterpolatedPercentiles)
	for rep := int64(0); rep < tr.numResamples; rep++ {
		replicateStats.Ingest(
			tr.computeStatistic(groupingKey, valueFieldName, accumulatorName, values, true),
		)
	}

	// E.g. for confidence level 0.95, the interval endpoints are the 2.5th
	// and 97.5th percentiles of the resampled statistics.
	alpha := (1.0 - tr.confidenceLevel) / 2.0
	lo := replicateStats.Emit(100.0 * alpha)
	hi := replicateStats.Emit(100.0 * (1.0 - alpha))

	outputFieldNameBase := valueFieldName + "_" + accumulatorName
	outrec.PutCopy(outputFieldNameBase, pointEstimate)
	outrec.PutCopy(outputFieldNameBase+"_lo", lo)
	outrec.PutCopy(outputFieldNameBase+"_hi", hi)
}

// computeStatistic computes a single stats1-style statistic over the given
// values -- either as-is, or over one same-length resample with replacement.
func (tr *TransformerBootstrapCI) computeStatistic(
	groupingKey string,
	valueFieldName string,
	accumulatorName string,
	values []*mlrval.Mlrval,
	resample bool,
) *mlrval.Mlrval {
	// Reset the factory so that percentile-statistic accumulators get a fresh
	// percentile-keeper for each computation, rather than sharing one.
	tr.accumulatorFactory.Reset()
	accumulator := tr.accumulatorFactory.MakeAccumulator(
		accumulatorName,
		groupingKey,
		valueFieldName,
		tr.doInterpolatedPercentiles,
	)
	// Accumulator names were pre-validated at construction time.
	lib.InternalCodingErrorIf(accumulator == nil)

	n := int64(len(values))
	if resample {
		for i := int64(0); i < n; i++ {
			accumulator.Ingest(values[lib.RandRange(0, n)])
		}
	} else {
		for _, value := range values {
			accumulator.Ingest(value)
		}
	}

	return accumulator.Emit()
}
