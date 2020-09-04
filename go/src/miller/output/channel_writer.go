package output

import (
	"os"

	"miller/lib"
)

func ChannelWriter(
	outrecsAndContexts <-chan *lib.RecordAndContext,
	recordWriter IRecordWriter,
	done chan<- bool,
	ostream *os.File,
) {
	for {
		recordAndContext := <-outrecsAndContexts
		record := recordAndContext.Record
		recordWriter.Write(record)
		if record == nil { // end of record stream
			done <- true
			break
		}
	}
}
