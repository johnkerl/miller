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
		if lrec == nil {
			done <- true
			break
		} else {
			recordWriter.Write(lrec)
		}
	}
}
