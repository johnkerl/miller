package input

import (
	"bufio"
	"container/list"
	"io"
	"strconv"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/types"
)

type RecordReaderDKVP struct {
	readerOptions   *cli.TReaderOptions
	recordsPerBatch int
}

func NewRecordReaderDKVP(
	readerOptions *cli.TReaderOptions,
	recordsPerBatch int,
) (*RecordReaderDKVP, error) {
	return &RecordReaderDKVP{
		readerOptions:   readerOptions,
		recordsPerBatch: recordsPerBatch,
	}, nil
}

func (reader *RecordReaderDKVP) Read(
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

func (reader *RecordReaderDKVP) processHandle(
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
	linesChannel := make(chan string, recordsPerBatch)
	go provideChannelizedLines(lineScanner, linesChannel, downstreamDoneChannel, recordsPerBatch)

	eof := false
	for !eof {
		var recordsAndContexts *list.List
		recordsAndContexts, eof = reader.getRecordBatch(linesChannel, recordsPerBatch, context)
		//fmt.Fprintf(os.Stderr, "GOT RECORD BATCH OF LENGTH %d\n", recordsAndContexts.Len())
		readerChannel <- recordsAndContexts
	}
}

// TODO: comment
func provideChannelizedLines(
	lineScanner *bufio.Scanner,
	linesChannel chan<- string,
	downstreamDoneChannel <-chan bool, // for mlr head
	recordsPerBatch int,
) {
	i := 0
	done := false
	for !done && lineScanner.Scan() {
		i++

		// See if downstream processors will be ignoring further data (e.g. mlr
		// head).  If so, stop reading. This makes 'mlr head hugefile' exit
		// quickly, as it should.
		if i%recordsPerBatch == 0 {
			select {
			case _ = <-downstreamDoneChannel:
				done = true
				break
			default:
				break
			}
			if done {
				break
			}
		}

		linesChannel <- lineScanner.Text()
	}
	close(linesChannel) // end-of-stream marker
}

//// TODO: productionalize this for the case no-head -- if profiling shows it to be worthwhile
//// TODO: comment
//func provideChannelizedLines(
//	lineScanner *bufio.Scanner,
//	linesChannel chan<- *list.List,
//	downstreamDoneChannel <-chan bool, // for mlr head
//	recordsPerBatch int,
//) {
//	i := 0
//	done := false
//
//	lines := list.New()
//
//	for !done && lineScanner.Scan() {
//		i++
//
//		lines.PushBack(lineScanner.Text())
//
//		// See if downstream processors will be ignoring further data (e.g. mlr
//		// head).  If so, stop reading. This makes 'mlr head hugefile' exit
//		// quickly, as it should.
//		if i%recordsPerBatch == 0 {
//			select {
//			case _ = <-downstreamDoneChannel:
//				done = true
//				break
//			default:
//				break
//			}
//			if done {
//				break
//			}
//			linesChannel <- lines
//			lines = list.New()
//		}
//
//		//linesChannel <- lineScanner.Text()
//	}
//	linesChannel <- lines
//	close(linesChannel) // end-of-stream marker
//}

//// TODO: productionalize this for the case no-head -- if profiling shows it to be worthwhile
//// TODO: comment copiously we're trying to handle slow/fast/short/long
//// reads: tail -f, smallfile, bigfile.
//func (reader *RecordReaderDKVP) getRecordBatch(
//	linesChannel <-chan *list.List,
//	maxBatchSize int,
//	context *types.Context,
//) (
//	recordsAndContexts *list.List,
//	eof bool,
//) {
//	//fmt.Printf("GRB ENTER\n")
//	recordsAndContexts = list.New()
//
//	lines, more := <-linesChannel
//	if !more {
//		return recordsAndContexts, true
//	}
//
//	for e := lines.Front(); e != nil; e = e.Next() {
//		line := e.Value.(string)
//
//		// Check for comments-in-data feature
//		if strings.HasPrefix(line, reader.readerOptions.CommentString) {
//			if reader.readerOptions.CommentHandling == cli.PassComments {
//				recordsAndContexts.PushBack(types.NewOutputStringList(line+"\n", context))
//				continue
//			} else if reader.readerOptions.CommentHandling == cli.SkipComments {
//				continue
//			}
//			// else comments are data
//		}
//
//		record := reader.recordFromDKVPLine(line)
//		context.UpdateForInputRecord()
//		recordAndContext := types.NewRecordAndContext(record, context)
//		recordsAndContexts.PushBack(recordAndContext)
//	}
//
//	//fmt.Printf("GRB EXIT\n")
//	return recordsAndContexts, false
//}

// TODO: comment copiously we're trying to handle slow/fast/short/long
// reads: tail -f, smallfile, bigfile.
func (reader *RecordReaderDKVP) getRecordBatch(
	linesChannel <-chan string,
	maxBatchSize int,
	context *types.Context,
) (
	recordsAndContexts *list.List,
	eof bool,
) {
	//fmt.Printf("GRB ENTER\n")
	recordsAndContexts = list.New()
	eof = false

	for i := 0; i < maxBatchSize; i++ {
		//fmt.Fprintf(os.Stderr, "-- %d/%d %d/%d\n", i, maxBatchSize, len(linesChannel), cap(linesChannel))
		if len(linesChannel) == 0 && i > 0 {
			//fmt.Println(" .. BREAK")
			break
		}
		//fmt.Println(" .. B:BLOCK")
		line, more := <-linesChannel
		//fmt.Printf(" .. E:BLOCK <<%s>> %v\n", line, more)
		if !more {
			eof = true
			break
		}

		// Check for comments-in-data feature
		if strings.HasPrefix(line, reader.readerOptions.CommentString) {
			if reader.readerOptions.CommentHandling == cli.PassComments {
				recordsAndContexts.PushBack(types.NewOutputStringList(line+"\n", context))
				continue
			} else if reader.readerOptions.CommentHandling == cli.SkipComments {
				continue
			}
			// else comments are data
		}

		record := reader.recordFromDKVPLine(line)
		context.UpdateForInputRecord()
		recordAndContext := types.NewRecordAndContext(record, context)
		recordsAndContexts.PushBack(recordAndContext)
	}

	//fmt.Printf("GRB EXIT\n")
	return recordsAndContexts, eof
}

func (reader *RecordReaderDKVP) recordFromDKVPLine(
	line string,
) *types.Mlrmap {
	record := types.NewMlrmapAsRecord()

	var pairs []string
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
