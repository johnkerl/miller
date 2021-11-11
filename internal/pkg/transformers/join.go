package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/input"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/transformers/utils"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameJoin = "join"

var JoinSetup = TransformerSetup{
	Verb:         verbNameJoin,
	UsageFunc:    transformerJoinUsage,
	ParseCLIFunc: transformerJoinParseCLI,
	IgnoresInput: false,
}

// ----------------------------------------------------------------
// Most transformers have option-variables as individual locals within the
// transformerXYZParseCLI function, which are passed as individual arguments to
// the NewTransformerXYZ function. For join, things are a bit more complex
// and we bag up the option-variables into this data structure.

type tJoinOptions struct {
	leftPrefix  string
	rightPrefix string

	outputJoinFieldNames []string
	leftJoinFieldNames   []string
	rightJoinFieldNames  []string

	allowUnsortedInput   bool
	emitPairables        bool
	emitLeftUnpairables  bool
	emitRightUnpairables bool

	leftFileName string
	prepipe      string
	prepipeIsRaw bool

	// These allow the joiner to have its own different format/delimiter for the left-file:
	joinFlagOptions cli.TOptions
}

func newJoinOptions() *tJoinOptions {
	return &tJoinOptions{
		leftPrefix:  "",
		rightPrefix: "",

		outputJoinFieldNames: nil,
		leftJoinFieldNames:   nil,
		rightJoinFieldNames:  nil,

		allowUnsortedInput:   true,
		emitPairables:        true,
		emitLeftUnpairables:  false,
		emitRightUnpairables: false,

		leftFileName: "",
		prepipe:      "",
		prepipeIsRaw: false,
	}
}

// ----------------------------------------------------------------
func transformerJoinUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameJoin)
	fmt.Fprintf(o, "Joins records from specified left file name with records from all file names\n")
	fmt.Fprintf(o, "at the end of the Miller argument list.\n")
	fmt.Fprintf(o, "Functionality is essentially the same as the system \"join\" command, but for\n")
	fmt.Fprintf(o, "record streams.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "  -f {left file name}\n")
	fmt.Fprintf(o, "  -j {a,b,c}   Comma-separated join-field names for output\n")
	fmt.Fprintf(o, "  -l {a,b,c}   Comma-separated join-field names for left input file;\n")
	fmt.Fprintf(o, "               defaults to -j values if omitted.\n")
	fmt.Fprintf(o, "  -r {a,b,c}   Comma-separated join-field names for right input file(s);\n")
	fmt.Fprintf(o, "               defaults to -j values if omitted.\n")
	fmt.Fprintf(o, "  --lp {text}  Additional prefix for non-join output field names from\n")
	fmt.Fprintf(o, "               the left file\n")
	fmt.Fprintf(o, "  --rp {text}  Additional prefix for non-join output field names from\n")
	fmt.Fprintf(o, "               the right file(s)\n")
	fmt.Fprintf(o, "  --np         Do not emit paired records\n")
	fmt.Fprintf(o, "  --ul         Emit unpaired records from the left file\n")
	fmt.Fprintf(o, "  --ur         Emit unpaired records from the right file(s)\n")
	fmt.Fprintf(o, "  -s|--sorted-input  Require sorted input: records must be sorted\n")
	fmt.Fprintf(o, "               lexically by their join-field names, else not all records will\n")
	fmt.Fprintf(o, "               be paired. The only likely use case for this is with a left\n")
	fmt.Fprintf(o, "               file which is too big to fit into system memory otherwise.\n")
	fmt.Fprintf(o, "  -u           Enable unsorted input. (This is the default even without -u.)\n")
	fmt.Fprintf(o, "               In this case, the entire left file will be loaded into memory.\n")
	fmt.Fprintf(o, "  --prepipe {command} As in main input options; see %s --help for details.\n",
		"mlr")
	fmt.Fprintf(o, "               If you wish to use a prepipe command for the main input as well\n")
	fmt.Fprintf(o, "               as here, it must be specified there as well as here.\n")
	fmt.Fprintf(o, "  --prepipex {command} Likewise.\n")
	fmt.Fprintf(o, "File-format options default to those for the right file names on the Miller\n")
	fmt.Fprintf(o, "argument list, but may be overridden for the left file as follows. Please see\n")
	fmt.Fprintf(o, "the main \"%s --help\" for more information on syntax for these arguments:\n", "mlr")
	fmt.Fprintf(o, "  -i {one of csv,dkvp,nidx,pprint,xtab}\n")
	fmt.Fprintf(o, "  --irs {record-separator character}\n")
	fmt.Fprintf(o, "  --ifs {field-separator character}\n")
	fmt.Fprintf(o, "  --ips {pair-separator character}\n")
	fmt.Fprintf(o, "  --repifs\n")
	fmt.Fprintf(o, "  --repips\n")
	fmt.Fprintf(o, "  --implicit-csv-header\n")
	fmt.Fprintf(o, "  --no-implicit-csv-header\n")
	fmt.Fprintf(o, "For example, if you have 'mlr --csv ... join -l foo ... ' then the left-file format will\n")
	fmt.Fprintf(o, "be specified CSV as well unless you override with 'mlr --csv ... join --ijson -l foo' etc.\n")
	fmt.Fprintf(o, "Likewise, if you have 'mlr --csv --implicit-csv-header ...' then the join-in file will be\n")
	fmt.Fprintf(o, "expected to be headerless as well unless you put '--no-implicit-csv-header' after 'join'.\n")
	fmt.Fprintf(o, "Please use \"%s --usage-separator-options\" for information on specifying separators.\n",
		"mlr")
	fmt.Fprintf(o, "Please see https://miller.readthedocs.io/en/latest/reference-verbs.html#join for more information\n")
	fmt.Fprintf(o, "including examples.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

// ----------------------------------------------------------------
func transformerJoinParseCLI(
	pargi *int,
	argc int,
	args []string,
	mainOptions *cli.TOptions, // Options for the right-files
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// Parse local flags
	opts := newJoinOptions()

	if mainOptions != nil { // for 'mlr --usage-all-verbs', it's nil
		// TODO: make sure this is a full nested-struct copy.
		opts.joinFlagOptions = *mainOptions // struct copy
	}

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerJoinUsage(os.Stdout, true, 0)

		} else if opt == "--prepipe" {
			opts.prepipe = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)
			opts.prepipeIsRaw = false

		} else if opt == "--prepipex" {
			opts.prepipe = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)
			opts.prepipeIsRaw = true

		} else if opt == "-f" {
			opts.leftFileName = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-j" {
			opts.outputJoinFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-l" {
			opts.leftJoinFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-r" {
			opts.rightJoinFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "--lp" {
			opts.leftPrefix = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "--rp" {
			opts.rightPrefix = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "--np" {
			opts.emitPairables = false

		} else if opt == "--ul" {
			opts.emitLeftUnpairables = true

		} else if opt == "--ur" {
			opts.emitRightUnpairables = true

		} else if opt == "-u" {
			opts.allowUnsortedInput = true

		} else if opt == "--sorted-input" || opt == "-s" {
			opts.allowUnsortedInput = false

		} else {
			// This is inelegant. For error-proofing we advance argi already in our
			// loop (so individual if-statements don't need to). However,
			// cli.Parse expects it unadvanced.
			largi := argi - 1
			if cli.FLAG_TABLE.Parse(args, argc, &largi, &opts.joinFlagOptions) {
				// This lets mlr main and mlr join have different input formats.
				// Nothing else to handle here.
				argi = largi
			} else {
				transformerJoinUsage(os.Stderr, true, 1)
			}
		}
	}

	cli.FinalizeReaderOptions(&opts.joinFlagOptions.ReaderOptions)

	if opts.leftFileName == "" {
		fmt.Fprintf(os.Stderr, "%s %s: need left file name\n", "mlr", verb)
		transformerJoinUsage(os.Stderr, true, 1)
		return nil
	}

	if !opts.emitPairables && !opts.emitLeftUnpairables && !opts.emitRightUnpairables {
		fmt.Fprintf(os.Stderr, "%s %s: all emit flags are unset; no output is possible.\n",
			"mlr", verb)
		transformerJoinUsage(os.Stderr, true, 1)
		return nil
	}

	if opts.outputJoinFieldNames == nil {
		fmt.Fprintf(os.Stderr, "%s %s: need output field names\n", "mlr", verb)
		transformerJoinUsage(os.Stderr, true, 1)
		return nil
	}

	if opts.leftJoinFieldNames == nil {
		opts.leftJoinFieldNames = opts.outputJoinFieldNames // array copy
	}
	if opts.rightJoinFieldNames == nil {
		opts.rightJoinFieldNames = opts.outputJoinFieldNames // array copy
	}

	llen := len(opts.leftJoinFieldNames)
	rlen := len(opts.rightJoinFieldNames)
	olen := len(opts.outputJoinFieldNames)
	if llen != rlen || llen != olen {
		fmt.Fprintf(os.Stderr,
			"%s %s: must have equal left,right,output field-name lists; got lengths %d,%d,%d.\n",
			"mlr", verb, llen, rlen, olen)
		os.Exit(1)
	}

	transformer, err := NewTransformerJoin(opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerJoin struct {
	opts *tJoinOptions

	leftFieldNameSet  map[string]bool
	rightFieldNameSet map[string]bool

	// For unsorted/half-streaming input
	ingested                         bool
	leftBucketsByJoinFieldValues     *lib.OrderedMap
	leftUnpairableRecordsAndContexts *list.List

	// For sorted/doubly-streaming input
	joinBucketKeeper *utils.JoinBucketKeeper

	recordTransformerFunc RecordTransformerFunc
}

// ----------------------------------------------------------------
func NewTransformerJoin(
	opts *tJoinOptions,
) (*TransformerJoin, error) {

	tr := &TransformerJoin{
		opts: opts,

		leftFieldNameSet:  lib.StringListToSet(opts.leftJoinFieldNames),
		rightFieldNameSet: lib.StringListToSet(opts.rightJoinFieldNames),

		ingested:                         false,
		leftBucketsByJoinFieldValues:     nil,
		leftUnpairableRecordsAndContexts: nil,
		joinBucketKeeper:                 nil,
	}

	if opts.allowUnsortedInput {
		// Half-streaming (default) case: ingest entire left file first.

		tr.leftUnpairableRecordsAndContexts = list.New()
		tr.leftBucketsByJoinFieldValues = lib.NewOrderedMap()
		tr.recordTransformerFunc = tr.transformHalfStreaming

	} else {
		// Doubly-streaming (non-default) case: step left/right files forward.
		// Requires both files be sorted on their join keys in order to not
		// miss anything. This lets people do joins that would otherwise take
		// too much RAM.

		tr.joinBucketKeeper = utils.NewJoinBucketKeeper(
			//		opts.prepipe,
			opts.leftFileName,
			&opts.joinFlagOptions.ReaderOptions,
			opts.leftJoinFieldNames,
		)

		tr.recordTransformerFunc = tr.transformDoublyStreaming
	}

	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerJoin) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, inputDownstreamDoneChannel, outputDownstreamDoneChannel, outputChannel)
}

// ----------------------------------------------------------------
// This is for the half-streaming case. We ingest the entire left file,
// matching each right record against those.
func (tr *TransformerJoin) transformHalfStreaming(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	// This can't be done in the CLI-parser since it requires information which
	// isn't known until after the CLI-parser is called.
	//
	// TODO: check if this is still true for the Go port, once everything else
	// is done.
	if !tr.ingested { // First call
		tr.ingestLeftFile()
		tr.ingested = true
	}

	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		groupingKey, hasAllJoinKeys := inrec.GetSelectedValuesJoined(
			tr.opts.rightJoinFieldNames,
		)
		if hasAllJoinKeys {
			iLeftBucket := tr.leftBucketsByJoinFieldValues.Get(groupingKey)
			if iLeftBucket == nil {
				if tr.opts.emitRightUnpairables {
					outputChannel <- inrecAndContext
				}
			} else {
				leftBucket := iLeftBucket.(*utils.JoinBucket)
				leftBucket.WasPaired = true
				if tr.opts.emitPairables {
					tr.formAndEmitPairs(
						leftBucket.RecordsAndContexts,
						inrecAndContext,
						outputChannel,
					)
				}
			}
		} else if tr.opts.emitRightUnpairables {
			outputChannel <- inrecAndContext
		}

	} else { // end of record stream
		if tr.opts.emitLeftUnpairables {
			tr.emitLeftUnpairedBuckets(outputChannel)
			tr.emitLeftUnpairables(outputChannel)
		}
		outputChannel <- inrecAndContext // emit end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (tr *TransformerJoin) transformDoublyStreaming(
	rightRecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	keeper := tr.joinBucketKeeper // keystroke-saver

	if !rightRecAndContext.EndOfStream {
		rightRec := rightRecAndContext.Record
		isPaired := false

		rightFieldValues, hasAllJoinKeys := rightRec.ReferenceSelectedValues(
			tr.opts.rightJoinFieldNames,
		)
		if hasAllJoinKeys {
			isPaired = keeper.FindJoinBucket(rightFieldValues)
		}
		if tr.opts.emitLeftUnpairables {
			keeper.OutputAndReleaseLeftUnpaireds(outputChannel)
		} else {
			keeper.ReleaseLeftUnpaireds(outputChannel)
		}

		lefts := keeper.JoinBucket.RecordsAndContexts // keystroke-saver

		if !isPaired && tr.opts.emitRightUnpairables {
			outputChannel <- rightRecAndContext
		}

		if isPaired && tr.opts.emitPairables && lefts != nil {
			tr.formAndEmitPairs(lefts, rightRecAndContext, outputChannel)
		}

	} else { // end of record stream
		keeper.FindJoinBucket(nil)

		if tr.opts.emitLeftUnpairables {
			keeper.OutputAndReleaseLeftUnpaireds(outputChannel)
		}

		outputChannel <- rightRecAndContext // emit end-of-stream marker
	}
}

// ----------------------------------------------------------------
// This is for the half-streaming case. We ingest the entire left file,
// matching each right record against those.
//
// Note: this logic is very similar to that in stream.go, which is what
// processes the main/right files.

func (tr *TransformerJoin) ingestLeftFile() {
	readerOpts := &tr.opts.joinFlagOptions.ReaderOptions

	// Instantiate the record-reader
	recordReader, err := input.Create(readerOpts)
	if recordReader == nil {
		fmt.Fprintf(os.Stderr, "mlr join: %v\n", err)
		os.Exit(1)
	}

	// Set the initial context for the left-file.
	//
	// Since Go is concurrent, the context struct needs to be duplicated and
	// passed through the channels along with each record.
	initialContext := types.NewNilContext()
	initialContext.UpdateForStartOfFile(tr.opts.leftFileName)

	// Set up channels for the record-reader.
	inputChannel := make(chan *types.RecordAndContext, 10)
	errorChannel := make(chan error, 1)
	downstreamDoneChannel := make(chan bool, 1)

	// Start the record reader.
	// TODO: prepipe
	leftFileNameArray := [1]string{tr.opts.leftFileName}
	go recordReader.Read(leftFileNameArray[:], *initialContext, inputChannel, errorChannel, downstreamDoneChannel)

	// Ingest parsed records and bucket them by their join-field values.  E.g.
	// if the join-field is "id" then put all records with id=1 in one bucket,
	// all those with id=2 in another bucket, etc. And any records lacking an
	// "id" field go into the unpairable list.
	done := false
	for !done {
		select {

		case err := <-errorChannel:
			fmt.Fprintln(os.Stderr, "mlr", ": ", err)
			os.Exit(1)

		case leftrecAndContext := <-inputChannel:
			if leftrecAndContext.EndOfStream {
				done = true
				break // breaks the switch, not the for, in Golang
			}
			leftrec := leftrecAndContext.Record

			groupingKey, leftFieldValues, ok := leftrec.GetSelectedValuesAndJoined(
				tr.opts.leftJoinFieldNames,
			)
			if ok {
				iBucket := tr.leftBucketsByJoinFieldValues.Get(groupingKey)
				if iBucket == nil { // New key-field-value: new bucket and hash-map entry
					bucket := utils.NewJoinBucket(leftFieldValues)
					bucket.RecordsAndContexts.PushBack(leftrecAndContext)
					tr.leftBucketsByJoinFieldValues.Put(groupingKey, bucket)
				} else { // Previously seen key-field-value: append record to bucket
					bucket := iBucket.(*utils.JoinBucket)
					bucket.RecordsAndContexts.PushBack(leftrecAndContext)
				}
			} else {
				tr.leftUnpairableRecordsAndContexts.PushBack(leftrecAndContext)
			}
		}
	}
}

// ----------------------------------------------------------------
// This helper method is used by the half-streaming/unsorted join, as well as
// the doubly-streaming/sorted join.

func (tr *TransformerJoin) formAndEmitPairs(
	leftRecordsAndContexts *list.List,
	rightRecordAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	////fmt.Println("-- pairs start") // VERBOSE
	// Loop over each to-be-paired-with record from the left file.
	for pe := leftRecordsAndContexts.Front(); pe != nil; pe = pe.Next() {
		////fmt.Println("-- pairs pe") // VERBOSE
		leftRecordAndContext := pe.Value.(*types.RecordAndContext)
		leftrec := leftRecordAndContext.Record
		rightrec := rightRecordAndContext.Record

		// Allocate a new output record which is the join of the left and right records.
		outrec := types.NewMlrmapAsRecord()

		// Add the joined-on fields to the new output record
		n := len(tr.opts.leftJoinFieldNames)
		for i := 0; i < n; i++ {
			// These arrays are already guaranteed same-length by CLI parser
			leftJoinFieldName := tr.opts.leftJoinFieldNames[i]
			outputJoinFieldName := tr.opts.outputJoinFieldNames[i]
			value := leftrec.Get(leftJoinFieldName)
			if value != nil {
				outrec.PutCopy(outputJoinFieldName, value)
			}
		}

		// Add the left-record fields not already added
		for pl := leftrec.Head; pl != nil; pl = pl.Next {
			_, ok := tr.leftFieldNameSet[pl.Key]
			if !ok {
				key := tr.opts.leftPrefix + pl.Key
				outrec.PutCopy(key, pl.Value)
			}
		}

		// Add the right-record fields not already added
		for pr := rightrec.Head; pr != nil; pr = pr.Next {
			_, ok := tr.rightFieldNameSet[pr.Key]
			if !ok {
				key := tr.opts.rightPrefix + pr.Key
				outrec.PutCopy(key, pr.Value)
			}
		}
		////fmt.Println("-- pairs outrec") // VERBOSE
		////outrec.Print() // VERBOSE

		// Clone the right record's context (NR, FILENAME, etc) to use for the new output record
		context := rightRecordAndContext.Context // struct copy
		outrecAndContext := types.NewRecordAndContext(outrec, &context)

		// Emit the new joined record on the downstream channel
		outputChannel <- outrecAndContext
	}
	////fmt.Println("-- pairs end") // VERBOSE
}

// ----------------------------------------------------------------
// There are two kinds of left non-pair records: (a) those lacking the
// specified join-keys -- can't possibly pair with anything on the right; (b)
// those having the join-keys but not matching with a record on the right.
//
// Example: join on "id" field. Records lacking an "id" field are in the first
// category.  Now suppose there's a left record with id=0, but there were three
// right-file records with id-field values 1,2,3. Then the id=0 left records is
// in the second category.

func (tr *TransformerJoin) emitLeftUnpairables(
	outputChannel chan<- *types.RecordAndContext,
) {
	// Loop over each to-be-paired-with record from the left file.
	for pe := tr.leftUnpairableRecordsAndContexts.Front(); pe != nil; pe = pe.Next() {
		leftRecordAndContext := pe.Value.(*types.RecordAndContext)
		outputChannel <- leftRecordAndContext
	}
}

func (tr *TransformerJoin) emitLeftUnpairedBuckets(
	outputChannel chan<- *types.RecordAndContext,
) {
	for pe := tr.leftBucketsByJoinFieldValues.Head; pe != nil; pe = pe.Next {
		bucket := pe.Value.(*utils.JoinBucket)
		if !bucket.WasPaired {
			for pf := bucket.RecordsAndContexts.Front(); pf != nil; pf = pf.Next() {
				recordAndContext := pf.Value.(*types.RecordAndContext)
				outputChannel <- recordAndContext
			}
		}
	}
}
