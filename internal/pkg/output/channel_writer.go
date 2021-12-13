package output

import (
	"bufio"
	"container/list"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/types"
)

func ChannelWriter(
	writerChannel <-chan *list.List, // list of *types.RecordAndContext
	recordWriter IRecordWriter,
	writerOptions *cli.TWriterOptions,
	doneChannel chan<- bool,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) {

	for {
		recordsAndContexts := <-writerChannel
		done := channelWriterHandleBatch(
			recordsAndContexts,
			recordWriter,
			writerOptions,
			bufferedOutputStream,
			outputIsStdout,
		)
		if done {
			doneChannel <- true
			break
		}
	}
}

// TODO: comment
// Returns true on end of record stream
func channelWriterHandleBatch(
	recordsAndContexts *list.List,
	recordWriter IRecordWriter,
	writerOptions *cli.TWriterOptions,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) bool {
	for e := recordsAndContexts.Front(); e != nil; e = e.Next() {
		recordAndContext := e.Value.(*types.RecordAndContext)

		// Three things can come through:
		// * End-of-stream marker
		// * Non-nil records to be printed
		// * Strings to be printed from put/filter DSL print/dump/etc
		//   statements. They are handled here rather than fmt.Println directly
		//   in the put/filter handlers since we want all print statements and
		//   record-output to be in the same goroutine, for deterministic
		//   output ordering.

		if !recordAndContext.EndOfStream {
			record := recordAndContext.Record
			if record != nil {
				recordWriter.Write(record, bufferedOutputStream, outputIsStdout)
			}

			outputString := recordAndContext.OutputString
			if outputString != "" {
				bufferedOutputStream.WriteString(outputString)
			}

			if writerOptions.FlushOnEveryRecord {
				bufferedOutputStream.Flush()
			}

		} else {
			// Let the record-writers drain their output, if they have any
			// queued up. For example, PPRINT needs to see all same-schema
			// records before printing any, since it needs to compute max width
			// down columns.
			recordWriter.Write(nil, bufferedOutputStream, outputIsStdout)
			return true
		}
	}
	return false
}
