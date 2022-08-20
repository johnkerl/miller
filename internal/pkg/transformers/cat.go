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
const verbNameCat = "cat"

var CatSetup = TransformerSetup{
	Verb:         verbNameCat,
	UsageFunc:    transformerCatUsage,
	ParseCLIFunc: transformerCatParseCLI,
	IgnoresInput: false,
}

func transformerCatUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameCat)
	fmt.Fprintf(o, "Passes input records directly to output. Most useful for format conversion.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-n         Prepend field \"n\" to each record with record-counter starting at 1.\n")
	fmt.Fprintf(o, "-N {name}  Prepend field {name} to each record with record-counter starting at 1.\n")
	fmt.Fprintf(o, "-g {a,b,c} Optional group-by-field names for counters, e.g. a,b,c\n")
	fmt.Fprintf(o, "--filename Prepend current filename to each record.\n")
	fmt.Fprintf(o, "--filenum  Prepend current filenum (1-up) to each record.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
}

func transformerCatParseCLI(
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

	// Parse local flags
	doCounters := false
	counterFieldName := ""
	var groupByFieldNames []string = nil
	doFileName := false
	doFileNum := false

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
			transformerCatUsage(os.Stdout)
			os.Exit(0)

		} else if opt == "-n" {
			counterFieldName = "n"

		} else if opt == "-N" {
			counterFieldName = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-g" {
			groupByFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "--filename" {
			doFileName = true

		} else if opt == "--filenum" {
			doFileNum = true

		} else {
			transformerCatUsage(os.Stderr)
			os.Exit(1)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerCat(
		doCounters,
		counterFieldName,
		groupByFieldNames,
		doFileName,
		doFileNum,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerCat struct {
	doCounters        bool
	groupByFieldNames []string

	counter          int64
	countsByGroup    map[string]int64
	counterFieldName string

	doFileName bool
	doFileNum  bool

	recordTransformerFunc RecordTransformerFunc
}

// ----------------------------------------------------------------
func NewTransformerCat(
	doCounters bool,
	counterFieldName string,
	groupByFieldNames []string,
	doFileName bool,
	doFileNum bool,
) (*TransformerCat, error) {

	if counterFieldName != "" {
		doCounters = true
	}

	tr := &TransformerCat{
		doCounters:        doCounters,
		groupByFieldNames: groupByFieldNames,
		counter:           0,
		countsByGroup:     make(map[string]int64),
		counterFieldName:  counterFieldName,
		doFileName:        doFileName,
		doFileNum:         doFileNum,
	}

	if !doCounters {
		tr.recordTransformerFunc = tr.simpleCat
	} else {
		if groupByFieldNames == nil {
			tr.recordTransformerFunc = tr.countersUngrouped
		} else {
			tr.recordTransformerFunc = tr.countersGrouped
		}
	}

	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerCat) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(
		inrecAndContext,
		outputRecordsAndContexts,
		inputDownstreamDoneChannel,
		outputDownstreamDoneChannel,
	)
}

// ----------------------------------------------------------------
func (tr *TransformerCat) simpleCat(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		if tr.doFileName {
			inrecAndContext.Record.PrependCopy("filename", mlrval.FromString(inrecAndContext.Context.FILENAME))
		}
		if tr.doFileNum {
			inrecAndContext.Record.PrependCopy("filenum", mlrval.FromInt(inrecAndContext.Context.FILENUM))
		}
	}
	outputRecordsAndContexts.PushBack(inrecAndContext)
}

// ----------------------------------------------------------------
func (tr *TransformerCat) countersUngrouped(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		tr.counter++
		key := tr.counterFieldName
		inrec.PrependCopy(key, mlrval.FromInt(tr.counter))

		if tr.doFileName {
			inrec.PrependCopy("filename", mlrval.FromString(inrecAndContext.Context.FILENAME))
		}
		if tr.doFileNum {
			inrec.PrependCopy("filenum", mlrval.FromInt(inrecAndContext.Context.FILENUM))
		}
	}
	outputRecordsAndContexts.PushBack(inrecAndContext)
}

// ----------------------------------------------------------------
func (tr *TransformerCat) countersGrouped(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		groupingKey, ok := inrec.GetSelectedValuesJoined(tr.groupByFieldNames)
		var counter int64 = 0
		if !ok {
			// Treat as unkeyed
			tr.counter++
			counter = tr.counter
		} else {
			counter, ok = tr.countsByGroup[groupingKey]
			if ok {
				counter++
			} else {
				counter = 1
			}
			tr.countsByGroup[groupingKey] = counter
		}

		key := tr.counterFieldName
		inrec.PrependCopy(key, mlrval.FromInt(counter))

		if tr.doFileName {
			inrec.PrependCopy("filename", mlrval.FromString(inrecAndContext.Context.FILENAME))
		}
		if tr.doFileNum {
			inrec.PrependCopy("filenum", mlrval.FromInt(inrecAndContext.Context.FILENUM))
		}
	}
	outputRecordsAndContexts.PushBack(inrecAndContext)
}
