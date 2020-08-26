package output

import (
	// System:
	"os"
	// Miller:
	"containers"
)

func ChannelWriter(
	ostream *os.File,
	outrecs <-chan *containers.Lrec,
	done chan<- bool,
) {
	for {
		lrec := <-outrecs
		if lrec == nil {
			done <- true
			break
		} else {
			lrec.Print(os.Stdout)
		}
	}
}
