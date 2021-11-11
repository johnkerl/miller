package output

import (
	"fmt"
	"io"

	"mlr/internal/pkg/types"
)

func ChannelWriter(
	outputChannel <-chan *types.RecordAndContext,
	recordWriter IRecordWriter,
	doneChannel chan<- bool,
	ostream io.WriteCloser,
	outputIsStdout bool,
) {
	for {
		recordAndContext := <-outputChannel

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
				recordWriter.Write(record, ostream, outputIsStdout)
			}

			outputString := recordAndContext.OutputString
			if outputString != "" {
				fmt.Print(outputString)
			}

		} else {
			// Let the record-writers drain their output, if they have any
			// queued up. For example, PPRINT needs to see all same-schema
			// records before printing any, since it needs to compute max width
			// down columns.
			recordWriter.Write(nil, ostream, outputIsStdout)
			doneChannel <- true
			break
		}
	}
}
