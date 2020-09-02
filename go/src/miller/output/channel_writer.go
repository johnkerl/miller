package output

import (
	"os"

	"miller/containers"
)

func ChannelWriter(
	outrecsAndContexts <-chan *containers.LrecAndContext,
	recordWriter IRecordWriter,
	done chan<- bool,
	ostream *os.File,
) {
	for {
		lrecAndContext := <-outrecsAndContexts
		lrec := lrecAndContext.Lrec
		recordWriter.Write(lrec)
		if lrec == nil { // end of record stream
			done <- true
			break
		}
	}
}
