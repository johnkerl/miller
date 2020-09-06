package output

import (
	"os"

	"miller/lib"
)

func ChannelWriter(
	outputChannel <-chan *lib.RecordAndContext,
	recordWriter IRecordWriter,
	done chan<- bool,
	ostream *os.File,
) {
	for {
		recordAndContext := <-outputChannel
		record := recordAndContext.Record
		recordWriter.Write(record)
		if record == nil { // end of record stream
			done <- true
			break
		}
	}
}
