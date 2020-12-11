// ================================================================
// JOIN_BUCKET_KEEPER
//
// This data structure supports Miller's sorted (double-streaming) join.  It is
// perhaps best explained by first comparing with the unsorted (half-streaming)
// case (see join.go).
//
// In both cases, we have left and right join keys. Suppose the left file has
// data with field name "L" to be joined with right-file(s) data with field
// name "R". For the unsorted case (see mapper_join.c) the entire left file is
// first loaded into buckets of record-lists, one for each distinct value of L.
// E.g. given the following:
//
//   +-----+-----+
//   |  L  |  R  |
//   + --- + --- +
//   |  a  |  a  |
//   |  c  |  b  |
//   |  a  |  f  |
//   |  b  |     |
//   |  c  |     |
//   |  d  |     |
//   |  a  |     |
//   +-----+-----+
//
// the left file is bucketed as
//
//   +-----+     +-----+     +-----+     +-----+
//   |  L  |     |  L  |     |  L  |     |  L  |
//   + --- +     + --- +     + --- +     + --- +
//   |  a  |     |  c  |     |  b  |     |  d  |
//   |  a  |     |  c  |     +-----+     +-----+
//   |  a  |     + --- +
//   + --- +
//
// Then the right file is processed one record at a time (hence
// "half-streaming"). The pairings are easy:
// * the right record with R=a is paired with the L=a bucket,
// * the right record with R=b is paired with the L=b bucket,
// * the right record with R=f is unpaired, and
// * the left records with L=c and L=d are unpaired.
//
// ----------------------------------------------------------------
// Now for the sorted (doubly-streaming) case. Here we require that the left
// and right files be already sorted (lexically ascending) by the join fields.
// Then the example inputs look like this:
//
//   +-----+-----+
//   |  L  |  R  |
//   + --- + --- +
//   |  a  |  a  |
//   |  a  |  b  |
//   |  a  |  f  |
//   |  b  |     |
//   |  c  |     |
//   |  c  |     |
//   |  d  |     |
//   +-----+-----+
//
// The right file is still read one record at a time. It's the job of this
// join-bucket-keeper class to keep track of the left-file buckets, one bucket at
// a time.  This includes all records with same values for the join field(s),
// e.g. the three L=a records, as well as a "peek" record which is either the
// next record with a different join value (e.g. the L=b record), or an
// end-of-file indicator.
//
// If a right-file record has join field matching the current left-file bucket,
// then it's paired with all records in that bucket. Otherwise the
// join-bucket-keeper needs to either stay with the current bucket or advance
// to the next one, depending whether the current right-file record's
// join-field values compare lexically with the the left-file bucket's
// join-field values.
//
// Examples:
//
// +-----------+-----------+-----------+-----------+-----------+-----------+
// |  L    R   |   L   R   |   L   R   |   L   R   |   L   R   |   L   R   |
// + ---  ---  + ---  ---  + ---  ---  + ---  ---  + ---  ---  + ---  ---  +
// |       a   |       a   |   e       |       a   |   e   e   |   e   e   |
// |       b   |   e       |   e       |   e   e   |   e       |   e   e   |
// |   e       |   e       |   e       |   e       |   e       |   e       |
// |   e       |   e       |       f   |   e       |       f   |   g   g   |
// |   e       |       f   |   g       |   g       |   g       |   g       |
// |   g       |   g       |   g       |   g       |   g       |           |
// |   g       |   g       |       h   |           |           |           |
// +-----------+-----------+-----------+-----------+-----------+-----------+
//
// In all these examples, the join-bucket-keeper goes through these steps:
// * bucket is empty, peek rec has L==e
// * bucket is L==e records, peek rec has L==g
// * bucket is L==g records, peek rec is null (due to EOF)
// * bucket is empty, peek rec is null (due to EOF)
//
// Example 1:
// * left-bucket is empty and left-peek has L==e
// * right record has R==a; join-bucket-keeper does not advance
// * right record has R==b; join-bucket-keeper does not advance
// * right end of file; all left records are unpaired.
//
// Example 2:
// * left-bucket is empty and left-peek has L==e
// * right record has R==a; join-bucket-keeper does not advance
// * right record has R==f; left records with L==e are unpaired.
// * etc.
//
// ================================================================

package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"miller/clitypes"
	"miller/input"
	"miller/lib"
	"miller/types"
)

// ----------------------------------------------------------------
// Data stored in this class
type tJoinBucketKeeper struct {
	// For streaming through the left-side file
	recordReader input.IRecordReader
	context      *types.Context
	inputChannel <-chan *types.RecordAndContext
	errorChannel chan error
	// TODO: merge with leof flag
	recordReaderDone bool

	leftJoinFieldNames []string

	peekRecordAndContext *types.RecordAndContext
	joinBucket           *tJoinBucket
	leftUnpaireds        *list.List

	leof  bool
	state tJoinBucketKeeperState
}

// ----------------------------------------------------------------
func newJoinBucketKeeper(
	// TODO prepipe string,
	leftFileName string,
	joinReaderOptions *clitypes.TReaderOptions,
	leftJoinFieldNames []string,
) *tJoinBucketKeeper {

	// Instantiate the record-reader
	recordReader := input.Create(joinReaderOptions)
	if recordReader == nil {
		fmt.Fprintf(
			os.Stderr,
			"%s join: Input format %s not found.\n",
			os.Args[0],
			joinReaderOptions.InputFileFormat,
		)
		os.Exit(1)
	}

	// Set the initial context for the left-file.  Since Go is concurrent, the
	// context struct needs to be duplicated and passed through the channels
	// along with each record.
	initialContext := types.NewContext()
	initialContext.UpdateForStartOfFile(leftFileName)

	// Set up channels for the record-reader
	inputChannel := make(chan *types.RecordAndContext, 10)
	errorChannel := make(chan error, 1)

	// Start the record-reader in its own goroutine.
	leftFileNameArray := [1]string{leftFileName}
	go recordReader.Read(leftFileNameArray[:], *initialContext, inputChannel, errorChannel)

	this := &tJoinBucketKeeper{
		recordReader:     recordReader,
		context:          initialContext,
		inputChannel:     inputChannel,
		errorChannel:     errorChannel,
		recordReaderDone: false,

		leftJoinFieldNames: leftJoinFieldNames,

		peekRecordAndContext: nil,
		joinBucket:           newJoinBucket(nil),
		leftUnpaireds:        nil,

		leof:  false,
		state: LEFT_STATE_0_PREFILL,
	}

	return this
}

// ----------------------------------------------------------------
// Returns nil at EOF

func (this *tJoinBucketKeeper) readRecord() *types.RecordAndContext {
	if this.recordReaderDone {
		return nil
	}

	select {

	case err := <-this.errorChannel:
		fmt.Fprintln(os.Stderr, os.Args[0], ": ", err)
		os.Exit(1)

	case leftrecAndContext := <-this.inputChannel:
		if leftrecAndContext.Record == nil { // end-of-stream marker
			this.recordReaderDone = true
			return nil
		} else {
			return leftrecAndContext
		}
	}

	return nil
}

// ----------------------------------------------------------------
// For tJoinBucketKeeper state machine
type tJoinBucketKeeperState int

// (0) pre-fill:    Lv == null, peek == null, leof = false
// (1) midstream:   Lv != null, peek != null, leof = false
// (2) last bucket: Lv != null, peek == null, leof = true
// (3) leof:        Lv == null, peek == null, leof = true
const (
	LEFT_STATE_0_PREFILL     = 0
	LEFT_STATE_1_FULL        = 1
	LEFT_STATE_2_LAST_BUCKET = 2
	LEFT_STATE_3_EOF         = 3
)

func (this *tJoinBucketKeeper) computeState() tJoinBucketKeeperState {
	if this.joinBucket.leftFieldValues == nil {
		if this.leof {
			return LEFT_STATE_3_EOF
		} else {
			return LEFT_STATE_0_PREFILL
		}
	} else {
		if this.peekRecordAndContext == nil {
			return LEFT_STATE_2_LAST_BUCKET
		} else {
			return LEFT_STATE_1_FULL
		}
	}
}

// ----------------------------------------------------------------
func (this *tJoinBucketKeeper) advanceOrDrain(
	rightFieldValues []*types.Mlrval, // nil means right-file EOF
) {
	this.leftUnpaireds = nil

	if this.state == LEFT_STATE_0_PREFILL {
		this.initialFill()
		this.state = this.computeState()
	}

	if rightFieldValues != nil { // Not right EOF
		if this.state == LEFT_STATE_1_FULL || this.state == LEFT_STATE_2_LAST_BUCKET {
			cmp := compareLexicallyVP(this.joinBucket.leftFieldValues, rightFieldValues)
			if cmp < 0 {
				// Advance left until match or left EOF.
				this.advanceTo(rightFieldValues)
			} else if cmp == 0 {
				this.joinBucket.wasPaired = true
			} else {
				// No match and no need to advance left; return null lists.
			}
		} else if this.state != LEFT_STATE_3_EOF {
			fmt.Fprintf(
				os.Stderr,
				"%s: internal coding error: failed transition from prefill state.\n",
				os.Args[0],
			)
			os.Exit(1)
		}

	} else { // Right EOF: return the final left-unpaireds.
		this.drain(rightFieldValues)
	}

	this.state = this.computeState()
}

// ----------------------------------------------------------------
func (this *tJoinBucketKeeper) initialFill() {
	for {
		// Skip over records not having the join keys. These go straight to the
		// left-unpaired list.
		this.peekRecordAndContext = this.readRecord()
		if this.peekRecordAndContext == nil { // left EOF
			break
		}
		if this.peekRecordAndContext.Record.HasSelectedKeys(this.leftJoinFieldNames) {
			break
		} else {
			if this.leftUnpaireds == nil {
				this.leftUnpaireds = list.New()
			}
			this.leftUnpaireds.PushBack(this.peekRecordAndContext)
		}
	}

	if this.peekRecordAndContext == nil {
		this.leof = true
		return
	}
	this.fill()
}

// ----------------------------------------------------------------
// Preconditions:
// * peekRecordAndContext != nil
// * peekRecordAndContext has the join keys

func (this *tJoinBucketKeeper) fill() {
	peekRec := this.peekRecordAndContext.Record
	leftFieldValues, hasAllJoinKeys := peekRec.ReferenceSelectedValues(
		this.leftJoinFieldNames,
	)

	if !hasAllJoinKeys {
		fmt.Fprintf(
			os.Stderr,
			"%s: internal coding error: peek record should have had join keys.\n",
			os.Args[0],
		)
		os.Exit(1)
	}

	this.joinBucket.leftFieldValues = types.CopyMlrvalPointerArray(leftFieldValues)
	this.joinBucket.recordsAndContexts.PushBack(this.peekRecordAndContext)
	this.peekRecordAndContext = nil
	this.joinBucket.wasPaired = false

	for {
		// Skip over records not having the join keys. These go straight to the
		// left-unpaired list.
		this.peekRecordAndContext = this.readRecord()
		if this.peekRecordAndContext == nil { // left EOF
			this.leof = true
			break
		}

		peekRec := this.peekRecordAndContext.Record
		peekFieldValues, hasAllJoinKeys := peekRec.ReferenceSelectedValues(
			this.leftJoinFieldNames,
		)

		if hasAllJoinKeys {
			cmp := compareLexicallyVP(
				this.joinBucket.leftFieldValues,
				peekFieldValues,
			)
			if cmp != 0 {
				break
			}
			this.joinBucket.recordsAndContexts.PushBack(this.peekRecordAndContext)
		} else {
			if this.leftUnpaireds == nil {
				this.leftUnpaireds = list.New()
			}
			this.leftUnpaireds.PushBack(this.peekRecordAndContext)
		}
		this.peekRecordAndContext = nil
	}
}

// ----------------------------------------------------------------
// Pre-conditions:
// * this.leftFieldValues < rightFieldValues.
// * currently in state 1 or 2 so there is a bucket but there may or may not be a peek-record
// * current bucket was/wasn't paired on previous emits but is not paired on this emit.
// Actions:
// * if bucket was never paired, return it to the caller; else discard.
// * consume left input stream, feeding into unpaired, for as long as leftvals < rightvals && !eof.
// * if there is leftrec with vals == rightvals: parallel initial_fill.
//   else, mimic initial_fill.

func (this *tJoinBucketKeeper) advanceTo(
	rightFieldValues []*types.Mlrval,
) {
	if this.joinBucket.wasPaired {
		this.joinBucket.recordsAndContexts = nil
	} else {
		if this.leftUnpaireds == nil {
			this.leftUnpaireds = this.joinBucket.recordsAndContexts
		} else {
			moveRecordsAndContexts(this.leftUnpaireds, this.joinBucket.recordsAndContexts)
		}
	}

	this.joinBucket.recordsAndContexts = list.New()
	if this.joinBucket.leftFieldValues != nil {
		this.joinBucket.leftFieldValues = nil
	}
	this.joinBucket.wasPaired = false

	if this.peekRecordAndContext == nil { // left EOF
		return
	}

	peekRec := this.peekRecordAndContext.Record
	peekFieldValues, hasAllJoinKeys := peekRec.ReferenceSelectedValues(
		this.leftJoinFieldNames,
	)
	// TODO: We need a double condition here ... the peek record is either
	// heterogeneous or homogeneous.  (Or, change that, and ensure elsewhere
	// the peek record is homogeneous.) The former is destined for
	// left-unpaired and shouldn't be lexically compared. The latter should be.
	lib.InternalCodingErrorIf(!hasAllJoinKeys)

	cmp := compareLexicallyPP(peekFieldValues, rightFieldValues)
	if cmp < 0 {
		// Keep seeking and filling the bucket until = or >; this may or may not
		// end up being a match.

		if this.leftUnpaireds == nil {
			this.leftUnpaireds = list.New()
		}

		for {
			this.leftUnpaireds.PushBack(this.peekRecordAndContext)
			this.peekRecordAndContext = nil

			for {
				// Skip over records not having the join keys. These go straight to the
				// left-unpaired list.

				this.peekRecordAndContext = this.readRecord()
				if this.peekRecordAndContext == nil {
					break
				}

				if this.peekRecordAndContext.Record.HasSelectedKeys(this.leftJoinFieldNames) {
					break
				} else {
					if this.leftUnpaireds == nil {
						this.leftUnpaireds = list.New()
					}
					this.leftUnpaireds.PushBack(this.peekRecordAndContext)
				}
			}

			// Double break from double for-loop
			if this.peekRecordAndContext == nil {
				this.leof = true
				break
			}

			peekRec := this.peekRecordAndContext.Record
			// The second return value is a has-all-join-keys indicator -- but
			// we already checked above, so we leave it as _.
			peekFieldValues, _ := peekRec.ReferenceSelectedValues(
				this.leftJoinFieldNames,
			)

			cmp := compareLexicallyVP(this.joinBucket.leftFieldValues, peekFieldValues)
			if cmp != 0 {
				break
			}
		}
	}

	if cmp == 0 {
		this.fill()
		this.joinBucket.wasPaired = true
	} else if cmp > 0 {
		this.fill()
	}
}

// ----------------------------------------------------------------
// TODO: rename?
func (this *tJoinBucketKeeper) drain(
	rightFieldValues []*types.Mlrval,
) {
	// 1. Any records already in this.joinBucket.records (current bucket)
	if this.joinBucket.wasPaired {
		if this.leftUnpaireds == nil {
			this.leftUnpaireds = list.New()
		}
	} else {
		if this.leftUnpaireds == nil {
			this.leftUnpaireds = this.joinBucket.recordsAndContexts
		} else {
			moveRecordsAndContexts(this.leftUnpaireds, this.joinBucket.recordsAndContexts)
		}
	}

	// 2. Peek-record, if any
	if this.peekRecordAndContext != nil {
		this.leftUnpaireds.PushBack(this.peekRecordAndContext)
		this.peekRecordAndContext = nil
	}

	// 3. Remainder of left input stream
	for {
		this.peekRecordAndContext = this.readRecord()
		if this.peekRecordAndContext == nil {
			break
		}
		this.leftUnpaireds.PushBack(this.peekRecordAndContext)
	}

	this.joinBucket.recordsAndContexts = nil
}

// ----------------------------------------------------------------
func (this *tJoinBucketKeeper) outputAndReleaseLeftUnpaireds(
	outputChannel chan<- *types.RecordAndContext,
) {
	if this.leftUnpaireds != nil {
		for {
			element := this.leftUnpaireds.Front()
			if element == nil {
				break
			}
			recordAndContext := element.Value.(*types.RecordAndContext)
			outputChannel <- recordAndContext
		}
	}
}

// ================================================================
// HELPER FUNCTIONS

// ----------------------------------------------------------------
// Pop everything off second-argument list and push to first-argument list.
func moveRecordsAndContexts(
	destination *list.List,
	source *list.List,
) {
	for {
		element := source.Front()
		if element == nil {
			break
		}
		destination.PushBack(element)
	}
}

// ----------------------------------------------------------------
// Returns -1, 0, 1 as left <, ==, > right, using lexical comparison only (even
// for numerical values).
//
// The "VP" in the name is because this variant takes an array of values on the
// left, array of pointers on the right.
//
// We assume (and do not check on each call) that both arrays have the same
// length.
func compareLexicallyVP(
	leftFieldValues []types.Mlrval,
	rightFieldValues []*types.Mlrval,
) int {
	n := len(leftFieldValues)
	for i := 0; i < n; i++ {
		left := leftFieldValues[i].String()
		right := rightFieldValues[i].String()
		// Returns -1, 0, 1 as left <, ==, > right
		cmp := strings.Compare(left, right)
		if cmp != 0 {
			return cmp
		}
	}
	return 0
}

// As previous, except the "PP" in the name is because this variant takes an
// array of pointers for both arguments.
func compareLexicallyPP(
	leftFieldValues []*types.Mlrval,
	rightFieldValues []*types.Mlrval,
) int {
	n := len(leftFieldValues)
	for i := 0; i < n; i++ {
		left := leftFieldValues[i].String()
		right := rightFieldValues[i].String()
		// Returns -1, 0, 1 as left <, ==, > right
		cmp := strings.Compare(left, right)
		if cmp != 0 {
			return cmp
		}
	}
	return 0
}

// TODO: make leftJoinFieldValues be an array of pointers; then simplify this.
