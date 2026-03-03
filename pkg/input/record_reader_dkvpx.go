// RecordReaderDKVPX reads DKVPX format: comma-delimited key=value pairs with
// CSV-style quoting. It uses the dkvpx package for parsing.
package input

import (
	"fmt"
	"io"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/dkvpx"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

type RecordReaderDKVPX struct {
	readerOptions   *cli.TReaderOptions
	recordsPerBatch int64
}

func NewRecordReaderDKVPX(
	readerOptions *cli.TReaderOptions,
	recordsPerBatch int64,
) (*RecordReaderDKVPX, error) {
	if readerOptions.IRS != "\n" && readerOptions.IRS != "\r\n" {
		return nil, fmt.Errorf("for DKVPX, IRS cannot be altered; LF vs CR/LF is autodetected")
	}
	return &RecordReaderDKVPX{
		readerOptions:   readerOptions,
		recordsPerBatch: recordsPerBatch,
	}, nil
}

func (reader *RecordReaderDKVPX) Read(
	filenames []string,
	context types.Context,
	readerChannel chan<- []*types.RecordAndContext,
	errorChannel chan error,
	downstreamDoneChannel <-chan bool,
) {
	if filenames != nil {
		if len(filenames) == 0 {
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

func (reader *RecordReaderDKVPX) processHandle(
	handle io.Reader,
	filename string,
	context *types.Context,
	readerChannel chan<- []*types.RecordAndContext,
	errorChannel chan<- error,
	downstreamDoneChannel <-chan bool,
) {
	context.UpdateForStartOfFile(filename)
	recordsPerBatch := reader.recordsPerBatch

	dkvpxReader := dkvpx.NewReader(NewBOMStrippingReader(handle))
	dkvpxReader.Comma = ','
	if reader.readerOptions.CommentHandling != cli.CommentsAreData &&
		len(reader.readerOptions.CommentString) == 1 {
		dkvpxReader.Comment = rune(reader.readerOptions.CommentString[0])
	}

	dkvpxRecordsChannel := make(chan []*lib.OrderedMap[string], recordsPerBatch)
	go channelizedDKVPXRecordScanner(dkvpxReader, dkvpxRecordsChannel, downstreamDoneChannel, errorChannel, recordsPerBatch)

	for {
		recordsAndContexts, eof := reader.getRecordBatch(dkvpxRecordsChannel, errorChannel, context)
		if len(recordsAndContexts) > 0 {
			readerChannel <- recordsAndContexts
		}
		if eof {
			break
		}
	}
}

func channelizedDKVPXRecordScanner(
	dkvpxReader *dkvpx.Reader,
	dkvpxRecordsChannel chan<- []*lib.OrderedMap[string],
	downstreamDoneChannel <-chan bool,
	errorChannel chan<- error,
	recordsPerBatch int64,
) {
	i := int64(0)
	done := false

	dkvpxRecords := make([]*lib.OrderedMap[string], 0, recordsPerBatch)

	for {
		i++

		dkvpxRecord, err := dkvpxReader.Read()
		if lib.IsEOF(err) {
			break
		}
		if err != nil {
			errorChannel <- err
			break
		}

		dkvpxRecords = append(dkvpxRecords, dkvpxRecord)

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
			dkvpxRecordsChannel <- dkvpxRecords
			dkvpxRecords = make([]*lib.OrderedMap[string], 0, recordsPerBatch)
		}

		if done {
			break
		}
	}
	dkvpxRecordsChannel <- dkvpxRecords
	close(dkvpxRecordsChannel)
}

func (reader *RecordReaderDKVPX) getRecordBatch(
	dkvpxRecordsChannel <-chan []*lib.OrderedMap[string],
	errorChannel chan<- error,
	context *types.Context,
) ([]*types.RecordAndContext, bool) {
	recordsAndContexts := []*types.RecordAndContext{}
	dedupeFieldNames := reader.readerOptions.DedupeFieldNames

	dkvpxRecords, more := <-dkvpxRecordsChannel
	if !more {
		return recordsAndContexts, true
	}

	for _, omap := range dkvpxRecords {
		record := mlrval.NewMlrmapAsRecord()

		for pe := omap.Head; pe != nil; pe = pe.Next {
			value := mlrval.FromDeferredType(pe.Value)
			_, err := record.PutReferenceMaybeDedupe(pe.Key, value, dedupeFieldNames)
			if err != nil {
				errorChannel <- err
				return nil, false
			}
		}

		context.UpdateForInputRecord()
		recordsAndContexts = append(recordsAndContexts, types.NewRecordAndContext(record, context))
	}

	return recordsAndContexts, false
}
