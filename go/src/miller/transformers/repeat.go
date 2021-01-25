package transformers

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"miller/clitypes"
	"miller/transforming"
	"miller/types"
)

type tRepeatCountSource int

const (
	repeatCountSourceUnspecified tRepeatCountSource = iota
	repeatCountFromInt
	repeatCountFromFieldName
)

// ----------------------------------------------------------------
var RepeatSetup = transforming.TransformerSetup{
	Verb:         "repeat",
	ParseCLIFunc: transformerRepeatParseCLI,
	IgnoresInput: false,
}

func transformerRepeatParseCLI(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) transforming.IRecordTransformer {

	repeatCountSource := repeatCountSourceUnspecified
	repeatCount := 0
	repeatCountFieldName := ""

	// Get the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		if !strings.HasPrefix(args[argi], "-") {
			break // No more flag options to process

		} else if args[argi] == "-h" || args[argi] == "--help" {
			transformerPutUsage(os.Stdout, 0, errorHandling, args[0], verb)
			return nil // help intentionally requested

		} else if args[argi] == "-n" {
			checkArgCountRepeat(verb, args, argi, argc, 2)
			n, err := fmt.Sscanf(args[argi+1], "%d", &repeatCount)
			if n != 1 || err != nil {
				transformerPutUsage(os.Stderr, 1, flag.ExitOnError, args[0], verb)
				os.Exit(1)
			}
			repeatCountSource = repeatCountFromInt
			argi += 2
		} else if args[argi] == "-f" {
			checkArgCountRepeat(verb, args, argi, argc, 2)
			repeatCountFieldName = args[argi+1]
			repeatCountSource = repeatCountFromFieldName
			argi += 2

		} else {
			transformerPutUsage(os.Stderr, 1, flag.ExitOnError, args[0], verb)
			os.Exit(1)
		}
	}

	if repeatCountSource == repeatCountSourceUnspecified {
		transformerRepeatUsage(os.Stderr, 1, flag.ExitOnError, args[0], verb)
		os.Exit(1)
	}

	transformer, _ := NewTransformerRepeat(
		repeatCountSource,
		repeatCount,
		repeatCountFieldName,
	)

	*pargi = argi
	return transformer
}

// For flags with values, e.g. ["-n" "10"], while we're looking at the "-n"
// this let us see if the "10" slot exists.
func checkArgCountRepeat(verb string, args []string, argi int, argc int, n int) {
	if (argc - argi) < n {
		fmt.Fprintf(os.Stderr, "%s %s: option \"%s\" missing argument(s).\n",
			args[0], verb, args[argi],
		)
		os.Exit(1)
	}
}

func transformerRepeatUsage(
	o *os.File,
	exitCode int,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	argv0 string,
	verb string,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprintf(o, "Copies input records to output records multiple times.\n")
	fmt.Fprintf(o, "Options must be exactly one of the following:\n")
	fmt.Fprintf(o, "  -n {repeat count}  Repeat each input record this many times.\n")
	fmt.Fprintf(o, "  -f {field name}    Same, but take the repeat count from the specified\n")
	fmt.Fprintf(o, "                     field name of each input record.\n")
	fmt.Fprintf(o, "Example:\n")
	fmt.Fprintf(o, "  echo x=0 | %s %s -n 4 then put '$x=urand()'\n", argv0, verb)
	fmt.Fprintf(o, "produces:\n")
	fmt.Fprintf(o, " x=0.488189\n")
	fmt.Fprintf(o, " x=0.484973\n")
	fmt.Fprintf(o, " x=0.704983\n")
	fmt.Fprintf(o, " x=0.147311\n")
	fmt.Fprintf(o, "Example:\n")
	fmt.Fprintf(o, "  echo a=1,b=2,c=3 | %s %s -f b\n", argv0, verb)
	fmt.Fprintf(o, "produces:\n")
	fmt.Fprintf(o, "  a=1,b=2,c=3\n")
	fmt.Fprintf(o, "  a=1,b=2,c=3\n")
	fmt.Fprintf(o, "Example:\n")
	fmt.Fprintf(o, "  echo a=1,b=2,c=3 | %s %s -f c\n", argv0, verb)
	fmt.Fprintf(o, "produces:\n")
	fmt.Fprintf(o, "  a=1,b=2,c=3\n")
	fmt.Fprintf(o, "  a=1,b=2,c=3\n")
	fmt.Fprintf(o, "  a=1,b=2,c=3\n")
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
