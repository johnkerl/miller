// This is mostly-identical code for the DKVP and NIDX record-readers.

package input

import (
	"container/list"
	"io"
	"strconv"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
	"github.com/johnkerl/miller/internal/pkg/types"
)

// splitter_DKVP_NIDX is a function type for the one bit of code differing
// between the DKVP reader and the NIDX reader, namely, how it splits lines.
type line_splitter_DKVP_NIDX func(reader *RecordReaderDKVPNIDX, line string) (*mlrval.Mlrmap, error)

type RecordReaderDKVPNIDX struct {
	readerOptions   *cli.TReaderOptions
	recordsPerBatch int64 // distinct from readerOptions.RecordsPerBatch for join/repl
	lineSplitter    line_splitter_DKVP_NIDX
	fieldSplitter   iFieldSplitter
	pairSplitter    iPairSplitter
}

func NewRecordReaderDKVP(
	readerOptions *cli.TReaderOptions,
	recordsPerBatch int64,
) (*RecordReaderDKVPNIDX, error) {
	return &RecordReaderDKVPNIDX{
		readerOptions:   readerOptions,
		recordsPerBatch: recordsPerBatch,
		lineSplitter:    recordFromDKVPLine,
		fieldSplitter:   newFieldSplitter(readerOptions),
		pairSplitter:    newPairSplitter(readerOptions),
	}, nil
}

func NewRecordReaderNIDX(
	readerOptions *cli.TReaderOptions,
	recordsPerBatch int64,
) (*RecordReaderDKVPNIDX, error) {
	return &RecordReaderDKVPNIDX{
		readerOptions:   readerOptions,
		recordsPerBatch: recordsPerBatch,
		lineSplitter:    recordFromNIDXLine,
		fieldSplitter:   newFieldSplitter(readerOptions),
		pairSplitter:    newPairSplitter(readerOptions),
	}, nil
}

func (reader *RecordReaderDKVPNIDX) Read(
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

func (reader *RecordReaderDKVPNIDX) processHandle(
	handle io.Reader,
	filename string,
	context *types.Context,
	readerChannel chan<- *list.List,
	errorChannel chan<- error,
	downstreamDoneChannel <-chan bool, // for mlr head
) {
	context.UpdateForStartOfFile(filename)
	recordsPerBatch := reader.recordsPerBatch

	lineScanner := NewLineScanner(handle, reader.readerOptions.IRS)
	linesChannel := make(chan *list.List, recordsPerBatch)
	go channelizedLineScanner(lineScanner, linesChannel, downstreamDoneChannel, recordsPerBatch)

	for {
		recordsAndContexts, eof := reader.getRecordBatch(linesChannel, errorChannel, context)
		if recordsAndContexts.Len() > 0 {
			readerChannel <- recordsAndContexts
		}
		if eof {
			break
		}
	}
}

// TODO: comment copiously we're trying to handle slow/fast/short/long reads: tail -f, smallfile, bigfile.
func (reader *RecordReaderDKVPNIDX) getRecordBatch(
	linesChannel <-chan *list.List,
	errorChannel chan<- error,
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

		record, err := reader.lineSplitter(reader, line)
		if err != nil {
			errorChannel <- err
			return
		}
		context.UpdateForInputRecord()
		recordAndContext := types.NewRecordAndContext(record, context)
		recordsAndContexts.PushBack(recordAndContext)
	}

	return recordsAndContexts, false
}

func recordFromDKVPLine(reader *RecordReaderDKVPNIDX, line string) (*mlrval.Mlrmap, error) {
	record := mlrval.NewMlrmapAsRecord()
	dedupeFieldNames := reader.readerOptions.DedupeFieldNames

	pairs := reader.fieldSplitter.Split(line)

	for i, pair := range pairs {
		kv := reader.pairSplitter.Split(pair)

		if len(kv) == 0 || (len(kv) == 1 && kv[0] == "") {
			// Ignore. This is expected when splitting with repeated IFS.
		} else if len(kv) == 1 {
			// E.g the pair has no equals sign: "a" rather than "a=1" or
			// "a=".  Here we use the positional index as the key. This way
			// DKVP is a generalization of NIDX.
			key := strconv.Itoa(i + 1) // Miller userspace indices are 1-up
			value := mlrval.FromDeferredType(kv[0])
			_, err := record.PutReferenceMaybeDedupe(key, value, dedupeFieldNames)
			if err != nil {
				return nil, err
			}
		} else {
			key := kv[0]
			value := mlrval.FromDeferredType(kv[1])
			_, err := record.PutReferenceMaybeDedupe(key, value, dedupeFieldNames)
			if err != nil {
				return nil, err
			}
		}
	}
	return record, nil
}

func recordFromNIDXLine(reader *RecordReaderDKVPNIDX, line string) (*mlrval.Mlrmap, error) {
	record := mlrval.NewMlrmapAsRecord()

	values := reader.fieldSplitter.Split(line)

	var i int = 0
	for _, value := range values {
		i++
		key := strconv.Itoa(i)
		mval := mlrval.FromDeferredType(value)
		record.PutReference(key, mval)
	}
	return record, nil
}
