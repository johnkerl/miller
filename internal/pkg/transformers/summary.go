package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
	"github.com/johnkerl/miller/internal/pkg/transformers/utils"
	"github.com/johnkerl/miller/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameSummary = "summary"

type tSummarizerType int

const (
	stFieldType = iota
	stAccumulator
	stPercentile
)

type tSummarizerInfo struct {
	name  string
	help  string
	stype tSummarizerType
}

var allSummarizerInfos = []tSummarizerInfo{
	{"field_type", "string, int, etc. -- if a column has mixed types, all encountered types are printed", stFieldType},

	{"count", "+1 for every instance of the field across all records in the input record stream", stAccumulator},
	{"null_count", "count of field values either empty string or JSON null", stAccumulator},
	{"distinct_count", "count of distinct values for the field", stAccumulator},
	{"mode", "most-frequently-occurring value for the field", stAccumulator},

	{"sum", "sum of field values", stAccumulator},
	{"mean", "mean of the field values", stAccumulator},
	{"stddev", "standard deviation of the field values", stAccumulator},
	{"var", "variance of the field values", stAccumulator},
	{"skewness", "skewness of the field values", stAccumulator},

	{"minlen", "length of shortest string representation for the field", stAccumulator},
	{"maxlen", "length of longest string representation for the field", stAccumulator},

	{"min", "minimum field value", stAccumulator},
	{"p25", "first-quartile field value", stPercentile},
	{"median", "median field value", stPercentile},
	{"p75", "third-quartile field value", stPercentile},
	{"max", "maximum field value", stAccumulator},
	{"iqr", "interquartile range: p75 - p25", stPercentile},
	{"lof", "lower outer fence: p25 - 3.0 * iqr", stPercentile},
	{"lif", "lower inner fence: p25 - 1.5 * iqr", stPercentile},
	{"uif", "upper inner fence: p75 + 1.5 * iqr", stPercentile},
	{"uof", "upper outer fence: p75 + 3.0 * iqr", stPercentile},
}

var summaryDefaultSummarizerNames = []string{
	"field_type",
	"count",
	"mean",
	"min",
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
	for _, info := range allSummarizerInfos {
		fmt.Fprintf(o, "  %-14s  %s\n", info.name, info.help)
	}

	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Default summarizers:\n")
	fmt.Fprintf(o, " ")
	for _, summarizerName := range summaryDefaultSummarizerNames {
		fmt.Fprintf(o, " %s", summarizerName)
	}
	fmt.Fprintf(o, "\n")

	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Notes:\n")
	fmt.Fprintf(o, "* min, p25, median, p75, and max work for strings as well as numbers\n")
	fmt.Fprintf(o, "* Distinct-counts are computed on string representations -- so 4.1 and 4.10 are counted as distinct here.\n")
	fmt.Fprintf(o, "* If the mode is not unique in the input data, the first-encountered value is reported as the mode.\n")
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

	summarizerNames := summaryDefaultSummarizerNames

	allSummarizerNamesList := make([]string, len(allSummarizerInfos))
	for i, info := range allSummarizerInfos {
		allSummarizerNamesList[i] = info.name
	}

	allSummarizerNamesSet := make(map[string]bool)
	for _, summarizerName := range allSummarizerNamesList {
		allSummarizerNamesSet[summarizerName] = true
	}

	transposeOutput := false

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
			summarizerNames = allSummarizerNamesList

		} else if opt == "-a" {
			summarizerNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)
			for _, summarizerName := range summarizerNames {
				if !allSummarizerNamesSet[summarizerName] {
					fmt.Fprintf(os.Stderr, "mlr %s: unrecognized summarizer name %s\n",
						verb, summarizerName,
					)
					os.Exit(1)
				}
			}

		} else if opt == "-x" {
			excludeSummarizerNames := cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)
			for _, excludeSummarizerName := range excludeSummarizerNames {
				if !allSummarizerNamesSet[excludeSummarizerName] {
					fmt.Fprintf(os.Stderr, "mlr %s: unrecognized Summarizer name %s\n",
						verb, excludeSummarizerName,
					)
					os.Exit(1)
				}
			}

			excludeSummarizerNamesSet := make(map[string]bool)
			for _, excludeSummarizerName := range excludeSummarizerNames {
				excludeSummarizerNamesSet[excludeSummarizerName] = true
			}

			summarizerNames = make([]string, 0)
			for _, summarizerName := range allSummarizerNamesList {
				if !excludeSummarizerNamesSet[summarizerName] {
					summarizerNames = append(summarizerNames, summarizerName)
				}
			}

		} else if opt == "--transpose" {
			transposeOutput = true

		} else {
			transformerSummaryUsage(os.Stderr)
			os.Exit(1)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerSummary(summarizerNames, transposeOutput)
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

	accumulators map[string]utils.IStats1Accumulator

	percentileKeeper *utils.PercentileKeeper
}

func newFieldSummary() *tFieldSummary {
	fieldSummary := &tFieldSummary{
		fieldTypesMap: lib.NewOrderedMap(),

		accumulators: make(map[string]utils.IStats1Accumulator),

		// Interpolated percentiles don't play well with string-valued input data
		percentileKeeper: utils.NewPercentileKeeper(false),
	}

	fieldSummary.accumulators["count"] = utils.NewStats1CountAccumulator()
	fieldSummary.accumulators["null_count"] = utils.NewStats1NullCountAccumulator()
	fieldSummary.accumulators["distinct_count"] = utils.NewStats1DistinctCountAccumulator()
	fieldSummary.accumulators["mode"] = utils.NewStats1ModeAccumulator()

	fieldSummary.accumulators["min"] = utils.NewStats1MinAccumulator()
	fieldSummary.accumulators["max"] = utils.NewStats1MaxAccumulator()

	fieldSummary.accumulators["sum"] = utils.NewStats1SumAccumulator()
	fieldSummary.accumulators["mean"] = utils.NewStats1MeanAccumulator()
	fieldSummary.accumulators["stddev"] = utils.NewStats1StddevAccumulator()
	fieldSummary.accumulators["var"] = utils.NewStats1VarAccumulator()
	fieldSummary.accumulators["skewness"] = utils.NewStats1SkewnessAccumulator()

	fieldSummary.accumulators["minlen"] = utils.NewStats1MinLenAccumulator()
	fieldSummary.accumulators["maxlen"] = utils.NewStats1MaxLenAccumulator()

	return fieldSummary
}

type TransformerSummary struct {
	fieldSummaries    *lib.OrderedMap
	summarizerNames   map[string]bool
	hasAnyPercentiles bool
	transposeOutput   bool
}

func NewTransformerSummary(
	summarizerNames []string,
	transposeOutput bool,
) (*TransformerSummary, error) {

	tr := &TransformerSummary{
		fieldSummaries:  lib.NewOrderedMap(),
		summarizerNames: make(map[string]bool),
		transposeOutput: transposeOutput,
	}

	for _, summarizerName := range summarizerNames {
		tr.summarizerNames[summarizerName] = true
	}

	// Different percentile summarizers share the same data structure.
	// If no percentile summarizers are requested, don't ingest percentiles.
	// This is to help running out of memory for large input data files.
	tr.hasAnyPercentiles = false
	for _, info := range allSummarizerInfos {
		if info.stype == stPercentile && tr.summarizerNames[info.name] {
			tr.hasAnyPercentiles = true
			break
		}
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
		tr.ingest(inrecAndContext)
	} else {
		if tr.transposeOutput {
			tr.emitTransposed(inrecAndContext, outputRecordsAndContexts)
		} else {
			tr.emit(inrecAndContext, outputRecordsAndContexts)
		}
	}
}

func (tr *TransformerSummary) ingest(
	inrecAndContext *types.RecordAndContext,
) {
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

		if tr.summarizerNames["field_type"] {
			// Go generics would be grand to put into lib.OrderedMap, but not all platforms
			// are on recent enough Go compiler versions.
			typeName := pe.Value.GetTypeName()
			iValue := fieldSummary.fieldTypesMap.Get(typeName)
			if iValue == nil {
				fieldSummary.fieldTypesMap.Put(typeName, int64(1))
			} else {
				fieldSummary.fieldTypesMap.Put(typeName, iValue.(int64)+1)
			}
		}

		for _, info := range allSummarizerInfos {
			if info.stype == stAccumulator && tr.summarizerNames[info.name] {
				fieldSummary.accumulators[info.name].Ingest(pe.Value)
			}
		}

		if tr.hasAnyPercentiles {
			fieldSummary.percentileKeeper.Ingest(pe.Value)
		}
	}
}

func (tr *TransformerSummary) emit(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
) {

	for pe := tr.fieldSummaries.Head; pe != nil; pe = pe.Next {
		newrec := mlrval.NewMlrmapAsRecord()
		fieldName := pe.Key
		fieldSummary := pe.Value.(*tFieldSummary)

		newrec.PutCopy("field_name", mlrval.FromString(fieldName))

		// Display field type(s) for this column as a list of string, hyphen-joined into a single string.
		if tr.summarizerNames["field_type"] {
			fieldTypesList := make([]string, fieldSummary.fieldTypesMap.FieldCount)
			i := 0
			for pf := fieldSummary.fieldTypesMap.Head; pf != nil; pf = pf.Next {
				fieldType := pf.Key
				fieldTypesList[i] = fieldType
				i++
			}
			newrec.PutCopy("field_type", mlrval.FromString(strings.Join(fieldTypesList, "-")))
		}

		for _, info := range allSummarizerInfos {
			if info.stype == stAccumulator {
				if tr.summarizerNames[info.name] {
					newrec.PutCopy(info.name, fieldSummary.accumulators[info.name].Emit())
				}
			} else if info.stype == stPercentile {
				if tr.summarizerNames[info.name] {
					newrec.PutCopy(info.name, fieldSummary.percentileKeeper.EmitNamed(info.name))
				}
			}
		}

		outputRecordsAndContexts.PushBack(types.NewRecordAndContext(newrec, &inrecAndContext.Context))
	}

	outputRecordsAndContexts.PushBack(inrecAndContext) // end-of-stream marker
}

func (tr *TransformerSummary) emitTransposed(
	inrecAndContext *types.RecordAndContext,
	oracs *list.List, // list of *types.RecordAndContext
) {
	octx := &inrecAndContext.Context

	// Display field type(s) for this column as a list of string, hyphen-joined into a single string.
	if tr.summarizerNames["field_type"] {
		newrec := mlrval.NewMlrmapAsRecord()
		newrec.PutCopy("field_name", mlrval.FromString("field_type"))
		for pe := tr.fieldSummaries.Head; pe != nil; pe = pe.Next {
			fieldSummary := pe.Value.(*tFieldSummary)
			fieldTypesList := make([]string, fieldSummary.fieldTypesMap.FieldCount)
			i := 0
			for pf := fieldSummary.fieldTypesMap.Head; pf != nil; pf = pf.Next {
				fieldType := pf.Key
				fieldTypesList[i] = fieldType
				i++
			}
			newrec.PutCopy(pe.Key, mlrval.FromString(strings.Join(fieldTypesList, "-")))
		}
		oracs.PushBack(types.NewRecordAndContext(newrec, &inrecAndContext.Context))
	}

	for _, info := range allSummarizerInfos {
		if info.stype == stAccumulator {
			tr.maybeEmitAccumulatorTransposed(oracs, octx, info.name)
		} else if info.stype == stPercentile {
			tr.maybeEmitPercentileNameTransposed(oracs, octx, info.name)
		}
	}

	oracs.PushBack(inrecAndContext) // end-of-stream marker
}

// ----------------------------------------------------------------

// maybeEmitPercentileNameTransposed is a helper method for emitTransposed,
// for "count", "sum", "mean", etc.
func (tr *TransformerSummary) maybeEmitAccumulatorTransposed(
	oracs *list.List, // list of *types.RecordAndContext
	octx *types.Context,
	summarizerName string,
) {
	if tr.summarizerNames[summarizerName] {
		newrec := mlrval.NewMlrmapAsRecord()
		newrec.PutCopy("field_name", mlrval.FromString(summarizerName))
		for pe := tr.fieldSummaries.Head; pe != nil; pe = pe.Next {
			fieldSummary := pe.Value.(*tFieldSummary)
			newrec.PutCopy(pe.Key, fieldSummary.accumulators[summarizerName].Emit())
		}
		oracs.PushBack(types.NewRecordAndContext(newrec, octx))
	}
}

// maybeEmitPercentileNameTransposed is a helper method for emitTransposed,
// for "median", "iqr", "uof", etc.
func (tr *TransformerSummary) maybeEmitPercentileNameTransposed(
	oracs *list.List, // list of *types.RecordAndContext
	octx *types.Context,
	summarizerName string,
) {
	if tr.summarizerNames[summarizerName] {
		newrec := mlrval.NewMlrmapAsRecord()
		newrec.PutCopy("field_name", mlrval.FromString(summarizerName))
		for pe := tr.fieldSummaries.Head; pe != nil; pe = pe.Next {
			fieldSummary := pe.Value.(*tFieldSummary)
			newrec.PutCopy(pe.Key, fieldSummary.percentileKeeper.EmitNamed(summarizerName))
		}
		oracs.PushBack(types.NewRecordAndContext(newrec, octx))
	}
}
