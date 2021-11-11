package input

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
type RecordReaderCSV struct {
	readerOptions *cli.TReaderOptions
	ifs0          byte // Go's CSV library only lets its 'Comma' be a single character
}

// ----------------------------------------------------------------
func NewRecordReaderCSV(readerOptions *cli.TReaderOptions) (*RecordReaderCSV, error) {
	if readerOptions.IRS != "\n" {
		return nil, errors.New("CSV IRS cannot be altered; LF vs CR/LF is autodetected.")
	}
	if len(readerOptions.IFS) != 1 {
		return nil, errors.New("CSV IFS can only be a single character")
	}
	return &RecordReaderCSV{
		readerOptions: readerOptions,
		ifs0:          readerOptions.IFS[0],
	}, nil
}

// ----------------------------------------------------------------
func (reader *RecordReaderCSV) Read(
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

// ----------------------------------------------------------------
func (reader *RecordReaderCSV) processHandle(
	handle io.Reader,
	filename string,
	context *types.Context,
	inputChannel chan<- *types.RecordAndContext,
	errorChannel chan error,
	downstreamDoneChannel <-chan bool, // for mlr head
) {
	context.UpdateForStartOfFile(filename)
	needHeader := !reader.readerOptions.UseImplicitCSVHeader
	var header []string = nil
	var rowNumber int = 0

	csvReader := csv.NewReader(NewBOMStrippingReader(handle))
	csvReader.Comma = rune(reader.ifs0)

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
				value := types.MlrvalFromInferredTypeForDataFiles(csvRecord[i])
				record.PutReference(key, value)
			}

		} else {
			if !reader.readerOptions.AllowRaggedCSVInput {
				err := errors.New(
					fmt.Sprintf(
						"mlr: CSV header/data length mismatch %d != %d "+
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
					value := types.MlrvalFromInferredTypeForDataFiles(csvRecord[i])
					record.PutReference(key, value)
				}
				if nh < nd {
					// if header shorter than data: use 1-up itoa keys
					key := strconv.Itoa(i + 1)
					value := types.MlrvalFromInferredTypeForDataFiles(csvRecord[i])
					record.PutCopy(key, value)
				}
				if nh > nd {
					// if header longer than data: use "" values
					for i = nd; i < nh; i++ {
						record.PutCopy(header[i], types.MLRVAL_VOID)
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
		csvWriter.Comma = rune(reader.ifs0)
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

// ----------------------------------------------------------------
// BOM-stripping
//
// Some CSVs start with a "byte-order mark" which is the 3-byte sequene
// \xef\xbb\xbf".  Any file with such contents trips up csv.Reader:
//
// * If a header line is not double-quoted then we can simply look at the first
//   record returned by csv.Reader and strip away the first three bytes if they
//   are the BOM.
//
// * But if a header line is double-quoted then csv.Reader will complain that
//   the header line has RFC-incompliant double-quoting (it would want the BOM
//   to be *inside* the double quotes).
//
// So we must wrap the io.Reader which is passed to csv.Reader. This
// BOMStrippingReader class does precisely that.

// BOMStrippingReader implements io.Reader to strip leading byte-order-mark
// characters off of CSV data.
type BOMStrippingReader struct {
	underlying io.Reader
	pastBOM    bool
}

func NewBOMStrippingReader(underlying io.Reader) *BOMStrippingReader {
	return &BOMStrippingReader{
		underlying: underlying,
		pastBOM:    false,
	}
}

func (bsr *BOMStrippingReader) Read(p []byte) (n int, err error) {
	if bsr.pastBOM {
		return bsr.underlying.Read(p)
	}
	bsr.pastBOM = true

	// Return error conditions right away.
	n, err = bsr.underlying.Read(p)
	if err != nil {
		return n, err
	}

	// Either this is a small file (maybe a zero-length file, which happens) or
	// it's a big file but we were invoked with a tiny buffer size of 1 or 2.
	// The latter case would be a bit messy to handle; but also we consider it
	// negligibly likely (we expect buffer lengths on the order of kilobytes).
	if n < 3 {
		return n, err
	}

	// If the BOM is present, slip the contents of the buffer down by three.
	if p[0] == CSV_BOM[0] && p[1] == CSV_BOM[1] && p[2] == CSV_BOM[2] {
		for i := 0; i < n-3; i++ {
			p[i] = p[i+3]
		}
		return n - 3, nil
	}

	// No BOM found (normal case).
	return n, err
}
