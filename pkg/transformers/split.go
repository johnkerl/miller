package transformers

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/output"
	"github.com/johnkerl/miller/v6/pkg/types"
)

const verbNameSplit = "split"
const splitDefaultOutputFileNamePrefix = "split"
const splitDefaultFileNamePartJoiner = "_"

var splitOptions = []OptionSpec{
	{Flag: "-n", Arg: "{n}", Type: "int", Desc: "Cap output file sizes at N records."},
	{Flag: "-m", Arg: "{m}", Type: "int", Desc: "Produce M files, round-robining records among them."},
	{Flag: "-g", Arg: "{a,b,c}", Type: "csv-list", Desc: "Write separate files with records having distinct values for the specified field names."},
	{Flag: "--prefix", Arg: "{p}", Type: "string", Desc: "Output filename prefix. Default \"split\"."},
	{Flag: "--suffix", Arg: "{s}", Type: "string", Desc: "Output filename suffix. Default is from the output format, e.g. \"csv\"."},
	{Flag: "--folder", Arg: "{f}", Type: "filename", Desc: "Output directory. Default is current directory."},
	{Flag: "-a", Type: "bool", Desc: "Append to existing files rather than overwriting."},
	{Flag: "-v", Type: "bool", Desc: "Send records downstream as well as splitting to files."},
	{Flag: "-e", Type: "bool", Desc: "Do NOT URL-escape names of output files."},
	{Flag: "-j", Arg: "{J}", Type: "string", Desc: "String used to join filename parts. Default \"_\"."},
}

var SplitSetup = TransformerSetup{
	Verb:         verbNameSplit,
	UsageFunc:    transformerSplitUsage,
	ParseCLIFunc: transformerSplitParseCLI,
	IgnoresInput: false,
	Options:      splitOptions,
}

func transformerSplitUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options] {filename}\n", "mlr", verbNameSplit)
	WriteVerbOptions(o, splitOptions)
	fmt.Fprintf(o,
		`Exactly one of -m, -n, or -g must be supplied.
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
Or using --folder:
  mlr --csv --from myfile.csv split -m 10 --folder /tmp --prefix test --suffix dat

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
) (RecordTransformer, error) {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	var n int64 = 0
	doMod := false
	doSize := false
	var groupByFieldNames []string = nil
	emitDownstream := false
	escapeFileNameCharacters := true
	fileNamePartJoiner := splitDefaultFileNamePartJoiner
	doAppend := false
	outputFileNamePrefix := splitDefaultOutputFileNamePrefix
	outputFileNameSuffix := "uninit"
	haveOutputFileNameSuffix := false
	outputFolder := ""

	var localOptions *cli.TOptions = nil
	if mainOptions != nil {
		copyThereof := *mainOptions // struct copy
		localOptions = &copyThereof
	}

	// Parse local flags.
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
			transformerSplitUsage(os.Stdout)
			return nil, cli.ErrHelpRequested

		case "-n":
			n, err = cli.VerbGetIntArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			doSize = true

		case "-m":
			n, err = cli.VerbGetIntArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			doMod = true

		case "-g":
			groupByFieldNames, err = cli.VerbGetStringArrayArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		case "--prefix":
			outputFileNamePrefix, err = cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		case "--suffix":
			outputFileNameSuffix, err = cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			haveOutputFileNameSuffix = true

		case "--folder":
			outputFolder, err = cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		case "-a":
			doAppend = true

		case "-v":
			emitDownstream = true

		case "-e":
			escapeFileNameCharacters = false

		case "-j":
			fileNamePartJoiner, err = cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		default:
			// This is inelegant. For error-proofing we advance argi already in our
			// loop (so individual if-statements don't need to). However,
			// ParseWriterOptions expects it unadvanced.
			largi := argi - 1
			if cli.FLAG_TABLE.Parse(args, argc, &largi, localOptions) {
				// This lets mlr main and mlr split have different output formats.
				// Nothing else to handle here.
				argi = largi
			} else {
				return nil, cli.VerbErrorf(verb, "output format not recognized")
			}
		}
	}

	doGroup := groupByFieldNames != nil
	if !doMod && !doSize && !doGroup {
		return nil, cli.VerbErrorf(verb, "-n, -g, or -s is required")
	}
	if (doMod && doSize) || (doMod && doGroup) || (doSize && doGroup) {
		return nil, cli.VerbErrorf(verb, "-n, -g, and -s are mutually exclusive")
	}

	if err := cli.FinalizeWriterOptions(&localOptions.WriterOptions); err != nil {
		return nil, cli.VerbErrorf(verb, "%v", err)
	}
	if !haveOutputFileNameSuffix {
		outputFileNameSuffix = localOptions.WriterOptions.OutputFileFormat
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
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
		outputFolder,
		&localOptions.WriterOptions,
	)
	if err != nil {
		// Error message already printed out
		os.Exit(1)
	}

	return transformer, nil
}

type TransformerSplit struct {
	n                        int64
	outputFileNamePrefix     string
	outputFileNameSuffix     string
	outputFolder             string
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
	outputFolder string,
	recordWriterOptions *cli.TWriterOptions,
) (*TransformerSplit, error) {

	if outputFolder != "" {
		err := os.MkdirAll(outputFolder, 0755)
		if err != nil {
			return nil, fmt.Errorf("mlr split: could not create output folder %s: %w", outputFolder, err)
		}
	}

	tr := &TransformerSplit{
		n:                        n,
		outputFileNamePrefix:     outputFileNamePrefix,
		outputFileNameSuffix:     outputFileNameSuffix,
		outputFolder:             outputFolder,
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
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, outputRecordsAndContexts, inputDownstreamDoneChannel,
		outputDownstreamDoneChannel)
}

func (tr *TransformerSplit) splitModUngrouped(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		remainder := 1 + (tr.ungroupedCounter % tr.n)
		filename := tr.makeUngroupedOutputFileName(remainder)

		err := tr.outputHandlerManager.WriteRecordAndContext(inrecAndContext, filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
			os.Exit(1)
		}

		if tr.emitDownstream {
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
		}

		tr.ungroupedCounter++

	} else {
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // end-of-stream marker
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
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
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
					fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
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
				fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
				os.Exit(1)
			}

			tr.previousQuotient = quotient
		}

		err = tr.outputHandler.WriteRecordAndContext(inrecAndContext)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
			os.Exit(1)
		}

		if tr.emitDownstream {
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
		}

		tr.ungroupedCounter++

	} else {
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // end-of-stream marker

		if tr.outputHandler != nil {
			err := tr.outputHandler.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
				os.Exit(1)
			}
		}
	}
}

func (tr *TransformerSplit) splitGrouped(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		var filename string
		groupByFieldValues, ok := inrecAndContext.Record.GetSelectedValues(tr.groupByFieldNames)
		if !ok {
			baseName := fmt.Sprintf("%s_ungrouped.%s", tr.outputFileNamePrefix, tr.outputFileNameSuffix)
			if tr.outputFolder != "" {
				filename = filepath.Join(tr.outputFolder, baseName)
			} else {
				filename = baseName
			}
		} else {
			filename = tr.makeGroupedOutputFileName(groupByFieldValues)
		}
		err := tr.outputHandlerManager.WriteRecordAndContext(inrecAndContext, filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
			os.Exit(1)
		}

		if tr.emitDownstream {
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
		}

	} else {
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // emit end-of-stream marker

		errs := tr.outputHandlerManager.Close()
		if len(errs) > 0 {
			for _, err := range errs {
				fmt.Fprintf(os.Stderr, "mlr: file-close error: %v\n", err)
			}
			os.Exit(1)
		}
	}
}

// makeUngroupedOutputFileName example: "split_53.csv" or "folder/split_53.csv" with --folder
func (tr *TransformerSplit) makeUngroupedOutputFileName(k int64) string {
	baseName := fmt.Sprintf("%s_%d.%s", tr.outputFileNamePrefix, k, tr.outputFileNameSuffix)
	if tr.outputFolder != "" {
		return filepath.Join(tr.outputFolder, baseName)
	}
	return baseName
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

	baseName := fileName + "." + tr.outputFileNameSuffix
	if tr.outputFolder != "" {
		return filepath.Join(tr.outputFolder, baseName)
	}
	return baseName
}
