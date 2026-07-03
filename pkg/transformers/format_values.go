package transformers

import (
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

const verbNameFormatValues = "format-values"

const defaultFormatValuesStringFormat = "%s"
const defaultFormatValuesIntFormat = "%d"
const defaultFormatValuesFloatFormat = "%f"

var formatValuesOptions = []OptionSpec{
	{Flag: "-i", Arg: "{integer format}", Type: "string", Desc: "Integer format string; defaults to \"%d\". Examples: \"%06lld\", \"%08llx\". Note that Miller integers are long long so you must use formats which apply to long long, e.g. with ll in them. Undefined behavior results otherwise."},
	{Flag: "-f", Arg: "{float format}", Type: "string", Desc: "Float format string; defaults to \"%f\". Examples: \"%8.3lf\", \"%.6le\". Note that Miller floats are double-precision so you must use formats which apply to double, e.g. with l[efg] in them. Undefined behavior results otherwise."},
	{Flag: "-s", Arg: "{string format}", Type: "string", Desc: "String format string; defaults to \"%s\". Examples: \"_%s\", \"%08s\". Note that you must use formats which apply to string, e.g. with s in them. Undefined behavior results otherwise."},
	{Flag: "-n", Type: "bool", Desc: "Coerce field values autodetected as int to float, and then apply the float format."},
}

var FormatValuesSetup = TransformerSetup{
	Verb:         verbNameFormatValues,
	UsageFunc:    transformerFormatValuesUsage,
	ParseCLIFunc: transformerFormatValuesParseCLI,
	IgnoresInput: false,
	Options:      formatValuesOptions,
}

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
	WriteVerbOptions(o, formatValuesOptions)
}

func transformerFormatValuesParseCLI(
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

	stringFormat := defaultFormatValuesStringFormat
	intFormat := defaultFormatValuesIntFormat
	floatFormat := defaultFormatValuesFloatFormat
	coerceIntToFloat := false

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
			transformerFormatValuesUsage(os.Stdout)
			return nil, cli.ErrHelpRequested

		case "-s":
			stringFormat, err = cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
		case "-i":
			intFormat, err = cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
		case "-f":
			floatFormat, err = cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
		case "-n":
			coerceIntToFloat = true

		default:
			return nil, cli.VerbErrorf(verb, "option \"%s\" not recognized", opt)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	transformer, err := NewTransformerFormatValues(
		stringFormat,
		intFormat,
		floatFormat,
		coerceIntToFloat,
	)
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

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

func (tr *TransformerFormatValues) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if inrecAndContext.EndOfStream {
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // emit end-of-stream marker
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

	*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
}
