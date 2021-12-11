// This is mostly-identical code for the DKVP and NIDX record-readers.

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

// splitter_DKVP_NIDX is a function type for the one bit of code differing
// between the DKVP reader and the NIDX reader, namely, how it splits lines.
type splitter_DKVP_NIDX func (reader *RecordReaderDKVPNIDX, line string) *types.Mlrmap

type RecordReaderDKVPNIDX struct {
	readerOptions   *cli.TReaderOptions
	recordsPerBatch int
	splitter splitter_DKVP_NIDX
}

func NewRecordReaderDKVP(
	readerOptions *cli.TReaderOptions,
	recordsPerBatch int,
) (*RecordReaderDKVPNIDX, error) {
	return &RecordReaderDKVPNIDX{
		readerOptions:   readerOptions,
		recordsPerBatch: recordsPerBatch,
		splitter: recordFromDKVPLine,
	}, nil
}

func NewRecordReaderNIDX(
	readerOptions *cli.TReaderOptions,
	recordsPerBatch int,
) (*RecordReaderDKVPNIDX, error) {
	return &RecordReaderDKVPNIDX{
		readerOptions:   readerOptions,
		recordsPerBatch: recordsPerBatch,
		splitter: recordFromNIDXLine,
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
	recordsPerBatch := reader.readerOptions.RecordsPerBatch

	lineScanner := NewLineScanner(handle, reader.readerOptions.IRS)
	linesChannel := make(chan *list.List, recordsPerBatch)
	go channelizedLineScanner(lineScanner, linesChannel, downstreamDoneChannel, recordsPerBatch)

	for {
		recordsAndContexts, eof := reader.getRecordBatch(linesChannel, context)
		readerChannel <- recordsAndContexts
		if eof {
			break
		}
	}
}

// TODO: comment copiously we're trying to handle slow/fast/short/long reads: tail -f, smallfile, bigfile.
func (reader *RecordReaderDKVPNIDX) getRecordBatch(
	linesChannel <-chan *list.List,
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

		record := reader.splitter(reader, line)
		context.UpdateForInputRecord()
		recordAndContext := types.NewRecordAndContext(record, context)
		recordsAndContexts.PushBack(recordAndContext)
	}

	return recordsAndContexts, false
}

func recordFromDKVPLine(reader *RecordReaderDKVPNIDX, line string) *types.Mlrmap {
	record := types.NewMlrmapAsRecord()

	var pairs []string
	// TODO: func-pointer this away
	if reader.readerOptions.IFSRegex == nil { // e.g. --no-ifs-regex
		pairs = lib.SplitString(line, reader.readerOptions.IFS)
	} else {
		pairs = lib.RegexSplitString(reader.readerOptions.IFSRegex, line, -1)
	}

	for i, pair := range pairs {
		var kv []string
		if reader.readerOptions.IPSRegex == nil { // e.g. --no-ips-regex
			kv = strings.SplitN(pair, reader.readerOptions.IPS, 2)
		} else {
			kv = lib.RegexSplitString(reader.readerOptions.IPSRegex, pair, 2)
		}

		if len(kv) == 0 || (len(kv) == 1 && kv[0] == "") {
			// Ignore. This is expected when splitting with repeated IFS.
		} else if len(kv) == 1 {
			// E.g the pair has no equals sign: "a" rather than "a=1" or
			// "a=".  Here we use the positional index as the key. This way
			// DKVP is a generalization of NIDX.
			key := strconv.Itoa(i + 1) // Miller userspace indices are 1-up
			value := types.MlrvalFromInferredTypeForDataFiles(kv[0])
			record.PutReference(key, value)
		} else {
			key := kv[0]
			value := types.MlrvalFromInferredTypeForDataFiles(kv[1])
			record.PutReference(key, value)
		}
	}
	return record
}

func recordFromNIDXLine(reader *RecordReaderDKVPNIDX, line string) *types.Mlrmap {
	record := types.NewMlrmapAsRecord()

	var values []string
	// TODO: func-pointer this away
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
