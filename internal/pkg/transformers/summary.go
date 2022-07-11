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

type tSummarizerInfo struct {
	name string
	help string
}

var summaryAllAccumulatorInfos = []tSummarizerInfo{
	{"field_type", "string, int, etc. -- if a column has mixed types, all encountered types are printed"},
	{"count", "+1 for every instance of the field across all records in the input record stream"},
	{"sum", "sum of field values"},
	{"mean", "mean of the field values"},
	{"stddev", "standard deviation of the field values"},
	{"var", "variance of the field values"},
	{"skewness", "skewness of the field values"},
	{"min", "minimum field value"},
	{"p25", "first-quartile field value"},
	{"median", "median field value"},
	{"p75", "third-quartile field value"},
	{"max", "maximum field value"},
	{"iqr", "interquartile range: p75 - p25"},
	{"lof", "lower outer fence: p25 - 3.0 * iqr"},
	{"lif", "lower inner fence: p25 - 1.5 * iqr"},
	{"uif", "upper inner fence: p75 + 1.5 * iqr"},
	{"uof", "upper outer fence: p75 + 3.0 * iqr"},
	{"null_count", "count of field values either empty string or JSON null"},
	{"distinct_count", "count of distinct values for the field"},
	{"mode", "most-frequently-occurring value for the field"},
	{"minlen", "length of shortest string representation for the field"},
	{"maxlen", "length of longest string representation for the field"},
}

var summaryAllAccumulatorNames = []string{
	"field_type",
	"count",
	"sum",
	"mean",
	"stddev",
	"var",
	"skewness",
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
	"minlen",
	"maxlen",
}

var summaryDefaultAccumulatorNames = []string{
	"field_type",
	"count",
	"mean",
	"min",
	"median",
	"max",
	"null_count",
	"distinct_count",
}

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
	fmt.Fprintf(o, "Show summary statistics about the input data.\n")

	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "All summarizers:\n")
	for _, info := range summaryAllAccumulatorInfos {
		fmt.Fprintf(o, "  %-14s  %s\n", info.name, info.help)
	}

	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Default summarizers:\n")
	fmt.Fprintf(o, " ")
	for _, accumulatorName := range summaryDefaultAccumulatorNames {
		fmt.Fprintf(o, " %s", accumulatorName)
	}
	fmt.Fprintf(o, "\n")

	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, `Notes:
* min, p25, median, p75, and max work for strings as well as numbers
* Distinct-counts are computed on string representations -- so 4.1 and 4.10 are counted as distinct here.
* If the mode is not unique in the input data, the first-encountered value is reported as the mode.
`)
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-a {mean,sum,etc.} Use only the specified summarizers.\n")
	fmt.Fprintf(o, "-x {mean,sum,etc.} Use all summarizers, except the specified ones.\n")
	fmt.Fprintf(o, "--all              Use all available summarizers.\n")
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

	accumulatorNames := summaryDefaultAccumulatorNames

	allAccumulatorNamesSet := make(map[string]bool)
	for _, accumulatorName := range summaryAllAccumulatorNames {
		allAccumulatorNamesSet[accumulatorName] = true
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
			accumulatorNames = summaryAllAccumulatorNames

		} else if opt == "-a" {
			accumulatorNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)
			for _, accumulatorName := range accumulatorNames {
				if !allAccumulatorNamesSet[accumulatorName] {
					fmt.Fprintf(os.Stderr, "mlr %s: unrecognized accumulator name %s\n",
						verb, accumulatorName,
					)
					os.Exit(1)
				}
			}

		} else if opt == "-x" {
			excludeAccumulatorNames := cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)
			for _, excludeAccumulatorName := range excludeAccumulatorNames {
				if !allAccumulatorNamesSet[excludeAccumulatorName] {
					fmt.Fprintf(os.Stderr, "mlr %s: unrecognized accumulator name %s\n",
						verb, excludeAccumulatorName,
					)
					os.Exit(1)
				}
			}

			excludeAccumulatorNamesSet := make(map[string]bool)
			for _, excludeAccumulatorName := range excludeAccumulatorNames {
				excludeAccumulatorNamesSet[excludeAccumulatorName] = true
			}

			accumulatorNames = make([]string, 0)
			for _, accumulatorName := range summaryAllAccumulatorNames {
				if !excludeAccumulatorNamesSet[accumulatorName] {
					accumulatorNames = append(accumulatorNames, accumulatorName)
				}
			}

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

	countAccumulator    utils.IStats1Accumulator
	sumAccumulator      utils.IStats1Accumulator
	meanAccumulator     utils.IStats1Accumulator
	stddevAccumulator   utils.IStats1Accumulator
	varAccumulator      utils.IStats1Accumulator
	skewnessAccumulator utils.IStats1Accumulator

	percentileKeeper *utils.PercentileKeeper

	nullCountAccumulator utils.IStats1Accumulator
	// Needs lib.OrderedMap, not map[string]int64, for deterministic regression-test output.
	// This is used for distinct_count (map length) as well as mode (max value in the map).
	distincts *lib.OrderedMap

	minlenAccumulator utils.IStats1Accumulator
	maxlenAccumulator utils.IStats1Accumulator
}

func newFieldSummary() *tFieldSummary {
	return &tFieldSummary{
		fieldTypesMap: lib.NewOrderedMap(),

		countAccumulator:    utils.NewStats1CountAccumulator(),
		sumAccumulator:      utils.NewStats1SumAccumulator(),
		meanAccumulator:     utils.NewStats1MeanAccumulator(),
		stddevAccumulator:   utils.NewStats1StddevAccumulator(),
		varAccumulator:      utils.NewStats1VarAccumulator(),
		skewnessAccumulator: utils.NewStats1SkewnessAccumulator(),

		// Interpolated percentiles don't play well with string-valued input data
		percentileKeeper: utils.NewPercentileKeeper(false),

		nullCountAccumulator: utils.NewStats1CountAccumulator(),
		distincts:            lib.NewOrderedMap(),

		minlenAccumulator: utils.NewStats1MinAccumulator(),
		maxlenAccumulator: utils.NewStats1MaxAccumulator(),
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

	// The accumulators we ingest are mostly the same as the ones we emit, except for dependencies.
	// If they asked for iqr but not p25 or p75, we still need to *ingest* p25 and p75 in order to
	// be able to *emit* iqr.
	for _, accumulatorName := range accumulatorNames {
		tr.accumulatorNamesToIngest[accumulatorName] = true

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
				if tr.accumulatorNamesToIngest["var"] {
					fieldSummary.varAccumulator.Ingest(pe.Value)
				}
				if tr.accumulatorNamesToIngest["skewness"] {
					fieldSummary.skewnessAccumulator.Ingest(pe.Value)
				}
			}

			if tr.accumulatorNamesToIngest["null_count"] {
				if pe.Value.IsNull() || pe.Value.IsVoid() {
					fieldSummary.nullCountAccumulator.Ingest(pe.Value)
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

			if tr.accumulatorNamesToIngest["minlen"] {
				fieldSummary.minlenAccumulator.Ingest(bifs.BIF_strlen(mlrval.FromString(pe.Value.OriginalString())))
			}
			if tr.accumulatorNamesToIngest["maxlen"] {
				fieldSummary.maxlenAccumulator.Ingest(bifs.BIF_strlen(mlrval.FromString(pe.Value.OriginalString())))
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
			if tr.accumulatorNamesToEmit["var"] {
				newrec.PutCopy("var", fieldSummary.varAccumulator.Emit())
			}
			if tr.accumulatorNamesToEmit["skewness"] {
				newrec.PutCopy("skewness", fieldSummary.skewnessAccumulator.Emit())
			}

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
			if tr.accumulatorNamesToEmit["null_count"] {
				newrec.PutCopy("null_count", fieldSummary.nullCountAccumulator.Emit())
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

			if tr.accumulatorNamesToEmit["minlen"] {
				newrec.PutCopy("minlen", fieldSummary.minlenAccumulator.Emit())
			}
			if tr.accumulatorNamesToEmit["maxlen"] {
				newrec.PutCopy("maxlen", fieldSummary.maxlenAccumulator.Emit())
			}

			outputRecordsAndContexts.PushBack(types.NewRecordAndContext(newrec, &inrecAndContext.Context))
		}

		// TODO: xout ?

		outputRecordsAndContexts.PushBack(inrecAndContext) // end-of-stream marker
	}
}
