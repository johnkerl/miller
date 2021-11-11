package input

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"encoding/json"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/types"
)

type RecordReaderJSON struct {
	readerOptions *cli.TReaderOptions
}

func NewRecordReaderJSON(readerOptions *cli.TReaderOptions) (*RecordReaderJSON, error) {
	return &RecordReaderJSON{
		readerOptions: readerOptions,
	}, nil
}

func (reader *RecordReaderJSON) Read(
	filenames []string,
	context types.Context,
	inputChannel chan<- *types.RecordAndContext,
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
			reader.processHandle(handle, "(stdin)", &context, inputChannel, errorChannel, downstreamDoneChannel)
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
					reader.processHandle(handle, filename, &context, inputChannel, errorChannel, downstreamDoneChannel)
					handle.Close()
				}
			}
		}
	}
	inputChannel <- types.NewEndOfStreamMarker(&context)
}

func (reader *RecordReaderJSON) processHandle(
	handle io.Reader,
	filename string,
	context *types.Context,
	inputChannel chan<- *types.RecordAndContext,
	errorChannel chan error,
	downstreamDoneChannel <-chan bool, // for mlr head
) {
	context.UpdateForStartOfFile(filename)

	if reader.readerOptions.CommentHandling != cli.CommentsAreData {
		handle = NewJSONCommentEnabledReader(handle, reader.readerOptions, inputChannel)
	}
	decoder := json.NewDecoder(handle)

	eof := false
	for {

		// See if downstream processors will be ignoring further data (e.g. mlr
		// head).  If so, stop reading. This makes 'mlr head hugefile' exit
		// quickly, as it should.
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

		mlrval, eof, err := types.MlrvalDecodeFromJSON(decoder)
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
				errorChannel <- errors.New("Internal coding error detected in JSON record-reader")
				return
			}
			context.UpdateForInputRecord()
			inputChannel <- types.NewRecordAndContext(
				record,
				context,
			)
		} else if mlrval.IsArray() {
			records := mlrval.GetArray()
			if records == nil {
				errorChannel <- errors.New("Internal coding error detected in JSON record-reader")
				return
			}
			for _, mlrval := range records {
				if !mlrval.IsMap() {
					// TODO: more context
					errorChannel <- errors.New(
						fmt.Sprintf(
							"Valid but unmillerable JSON. Expected map (JSON object); got %s.",
							mlrval.GetTypeName(),
						),
					)
					return
				}
				record := mlrval.GetMap()
				if record == nil {
					errorChannel <- errors.New("Internal coding error detected in JSON record-reader")
					return
				}
				context.UpdateForInputRecord()
				inputChannel <- types.NewRecordAndContext(
					record,
					context,
				)

			}

		} else {
			errorChannel <- errors.New(
				fmt.Sprintf(
					"Valid but unmillerable JSON. Expected map (JSON object); got %s.",
					mlrval.GetTypeName(),
				),
			)
			return
		}
	}
}

// ================================================================
// JSON comment-stripping
//
// Miller lets users (on an opt-in basis) have comments in their data files,
// for all formats including JSON. Comments are only honored at start of line.
// Users can have them be printed to stdout straightaway, or simply discarded.
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
	context       *types.Context // Needed for channelized stdout-printing logic
	inputChannel  chan<- *types.RecordAndContext

	// In case a line was ingested which was longer than the read-buffer passed
	// to us, in which case we need to split up that line and return it over
	// the course of two or more calls.
	lineBytes []byte
}

func NewJSONCommentEnabledReader(
	underlying io.Reader,
	readerOptions *cli.TReaderOptions,
	inputChannel chan<- *types.RecordAndContext,
) *JSONCommentEnabledReader {
	return &JSONCommentEnabledReader{
		lineScanner:   bufio.NewScanner(underlying),
		readerOptions: readerOptions,
		context:       types.NewNilContext(),
		inputChannel:  inputChannel,

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
			bsr.inputChannel <- types.NewOutputString(line+"\n", bsr.context)
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
