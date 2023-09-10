package output

import (
	"bufio"
	"container/list"
	"fmt"
	"os"

	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/types"
)

func ChannelWriter(
	writerChannel <-chan *list.List, // list of *types.RecordAndContext
	recordWriter IRecordWriter,
	writerOptions *cli.TWriterOptions,
	doneChannel chan<- bool,
	dataProcessingErrorChannel chan<- bool,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) {

	for {
		recordsAndContexts := <-writerChannel
		done, errored := channelWriterHandleBatch(
			recordsAndContexts,
			recordWriter,
			writerOptions,
			dataProcessingErrorChannel,
			bufferedOutputStream,
			outputIsStdout,
		)
		if errored {
			dataProcessingErrorChannel <- true
			doneChannel <- true
			break
		}
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
	dataProcessingErrorChannel chan<- bool,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) (done bool, errored bool) {
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

			// XXX more
			// XXX also make sure this results in exit 1 & goroutine cleanup
			if writerOptions.FailOnDataError {
				ok := true
				for pe := record.Head; pe != nil; pe = pe.Next {
					if pe.Value.IsError() {
						context := recordAndContext.Context
						fmt.Fprintf(os.Stderr, "mlr: data error at NR=%d FNR=%d FILENAME=%s\n",
							context.NR, context.FNR, context.FILENAME,
						)
						is, err := pe.Value.GetError()
						if is {
							if err != nil {
								fmt.Fprintf(os.Stderr, "mlr: field %s: %v\n", pe.Key, err)
							} else {
								fmt.Fprintf(os.Stderr, "mlr: field %s\n", pe.Key)
							}
							ok = false
						}
					}
				}
				if !ok {
					return true, true
				}
			}

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
			return true, false
		}
	}
	return false, false
}
