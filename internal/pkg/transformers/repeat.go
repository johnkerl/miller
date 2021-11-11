package transformers

import (
	"fmt"
	"os"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/types"
)

type tRepeatCountSource int

const (
	repeatCountSourceUnspecified tRepeatCountSource = iota
	repeatCountFromInt
	repeatCountFromFieldName
)

// ----------------------------------------------------------------
const verbNameRepeat = "repeat"

var RepeatSetup = TransformerSetup{
	Verb:         verbNameRepeat,
	UsageFunc:    transformerRepeatUsage,
	ParseCLIFunc: transformerRepeatParseCLI,

	IgnoresInput: false,
}

func transformerRepeatUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameRepeat)
	fmt.Fprintf(o, "Copies input records to output records multiple times.\n")
	fmt.Fprintf(o, "Options must be exactly one of the following:\n")
	fmt.Fprintf(o, "-n {repeat count}  Repeat each input record this many times.\n")
	fmt.Fprintf(o, "-f {field name}    Same, but take the repeat count from the specified\n")
	fmt.Fprintf(o, "                   field name of each input record.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
	fmt.Fprintf(o, "Example:\n")
	fmt.Fprintf(o, "  echo x=0 | %s %s -n 4 then put '$x=urand()'\n", "mlr", verbNameRepeat)
	fmt.Fprintf(o, "produces:\n")
	fmt.Fprintf(o, " x=0.488189\n")
	fmt.Fprintf(o, " x=0.484973\n")
	fmt.Fprintf(o, " x=0.704983\n")
	fmt.Fprintf(o, " x=0.147311\n")
	fmt.Fprintf(o, "Example:\n")
	fmt.Fprintf(o, "  echo a=1,b=2,c=3 | %s %s -f b\n", "mlr", verbNameRepeat)
	fmt.Fprintf(o, "produces:\n")
	fmt.Fprintf(o, "  a=1,b=2,c=3\n")
	fmt.Fprintf(o, "  a=1,b=2,c=3\n")
	fmt.Fprintf(o, "Example:\n")
	fmt.Fprintf(o, "  echo a=1,b=2,c=3 | %s %s -f c\n", "mlr", verbNameRepeat)
	fmt.Fprintf(o, "produces:\n")
	fmt.Fprintf(o, "  a=1,b=2,c=3\n")
	fmt.Fprintf(o, "  a=1,b=2,c=3\n")
	fmt.Fprintf(o, "  a=1,b=2,c=3\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerRepeatParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
) IRecordTransformer {

	repeatCountSource := repeatCountSourceUnspecified
	repeatCount := 0
	repeatCountFieldName := ""

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerRepeatUsage(os.Stdout, true, 0)

		} else if opt == "-n" {
			repeatCount = cli.VerbGetIntArgOrDie(verb, opt, args, &argi, argc)
			repeatCountSource = repeatCountFromInt

		} else if opt == "-f" {
			repeatCountFieldName = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)
			repeatCountSource = repeatCountFromFieldName

		} else {
			transformerRepeatUsage(os.Stderr, true, 1)
		}
	}

	if repeatCountSource == repeatCountSourceUnspecified {
		transformerRepeatUsage(os.Stderr, true, 1)
	}

	transformer, err := NewTransformerRepeat(
		repeatCountSource,
		repeatCount,
		repeatCountFieldName,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerRepeat struct {
	repeatCount           int
	repeatCountFieldName  string
	recordTransformerFunc RecordTransformerFunc
}

// ----------------------------------------------------------------
func NewTransformerRepeat(
	repeatCountSource tRepeatCountSource,
	repeatCount int,
	repeatCountFieldName string,
) (*TransformerRepeat, error) {

	tr := &TransformerRepeat{
		repeatCount:          repeatCount,
		repeatCountFieldName: repeatCountFieldName,
	}

	if repeatCountSource == repeatCountFromInt {
		tr.recordTransformerFunc = tr.repeatByCount
	} else {
		tr.recordTransformerFunc = tr.repeatByFieldName
	}

	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerRepeat) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, inputDownstreamDoneChannel, outputDownstreamDoneChannel, outputChannel)
}

// ----------------------------------------------------------------
func (tr *TransformerRepeat) repeatByCount(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		for i := 0; i < tr.repeatCount; i++ {
			outputChannel <- types.NewRecordAndContext(
				inrecAndContext.Record.Copy(),
				&inrecAndContext.Context,
			)
		}
	} else {
		outputChannel <- inrecAndContext
	}
}

// ----------------------------------------------------------------
func (tr *TransformerRepeat) repeatByFieldName(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		fieldValue := inrecAndContext.Record.Get(tr.repeatCountFieldName)
		if fieldValue == nil {
			return
		}
		repeatCount, ok := fieldValue.GetIntValue()
		if !ok {
			return
		}
		for i := 0; i < int(repeatCount); i++ {
			outputChannel <- types.NewRecordAndContext(
				inrecAndContext.Record.Copy(),
				&inrecAndContext.Context,
			)
		}

	} else {
		outputChannel <- inrecAndContext
	}
}
