package output

import (
	// System:
	"os"
	// Miller:
	"miller/containers"
)

func ChannelWriter(
	outrecs <-chan *containers.Lrec,
	recordWriter RecordWriter,
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
