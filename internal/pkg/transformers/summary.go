package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/bifs"
	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
	"github.com/johnkerl/miller/internal/pkg/transformers/utils"
	"github.com/johnkerl/miller/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameSummary = "summary"

var SummarySetup = TransformerSetup{
	Verb:         verbNameSummary,
	UsageFunc:    transformerSummaryUsage,
	ParseCLIFunc: transformerSummaryParseCLI,
	IgnoresInput: false,
}

func transformerSummaryUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameSummary)
	fmt.Fprintf(o, `Show summary statistics about the input data:
  field_name field_type
  min        p25            median   p75   max
  count      mean           stddev
  null_count distinct_count mode

  TODO: defaults
  TODO: -a foo,bar
  TODO: --all

Notes:
* Distinct-counts are computed on string representations -- so 4.1 and 4.10 are counted as distinct here.
* If the mode is not unique in the input data, the first-encountered value is reported as the mode.
`)
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
}

func transformerSummaryParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// defaults
	accumulatorNames := []string{
		"field_type",
		"count",
		"mean",
		"min",
		"median",
		"max",
		"null_count",
		"distinct_count",
	}

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		if args[argi] == "--" {
			break // All transformers must do this so main-flags can follow verb-flags
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerSummaryUsage(os.Stdout)
			os.Exit(0)

		} else if opt == "--all" {
			// XXX to top
			accumulatorNames = []string{
				"field_type",
				"count",
				"sum",
				"mean",
				"stddev",
				"min",
				"p25",
				"median",
				"p75",
				"max",
				"iqr",
				"lof",
				"lif",
				"uif",
				"uof",
				"null_count",
				"distinct_count",
				"mode",
			}

		} else if opt == "-a" {
			accumulatorNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)
			// xxx: abend if not contained in all-list

		} else {
			transformerSummaryUsage(os.Stderr)
			os.Exit(1)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerSummary(accumulatorNames)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type tFieldSummary struct {
	// Needs lib.OrderedMap, not map[string]int64, for deterministic regression-test output.
	// This is a map (a set really) rather than a single value in case of heterogeneous data.
	fieldTypesMap *lib.OrderedMap

	countAccumulator  utils.IStats1Accumulator
	sumAccumulator    utils.IStats1Accumulator
	meanAccumulator   utils.IStats1Accumulator
	stddevAccumulator utils.IStats1Accumulator

	percentileKeeper *utils.PercentileKeeper

	nullCount int64
	// Needs lib.OrderedMap, not map[string]int64, for deterministic regression-test output.
	// This is used for distinct_count (map length) as well as mode (max value in the map).
	distincts *lib.OrderedMap
}

func newFieldSummary() *tFieldSummary {
	return &tFieldSummary{
		fieldTypesMap: lib.NewOrderedMap(),

		countAccumulator:  utils.NewStats1CountAccumulator(),
		sumAccumulator:    utils.NewStats1SumAccumulator(),
		meanAccumulator:   utils.NewStats1MeanAccumulator(),
		stddevAccumulator: utils.NewStats1StddevAccumulator(),

		// Interpolated percentiles don't play well with string-valued input data
		percentileKeeper: utils.NewPercentileKeeper(false),

		nullCount: 0,
		distincts: lib.NewOrderedMap(),
	}
}

type TransformerSummary struct {
	fieldSummaries           *lib.OrderedMap
	accumulatorNamesToIngest map[string]bool
	accumulatorNamesToEmit   map[string]bool
}

func NewTransformerSummary(
	accumulatorNames []string,
) (*TransformerSummary, error) {

	tr := &TransformerSummary{
		fieldSummaries:           lib.NewOrderedMap(),
		accumulatorNamesToIngest: make(map[string]bool),
		accumulatorNamesToEmit:   make(map[string]bool),
	}

	// xxx
	for _, accumulatorName := range accumulatorNames {
		tr.accumulatorNamesToIngest[accumulatorName] = true
		// Dependencies:
		if accumulatorName == "iqr" {
			tr.accumulatorNamesToIngest["p25"] = true
			tr.accumulatorNamesToIngest["p75"] = true
		} else if accumulatorName == "lof" || accumulatorName == "lif" || accumulatorName == "uif" || accumulatorName == "uof" {
			tr.accumulatorNamesToIngest["p25"] = true
			tr.accumulatorNamesToIngest["p75"] = true
			tr.accumulatorNamesToIngest["iqr"] = true
		}

		tr.accumulatorNamesToEmit[accumulatorName] = true
	}

	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerSummary) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)

	if !inrecAndContext.EndOfStream {
		// Ingest another record
		inrec := inrecAndContext.Record

		for pe := inrec.Head; pe != nil; pe = pe.Next {
			fieldName := pe.Key

			iFieldSummary := tr.fieldSummaries.Get(fieldName)
			var fieldSummary *tFieldSummary
			if iFieldSummary == nil {
				fieldSummary = newFieldSummary()
				tr.fieldSummaries.Put(fieldName, fieldSummary)
			} else {
				fieldSummary = iFieldSummary.(*tFieldSummary)
			}

			// Go generics would be grand to put into lib.OrderedMap, but not all platforms
			// are on recent enough Go compiler versions.
			if tr.accumulatorNamesToIngest["field_type"] {
				typeName := pe.Value.GetTypeName()
				iValue := fieldSummary.fieldTypesMap.Get(typeName)
				if iValue == nil {
					fieldSummary.fieldTypesMap.Put(typeName, int64(1))
				} else {
					fieldSummary.fieldTypesMap.Put(typeName, iValue.(int64)+1)
				}
			}

			fieldSummary.percentileKeeper.Ingest(pe.Value)

			if pe.Value.IsNumeric() {
				if tr.accumulatorNamesToIngest["count"] {
					fieldSummary.countAccumulator.Ingest(pe.Value)
				}
				if tr.accumulatorNamesToIngest["sum"] {
					fieldSummary.sumAccumulator.Ingest(pe.Value)
				}
				if tr.accumulatorNamesToIngest["mean"] {
					fieldSummary.meanAccumulator.Ingest(pe.Value)
				}
				if tr.accumulatorNamesToIngest["stddev"] {
					fieldSummary.stddevAccumulator.Ingest(pe.Value)
				}
			}

			if tr.accumulatorNamesToIngest["null_count"] {
				if pe.Value.IsNull() || pe.Value.IsVoid() {
					fieldSummary.nullCount++
				}
			}

			if tr.accumulatorNamesToIngest["distinct_count"] || tr.accumulatorNamesToIngest["mode"] {
				valueString := pe.Value.String()
				iValue := fieldSummary.distincts.Get(valueString)
				if iValue == nil {
					fieldSummary.distincts.Put(valueString, int64(1))
				} else {
					fieldSummary.distincts.Put(valueString, iValue.(int64)+1)
				}
			}
		}

	} else {
		// Emit results at end of record stream
		for pe := tr.fieldSummaries.Head; pe != nil; pe = pe.Next {
			newrec := mlrval.NewMlrmapAsRecord()
			fieldName := pe.Key
			fieldSummary := pe.Value.(*tFieldSummary)

			newrec.PutCopy("field_name", mlrval.FromString(fieldName))

			// Display field type(s) for this column as a list of string, hyphen-joined into a single string.
			if tr.accumulatorNamesToEmit["field_type"] {
				fieldTypesList := make([]string, fieldSummary.fieldTypesMap.FieldCount)
				i := 0
				for pf := fieldSummary.fieldTypesMap.Head; pf != nil; pf = pf.Next {
					fieldType := pf.Key
					fieldTypesList[i] = fieldType
					i++
				}
				newrec.PutCopy("field_type", mlrval.FromString(strings.Join(fieldTypesList, "-")))
			}

			if tr.accumulatorNamesToEmit["count"] {
				newrec.PutCopy("count", fieldSummary.countAccumulator.Emit())
			}
			if tr.accumulatorNamesToEmit["sum"] {
				newrec.PutCopy("sum", fieldSummary.sumAccumulator.Emit())
			}
			if tr.accumulatorNamesToEmit["mean"] {
				newrec.PutCopy("mean", fieldSummary.meanAccumulator.Emit())
			}
			if tr.accumulatorNamesToEmit["stddev"] {
				newrec.PutCopy("stddev", fieldSummary.stddevAccumulator.Emit())
			}
			// TODO: variance
			// TODO: skewness

			// minlen/maxlen

			q1 := fieldSummary.percentileKeeper.Emit(25.0)
			median := fieldSummary.percentileKeeper.Emit(50.0)
			q3 := fieldSummary.percentileKeeper.Emit(75.0)
			max := fieldSummary.percentileKeeper.Emit(100.0)

			iqr := bifs.BIF_minus_binary(q3, q1)
			inner_k := mlrval.FromFloat(1.5)
			outer_k := mlrval.FromFloat(3.0)

			if tr.accumulatorNamesToEmit["min"] {
				min := fieldSummary.percentileKeeper.Emit(0.0)
				newrec.PutCopy("min", min)
			}

			if tr.accumulatorNamesToEmit["p25"] {
				newrec.PutCopy("p25", q1)
			}
			if tr.accumulatorNamesToEmit["median"] {
				newrec.PutCopy("median", median)
			}
			if tr.accumulatorNamesToEmit["p75"] {
				newrec.PutCopy("p75", q3)
			}
			if tr.accumulatorNamesToEmit["max"] {
				newrec.PutCopy("max", max)
			}

			// TODO: leave "" if float-count is zero ...

			if tr.accumulatorNamesToEmit["iqr"] {
				newrec.PutCopy("iqr", iqr)
			}
			if tr.accumulatorNamesToEmit["lof"] {
				lof := bifs.BIF_minus_binary(q1, bifs.BIF_times(outer_k, iqr))
				newrec.PutCopy("lof", lof)
			}
			if tr.accumulatorNamesToEmit["lif"] {

				lif := bifs.BIF_minus_binary(q1, bifs.BIF_times(inner_k, iqr))
				newrec.PutCopy("lif", lif)
			}
			if tr.accumulatorNamesToEmit["uif"] {
				uif := bifs.BIF_plus_binary(q3, bifs.BIF_times(inner_k, iqr))
				newrec.PutCopy("uif", uif)
			}
			if tr.accumulatorNamesToEmit["uof"] {
				uof := bifs.BIF_plus_binary(q3, bifs.BIF_times(outer_k, iqr))
				newrec.PutCopy("uof", uof)
			}

			// TODO: MOVE COMMENT
			// iqr = q3 - q1
			// lof/lif
			// q1 - k*iqr, q3 + k*iqr k=1.5
			// q1 - k*iqr, q3 + k*iqr k=3.0

			if tr.accumulatorNamesToEmit["null_count"] {
				newrec.PutCopy("null_count", mlrval.FromInt(fieldSummary.nullCount))
			}
			if tr.accumulatorNamesToEmit["distinct_count"] {
				newrec.PutCopy("distinct_count", mlrval.FromInt(fieldSummary.distincts.FieldCount))
			}

			// The mode is the most-occurring value for this column. In case of ties, use the first
			// found. We need OrderedMap so regression-test outputs are deterministic in case of
			// ties.
			if tr.accumulatorNamesToEmit["mode"] {
				mode := ""
				var maxCount int64 = 0
				for pf := fieldSummary.distincts.Head; pf != nil; pf = pf.Next {
					distinctValue := pf.Key
					distinctCount := pf.Value.(int64)
					if distinctCount > maxCount {
						maxCount = distinctCount
						mode = distinctValue
					}
				}
				newrec.PutCopy("mode", mlrval.FromString(mode))
			}

			outputRecordsAndContexts.PushBack(types.NewRecordAndContext(newrec, &inrecAndContext.Context))
		}

		outputRecordsAndContexts.PushBack(inrecAndContext) // end-of-stream marker
	}
}

// field_name     a       b       i          x        y
// field_type     string  string  int        float    float

// count          0       0       10000      10000    10000
// sum            0       0       50005000   4986.020 5062.057
// mean           -       -       5000.500   0.499    0.506
// stddev         -       -       2886.896   0.290    0.291
// min            eks     eks     1          0.000    0.000
// p25            hat     hat     2501       0.247    0.252
// median         pan     pan     5001       0.501    0.506
// p75            wye     wye     7501       0.748    0.764
// max            zee     zee     10000      1.000    1.000
// iqr            (error) (error) 5000       0.502    0.512
// lof            (error) (error) -12499.000 -1.258   -1.283
// lif            (error) (error) -4999.000  -0.506   -0.516
// uif            (error) (error) 15001.000  1.500    1.532
// uof            (error) (error) 22501.000  2.253    2.300
// null_count     0       0       0          0        0
// distinct_count 1       1       1          1        1
// mode           string  string  int        float    float

// defaults:

// field_name
// field_type
// count
// mean
// min
// median
// max
// null_count
// distinct_count

// all:
// d field_name
// d field_type
// d count
//   sum
// d mean
//   stddev
// d min
//   p25
// d median
//   p75
// d max
//   iqr -- needs p25 p75
//   lof -- needs p25 p75 iqr
//   lif -- needs p25 p75 iqr
//   uif -- needs p25 p75 iqr
//   uof -- needs p25 p75 iqr
// d null_count
// d distinct_count -- needs distincts
//   mode -- needs distincts
