package transformers

import (
	"container/list"
	"errors"
	"fmt"
	"os"
	"strings"

	"miller/src/cliutil"
	"miller/src/input"
	"miller/src/lib"
	"miller/src/transformers/utils"
	"miller/src/transforming"
	"miller/src/types"
)

// ----------------------------------------------------------------
const verbNameJoin = "join"

var JoinSetup = transforming.TransformerSetup{
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

	// These allow the joiner to have its own different format/delimiter for the left-file:
	joinReaderOptions cliutil.TReaderOptions
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

		// TODO
		// readerOptions: readerOptions,
	}
}

// ----------------------------------------------------------------
func transformerJoinUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", lib.MlrExeName(), verbNameJoin)
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
		lib.MlrExeName())
	fmt.Fprintf(o, "               If you wish to use a prepipe command for the main input as well\n")
	fmt.Fprintf(o, "               as here, it must be specified there as well as here.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "File-format options default to those for the right file names on the Miller\n")
	fmt.Fprintf(o, "argument list, but may be overridden for the left file as follows. Please see\n")
	fmt.Fprintf(o, "the main \"%s --help\" for more information on syntax for these arguments.\n", lib.MlrExeName())
	fmt.Fprintf(o, "  -i {one of csv,dkvp,nidx,pprint,xtab}\n")
	fmt.Fprintf(o, "  --irs {record-separator character}\n")
	fmt.Fprintf(o, "  --ifs {field-separator character}\n")
	fmt.Fprintf(o, "  --ips {pair-separator character}\n")
	fmt.Fprintf(o, "  --repifs\n")
	fmt.Fprintf(o, "  --repips\n")
	fmt.Fprintf(o, "Please use \"%s --usage-separator-options\" for information on specifying separators.\n",
		lib.MlrExeName())
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
	mainReaderOptions *cliutil.TReaderOptions, // Options for the right-files
	__ *cliutil.TWriterOptions,
) transforming.IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// Parse local flags
	opts := newJoinOptions()

	if mainReaderOptions != nil { // for 'mlr --usage-all-verbs', it's nil
		// TODO: make sure this is a full nested-struct copy.
		opts.joinReaderOptions = *mainReaderOptions // struct copy
	}

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerSortUsage(os.Stdout, true, 0)

		} else if opt == "--prepipe" {
			opts.prepipe = cliutil.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-f" {
			opts.leftFileName = cliutil.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-j" {
			opts.outputJoinFieldNames = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-l" {
			opts.leftJoinFieldNames = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-r" {
			opts.rightJoinFieldNames = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "--lp" {
			opts.leftPrefix = cliutil.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "--rp" {
			opts.rightPrefix = cliutil.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

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
			// ParseReaderOptions expects it unadvanced.
			rargi := argi - 1
			if cliutil.ParseReaderOptions(args, argc, &rargi, &opts.joinReaderOptions) {
				// This lets mlr main and mlr join have different input formats.
				// Nothing else to handle here.
				argi = rargi
			} else {
				transformerJoinUsage(os.Stderr, true, 1)
			}
		}
	}

	if opts.leftFileName == "" {
		fmt.Fprintf(os.Stderr, "%s %s: need left file name\n", lib.MlrExeName(), verb)
		transformerSortUsage(os.Stderr, true, 1)
		return nil
	}

	if !opts.emitPairables && !opts.emitLeftUnpairables && !opts.emitRightUnpairables {
		fmt.Fprintf(os.Stderr, "%s %s: all emit flags are unset; no output is possible.\n",
			lib.MlrExeName(), verb)
		transformerSortUsage(os.Stderr, true, 1)
		return nil
	}

	if opts.outputJoinFieldNames == nil {
		fmt.Fprintf(os.Stderr, "%s %s: need output field names\n", lib.MlrExeName(), verb)
		transformerSortUsage(os.Stderr, true, 1)
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
			lib.MlrExeName(), verb, llen, rlen, olen)
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

	recordTransformerFunc transforming.RecordTransformerFunc
}

// ----------------------------------------------------------------
func NewTransformerJoin(
	opts *tJoinOptions,
) (*TransformerJoin, error) {

	this := &TransformerJoin{
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

		this.leftUnpairableRecordsAndContexts = list.New()
		this.leftBucketsByJoinFieldValues = lib.NewOrderedMap()
		this.recordTransformerFunc = this.transformHalfStreaming

	} else {
		// Doubly-streaming (non-default) case: step left/right files forward.
		// Requires both files be sorted on their join keys in order to not
		// miss anything. This lets people do joins that would otherwise take
		// too much RAM.

		this.joinBucketKeeper = utils.NewJoinBucketKeeper(
			//		opts.prepipe,
			opts.leftFileName,
			&opts.joinReaderOptions,
			opts.leftJoinFieldNames,
		)

		this.recordTransformerFunc = this.transformDoublyStreaming
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerJoin) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	this.recordTransformerFunc(inrecAndContext, outputChannel)
}

// ----------------------------------------------------------------
// This is for the half-streaming case. We ingest the entire left file,
// matching each right record against those.
func (this *TransformerJoin) transformHalfStreaming(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	// This can't be done in the CLI-parser since it requires information which
	// isn't known until after the CLI-parser is called.
	//
	// TODO: check if this is still true for the Go port, once everything else
	// is done.
	if !this.ingested { // First call
		this.ingestLeftFile()
		this.ingested = true
	}

	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		groupingKey, hasAllJoinKeys := inrec.GetSelectedValuesJoined(
			this.opts.rightJoinFieldNames,
		)
		if hasAllJoinKeys {
			iLeftBucket := this.leftBucketsByJoinFieldValues.Get(groupingKey)
			if iLeftBucket == nil {
				if this.opts.emitRightUnpairables {
					outputChannel <- inrecAndContext
				}
			} else {
				leftBucket := iLeftBucket.(*utils.JoinBucket)
				leftBucket.WasPaired = true
				if this.opts.emitPairables {
					this.formAndEmitPairs(
						leftBucket.RecordsAndContexts,
						inrecAndContext,
						outputChannel,
					)
				}
			}
		} else if this.opts.emitRightUnpairables {
			outputChannel <- inrecAndContext
		}

	} else { // end of record stream
		if this.opts.emitLeftUnpairables {
			this.emitLeftUnpairedBuckets(outputChannel)
			this.emitLeftUnpairables(outputChannel)
		}
		outputChannel <- inrecAndContext // emit end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (this *TransformerJoin) transformDoublyStreaming(
	rightRecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	keeper := this.joinBucketKeeper // keystroke-saver

	if !rightRecAndContext.EndOfStream {
		rightRec := rightRecAndContext.Record
		isPaired := false

		rightFieldValues, hasAllJoinKeys := rightRec.ReferenceSelectedValues(
			this.opts.rightJoinFieldNames,
		)
		if hasAllJoinKeys {
			isPaired = keeper.FindJoinBucket(rightFieldValues)
		}
		if this.opts.emitLeftUnpairables {
			keeper.OutputAndReleaseLeftUnpaireds(outputChannel)
		} else {
			keeper.ReleaseLeftUnpaireds(outputChannel)
		}

		lefts := keeper.JoinBucket.RecordsAndContexts // keystroke-saver

		if !isPaired && this.opts.emitRightUnpairables {
			outputChannel <- rightRecAndContext
		}

		if isPaired && this.opts.emitPairables && lefts != nil {
			this.formAndEmitPairs(lefts, rightRecAndContext, outputChannel)
		}

	} else { // end of record stream
		keeper.FindJoinBucket(nil)

		if this.opts.emitLeftUnpairables {
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

func (this *TransformerJoin) ingestLeftFile() {
	readerOpts := &this.opts.joinReaderOptions

	// Instantiate the record-reader
	recordReader := input.Create(readerOpts)
	if recordReader == nil {
		fmt.Fprintln(
			os.Stderr,
			errors.New("Input format not found: "+readerOpts.InputFileFormat),
		)
		os.Exit(1)
	}

	// Set the initial context for the left-file.
	//
	// Since Go is concurrent, the context struct needs to be duplicated and
	// passed through the channels along with each record.
	initialContext := types.NewContext(nil)
	initialContext.UpdateForStartOfFile(this.opts.leftFileName)

	// Set up channels for the record-reader.
	inputChannel := make(chan *types.RecordAndContext, 10)
	errorChannel := make(chan error, 1)

	// Start the record reader.
	// TODO: prepipe
	leftFileNameArray := [1]string{this.opts.leftFileName}
	go recordReader.Read(leftFileNameArray[:], *initialContext, inputChannel, errorChannel)

	// Ingest parsed records and bucket them by their join-field values.  E.g.
	// if the join-field is "id" then put all records with id=1 in one bucket,
	// all those with id=2 in another bucket, etc. And any records lacking an
	// "id" field go into the unpairable list.
	done := false
	for !done {
		select {

		case err := <-errorChannel:
			fmt.Fprintln(os.Stderr, lib.MlrExeName(), ": ", err)
			os.Exit(1)

		case leftrecAndContext := <-inputChannel:
			if leftrecAndContext.EndOfStream {
				done = true
				break // breaks the switch, not the for, in Golang
			}
			leftrec := leftrecAndContext.Record

			groupingKey, leftFieldValues, ok := leftrec.GetSelectedValuesAndJoined(
				this.opts.leftJoinFieldNames,
			)
			if ok {
				iBucket := this.leftBucketsByJoinFieldValues.Get(groupingKey)
				if iBucket == nil { // New key-field-value: new bucket and hash-map entry
					bucket := utils.NewJoinBucket(leftFieldValues)
					bucket.RecordsAndContexts.PushBack(leftrecAndContext)
					this.leftBucketsByJoinFieldValues.Put(groupingKey, bucket)
				} else { // Previously seen key-field-value: append record to bucket
					bucket := iBucket.(*utils.JoinBucket)
					bucket.RecordsAndContexts.PushBack(leftrecAndContext)
				}
			} else {
				this.leftUnpairableRecordsAndContexts.PushBack(leftrecAndContext)
			}
		}
	}
}

// ----------------------------------------------------------------
// This helper method is used by the half-streaming/unsorted join, as well as
// the doubly-streaming/sorted join.

func (this *TransformerJoin) formAndEmitPairs(
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
		n := len(this.opts.leftJoinFieldNames)
		for i := 0; i < n; i++ {
			// These arrays are already guaranteed same-length by CLI parser
			leftJoinFieldName := this.opts.leftJoinFieldNames[i]
			outputJoinFieldName := this.opts.outputJoinFieldNames[i]
			value := leftrec.Get(leftJoinFieldName)
			if value != nil {
				outrec.PutCopy(outputJoinFieldName, value)
			}
		}

		// Add the left-record fields not already added
		for pl := leftrec.Head; pl != nil; pl = pl.Next {
			_, ok := this.leftFieldNameSet[pl.Key]
			if !ok {
				key := this.opts.leftPrefix + pl.Key
				outrec.PutCopy(key, pl.Value)
			}
		}

		// Add the right-record fields not already added
		for pr := rightrec.Head; pr != nil; pr = pr.Next {
			_, ok := this.rightFieldNameSet[pr.Key]
			if !ok {
				key := this.opts.rightPrefix + pr.Key
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

func (this *TransformerJoin) emitLeftUnpairables(
	outputChannel chan<- *types.RecordAndContext,
) {
	// Loop over each to-be-paired-with record from the left file.
	for pe := this.leftUnpairableRecordsAndContexts.Front(); pe != nil; pe = pe.Next() {
		leftRecordAndContext := pe.Value.(*types.RecordAndContext)
		outputChannel <- leftRecordAndContext
	}
}

func (this *TransformerJoin) emitLeftUnpairedBuckets(
	outputChannel chan<- *types.RecordAndContext,
) {
	for pe := this.leftBucketsByJoinFieldValues.Head; pe != nil; pe = pe.Next {
		bucket := pe.Value.(*utils.JoinBucket)
		if !bucket.WasPaired {
			for pf := bucket.RecordsAndContexts.Front(); pf != nil; pf = pf.Next() {
				recordAndContext := pf.Value.(*types.RecordAndContext)
				outputChannel <- recordAndContext
			}
		}
	}
}
