package input

import (
	// System:
	"bufio"
	"io"
	"strings"
	// Miller:
	"containers"
)

func ChannelReader(
	reader *bufio.Reader,
	inrecs chan<- *containers.Lrec,
	echan chan error,
) {
	eof := false

	for !eof {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			err = nil
			eof = true
		} else if err != nil {
			echan <- err
		} else {

			// This is how to do a chomp:
			line = strings.TrimRight(line, "\n")

			// xxx temp
			ifs := ","
			ips := "="
			lrec := LrecFromDKVPLine(&line, &ifs, &ips)
			inrecs <- lrec
		}
	}

	inrecs <- nil // signals end of input record stream
}
