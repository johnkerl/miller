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
	argi++

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

		} else {
			transformerSummaryUsage(os.Stderr)
			os.Exit(1)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerSummary()
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
	count         int64

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
		count:         0,

		// Interpolated percentiles don't play well with string-valued input data
		percentileKeeper:  utils.NewPercentileKeeper(false),
		meanAccumulator:   utils.NewStats1MeanAccumulator(),
		stddevAccumulator: utils.NewStats1StddevAccumulator(),

		nullCount: 0,
		distincts: lib.NewOrderedMap(),
	}
}

type TransformerSummary struct {
	fieldSummaries *lib.OrderedMap
}

func NewTransformerSummary() (*TransformerSummary, error) {
	tr := &TransformerSummary{
		fieldSummaries: lib.NewOrderedMap(),
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
			typeName := pe.Value.GetTypeName()
			iValue := fieldSummary.fieldTypesMap.Get(typeName)
			if iValue == nil {
				fieldSummary.fieldTypesMap.Put(typeName, int64(1))
			} else {
				fieldSummary.fieldTypesMap.Put(typeName, iValue.(int64)+1)
			}
			fieldSummary.count++

			fieldSummary.percentileKeeper.Ingest(pe.Value)
			if pe.Value.IsNumeric() {
				fieldSummary.meanAccumulator.Ingest(pe.Value)
				fieldSummary.stddevAccumulator.Ingest(pe.Value)
			}

			if pe.Value.IsNull() || pe.Value.IsVoid() {
				fieldSummary.nullCount++
			}
			valueString := pe.Value.String()
			iValue = fieldSummary.distincts.Get(valueString)
			if iValue == nil {
				fieldSummary.distincts.Put(typeName, int64(1))
			} else {
				fieldSummary.distincts.Put(typeName, iValue.(int64)+1)
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
			fieldTypesList := make([]string, fieldSummary.fieldTypesMap.FieldCount)
			i := 0
			for pf := fieldSummary.fieldTypesMap.Head; pf != nil; pf = pf.Next {
				fieldType := pf.Key
				fieldTypesList[i] = fieldType
				i++
			}
			newrec.PutCopy("field_type", mlrval.FromString(strings.Join(fieldTypesList, "-")))

			newrec.PutCopy("count", mlrval.FromInt(fieldSummary.count))
			newrec.PutCopy("min", fieldSummary.percentileKeeper.Emit(0.0))
			newrec.PutCopy("p25", fieldSummary.percentileKeeper.Emit(25.0))
			newrec.PutCopy("median", fieldSummary.percentileKeeper.Emit(50.0))
			newrec.PutCopy("p75", fieldSummary.percentileKeeper.Emit(75.0))
			newrec.PutCopy("max", fieldSummary.percentileKeeper.Emit(100.0))

			newrec.PutCopy("mean", fieldSummary.meanAccumulator.Emit())
			newrec.PutCopy("stddev", fieldSummary.stddevAccumulator.Emit())

			newrec.PutCopy("null_count", mlrval.FromInt(fieldSummary.nullCount))
			newrec.PutCopy("distinct_count", mlrval.FromInt(fieldSummary.distincts.FieldCount))

			// The mode is the most-occurringe value for this column. In case of ties, use the first
			// found. We need OrderedMap so regression-test outputs are deterministic in case of
			// ties.
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

			outputRecordsAndContexts.PushBack(types.NewRecordAndContext(newrec, &inrecAndContext.Context))
		}

		outputRecordsAndContexts.PushBack(inrecAndContext) // end-of-stream marker
	}
}
