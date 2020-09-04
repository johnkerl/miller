package output

import (
	"os"

	"miller/lib"
)

func ChannelWriter(
	outrecsAndContexts <-chan *lib.LrecAndContext,
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
