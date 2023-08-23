package transformers

import (
	"container/list"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
	"github.com/johnkerl/miller/internal/pkg/output"
	"github.com/johnkerl/miller/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameSplit = "split"
const splitDefaultOutputFileNamePrefix = "split"
const splitDefaultFileNamePartJoiner = "_"

var SplitSetup = TransformerSetup{
	Verb:         verbNameSplit,
	UsageFunc:    transformerSplitUsage,
	ParseCLIFunc: transformerSplitParseCLI,
	IgnoresInput: false,
}

func transformerSplitUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options] {filename}\n", "mlr", verbNameSplit)
	fmt.Fprintf(o,
		`Options:
-n {n}:      Cap file sizes at N records.
-m {m}:      Produce M files, round-robining records among them.
-g {a,b,c}:  Write separate files with records having distinct values for fields named a,b,c.
Exactly one  of -m, -n, or -g must be supplied.
--prefix {p} Specify filename prefix; default "`+splitDefaultOutputFileNamePrefix+`".
--suffix {s} Specify filename suffix; default is from mlr output format, e.g. "csv".
-a           Append to existing file(s), if any, rather than overwriting.
-v           Send records along to downstream verbs as well as splitting to files.
-e           Do NOT URL-escape names of output files.
-j {J}       Use string J to join filename parts; default "`+splitDefaultFileNamePartJoiner+`".
-h|--help    Show this message.
Any of the output-format command-line flags (see mlr -h). For example, using
  mlr --icsv --from myfile.csv split --ojson -n 1000
the input is CSV, but the output files are JSON.

Examples: Suppose myfile.csv has 1,000,000 records.

100 output files, 10,000 records each. First 10,000 records in split_1.csv, next in split_2.csv, etc.
  mlr --csv --from myfile.csv split -n 10000

10 output files, 100,000 records each. Records 1,11,21,etc in split_1.csv, records 2,12,22, etc in split_2.csv, etc.
  mlr --csv --from myfile.csv split -m 10
Same, but with JSON output.
  mlr --csv --from myfile.csv split -m 10 -o json

Same but instead of split_1.csv, split_2.csv, etc. there are test_1.dat, test_2.dat, etc.
  mlr --csv --from myfile.csv split -m 10 --prefix test --suffix dat
Same, but written to the /tmp/ directory.
  mlr --csv --from myfile.csv split -m 10 --prefix /tmp/test --suffix dat

If the shape field has values triangle and square, then there will be split_triangle.csv and split_square.csv.
  mlr --csv --from myfile.csv split -g shape

If the color field has values yellow and green, and the shape field has values triangle and square,
then there will be split_yellow_triangle.csv, split_yellow_square.csv, etc.
  mlr --csv --from myfile.csv split -g color,shape

See also the "tee" DSL function which lets you do more ad-hoc customization.
`)
}

func transformerSplitParseCLI(
	pargi *int,
	argc int,
	args []string,
	mainOptions *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	var n int64 = 0
	var doMod bool = false
	var doSize bool = false
	var groupByFieldNames []string = nil
	var emitDownstream bool = false
	var escapeFileNameCharacters bool = true
	var fileNamePartJoiner string = splitDefaultFileNamePartJoiner
	var doAppend bool = false
	var outputFileNamePrefix string = splitDefaultOutputFileNamePrefix
	var outputFileNameSuffix string = "uninit"
	haveOutputFileNameSuffix := false

	var localOptions *cli.TOptions = nil
	if mainOptions != nil {
		copyThereof := *mainOptions // struct copy
		localOptions = &copyThereof
	}

	// Parse local flags.
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
			transformerSplitUsage(os.Stdout)
			os.Exit(0)

		} else if opt == "-n" {
			n = cli.VerbGetIntArgOrDie(verb, opt, args, &argi, argc)
			doSize = true

		} else if opt == "-m" {
			n = cli.VerbGetIntArgOrDie(verb, opt, args, &argi, argc)
			doMod = true

		} else if opt == "-g" {
			groupByFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "--prefix" {
			outputFileNamePrefix = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "--suffix" {
			outputFileNameSuffix = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)
			haveOutputFileNameSuffix = true

		} else if opt == "-a" {
			doAppend = true

		} else if opt == "-v" {
			emitDownstream = true

		} else if opt == "-e" {
			escapeFileNameCharacters = false

		} else if opt == "-j" {
			fileNamePartJoiner = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else {
			// This is inelegant. For error-proofing we advance argi already in our
			// loop (so individual if-statements don't need to). However,
			// ParseWriterOptions expects it unadvanced.
			largi := argi - 1
			if cli.FLAG_TABLE.Parse(args, argc, &largi, localOptions) {
				// This lets mlr main and mlr split have different output formats.
				// Nothing else to handle here.
				argi = largi
			} else {
				transformerSplitUsage(os.Stderr)
				os.Exit(1)
			}
		}
	}

	doGroup := groupByFieldNames != nil
	if !doMod && !doSize && !doGroup {
		fmt.Fprintf(os.Stderr, "mlr %s: At least one of -m, -n, or -g is required.\n", verb)
		os.Exit(1)
	}
	if (doMod && doSize) || (doMod && doGroup) || (doSize && doGroup) {
		fmt.Fprintf(os.Stderr, "mlr %s: Only one of -m, -n, or -g is required.\n", verb)
		os.Exit(1)
	}

	cli.FinalizeWriterOptions(&localOptions.WriterOptions)
	if !haveOutputFileNameSuffix {
		outputFileNameSuffix = localOptions.WriterOptions.OutputFileFormat
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerSplit(
		n,
		doMod,
		doSize,
		groupByFieldNames,
		emitDownstream,
		escapeFileNameCharacters,
		fileNamePartJoiner,
		doAppend,
		outputFileNamePrefix,
		outputFileNameSuffix,
		&localOptions.WriterOptions,
	)
	if err != nil {
		// Error message already printed out
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerSplit struct {
	n                        int64
	outputFileNamePrefix     string
	outputFileNameSuffix     string
	emitDownstream           bool
	escapeFileNameCharacters bool
	fileNamePartJoiner       string
	ungroupedCounter         int64
	groupByFieldNames        []string
	recordWriterOptions      *cli.TWriterOptions
	doAppend                 bool

	// For doSize ungrouped: only one file open at a time
	outputHandler    output.OutputHandler
	previousQuotient int64

	// For all other cases: multiple files open at a time
	outputHandlerManager output.OutputHandlerManager

	recordTransformerFunc RecordTransformerFunc
}

func NewTransformerSplit(
	n int64,
	doMod bool,
	doSize bool,
	groupByFieldNames []string,
	emitDownstream bool,
	escapeFileNameCharacters bool,
	fileNamePartJoiner string,
	doAppend bool,
	outputFileNamePrefix string,
	outputFileNameSuffix string,
	recordWriterOptions *cli.TWriterOptions,
) (*TransformerSplit, error) {

	tr := &TransformerSplit{
		n:                        n,
		outputFileNamePrefix:     outputFileNamePrefix,
		outputFileNameSuffix:     outputFileNameSuffix,
		emitDownstream:           emitDownstream,
		escapeFileNameCharacters: escapeFileNameCharacters,
		fileNamePartJoiner:       fileNamePartJoiner,
		ungroupedCounter:         0,
		groupByFieldNames:        groupByFieldNames,
		recordWriterOptions:      recordWriterOptions,
		doAppend:                 doAppend,

		outputHandler:    nil,
		previousQuotient: -1,
	}

	tr.outputHandlerManager = output.NewFileOutputHandlerManager(recordWriterOptions, doAppend)

	if groupByFieldNames != nil {
		tr.recordTransformerFunc = tr.splitGrouped
	} else if doMod {
		tr.recordTransformerFunc = tr.splitModUngrouped
	} else {
		tr.recordTransformerFunc = tr.splitSizeUngrouped
	}

	return tr, nil
}

func (tr *TransformerSplit) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, outputRecordsAndContexts, inputDownstreamDoneChannel,
		outputDownstreamDoneChannel)
}

func (tr *TransformerSplit) splitModUngrouped(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		remainder := 1 + (tr.ungroupedCounter % tr.n)
		filename := tr.makeUngroupedOutputFileName(remainder)

		err := tr.outputHandlerManager.WriteRecordAndContext(inrecAndContext, filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mlr: file-write error: %v\n", err)
			os.Exit(1)
		}

		if tr.emitDownstream {
			outputRecordsAndContexts.PushBack(inrecAndContext)
		}

		tr.ungroupedCounter++

	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext) // end-of-stream marker
		errs := tr.outputHandlerManager.Close()
		if len(errs) > 0 {
			for _, err := range errs {
				fmt.Fprintf(os.Stderr, "mlr: file-close error: %v\n", err)
			}
			os.Exit(1)
		}
	}
}

func (tr *TransformerSplit) splitSizeUngrouped(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	var err error
	if !inrecAndContext.EndOfStream {
		quotient := 1 + (tr.ungroupedCounter / tr.n)

		if quotient != tr.previousQuotient {
			if tr.outputHandler != nil {
				err = tr.outputHandler.Close()
				if err != nil {
					fmt.Fprintf(os.Stderr, "mlr: file-close error: %v\n", err)
					os.Exit(1)
				}
			}

			filename := tr.makeUngroupedOutputFileName(quotient)
			tr.outputHandler, err = output.NewFileOutputHandler(
				filename,
				tr.recordWriterOptions,
				tr.doAppend,
			)
			if err != nil {
				fmt.Fprintf(os.Stderr, "mlr: file-open error: %v\n", err)
				os.Exit(1)
			}

			tr.previousQuotient = quotient
		}

		err = tr.outputHandler.WriteRecordAndContext(inrecAndContext)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mlr: file-write error: %v\n", err)
			os.Exit(1)
		}

		if tr.emitDownstream {
			outputRecordsAndContexts.PushBack(inrecAndContext)
		}

		tr.ungroupedCounter++

	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext) // end-of-stream marker

		if tr.outputHandler != nil {
			err := tr.outputHandler.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "mlr: file-close error: %v\n", err)
				os.Exit(1)
			}
		}
	}
}

func (tr *TransformerSplit) splitGrouped(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		var filename string
		groupByFieldValues, ok := inrecAndContext.Record.GetSelectedValues(tr.groupByFieldNames)
		if !ok {
			filename = fmt.Sprintf("%s_ungrouped.%s", tr.outputFileNamePrefix, tr.outputFileNameSuffix)
		} else {
			filename = tr.makeGroupedOutputFileName(groupByFieldValues)
		}
		err := tr.outputHandlerManager.WriteRecordAndContext(inrecAndContext, filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
			os.Exit(1)
		}

		if tr.emitDownstream {
			outputRecordsAndContexts.PushBack(inrecAndContext)
		}

	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext) // emit end-of-stream marker

		errs := tr.outputHandlerManager.Close()
		if len(errs) > 0 {
			for _, err := range errs {
				fmt.Fprintf(os.Stderr, "mlr: file-close error: %v\n", err)
			}
			os.Exit(1)
		}
	}
}

// makeUngroupedOutputFileName example: "split_53.csv"
func (tr *TransformerSplit) makeUngroupedOutputFileName(k int64) string {
	return fmt.Sprintf("%s_%d.%s", tr.outputFileNamePrefix, k, tr.outputFileNameSuffix)
}

// makeGroupedOutputFileName example: "split_orange.csv"
func (tr *TransformerSplit) makeGroupedOutputFileName(
	groupByFieldValues []*mlrval.Mlrval,
) string {
	var fileNameParts []string

	for _, groupByFieldValue := range groupByFieldValues {
		fileNameParts = append(fileNameParts, groupByFieldValue.String())
	}

	fileName := strings.Join(fileNameParts, tr.fileNamePartJoiner)

	if tr.escapeFileNameCharacters {
		fileName = url.QueryEscape(fileName)
	}

	if tr.outputFileNamePrefix != "" {
		fileName = tr.outputFileNamePrefix + tr.fileNamePartJoiner + fileName
	}

	return fileName + "." + tr.outputFileNameSuffix
}
