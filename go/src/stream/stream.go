package stream

import (
	// System:
	"bufio"
	"log"
	"os"
	// Miller:
	"containers"
	"input"
	"mapping"
	"output"
)

// ----------------------------------------------------------------
func Stream(filenames []string) error {
	istream, err := Argf(filenames)
	if err != nil {
		return err
		os.Exit(1)
	}
	reader := bufio.NewReader(istream)

	inrecs := make(chan *containers.Lrec, 10)
	echan := make(chan error, 1)
	outrecs := make(chan *containers.Lrec, 1)
	donechan := make(chan bool, 1)

	go input.ChannelReader(reader, inrecs, echan)
	go mapping.ChannelMapper(inrecs, outrecs)
	go output.ChannelWriter(outrecs, donechan, os.Stdout)

	done := false
	for !done {
		select {
		case err := <-echan:
			log.Fatal(err)
		case _ = <-donechan:
			done = true
			break
		}
	}

	return nil
}
