package input

import (
	"container/list"
	"io"
	"strconv"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/types"
)

type RecordReaderNIDX struct {
	readerOptions   *cli.TReaderOptions
	recordsPerBatch int
}

func NewRecordReaderNIDX(
	readerOptions *cli.TReaderOptions,
	recordsPerBatch int,
) (*RecordReaderNIDX, error) {
	return &RecordReaderNIDX{
		readerOptions:   readerOptions,
		recordsPerBatch: recordsPerBatch,
	}, nil
}

func (reader *RecordReaderNIDX) Read(
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

func (reader *RecordReaderNIDX) processHandle(
	handle io.Reader,
	filename string,
	context *types.Context,
	readerChannel chan<- *list.List, // list of *types.RecordAndContext
	errorChannel chan error,
	downstreamDoneChannel <-chan bool, // for mlr head
) {
	context.UpdateForStartOfFile(filename)
	recordsPerBatch := reader.readerOptions.RecordsPerBatch

	lineScanner := NewLineScanner(handle, reader.readerOptions.IRS)
	linesChannel := make(chan *list.List, recordsPerBatch)
	go channelizedLineScanner(lineScanner, linesChannel, downstreamDoneChannel, recordsPerBatch)

	for {
		recordsAndContexts, eof := reader.getRecordBatch(linesChannel, recordsPerBatch, context)
		readerChannel <- recordsAndContexts
		if eof {
			break
		}
	}
}

// TODO: comment copiously we're trying to handle slow/fast/short/long reads: tail -f, smallfile, bigfile.
func (reader *RecordReaderNIDX) getRecordBatch(
	linesChannel <-chan *list.List,
	maxBatchSize int,
	context *types.Context,
) (
	recordsAndContexts *list.List,
	eof bool,
) {
	recordsAndContexts = list.New()

	lines, more := <-linesChannel
	if !more {
		return recordsAndContexts, true
	}

	for e := lines.Front(); e != nil; e = e.Next() {
		line := e.Value.(string)

		// Check for comments-in-data feature
		// TODO: function-pointer this away
		if reader.readerOptions.CommentHandling != cli.CommentsAreData {
			if strings.HasPrefix(line, reader.readerOptions.CommentString) {
				if reader.readerOptions.CommentHandling == cli.PassComments {
					recordsAndContexts.PushBack(types.NewOutputString(line+"\n", context))
					continue
				} else if reader.readerOptions.CommentHandling == cli.SkipComments {
					continue
				}
				// else comments are data
			}
		}

		record := reader.recordFromNIDXLine(line)

		context.UpdateForInputRecord()
		recordAndContext := types.NewRecordAndContext(record, context)
		recordsAndContexts.PushBack(recordAndContext)
	}

	return recordsAndContexts, false
}

// ----------------------------------------------------------------
func (reader *RecordReaderNIDX) recordFromNIDXLine(
	line string,
) *types.Mlrmap {
	record := types.NewMlrmapAsRecord()

	var values []string
	if reader.readerOptions.IFSRegex == nil { // e.g. --no-ifs-regex
		values = lib.SplitString(line, reader.readerOptions.IFS)
	} else {
		values = lib.RegexSplitString(reader.readerOptions.IFSRegex, line, -1)
	}
	if reader.readerOptions.AllowRepeatIFS {
		values = lib.StripEmpties(values) // left/right trim
	}

	var i int = 0
	for _, value := range values {
		i++
		key := strconv.Itoa(i)
		mval := types.MlrvalFromInferredTypeForDataFiles(value)
		record.PutReference(key, mval)
	}
	return record
}
