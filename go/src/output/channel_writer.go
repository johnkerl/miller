package output

import (
	// System:
	"os"
	// Miller:
	"containers"
)

func ChannelWriter(
	outrecs <-chan *containers.Lrec,
	done chan<- bool,
	ostream *os.File,
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
