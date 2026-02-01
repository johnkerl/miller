// This is mostly-identical code for the DKVP and NIDX record-readers.

package input

import (
	"io"
	"strconv"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

// line_splitter_DKVP_NIDX is a function type for the one bit of code differing
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

func (reader *RecordReaderDKVPNIDX) processHandle(
	handle io.Reader,
	filename string,
	context *types.Context,
	readerChannel chan<- []*types.RecordAndContext,
	errorChannel chan<- error,
	downstreamDoneChannel <-chan bool, // for mlr head
) {
	context.UpdateForStartOfFile(filename)
	recordsPerBatch := reader.recordsPerBatch

	lineReader := NewLineReader(handle, reader.readerOptions.IRS)
	linesChannel := make(chan []string, recordsPerBatch)
	go channelizedLineReader(lineReader, linesChannel, downstreamDoneChannel, recordsPerBatch)

	for {
		recordsAndContexts, eof := reader.getRecordBatch(linesChannel, errorChannel, context)
		if len(recordsAndContexts) > 0 {
			readerChannel <- recordsAndContexts
		}
		if eof {
			break
		}
	}
}

// TODO: comment copiously we're trying to handle slow/fast/short/long reads: tail -f, smallfile, bigfile.
func (reader *RecordReaderDKVPNIDX) getRecordBatch(
	linesChannel <-chan []string,
	errorChannel chan<- error,
	context *types.Context,
) (
	recordsAndContexts []*types.RecordAndContext,
	eof bool,
) {
	recordsAndContexts = make([]*types.RecordAndContext, 0)

	lines, more := <-linesChannel
	if !more {
		return recordsAndContexts, true
	}

	for _, line := range lines {

		// Check for comments-in-data feature
		// TODO: function-pointer this away
		if reader.readerOptions.CommentHandling != cli.CommentsAreData {
			if strings.HasPrefix(line, reader.readerOptions.CommentString) {
				if reader.readerOptions.CommentHandling == cli.PassComments {
					recordsAndContexts = append(recordsAndContexts, types.NewOutputString(line+"\n", context))
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
		recordsAndContexts = append(recordsAndContexts, recordAndContext)
	}

	return recordsAndContexts, false
}

func recordFromDKVPLine(reader *RecordReaderDKVPNIDX, line string) (*mlrval.Mlrmap, error) {
	record := mlrval.NewMlrmapAsRecord()
	dedupeFieldNames := reader.readerOptions.DedupeFieldNames

	pairs := reader.fieldSplitter.Split(line)

	// Without --incr-key:
	//   echo 'a,z=b,c' | mlr cat gives 1=a,z=b,3=c
	//   I.e. implicit keys are taken from the 1-up field counter.
	// With it:
	//   echo 'a,z=b,c' | mlr cat gives 1=a,z=b,2=c
	//   I.e. implicit keys are taken from a 1-up count of fields lacking explicit keys.
	incr_key := 0

	for i, pair := range pairs {
		kv := reader.pairSplitter.Split(pair)

		if len(kv) == 0 || (len(kv) == 1 && kv[0] == "") {
			// Ignore. This is expected when splitting with repeated IFS.
		} else if len(kv) == 1 {
			// E.g. the pair has no equals sign: "a" rather than "a=1" or
			// "a=".  Here we use the positional index as the key. This way
			// DKVP is a generalization of NIDX.
			//
			// Also: recall that Miller userspace indices are 1-up.
			var int_key int
			if reader.readerOptions.IncrementImplicitKey {
				int_key = incr_key
			} else {
				int_key = i
			}
			str_key := strconv.Itoa(int_key + 1)
			incr_key++
			value := mlrval.FromDeferredType(kv[0])
			_, err := record.PutReferenceMaybeDedupe(str_key, value, dedupeFieldNames)
			if err != nil {
				return nil, err
			}
		} else {
			str_key := kv[0]
			value := mlrval.FromDeferredType(kv[1])
			_, err := record.PutReferenceMaybeDedupe(str_key, value, dedupeFieldNames)
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
		str_key := strconv.Itoa(i)
		mval := mlrval.FromDeferredType(value)
		record.PutReference(str_key, mval)
	}
	return record, nil
}
