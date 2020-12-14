package transformers

import (
	"container/list"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"miller/clitypes"
	"miller/input"
	"miller/lib"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
var JoinSetup = transforming.TransformerSetup{
	Verb:         "join",
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
	joinReaderOptions clitypes.TReaderOptions
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
func transformerJoinParseCLI(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	mainReaderOptions *clitypes.TReaderOptions, // Options for the right-files
	__ *clitypes.TWriterOptions,
) transforming.IRecordTransformer {

	// Get the verb name from the current spot in the mlr command line
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
		if !strings.HasPrefix(args[argi], "-") {
			break // No more flag options to process
		} else if args[argi] == "-h" || args[argi] == "--help" {
			transformerSortUsage(os.Stdout, 0, errorHandling, args[0], verb)
			return nil // help intentionally requested

		} else if clitypes.ParseReaderOptions(args, argc, &argi, &opts.joinReaderOptions) {
			// handled

		} else if args[argi] == "--prepipe" {
			if (argc - argi) < 2 {
				transformerSortUsage(os.Stdout, 0, errorHandling, args[0], verb)
				return nil
			}
			opts.prepipe = args[argi+1]
			argi += 2

		} else if args[argi] == "-f" {
			if (argc - argi) < 2 {
				transformerSortUsage(os.Stdout, 0, errorHandling, args[0], verb)
				return nil
			}
			opts.leftFileName = args[argi+1]
			argi += 2

		} else if args[argi] == "-j" {
			if (argc - argi) < 2 {
				transformerSortUsage(os.Stdout, 0, errorHandling, args[0], verb)
				return nil
			}
			opts.outputJoinFieldNames = lib.SplitString(args[argi+1], ",")
			argi += 2

		} else if args[argi] == "-l" {
			if (argc - argi) < 2 {
				transformerSortUsage(os.Stdout, 0, errorHandling, args[0], verb)
				return nil
			}
			opts.leftJoinFieldNames = lib.SplitString(args[argi+1], ",")
			argi += 2

		} else if args[argi] == "-r" {
			if (argc - argi) < 2 {
				transformerSortUsage(os.Stdout, 0, errorHandling, args[0], verb)
				return nil
			}
			opts.rightJoinFieldNames = lib.SplitString(args[argi+1], ",")
			argi += 2

		} else if args[argi] == "--lp" {
			if (argc - argi) < 2 {
				transformerSortUsage(os.Stdout, 0, errorHandling, args[0], verb)
				return nil
			}
			opts.leftPrefix = args[argi+1]
			argi += 2

		} else if args[argi] == "--rp" {
			if (argc - argi) < 2 {
				transformerSortUsage(os.Stdout, 0, errorHandling, args[0], verb)
				return nil
			}
			opts.rightPrefix = args[argi+1]
			argi += 2

		} else if args[argi] == "--np" {
			opts.emitPairables = false
			argi += 1

		} else if args[argi] == "--ul" {
			opts.emitLeftUnpairables = true
			argi += 1

		} else if args[argi] == "--ur" {
			opts.emitRightUnpairables = true
			argi += 1

		} else if args[argi] == "-u" {
			opts.allowUnsortedInput = true
			argi += 1

		} else if args[argi] == "--sorted-input" || args[argi] == "-s" {
			opts.allowUnsortedInput = false
			argi += 1

		} else {
			transformerSortUsage(os.Stderr, 1, flag.ExitOnError, args[0], verb)
			os.Exit(1)
		}
	}

	if opts.leftFileName == "" {
		fmt.Fprintf(os.Stderr, "%s %s: need left file name\n", os.Args[0], verb)
		transformerSortUsage(os.Stderr, 1, flag.ExitOnError, args[0], verb)
		return nil
	}

	if !opts.emitPairables && !opts.emitLeftUnpairables && !opts.emitRightUnpairables {
		fmt.Fprintf(os.Stderr, "%s %s: all emit flags are unset; no output is possible.\n",
			os.Args[0], verb)
		transformerSortUsage(os.Stderr, 1, flag.ExitOnError, args[0], verb)
		return nil
	}

	if opts.outputJoinFieldNames == nil {
		fmt.Fprintf(os.Stderr, "%s %s: need output field names\n", os.Args[0], verb)
		transformerSortUsage(os.Stderr, 1, flag.ExitOnError, args[0], verb)
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
			os.Args[0], verb, llen, rlen, olen)
		os.Exit(1)
	}

	transformer, _ := NewTransformerJoin(opts)

	*pargi = argi
	return transformer
}

func transformerJoinUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
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
		os.Args[0])
	fmt.Fprintf(o, "               If you wish to use a prepipe command for the main input as well\n")
	fmt.Fprintf(o, "               as here, it must be specified there as well as here.\n")
	fmt.Fprintf(o, "File-format options default to those for the right file names on the Miller\n")
	fmt.Fprintf(o, "argument list, but may be overridden for the left file as follows. Please see\n")
	fmt.Fprintf(o, "the main \"%s --help\" for more information on syntax for these arguments.\n", argv0)
	fmt.Fprintf(o, "  -i {one of csv,dkvp,nidx,pprint,xtab}\n")
	fmt.Fprintf(o, "  --irs {record-separator character}\n")
	fmt.Fprintf(o, "  --ifs {field-separator character}\n")
	fmt.Fprintf(o, "  --ips {pair-separator character}\n")
	fmt.Fprintf(o, "  --repifs\n")
	fmt.Fprintf(o, "  --repips\n")
	fmt.Fprintf(o, "Please use \"%s --usage-separator-options\" for information on specifying separators.\n",
		argv0)
	fmt.Fprintf(o, "Please see https://miller.readthedocs.io/en/latest/reference-verbs.html#join for more information\n")
	fmt.Fprintf(o, "including examples.\n")
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
	joinBucketKeeper *tJoinBucketKeeper

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

		this.joinBucketKeeper = newJoinBucketKeeper(
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

	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream
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
				leftBucket := iLeftBucket.(*tJoinBucket)
				leftBucket.wasPaired = true
				if this.opts.emitPairables {
					this.formAndEmitPairs(
						leftBucket.recordsAndContexts,
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

	////fmt.Println() // VERBOSE
	////keeper.dump("pre") // VERBOSE

	rightRec := rightRecAndContext.Record

	if rightRec != nil { // not end of record stream

		////fmt.Println("RIGHT REC", rightRec.ToDKVPString()) // VERBOSE

		isPaired := false

		rightFieldValues, hasAllJoinKeys := rightRec.ReferenceSelectedValues(
			this.opts.rightJoinFieldNames,
		)
		if hasAllJoinKeys {
			isPaired = keeper.findJoinBucket(rightFieldValues)
		}
		////fmt.Println("IS_PAIRED", isPaired) // VERBOSE
		////keeper.dump("post") // VERBOSE
		if this.opts.emitLeftUnpairables {
			keeper.outputAndReleaseLeftUnpaireds(outputChannel)
		} else {
			keeper.releaseLeftUnpaireds(outputChannel)
		}

		lefts := keeper.joinBucket.recordsAndContexts // keystroke-saver

		if !isPaired && this.opts.emitRightUnpairables {
			outputChannel <- rightRecAndContext
		}

		if isPaired && this.opts.emitPairables && lefts != nil {
			this.formAndEmitPairs(lefts, rightRecAndContext, outputChannel)
		}

	} else { // end of record stream
		keeper.findJoinBucket(nil)

		if this.opts.emitLeftUnpairables {
			keeper.outputAndReleaseLeftUnpaireds(outputChannel)
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
	initialContext := types.NewContext()
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
			fmt.Fprintln(os.Stderr, os.Args[0], ": ", err)
			os.Exit(1)

		case leftrecAndContext := <-inputChannel:
			leftrec := leftrecAndContext.Record
			if leftrec == nil { // end-of-stream marker
				done = true
				break // breaks the switch, not the for, in Golang
			}

			groupingKey, leftFieldValues, ok := leftrec.GetSelectedValuesAndJoined(
				this.opts.leftJoinFieldNames,
			)
			if ok {
				iBucket := this.leftBucketsByJoinFieldValues.Get(groupingKey)
				if iBucket == nil { // New key-field-value: new bucket and hash-map entry
					bucket := newJoinBucket(leftFieldValues)
					bucket.recordsAndContexts.PushBack(leftrecAndContext)
					this.leftBucketsByJoinFieldValues.Put(groupingKey, bucket)
				} else { // Previously seen key-field-value: append record to bucket
					bucket := iBucket.(*tJoinBucket)
					bucket.recordsAndContexts.PushBack(leftrecAndContext)
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
	//fmt.Println("-- pairs start")
	// Loop over each to-be-paired-with record from the left file.
	for pe := leftRecordsAndContexts.Front(); pe != nil; pe = pe.Next() {
		//fmt.Println("-- pairs pe")
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
			value := leftrec.Get(&leftJoinFieldName)
			if value != nil {
				outrec.PutCopy(&outputJoinFieldName, value)
			}
		}

		// Add the left-record fields not already added
		for pl := leftrec.Head; pl != nil; pl = pl.Next {
			_, ok := this.leftFieldNameSet[*pl.Key]
			if !ok {
				key := this.opts.leftPrefix + *pl.Key
				outrec.PutCopy(&key, pl.Value)
			}
		}

		// Add the right-record fields not already added
		for pr := rightrec.Head; pr != nil; pr = pr.Next {
			_, ok := this.rightFieldNameSet[*pr.Key]
			if !ok {
				key := this.opts.rightPrefix + *pr.Key
				outrec.PutCopy(&key, pr.Value)
			}
		}
		//fmt.Println("-- pairs outrec")
		//outrec.Print()

		// Clone the right record's context (NR, FILENAME, etc) to use for the new output record
		context := rightRecordAndContext.Context // struct copy
		outrecAndContext := types.NewRecordAndContext(outrec, &context)

		// Emit the new joined record on the downstream channel
		outputChannel <- outrecAndContext
	}
	//fmt.Println("-- pairs end")
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
		bucket := pe.Value.(*tJoinBucket)
		if !bucket.wasPaired {
			for pf := bucket.recordsAndContexts.Front(); pf != nil; pf = pf.Next() {
				recordAndContext := pf.Value.(*types.RecordAndContext)
				outputChannel <- recordAndContext
			}
		}
	}
}
