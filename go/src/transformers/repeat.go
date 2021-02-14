package transformers

import (
	"fmt"
	"os"
	"strings"

	"miller/src/cliutil"
	"miller/src/lib"
	"miller/src/transforming"
	"miller/src/types"
)

type tRepeatCountSource int

const (
	repeatCountSourceUnspecified tRepeatCountSource = iota
	repeatCountFromInt
	repeatCountFromFieldName
)

// ----------------------------------------------------------------
const verbNameRepeat = "repeat"

var RepeatSetup = transforming.TransformerSetup{
	Verb:         verbNameRepeat,
	ParseCLIFunc: transformerRepeatParseCLI,
	UsageFunc:    transformerRepeatUsage,

	IgnoresInput: false,
}

func transformerRepeatParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cliutil.TReaderOptions,
	__ *cliutil.TWriterOptions,
) transforming.IRecordTransformer {

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
			repeatCount = cliutil.VerbGetIntArgOrDie(verb, opt, args, &argi, argc)
			repeatCountSource = repeatCountFromInt

		} else if opt == "-f" {
			repeatCountFieldName = cliutil.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)
			repeatCountSource = repeatCountFromFieldName

		} else {
			transformerRepeatUsage(os.Stderr, true, 1)
		}
	}

	if repeatCountSource == repeatCountSourceUnspecified {
		transformerRepeatUsage(os.Stderr, true, 1)
	}

	transformer, _ := NewTransformerRepeat(
		repeatCountSource,
		repeatCount,
		repeatCountFieldName,
	)

	*pargi = argi
	return transformer
}

func transformerRepeatUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", lib.MlrExeName(), verbNameRepeat)
	fmt.Fprintf(o, "Copies input records to output records multiple times.\n")
	fmt.Fprintf(o, "Options must be exactly one of the following:\n")
	fmt.Fprintf(o, "-n {repeat count}  Repeat each input record this many times.\n")
	fmt.Fprintf(o, "-f {field name}    Same, but take the repeat count from the specified\n")
	fmt.Fprintf(o, "                   field name of each input record.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
	fmt.Fprintf(o, "Example:\n")
	fmt.Fprintf(o, "  echo x=0 | %s %s -n 4 then put '$x=urand()'\n", lib.MlrExeName(), verbNameRepeat)
	fmt.Fprintf(o, "produces:\n")
	fmt.Fprintf(o, " x=0.488189\n")
	fmt.Fprintf(o, " x=0.484973\n")
	fmt.Fprintf(o, " x=0.704983\n")
	fmt.Fprintf(o, " x=0.147311\n")
	fmt.Fprintf(o, "Example:\n")
	fmt.Fprintf(o, "  echo a=1,b=2,c=3 | %s %s -f b\n", lib.MlrExeName(), verbNameRepeat)
	fmt.Fprintf(o, "produces:\n")
	fmt.Fprintf(o, "  a=1,b=2,c=3\n")
	fmt.Fprintf(o, "  a=1,b=2,c=3\n")
	fmt.Fprintf(o, "Example:\n")
	fmt.Fprintf(o, "  echo a=1,b=2,c=3 | %s %s -f c\n", lib.MlrExeName(), verbNameRepeat)
	fmt.Fprintf(o, "produces:\n")
	fmt.Fprintf(o, "  a=1,b=2,c=3\n")
	fmt.Fprintf(o, "  a=1,b=2,c=3\n")
	fmt.Fprintf(o, "  a=1,b=2,c=3\n")
	if doExit {
		os.Exit(exitCode)
	}
}

// ----------------------------------------------------------------
type TransformerRepeat struct {
	repeatCount           int
	repeatCountFieldName  string
	recordTransformerFunc transforming.RecordTransformerFunc
}

// ----------------------------------------------------------------
func NewTransformerRepeat(
	repeatCountSource tRepeatCountSource,
	repeatCount int,
	repeatCountFieldName string,
) (*TransformerRepeat, error) {

	this := &TransformerRepeat{
		repeatCount:          repeatCount,
		repeatCountFieldName: repeatCountFieldName,
	}

	if repeatCountSource == repeatCountFromInt {
		this.recordTransformerFunc = this.repeatByCount
	} else {
		this.recordTransformerFunc = this.repeatByFieldName
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerRepeat) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	this.recordTransformerFunc(inrecAndContext, outputChannel)
}

// ----------------------------------------------------------------
func (this *TransformerRepeat) repeatByCount(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		for i := 0; i < this.repeatCount; i++ {
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
func (this *TransformerRepeat) repeatByFieldName(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		fieldValue := inrecAndContext.Record.Get(this.repeatCountFieldName)
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
