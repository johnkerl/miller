package input

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"mlr/src/cli"
	"mlr/src/lib"
	"mlr/src/types"
)

// ----------------------------------------------------------------
type RecordReaderCSV struct {
	readerOptions     *cli.TReaderOptions
	emptyStringMlrval types.Mlrval
}

// ----------------------------------------------------------------
func NewRecordReaderCSV(readerOptions *cli.TReaderOptions) *RecordReaderCSV {
	return &RecordReaderCSV{
		readerOptions:     readerOptions,
		emptyStringMlrval: types.MlrvalFromString(""),
	}
}

// ----------------------------------------------------------------
func (reader *RecordReaderCSV) Read(
	filenames []string,
	context types.Context,
	inputChannel chan<- *types.RecordAndContext,
	errorChannel chan error,
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
			reader.processHandle(handle, "(stdin)", &context, inputChannel, errorChannel)
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
					reader.processHandle(handle, filename, &context, inputChannel, errorChannel)
					handle.Close()
				}
			}
		}
	}
	inputChannel <- types.NewEndOfStreamMarker(&context)
}

// ----------------------------------------------------------------
func (reader *RecordReaderCSV) processHandle(
	handle io.Reader,
	filename string,
	context *types.Context,
	inputChannel chan<- *types.RecordAndContext,
	errorChannel chan error,
) {
	context.UpdateForStartOfFile(filename)
	needHeader := !reader.readerOptions.UseImplicitCSVHeader
	var header []string = nil
	var rowNumber int = 0

	csvReader := csv.NewReader(handle)
	csvReader.Comma = rune(reader.readerOptions.IFS[0]) // xxx temp

	for {
		if needHeader {
			// TODO: make this a helper function
			csvRecord, err := csvReader.Read()
			if lib.IsEOF(err) {
				break
			}
			if err != nil && csvRecord == nil {
				// See https://golang.org/pkg/encoding/csv.
				// We handle field-count ourselves.
				errorChannel <- err
				return
			}

			isData := reader.maybeConsumeComment(csvRecord, context, inputChannel)
			if !isData {
				continue
			}

			header = csvRecord
			rowNumber++

			needHeader = false
		}

		csvRecord, err := csvReader.Read()
		if lib.IsEOF(err) {
			break
		}
		if err != nil && csvRecord == nil {
			// See https://golang.org/pkg/encoding/csv.
			// We handle field-count ourselves.
			errorChannel <- err
			return
		}
		rowNumber++

		isData := reader.maybeConsumeComment(csvRecord, context, inputChannel)
		if !isData {
			continue
		}

		if header == nil { // implicit CSV header
			n := len(csvRecord)
			header = make([]string, n)
			for i := 0; i < n; i++ {
				header[i] = strconv.Itoa(i + 1)
			}
		}

		record := types.NewMlrmap()

		nh := len(header)
		nd := len(csvRecord)

		if nh == nd {
			for i := 0; i < nh; i++ {
				key := header[i]
				value := types.MlrvalPointerFromInferredTypeForDataFiles(csvRecord[i])
				record.PutReference(key, value)
			}

		} else {
			if !reader.readerOptions.AllowRaggedCSVInput {
				err := errors.New(
					fmt.Sprintf(
						"Miller: CSV header/data length mismatch %d != %d "+
							"at filename %s row %d.\n",
						nh, nd, filename, rowNumber,
					),
				)
				errorChannel <- err
				return
			} else {
				i := 0
				n := lib.IntMin2(nh, nd)
				for i = 0; i < n; i++ {
					key := header[i]
					value := types.MlrvalPointerFromInferredTypeForDataFiles(csvRecord[i])
					record.PutReference(key, value)
				}
				if nh < nd {
					// if header shorter than data: use 1-up itoa keys
					key := strconv.Itoa(i + 1)
					value := types.MlrvalPointerFromInferredTypeForDataFiles(csvRecord[i])
					record.PutCopy(key, value)
				}
				if nh > nd {
					// if header longer than data: use "" values
					for i = nd; i < nh; i++ {
						record.PutCopy(header[i], &reader.emptyStringMlrval)
					}
				}
			}
		}

		context.UpdateForInputRecord()

		inputChannel <- types.NewRecordAndContext(
			record,
			context,
		)
	}
}

// maybeConsumeComment returns true if the CSV record should be processed as
// data, false otherwise.
func (reader *RecordReaderCSV) maybeConsumeComment(
	csvRecord []string,
	context *types.Context,
	inputChannel chan<- *types.RecordAndContext,
) bool {
	if reader.readerOptions.CommentHandling == cli.CommentsAreData {
		// Nothing is to be construed as a comment
		return true
	}

	if len(csvRecord) < 1 {
		// Not a comment
		return true
	}
	leader := csvRecord[0]

	if !strings.HasPrefix(leader, reader.readerOptions.CommentString) {
		// Not a comment
		return true
	}

	// Is a comment
	if reader.readerOptions.CommentHandling == cli.PassComments {
		// What we want to do here is simple enough: write the record back into
		// a buffer -- basically string-join on IFS but with csvWriter's
		// double-quote handling -- and then pass the resulting string along
		// down-channel to the goroutine which writes strings.
		//
		// However, sadly, bytes.Buffer does not implement io.Writer because
		// its Write method has pointer receiver. So we have a WorkaroundBuffer
		// struct below which has non-pointer receiver.
		buffer := NewWorkaroundBuffer()
		csvWriter := csv.NewWriter(buffer)
		csvWriter.Comma = rune(reader.readerOptions.IFS[0]) // xxx temp
		csvWriter.Write(csvRecord)
		csvWriter.Flush()
		inputChannel <- types.NewOutputString(buffer.String(), context)
	} else /* reader.readerOptions.CommentHandling == cli.SkipComments */ {
		// discard entirely
	}
	return false
}

// ----------------------------------------------------------------
// As noted above: wraps a bytes.Buffer, whose Write method has pointer
// receiver, in a struct with non-pointer receiver so that it implements
// io.Writer.
type WorkaroundBuffer struct {
	pbuffer *bytes.Buffer
}

func NewWorkaroundBuffer() WorkaroundBuffer {
	var buffer bytes.Buffer
	return WorkaroundBuffer{
		pbuffer: &buffer,
	}
}

func (wb WorkaroundBuffer) Write(p []byte) (n int, err error) {
	return wb.pbuffer.Write(p)
}

func (wb WorkaroundBuffer) String() string {
	return wb.pbuffer.String()
}
