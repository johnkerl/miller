package stream

import (
	// System:
	"errors"
	"log"
	"os"
	// Miller:
	"containers"
	"input"
	"mapping"
	"output"
)

// ----------------------------------------------------------------
func Stream(
	filenames []string,
	inputFormatName string,
	mapperName string,
	outputFormatName string,
) error {

	recordReader := input.Create(inputFormatName)
	if recordReader == nil {
		return errors.New("Input format not found: " + inputFormatName)
	}

	recordMapper := mapping.Create(mapperName)
	if recordMapper == nil {
		return errors.New("Mapper not found: " + mapperName)
	}

	inrecs := make(chan *containers.Lrec, 10)
	echan := make(chan error, 1)
	outrecs := make(chan *containers.Lrec, 1)
	donechan := make(chan bool, 1)

	go recordReader.Read(filenames, inrecs, echan)
	go mapping.ChannelMapper(inrecs, recordMapper, outrecs)
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
