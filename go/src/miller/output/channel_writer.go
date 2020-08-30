package output

import (
	"os"

	"miller/containers"
)

func ChannelWriter(
	outrecs <-chan *containers.Lrec,
	recordWriter IRecordWriter,
	done chan<- bool,
	ostream *os.File,
) {
	for {
		lrec := <-outrecs
		recordWriter.Write(lrec)
		if lrec == nil {
			done <- true
			break
		}
	}
}
