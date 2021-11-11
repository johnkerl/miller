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

package utils

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/input"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
// Data stored in this class
type JoinBucketKeeper struct {
	// For streaming through the left-side file
	recordReader input.IRecordReader
	context      *types.Context
	inputChannel <-chan *types.RecordAndContext
	errorChannel chan error
	// TODO: merge with leof flag
	recordReaderDone bool

	leftJoinFieldNames []string

	// Given a left-file of the following form (with left-join-field name "L"):
	//   +-----+
	//   |  L  |
	//   + --- +
	//   |  a  |
	//   |  a  |
	//   |  a  |
	//   |  b  |
	//   |  c  |
	//   |  c  |
	//   |  d  |
	//   +-----+
	// then the join-bucket points to a full list of records with same
	// left-join-field value, and the peek record is the next one (if any --
	// nil at left EOF) with a different value.

	peekRecordAndContext *types.RecordAndContext
	JoinBucket           *JoinBucket
	leftUnpaireds        *list.List

	leof  bool
	state tJoinBucketKeeperState
}

// ----------------------------------------------------------------
func NewJoinBucketKeeper(
	// TODO prepipe string,
	leftFileName string,
	joinReaderOptions *cli.TReaderOptions,
	leftJoinFieldNames []string,
) *JoinBucketKeeper {

	// Instantiate the record-reader
	recordReader, err := input.Create(joinReaderOptions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mlr join: %v", err)
		os.Exit(1)
	}

	// Set the initial context for the left-file.  Since Go is concurrent, the
	// context struct needs to be duplicated and passed through the channels
	// along with each record.
	initialContext := types.NewNilContext()
	initialContext.UpdateForStartOfFile(leftFileName)

	// Set up channels for the record-reader
	inputChannel := make(chan *types.RecordAndContext, 10)
	errorChannel := make(chan error, 1)
	downstreamDoneChannel := make(chan bool, 1)

	// Start the record-reader in its own goroutine.
	leftFileNameArray := [1]string{leftFileName}
	go recordReader.Read(leftFileNameArray[:], *initialContext, inputChannel, errorChannel, downstreamDoneChannel)

	keeper := &JoinBucketKeeper{
		recordReader:     recordReader,
		context:          initialContext,
		inputChannel:     inputChannel,
		errorChannel:     errorChannel,
		recordReaderDone: false,

		leftJoinFieldNames: leftJoinFieldNames,

		JoinBucket:           NewJoinBucket(nil),
		peekRecordAndContext: nil,
		leftUnpaireds:        list.New(),

		leof:  false,
		state: LEFT_STATE_0_PREFILL,
	}

	return keeper
}

// ----------------------------------------------------------------
// For JoinBucketKeeper state machine
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

func (keeper *JoinBucketKeeper) computeState() tJoinBucketKeeperState {
	if keeper.JoinBucket.leftFieldValues == nil {
		if !keeper.leof {
			return LEFT_STATE_0_PREFILL
		} else {
			return LEFT_STATE_3_EOF
		}
	} else {
		if keeper.peekRecordAndContext == nil {
			return LEFT_STATE_2_LAST_BUCKET
		} else {
			return LEFT_STATE_1_FULL
		}
	}
}

// ----------------------------------------------------------------
// This is the main entry point for the join verb.  Given the right-field
// values from the current right-file record, this method finds left-file
// join-bucket (if any) and points keeper.JoinBucket at it.
//
// If the join-keys have changed since the last right record, and if the
// previous join-bucket wasn't ever paired with a right record, then it will be
// moved to keeper.leftUnpaired.
//
// Also, if it's time to seek to a new left-side join bucket, then any
// left-file records found along the way lacking the specified join-field names
// will also be moved to keeper.leftUnpaired.

func (keeper *JoinBucketKeeper) FindJoinBucket(
	rightFieldValues []*types.Mlrval, // nil means right-file EOF
) bool {
	// TODO: comment me
	isPaired := false

	// This will produce a join bucket on the left side (if there is any at all
	// to be had) but it may or may not make the join keys from the current
	// right record.
	if keeper.state == LEFT_STATE_0_PREFILL {
		keeper.prepareForFirstJoinBucket()
		if keeper.peekRecordAndContext != nil {
			keeper.fillNextJoinBucket()
		}
		keeper.state = keeper.computeState()
	}

	if rightFieldValues != nil { // Not right EOF
		if keeper.state == LEFT_STATE_1_FULL || keeper.state == LEFT_STATE_2_LAST_BUCKET {

			cmp := compareLexically(keeper.JoinBucket.leftFieldValues, rightFieldValues)

			if cmp < 0 {
				// Advance left until match or left EOF.  This might find a
				// matching join bucket for the current record, or not.
				// Example: joining on "id" column and left file has several
				// join-field records with id=3, then several with id=7, but
				// the current right record has id=5.
				keeper.prepareForNewJoinBucket(rightFieldValues)

				if keeper.peekRecordAndContext != nil {
					keeper.fillNextJoinBucket()
				}

				// TODO: privatize more
				if keeper.JoinBucket.RecordsAndContexts.Len() > 0 {
					cmp := compareLexically(
						keeper.JoinBucket.leftFieldValues,
						rightFieldValues,
					)
					if cmp == 0 {
						isPaired = true
						keeper.JoinBucket.WasPaired = true
					}
				}

			} else if cmp == 0 {
				// Stay on current bucket
				keeper.JoinBucket.WasPaired = true
				isPaired = true
			} else {
				// E.g. joining on "id", current right-record has id=5,
				// previous join-bucket had id=4, new one has id=6.  No match
				// and no need to advance left.
				isPaired = false
			}
		} else if keeper.state != LEFT_STATE_3_EOF {
			fmt.Fprintf(
				os.Stderr,
				"%s: internal coding error: failed transition from prefill state.\n",
				"mlr",
			)
			os.Exit(1)
		}

	} else { // Right EOF
		keeper.markRemainingsAsUnpaired()
	}

	keeper.state = keeper.computeState()

	return isPaired
}

// ----------------------------------------------------------------
// This finds the first peek record which posseses all the necessary join-field
// keys.  Any other records found along the way, lacking the necessary
// join-field keys, are moved to the left-unpaired list.

func (keeper *JoinBucketKeeper) prepareForFirstJoinBucket() {
	for {
		// Skip over records not having the join keys. These go straight to the
		// left-unpaired list.
		keeper.peekRecordAndContext = keeper.readRecord()
		if keeper.peekRecordAndContext == nil { // left EOF
			break
		}
		if keeper.peekRecordAndContext.Record.HasSelectedKeys(keeper.leftJoinFieldNames) {
			break
		}
		keeper.leftUnpaireds.PushBack(keeper.peekRecordAndContext)
	}

	if keeper.peekRecordAndContext == nil {
		keeper.leof = true
		return
	}
}

// ----------------------------------------------------------------
// After right-file input has moved past the current join-bucket, this finds
// the next peek record which possesses all the necessary join-field keys.  Any
// other records found along the way, lacking the necessary join-field keys,
// are moved to the left-unpaired list.
//
// Pre-conditions:
// * Our keeper.JoinBucket.leftFieldValues < rightFieldValues (with lexical
//   comparison, even for numeric values).
// * Currently in state 1 or 2 so there is a bucket but there may or may not be
//   a peek-record.
// * Current bucket was/wasn't paired on previous emits but is not paired on this emit.
// Actions:
// * If the current bucket was never paired, move it to the left-unpaired list.
// * Consume the left input stream, feeding into unpaired, for as long as
//   leftvals < rightvals && !eof.

func (keeper *JoinBucketKeeper) prepareForNewJoinBucket(
	rightFieldValues []*types.Mlrval,
) {
	if !keeper.JoinBucket.WasPaired {
		moveRecordsAndContexts(keeper.leftUnpaireds, keeper.JoinBucket.RecordsAndContexts)
	}
	keeper.JoinBucket = NewJoinBucket(nil)

	if keeper.peekRecordAndContext == nil { // left EOF
		return
	}

	peekRec := keeper.peekRecordAndContext.Record
	peekFieldValues, hasAllJoinKeys := peekRec.ReferenceSelectedValues(
		keeper.leftJoinFieldNames,
	)
	lib.InternalCodingErrorIf(!hasAllJoinKeys)

	// We use a double condition here, implemented as a double for-loop. The
	// peek record is either heterogeneous or homogeneous. The former is
	// destined for left-unpaired and shouldn't be lexically compared. The
	// latter should be.

	cmp := compareLexically(peekFieldValues, rightFieldValues)
	if cmp >= 0 {
		return
	}

	// Keep seeking and filling the bucket until = or >; this may or may not
	// end up being a match.
	for {
		keeper.leftUnpaireds.PushBack(keeper.peekRecordAndContext)
		keeper.peekRecordAndContext = nil

		for {
			// Skip over records not having the join keys. These go straight to the
			// left-unpaired list.
			keeper.peekRecordAndContext = keeper.readRecord()
			if keeper.peekRecordAndContext == nil {
				break
			}
			peekRec := keeper.peekRecordAndContext.Record

			if peekRec.HasSelectedKeys(keeper.leftJoinFieldNames) {
				break
			}
			keeper.leftUnpaireds.PushBack(keeper.peekRecordAndContext)
		}

		// Double break from double for-loop
		if keeper.peekRecordAndContext == nil {
			keeper.leof = true
			break
		}

		peekRec := keeper.peekRecordAndContext.Record
		// The second return value is a has-all-join-keys indicator -- but
		// we already checked above, so we leave it as _.
		peekFieldValues, _ := peekRec.ReferenceSelectedValues(
			keeper.leftJoinFieldNames,
		)

		cmp = compareLexically(peekFieldValues, rightFieldValues)
		if cmp >= 0 {
			break
		}
	}
}

// ----------------------------------------------------------------
// This takes the peek record and forms a complete join-bucket with all records
// having its join-field values. E.g. if the join-field is "id" and the peek
// record has id=5, it's moved to a new join bucket with id=5 and all other
// left-file records with id=5 are put there as well. To get *all* such
// requires that we read until we have one too many, which becomes the next
// peek record (maybe having id=6, for example).
//
// It moves the previous join-bucket to the left-unpaired list, if that was never
// paired with a right-file record.
//
// Preconditions:
// * peekRecordAndContext != nil
// * peekRecordAndContext has the join keys

func (keeper *JoinBucketKeeper) fillNextJoinBucket() {
	peekRec := keeper.peekRecordAndContext.Record
	peekFieldValues, hasAllJoinKeys := peekRec.ReferenceSelectedValues(
		keeper.leftJoinFieldNames,
	)

	if !hasAllJoinKeys {
		fmt.Fprintf(
			os.Stderr,
			"%s: internal coding error: peek record should have had join keys.\n",
			"mlr",
		)
		os.Exit(1)
	}

	keeper.JoinBucket.leftFieldValues = types.CopyMlrvalPointerArray(peekFieldValues)
	keeper.JoinBucket.RecordsAndContexts.PushBack(keeper.peekRecordAndContext)
	keeper.JoinBucket.WasPaired = false

	keeper.peekRecordAndContext = nil

	for {
		// Skip over records not having the join keys. These go straight to the
		// left-unpaired list.
		keeper.peekRecordAndContext = keeper.readRecord()
		if keeper.peekRecordAndContext == nil { // left EOF
			keeper.leof = true
			break
		}

		peekRec := keeper.peekRecordAndContext.Record
		peekFieldValues, hasAllJoinKeys := peekRec.ReferenceSelectedValues(
			keeper.leftJoinFieldNames,
		)

		if hasAllJoinKeys {
			cmp := compareLexically(
				keeper.JoinBucket.leftFieldValues,
				peekFieldValues,
			)
			if cmp != 0 {
				break
			}
			keeper.JoinBucket.RecordsAndContexts.PushBack(keeper.peekRecordAndContext)
		} else {
			keeper.leftUnpaireds.PushBack(keeper.peekRecordAndContext)
		}
		keeper.peekRecordAndContext = nil
	}
}

// ----------------------------------------------------------------
// TODO: comment
func (keeper *JoinBucketKeeper) markRemainingsAsUnpaired() {
	// 1. Any records already in keeper.JoinBucket.records (current bucket)
	if !keeper.JoinBucket.WasPaired {
		moveRecordsAndContexts(keeper.leftUnpaireds, keeper.JoinBucket.RecordsAndContexts)
	}
	keeper.JoinBucket.RecordsAndContexts = nil

	// 2. Peek-record, if any
	if keeper.peekRecordAndContext != nil {
		keeper.leftUnpaireds.PushBack(keeper.peekRecordAndContext)
		keeper.peekRecordAndContext = nil
	}

	// 3. Remainder of left input stream
	for {
		keeper.peekRecordAndContext = keeper.readRecord()
		if keeper.peekRecordAndContext == nil {
			break
		}
		keeper.leftUnpaireds.PushBack(keeper.peekRecordAndContext)
	}
}

// ----------------------------------------------------------------
// TODO: comment
func (keeper *JoinBucketKeeper) OutputAndReleaseLeftUnpaireds(
	outputChannel chan<- *types.RecordAndContext,
) {
	for {
		element := keeper.leftUnpaireds.Front()
		if element == nil {
			break
		}
		recordAndContext := element.Value.(*types.RecordAndContext)
		outputChannel <- recordAndContext
		keeper.leftUnpaireds.Remove(element)
	}
}

func (keeper *JoinBucketKeeper) ReleaseLeftUnpaireds(
	outputChannel chan<- *types.RecordAndContext,
) {
	for {
		element := keeper.leftUnpaireds.Front()
		if element == nil {
			break
		}
		keeper.leftUnpaireds.Remove(element)
	}
}

// ================================================================
// HELPER FUNCTIONS

// ----------------------------------------------------------------
// Method to get the next left-file record from the record-reader goroutine.
// Returns nil at EOF.

func (keeper *JoinBucketKeeper) readRecord() *types.RecordAndContext {
	if keeper.recordReaderDone {
		return nil
	}

	select {
	case err := <-keeper.errorChannel:
		fmt.Fprintln(os.Stderr, "mlr", ": ", err)
		os.Exit(1)
	case leftrecAndContext := <-keeper.inputChannel:
		if leftrecAndContext.EndOfStream { // end-of-stream marker
			keeper.recordReaderDone = true
			return nil
		} else {
			return leftrecAndContext
		}
	}

	return nil
}

// ----------------------------------------------------------------
// Pops everything off second-argument list and push to first-argument list.

func moveRecordsAndContexts(
	destination *list.List,
	source *list.List,
) {
	for {
		element := source.Front()
		if element == nil {
			break
		}
		destination.PushBack(element.Value.(*types.RecordAndContext))
		source.Remove(element)
	}
}

// ----------------------------------------------------------------
// Returns -1, 0, 1 as left <, ==, > right, using lexical comparison only (even
// for numerical values).

func compareLexically(
	leftFieldValues []*types.Mlrval,
	rightFieldValues []*types.Mlrval,
) int {
	lib.InternalCodingErrorIf(len(leftFieldValues) != len(rightFieldValues))
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

// ================================================================
func (keeper *JoinBucketKeeper) dump(prefix string) {
	fmt.Printf("+----------------------------------------------------- %s\n", prefix)
	fmt.Println("| recordReaderDone     [", keeper.recordReaderDone, "]")
	fmt.Println("| leof                 [", keeper.leof, "]")
	fmt.Println("| stateCode            [", keeper.state, "]")
	fmt.Println("| leftJoinFieldNames   [", strings.Join(keeper.leftJoinFieldNames, ","), "]")

	fmt.Println("| JoinBucket:")
	// TODO: make utility method
	leftFieldValuesString := make([]string, len(keeper.JoinBucket.leftFieldValues))
	for i, leftFieldValue := range keeper.JoinBucket.leftFieldValues {
		leftFieldValuesString[i] = leftFieldValue.String()
	}
	fmt.Printf("|   leftFieldValues    [%s]\n", strings.Join(leftFieldValuesString, ","))
	fmt.Printf("|   RecordsAndContexts (%d)\n", keeper.JoinBucket.RecordsAndContexts.Len())
	for element := keeper.JoinBucket.RecordsAndContexts.Front(); element != nil; element = element.Next() {
		fmt.Println("|    ", element.Value.(*types.RecordAndContext).Record.ToDKVPString())
	}
	fmt.Println("|   WasPaired         ", keeper.JoinBucket.WasPaired)

	if keeper.peekRecordAndContext == nil || keeper.peekRecordAndContext.Record == nil {
		fmt.Println("| peekRecordAndContext [nil]")
	} else {
		fmt.Println("| peekRecordAndContext [", keeper.peekRecordAndContext.Record.ToDKVPString(), "]")
	}

	fmt.Printf("| leftUnpaireds        (%d)\n", keeper.leftUnpaireds.Len())
	for element := keeper.leftUnpaireds.Front(); element != nil; element = element.Next() {
		fmt.Println("|   ", element.Value.(*types.RecordAndContext).Record.ToDKVPString())
	}

	fmt.Printf("------------------------------------------------------\n")
}

func dumpFieldValues(name string, values []*types.Mlrval) {
	for i, value := range values {
		fmt.Printf("-- %s[%d] = %s\n", name, i, value.String())
	}
}
