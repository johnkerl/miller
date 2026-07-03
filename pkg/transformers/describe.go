// The describe verb reports a compact schema/shape for the input data: one
// output record per input field, with the types seen, occurrence counts,
// cardinality, min/max, and -- for low-cardinality fields -- the complete set
// of distinct values.
//
// This is PR6 of the AI-friendly roadmap (issue #2098): `mlr summary` is the
// statistician's view (means, percentiles, outlier fences); `describe` is the
// agent's view -- just enough shape for a person or an LLM agent to construct
// a correct next command. In particular the distinct-value list for
// low-cardinality fields gives the actual domain for flags like `-g` and for
// DSL comparisons, attacking value-hallucination with data-derived values the
// static help catalog cannot know.
//
// Use `mlr --ojson describe ...` for a machine-readable JSON document: the
// per-field `types` map and `values` array nest naturally there, and flatten
// in tabular output formats.

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

const verbNameDescribe = "describe"

const describeDefaultMaxValues = 20

var describeOptions = []OptionSpec{
	{Flag: "-n", Aliases: []string{"--max-values"}, Arg: "{n}", Type: "int",
		Desc: fmt.Sprintf("List a field's distinct values only if it has at most {n} of them; 0 suppresses the values array entirely. Defaults to %d.", describeDefaultMaxValues)},
}

var DescribeSetup = TransformerSetup{
	Verb:         verbNameDescribe,
	UsageFunc:    transformerDescribeUsage,
	ParseCLIFunc: transformerDescribeParseCLI,
	IgnoresInput: false,
	Options:      describeOptions,
}

func transformerDescribeUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameDescribe)
	fmt.Fprintf(o, "Shows a compact schema for the input data: field names, types, and value shape.\n")
	fmt.Fprintf(o, "Emits one output record per input field, with types seen, counts, cardinality,\n")
	fmt.Fprintf(o, "min/max, and (for low-cardinality fields) the complete set of distinct values.\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Output fields, one record per input field:\n")
	fmt.Fprintf(o, "  field_name      name of the input field\n")
	fmt.Fprintf(o, "  types           map from type name (int, string, etc.) to occurrence count\n")
	fmt.Fprintf(o, "  count           number of records in which the field appears\n")
	fmt.Fprintf(o, "  null_count      count of field values either empty string or JSON null\n")
	fmt.Fprintf(o, "  distinct_count  count of distinct values for the field\n")
	fmt.Fprintf(o, "  min, max        minimum/maximum field value (works for strings as well as numbers)\n")
	fmt.Fprintf(o, "  values          all distinct values, in order first seen -- only for fields\n")
	fmt.Fprintf(o, "                  whose distinct_count is within the -n limit\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Notes:\n")
	fmt.Fprintf(o, "* Distinctness is computed on string representations -- so 4.1 and 4.10 are counted as distinct here.\n")
	fmt.Fprintf(o, "* Use `mlr --ojson describe` for a machine-readable JSON document; in tabular\n")
	fmt.Fprintf(o, "  output formats the types map and values array are flattened.\n")
	fmt.Fprintf(o, "* See also the summary verb, which reports summary statistics (mean, percentiles, etc.).\n")
	fmt.Fprintf(o, "\n")
	WriteVerbOptions(o, describeOptions)
}

func transformerDescribeParseCLI(
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

	maxValues := int64(describeDefaultMaxValues)

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
			transformerDescribeUsage(os.Stdout)
			return nil, cli.ErrHelpRequested

		case "-n", "--max-values":
			n, err := cli.VerbGetIntArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			if n < 0 {
				return nil, cli.VerbErrorf(verb, "option \"%s\" requires a non-negative value", opt)
			}
			maxValues = n

		default:
			return nil, cli.VerbErrorf(verb, "option \"%s\" not recognized", opt)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	return NewTransformerDescribe(maxValues)
}

// tFieldDescription accumulates per-field shape information across the record
// stream.
type tFieldDescription struct {
	// Needs lib.OrderedMap, not map[string]int64, for deterministic output.
	// Type name -> occurrence count; a map in case of heterogeneous data.
	typesMap *lib.OrderedMap[int64]

	count     int64
	nullCount int64

	// String representation -> first-encountered value, in order first seen.
	// Distinctness is on string representations, matching the summary verb's
	// distinct_count; keeping the original mlrval preserves int/float/etc.
	// types in the emitted values array.
	distinctValues *lib.OrderedMap[*mlrval.Mlrval]

	minAccumulator utils.IStats1Accumulator
	maxAccumulator utils.IStats1Accumulator
}

func newFieldDescription() *tFieldDescription {
	return &tFieldDescription{
		typesMap:       lib.NewOrderedMap[int64](),
		distinctValues: lib.NewOrderedMap[*mlrval.Mlrval](),
		minAccumulator: utils.NewStats1MinAccumulator(),
		maxAccumulator: utils.NewStats1MaxAccumulator(),
	}
}

type TransformerDescribe struct {
	fieldDescriptions *lib.OrderedMap[*tFieldDescription]
	maxValues         int64
}

func NewTransformerDescribe(
	maxValues int64,
) (*TransformerDescribe, error) {
	return &TransformerDescribe{
		fieldDescriptions: lib.NewOrderedMap[*tFieldDescription](),
		maxValues:         maxValues,
	}, nil
}

func (tr *TransformerDescribe) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		tr.ingest(inrecAndContext)
	} else {
		tr.emit(inrecAndContext, outputRecordsAndContexts)
	}
}

func (tr *TransformerDescribe) ingest(
	inrecAndContext *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record

	for pe := inrec.Head; pe != nil; pe = pe.Next {
		fieldDescription := tr.fieldDescriptions.Get(pe.Key)
		if fieldDescription == nil {
			fieldDescription = newFieldDescription()
			tr.fieldDescriptions.Put(pe.Key, fieldDescription)
		}

		typeName := pe.Value.GetTypeName()
		typeCount, ok := fieldDescription.typesMap.GetWithCheck(typeName)
		if !ok {
			fieldDescription.typesMap.Put(typeName, int64(1))
		} else {
			fieldDescription.typesMap.Put(typeName, typeCount+1)
		}

		fieldDescription.count++
		if pe.Value.IsVoid() || pe.Value.IsNull() {
			fieldDescription.nullCount++
		}

		valueString := pe.Value.String()
		if !fieldDescription.distinctValues.Has(valueString) {
			fieldDescription.distinctValues.Put(valueString, pe.Value.Copy())
		}

		fieldDescription.minAccumulator.Ingest(pe.Value)
		fieldDescription.maxAccumulator.Ingest(pe.Value)
	}
}

func (tr *TransformerDescribe) emit(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
) {
	for pe := tr.fieldDescriptions.Head; pe != nil; pe = pe.Next {
		fieldDescription := pe.Value
		newrec := mlrval.NewMlrmapAsRecord()

		newrec.PutCopy("field_name", mlrval.FromString(pe.Key))

		typesMap := mlrval.NewMlrmap()
		for pf := fieldDescription.typesMap.Head; pf != nil; pf = pf.Next {
			typesMap.PutReference(pf.Key, mlrval.FromInt(pf.Value))
		}
		newrec.PutReference("types", mlrval.FromMap(typesMap))

		newrec.PutReference("count", mlrval.FromInt(fieldDescription.count))
		newrec.PutReference("null_count", mlrval.FromInt(fieldDescription.nullCount))
		newrec.PutReference("distinct_count", mlrval.FromInt(fieldDescription.distinctValues.FieldCount))
		newrec.PutCopy("min", fieldDescription.minAccumulator.Emit())
		newrec.PutCopy("max", fieldDescription.maxAccumulator.Emit())

		if tr.maxValues > 0 && fieldDescription.distinctValues.FieldCount <= tr.maxValues {
			values := make([]*mlrval.Mlrval, 0, fieldDescription.distinctValues.FieldCount)
			for pf := fieldDescription.distinctValues.Head; pf != nil; pf = pf.Next {
				values = append(values, pf.Value)
			}
			newrec.PutReference("values", mlrval.FromArray(values))
		}

		*outputRecordsAndContexts = append(*outputRecordsAndContexts, types.NewRecordAndContext(newrec, &inrecAndContext.Context))
	}

	*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // end-of-stream marker
}
