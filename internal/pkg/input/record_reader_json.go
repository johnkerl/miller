package input

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"strings"

	"encoding/json"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
	"github.com/johnkerl/miller/internal/pkg/types"
)

type RecordReaderJSON struct {
	readerOptions   *cli.TReaderOptions
	recordsPerBatch int64 // distinct from readerOptions.RecordsPerBatch for join/repl
}

func NewRecordReaderJSON(
	readerOptions *cli.TReaderOptions,
	recordsPerBatch int64,
) (*RecordReaderJSON, error) {
	return &RecordReaderJSON{
		readerOptions:   readerOptions,
		recordsPerBatch: recordsPerBatch,
	}, nil
}

func (reader *RecordReaderJSON) Read(
	filenames []string,
	context types.Context,
	readerChannel chan<- *list.List, // list of *types.RecordAndContext
	errorChannel chan error,
	downstreamDoneChannel <-chan bool, // for mlr head
) {
	if filenames != nil { // nil for mlr -n
		if len(filenames) == 0 { // read from stdin
			handle, err := lib.OpenStdin(
				reader.readerOptions.Prepipe,
				reader.readerOptions.PrepipeIsRaw,
				reader.readerOptions.FileInputEncoding,
			)
			if err != nil {
				errorChannel <- err
			}
			reader.processHandle(handle, "(stdin)", &context, readerChannel, errorChannel, downstreamDoneChannel)
		} else {
			for _, filename := range filenames {
				handle, err := lib.OpenFileForRead(
					filename,
					reader.readerOptions.Prepipe,
					reader.readerOptions.PrepipeIsRaw,
					reader.readerOptions.FileInputEncoding,
				)
				if err != nil {
					errorChannel <- err
				} else {
					reader.processHandle(handle, filename, &context, readerChannel, errorChannel, downstreamDoneChannel)
					handle.Close()
				}
			}
		}
	}
	readerChannel <- types.NewEndOfStreamMarkerList(&context)
}

func (reader *RecordReaderJSON) processHandle(
	handle io.Reader,
	filename string,
	context *types.Context,
	readerChannel chan<- *list.List, // list of *types.RecordAndContext
	errorChannel chan error,
	downstreamDoneChannel <-chan bool, // for mlr head
) {
	context.UpdateForStartOfFile(filename)
	// TODO: comment
	recordsPerBatch := reader.recordsPerBatch

	if reader.readerOptions.CommentHandling != cli.CommentsAreData {
		handle = NewJSONCommentEnabledReader(handle, reader.readerOptions, readerChannel)
	}
	decoder := json.NewDecoder(handle)
	recordsAndContexts := list.New()

	eof := false
	i := int64(0)
	for {
		// See if downstream processors will be ignoring further data (e.g. mlr
		// head).  If so, stop reading. This makes 'mlr head hugefile' exit
		// quickly, as it should. Do this channel-check every so often to avoid
		// scheduler overhead.
		i++
		if i%recordsPerBatch == 0 {
			select {
			case _ = <-downstreamDoneChannel:
				eof = true
				break
			default:
				break
			}
			if eof {
				break
			}
		}

		mlrval, eof, err := mlrval.MlrvalDecodeFromJSON(decoder)
		if eof {
			break
		}
		if err != nil {
			errorChannel <- err
			return
		}

		// Find out what we got.
		// * Map is an input record: deliver it.
		// * Array is OK if it's array of input record: deliver them.
		// * Non-collection types are valid but unmillerable JSON.

		if mlrval.IsMap() {
			// TODO: make a helper method
			record := mlrval.GetMap()
			if record == nil {
				errorChannel <- fmt.Errorf("internal coding error detected in JSON record-reader")
				return
			}
			context.UpdateForInputRecord()
			recordsAndContexts.PushBack(types.NewRecordAndContext(record, context))

			if int64(recordsAndContexts.Len()) >= recordsPerBatch {
				readerChannel <- recordsAndContexts
				recordsAndContexts = list.New()
			}

		} else if mlrval.IsArray() {
			records := mlrval.GetArray()
			if records == nil {
				errorChannel <- fmt.Errorf("internal coding error detected in JSON record-reader")
				return
			}

			for _, mlrval := range records {
				if !mlrval.IsMap() {
					// TODO: more context
					errorChannel <- fmt.Errorf(
						"valid but unmillerable JSON. Expected map (JSON object); got %s.",
						mlrval.GetTypeName(),
					)
					return
				}
				record := mlrval.GetMap()
				if record == nil {
					errorChannel <- fmt.Errorf("internal coding error detected in JSON record-reader")
					return
				}
				context.UpdateForInputRecord()
				recordsAndContexts.PushBack(types.NewRecordAndContext(record, context))

				if int64(recordsAndContexts.Len()) >= recordsPerBatch {
					readerChannel <- recordsAndContexts
					recordsAndContexts = list.New()
				}
			}

		} else {
			errorChannel <- fmt.Errorf(
				"valid but unmillerable JSON. Expected map (JSON object); got %s.",
				mlrval.GetTypeName(),
			)
			return
		}
	}

	if recordsAndContexts.Len() > 0 {
		readerChannel <- recordsAndContexts
	}
}

// ================================================================
// JSON comment-stripping
//
// Miller lets users (on an opt-in basis) have comments in their data files,
// for all formats including JSON. Comments are only honored at start of line.
// Users can have them be printed to stdout straight away, or simply discarded.
//
// For most file formats Miller is doing line-based I/O and can deal with
// comment lines easily and simply. But for JSON, the Go library needs an
// io.Reader object which we implement here.
//
// This could be done by peeking into the return value from the underlying
// io.Reader, detecting comment-line starts and line-endings within the byte
// array that io.Reader deals with. That's an appealing plan of action, but it
// gets messy if the comment-string is multi-character since then a comment
// string could be split between successive calls to Read() on the underlying
// handle.
//
// Instead we use a line-oriented scanner to do line-splitting for us.

// JSONCommentEnabledReader implements io.Reader to strip comment lines
// off of CSV data.
type JSONCommentEnabledReader struct {
	lineScanner   *bufio.Scanner
	readerOptions *cli.TReaderOptions
	context       *types.Context    // Needed for channelized stdout-printing logic
	readerChannel chan<- *list.List // list of *types.RecordAndContext

	// In case a line was ingested which was longer than the read-buffer passed
	// to us, in which case we need to split up that line and return it over
	// the course of two or more calls.
	lineBytes []byte
}

func NewJSONCommentEnabledReader(
	underlying io.Reader,
	readerOptions *cli.TReaderOptions,
	readerChannel chan<- *list.List, // list of *types.RecordAndContext
) *JSONCommentEnabledReader {
	return &JSONCommentEnabledReader{
		lineScanner:   bufio.NewScanner(underlying),
		readerOptions: readerOptions,
		context:       types.NewNilContext(),
		readerChannel: readerChannel,

		lineBytes: nil,
	}
}

func (bsr *JSONCommentEnabledReader) Read(p []byte) (n int, err error) {
	if bsr.lineBytes != nil {
		return bsr.populateFromLine(p), nil
	}

	// Loop until we can get a non-comment line to pass on, or end of file.
	for {
		// EOF
		if !bsr.lineScanner.Scan() {
			return 0, io.EOF
		}
		line := bsr.lineScanner.Text()

		// Non-comment line
		if !strings.HasPrefix(line, bsr.readerOptions.CommentString) {
			bsr.lineBytes = []byte(line)
			return bsr.populateFromLine(p), nil
		}

		// Comment line
		if bsr.readerOptions.CommentHandling == cli.PassComments {
			// Insert the string into the record-output stream, so that goroutine can
			// print it, resulting in deterministic output-ordering.
			ell := list.New()
			ell.PushBack(types.NewOutputString(line+"\n", bsr.context))
			bsr.readerChannel <- ell
		}
	}
}

// populateFromLine is a helper for Read. It takes a full line from the
// bufio.Scanner, and writes as much as it can to the caller's p-buffer.  If
// the entirety is written, the line is marked as done so a subsequent call to
// Read will retrieve the next line from the input file. Otherwise, as much as
// possible is transferred, and the rest is marked for transfer on a subsequent
// call.
func (bsr *JSONCommentEnabledReader) populateFromLine(p []byte) int {
	numBytesWritten := 0
	if len(bsr.lineBytes) < len(p) {
		for i := 0; i < len(bsr.lineBytes); i++ {
			p[i] = bsr.lineBytes[i]
		}
		numBytesWritten = len(bsr.lineBytes)
		bsr.lineBytes = nil
	} else {
		for i := 0; i < len(p); i++ {
			p[i] = bsr.lineBytes[i]
		}
		numBytesWritten = len(p)
		bsr.lineBytes = bsr.lineBytes[len(p):]
	}
	return numBytesWritten
}
