package input

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	csv "github.com/johnkerl/miller/v6/pkg/go-csv"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

type RecordReaderCSV struct {
	readerOptions       *cli.TReaderOptions
	recordsPerBatch     int64 // distinct from readerOptions.RecordsPerBatch for join/repl
	ifs0                byte  // Go's CSV library only lets its 'Comma' be a single character
	csvLazyQuotes       bool  // Maps directly to Go's CSV library's LazyQuotes
	csvTrimLeadingSpace bool  // Maps directly to Go's CSV library's TrimLeadingSpace

	filename   string
	rowNumber  int64
	needHeader bool
	header     []string
}

func NewRecordReaderCSV(
	readerOptions *cli.TReaderOptions,
	recordsPerBatch int64,
) (*RecordReaderCSV, error) {
	if readerOptions.IRS != "\n" && readerOptions.IRS != "\r\n" {
		return nil, fmt.Errorf("for CSV, IRS cannot be altered; LF vs CR/LF is autodetected")
	}
	if len(readerOptions.IFS) != 1 {
		return nil, fmt.Errorf("for CSV, IFS can only be a single character")
	}
	if readerOptions.CommentHandling != cli.CommentsAreData {
		if len(readerOptions.CommentString) != 1 {
			return nil, fmt.Errorf("for CSV, the comment prefix must be a single character")
		}
	}
	return &RecordReaderCSV{
		readerOptions:       readerOptions,
		ifs0:                readerOptions.IFS[0],
		recordsPerBatch:     recordsPerBatch,
		csvLazyQuotes:       readerOptions.CSVLazyQuotes,
		csvTrimLeadingSpace: readerOptions.CSVTrimLeadingSpace,
	}, nil
}

func (reader *RecordReaderCSV) Read(
	filenames []string,
	context types.Context,
	readerChannel chan<- []*types.RecordAndContext, // list of *types.RecordAndContext
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
			} else {
				reader.processHandle(handle, "(stdin)", &context, readerChannel, errorChannel, downstreamDoneChannel)
			}
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

func (reader *RecordReaderCSV) processHandle(
	handle io.Reader,
	filename string,
	context *types.Context,
	readerChannel chan<- []*types.RecordAndContext, // list of *types.RecordAndContext
	errorChannel chan error,
	downstreamDoneChannel <-chan bool, // for mlr head
) {
	context.UpdateForStartOfFile(filename)
	recordsPerBatch := reader.recordsPerBatch

	// Reset state for start of next input file
	reader.filename = filename
	reader.rowNumber = 0
	reader.needHeader = !reader.readerOptions.UseImplicitHeader
	reader.header = nil

	csvReader := csv.NewReader(NewBOMStrippingReader(handle))
	csvReader.Comma = rune(reader.ifs0)
	csvReader.LazyQuotes = reader.csvLazyQuotes
	csvReader.TrimLeadingSpace = reader.csvTrimLeadingSpace

	if reader.readerOptions.CommentHandling != cli.CommentsAreData {
		if len(reader.readerOptions.CommentString) == 1 {
			// Use our modified fork of the go-csv package
			csvReader.Comment = rune(reader.readerOptions.CommentString[0])
		}
	}

	csvRecordsChannel := make(chan [][]string, recordsPerBatch)
	go channelizedCSVRecordScanner(csvReader, csvRecordsChannel, downstreamDoneChannel, errorChannel,
		recordsPerBatch)

	for {
		recordsAndContexts, eof := reader.getRecordBatch(csvRecordsChannel, errorChannel, context)
		if len(recordsAndContexts) > 0 {
			readerChannel <- recordsAndContexts
		}
		if eof {
			break
		}
	}
}

// TODO: comment
func channelizedCSVRecordScanner(
	csvReader *csv.Reader,
	csvRecordsChannel chan<- [][]string,
	downstreamDoneChannel <-chan bool, // for mlr head
	errorChannel chan error,
	recordsPerBatch int64,
) {
	i := int64(0)
	done := false

	csvRecords := make([][]string, 0, recordsPerBatch)

	for {
		i++

		csvRecord, err := csvReader.Read()
		if lib.IsEOF(err) {
			break
		}
		if err != nil && csvRecord == nil {
			// See https://golang.org/pkg/encoding/csv.
			// We handle field-count ourselves.
			errorChannel <- err
			break
		}

		csvRecords = append(csvRecords, csvRecord)

		// See if downstream processors will be ignoring further data (e.g. mlr
		// head).  If so, stop reading. This makes 'mlr head hugefile' exit
		// quickly, as it should.
		if i%recordsPerBatch == 0 {
			select {
			case <-downstreamDoneChannel:
				done = true
				break
			default:
				break
			}
			if done {
				break
			}
			csvRecordsChannel <- csvRecords
			csvRecords = make([][]string, 0, recordsPerBatch)
		}

		if done {
			break
		}
	}
	csvRecordsChannel <- csvRecords
	close(csvRecordsChannel) // end-of-stream marker
}

// TODO: comment copiously we're trying to handle slow/fast/short/long reads: tail -f, smallfile, bigfile.
func (reader *RecordReaderCSV) getRecordBatch(
	csvRecordsChannel <-chan [][]string,
	errorChannel chan error,
	context *types.Context,
) (
	recordsAndContexts []*types.RecordAndContext,
	eof bool,
) {
	recordsAndContexts = make([]*types.RecordAndContext, 0)
	dedupeFieldNames := reader.readerOptions.DedupeFieldNames

	csvRecords, more := <-csvRecordsChannel
	if !more {
		return recordsAndContexts, true
	}

	for _, csvRecord := range csvRecords {

		if reader.needHeader {
			isData := reader.maybeConsumeComment(csvRecord, context, &recordsAndContexts)
			if !isData {
				continue
			}

			reader.header = csvRecord
			reader.rowNumber++
			reader.needHeader = false
			continue
		}

		isData := reader.maybeConsumeComment(csvRecord, context, &recordsAndContexts)
		if !isData {
			continue
		}
		reader.rowNumber++

		if reader.header == nil { // implicit CSV header
			n := len(csvRecord)
			reader.header = make([]string, n)
			for i := 0; i < n; i++ {
				reader.header[i] = strconv.Itoa(i + 1)
			}
		}

		record := mlrval.NewMlrmapAsRecord()

		nh := int64(len(reader.header))
		nd := int64(len(csvRecord))

		if nh == nd {
			for i := int64(0); i < nh; i++ {
				key := reader.header[i]
				value := mlrval.FromDeferredType(csvRecord[i])
				_, err := record.PutReferenceMaybeDedupe(key, value, dedupeFieldNames)
				if err != nil {
					errorChannel <- err
					return
				}
			}

		} else {
			if !reader.readerOptions.AllowRaggedCSVInput {
				err := fmt.Errorf(
					"mlr: CSV header/data length mismatch %d != %d at filename %s row %d",
					nh, nd, reader.filename, reader.rowNumber,
				)
				errorChannel <- err
				return
			}

			i := int64(0)
			n := lib.IntMin2(nh, nd)
			for i = 0; i < n; i++ {
				key := reader.header[i]
				value := mlrval.FromDeferredType(csvRecord[i])
				_, err := record.PutReferenceMaybeDedupe(key, value, dedupeFieldNames)
				if err != nil {
					errorChannel <- err
					return
				}
			}
			if nh < nd {
				// if header shorter than data: use 1-up itoa keys
				for i = nh; i < nd; i++ {
					key := strconv.FormatInt(i+1, 10)
					value := mlrval.FromDeferredType(csvRecord[i])
					_, err := record.PutReferenceMaybeDedupe(key, value, dedupeFieldNames)
					if err != nil {
						errorChannel <- err
						return
					}
				}
			}
			// if nh > nd: leave it short. This is a job for unsparsify.
		}

		context.UpdateForInputRecord()

		recordsAndContexts = append(recordsAndContexts, types.NewRecordAndContext(record, context))
	}

	return recordsAndContexts, false
}

// maybeConsumeComment returns true if the CSV record should be processed as
// data, false otherwise.
func (reader *RecordReaderCSV) maybeConsumeComment(
	csvRecord []string,
	context *types.Context,
	recordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
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

		// Contract with our fork of the go-csv CSV Reader, and, our own constructor.
		lib.InternalCodingErrorIf(len(csvRecord) != 1)
		*recordsAndContexts = append(*recordsAndContexts, types.NewOutputString(csvRecord[0], context))

	} else /* reader.readerOptions.CommentHandling == cli.SkipComments */ {
		// discard entirely
	}
	return false
}

// ----------------------------------------------------------------
// BOM-stripping
//
// Some CSVs start with a "byte-order mark" which is the 3-byte sequence
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
