package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
	"github.com/johnkerl/miller/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameFormatValues = "format-values"

const defaultFormatValuesStringFormat = "%s"
const defaultFormatValuesIntFormat = "%d"
const defaultFormatValuesFloatFormat = "%f"

var FormatValuesSetup = TransformerSetup{
	Verb:         verbNameFormatValues,
	UsageFunc:    transformerFormatValuesUsage,
	ParseCLIFunc: transformerFormatValuesParseCLI,
	IgnoresInput: false,
}

// ----------------------------------------------------------------
func transformerFormatValuesUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameFormatValues)
	fmt.Fprintf(o, "Applies format strings to all field values, depending on autodetected type.\n")
	fmt.Fprintf(o, "* If a field value is detected to be integer, applies integer format.\n")
	fmt.Fprintf(o, "* Else, if a field value is detected to be float, applies float format.\n")
	fmt.Fprintf(o, "* Else, applies string format.\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Note: this is a low-keystroke way to apply formatting to many fields. To get\n")
	fmt.Fprintf(o, "finer control, please see the fmtnum function within the mlr put DSL.\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Note: this verb lets you apply arbitrary format strings, which can produce\n")
	fmt.Fprintf(o, "undefined behavior and/or program crashes.  See your system's \"man printf\".\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-i {integer format} Defaults to \"%s\".\n", defaultFormatValuesIntFormat)
	fmt.Fprintf(o, "                    Examples: \"%%06lld\", \"%%08llx\".\n")
	fmt.Fprintf(o, "                    Note that Miller integers are long long so you must use\n")
	fmt.Fprintf(o, "                    formats which apply to long long, e.g. with ll in them.\n")
	fmt.Fprintf(o, "                    Undefined behavior results otherwise.\n")
	fmt.Fprintf(o, "-f {float format}   Defaults to \"%s\".\n", defaultFormatValuesFloatFormat)
	fmt.Fprintf(o, "                    Examples: \"%%8.3lf\", \"%%.6le\".\n")
	fmt.Fprintf(o, "                    Note that Miller floats are double-precision so you must\n")
	fmt.Fprintf(o, "                    use formats which apply to double, e.g. with l[efg] in them.\n")
	fmt.Fprintf(o, "                    Undefined behavior results otherwise.\n")
	fmt.Fprintf(o, "-s {string format}  Defaults to \"%s\".\n", defaultFormatValuesStringFormat)
	fmt.Fprintf(o, "                    Examples: \"_%%s\", \"%%08s\".\n")
	fmt.Fprintf(o, "                    Note that you must use formats which apply to string, e.g.\n")
	fmt.Fprintf(o, "                    with s in them. Undefined behavior results otherwise.\n")
	fmt.Fprintf(o, "-n                  Coerce field values autodetected as int to float, and then\n")
	fmt.Fprintf(o, "                    apply the float format.\n")
}

func transformerFormatValuesParseCLI(
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

	stringFormat := defaultFormatValuesStringFormat
	intFormat := defaultFormatValuesIntFormat
	floatFormat := defaultFormatValuesFloatFormat
	coerceIntToFloat := false

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
			transformerFormatValuesUsage(os.Stdout)
			os.Exit(0)

		} else if opt == "-s" {
			stringFormat = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)
		} else if opt == "-i" {
			intFormat = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)
		} else if opt == "-f" {
			floatFormat = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)
		} else if opt == "-n" {
			coerceIntToFloat = true

		} else {
			transformerFormatValuesUsage(os.Stderr)
			os.Exit(1)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerFormatValues(
		stringFormat,
		intFormat,
		floatFormat,
		coerceIntToFloat,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerFormatValues struct {
	stringFormatter  mlrval.IFormatter
	intFormatter     mlrval.IFormatter
	floatFormatter   mlrval.IFormatter
	coerceIntToFloat bool
}

func NewTransformerFormatValues(
	stringFormat string,
	intFormat string,
	floatFormat string,
	coerceIntToFloat bool,
) (*TransformerFormatValues, error) {
	stringFormatter, err := mlrval.GetFormatter(stringFormat)
	if err != nil {
		return nil, err
	}

	intFormatter, err := mlrval.GetFormatter(intFormat)
	if err != nil {
		return nil, err
	}

	floatFormatter, err := mlrval.GetFormatter(floatFormat)
	if err != nil {
		return nil, err
	}

	tr := &TransformerFormatValues{
		stringFormatter:  stringFormatter,
		intFormatter:     intFormatter,
		floatFormatter:   floatFormatter,
		coerceIntToFloat: coerceIntToFloat,
	}
	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerFormatValues) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if inrecAndContext.EndOfStream {
		outputRecordsAndContexts.PushBack(inrecAndContext) // emit end-of-stream marker
		return
	}

	for pe := inrecAndContext.Record.Head; pe != nil; pe = pe.Next {
		if tr.coerceIntToFloat {
			_, isNumeric := pe.Value.GetNumericToFloatValue()
			if isNumeric {
				pe.Value = tr.floatFormatter.Format(pe.Value)
			} else if pe.Value.IsStringOrVoid() {
				pe.Value = tr.stringFormatter.Format(pe.Value)
			} // else, don't rewrite booleans, arrays, maps, etc.
		} else {
			_, isInt := pe.Value.GetIntValue()
			_, isFloat := pe.Value.GetFloatValue()
			if isInt {
				pe.Value = tr.intFormatter.Format(pe.Value)
			} else if isFloat {
				pe.Value = tr.floatFormatter.Format(pe.Value)
			} else if pe.Value.IsStringOrVoid() {
				pe.Value = tr.stringFormatter.Format(pe.Value)
			} // else, don't rewrite booleans, arrays, maps, etc.
		}
	}

	outputRecordsAndContexts.PushBack(inrecAndContext)
}
