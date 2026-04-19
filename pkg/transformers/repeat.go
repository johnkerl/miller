package transformers

import (
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/types"
)

type tRepeatCountSource int

const (
	repeatCountSourceUnspecified tRepeatCountSource = iota
	repeatCountFromInt
	repeatCountFromFieldName
)

const verbNameRepeat = "repeat"

var RepeatSetup = TransformerSetup{
	Verb:         verbNameRepeat,
	UsageFunc:    transformerRepeatUsage,
	ParseCLIFunc: transformerRepeatParseCLI,

	IgnoresInput: false,
}

func transformerRepeatUsage(
	o *os.File,
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
}

func transformerRepeatParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) (RecordTransformer, error) {

	repeatCountSource := repeatCountSourceUnspecified
	repeatCount := int64(0)
	repeatCountFieldName := ""

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

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

		if opt == "-h" || opt == "--help" {
			transformerRepeatUsage(os.Stdout)
			return nil, cli.ErrHelpRequested

		} else if opt == "-n" {
			repeatCount, err = cli.VerbGetIntArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			repeatCountSource = repeatCountFromInt

		} else if opt == "-f" {
			repeatCountFieldName, err = cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			repeatCountSource = repeatCountFromFieldName

		} else {
			return nil, cli.VerbErrorf(verb, "option \"%s\" not recognized", opt)
		}
	}

	if repeatCountSource == repeatCountSourceUnspecified {
		return nil, cli.VerbErrorf(verb, "-n or -f is required")
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	transformer, err := NewTransformerRepeat(
		repeatCountSource,
		repeatCount,
		repeatCountFieldName,
	)
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

type TransformerRepeat struct {
	repeatCount           int64
	repeatCountFieldName  string
	recordTransformerFunc RecordTransformerFunc
}

func NewTransformerRepeat(
	repeatCountSource tRepeatCountSource,
	repeatCount int64,
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

func (tr *TransformerRepeat) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, outputRecordsAndContexts, inputDownstreamDoneChannel, outputDownstreamDoneChannel)
}

func (tr *TransformerRepeat) repeatByCount(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		for i := int64(0); i < tr.repeatCount; i++ {
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, types.NewRecordAndContext(
				inrecAndContext.Record.Copy(),
				&inrecAndContext.Context,
			))
		}
	} else {
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
	}
}

func (tr *TransformerRepeat) repeatByFieldName(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
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
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, types.NewRecordAndContext(
				inrecAndContext.Record.Copy(),
				&inrecAndContext.Context,
			))
		}

	} else {
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
	}
}
