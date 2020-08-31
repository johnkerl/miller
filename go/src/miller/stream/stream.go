package stream

import (
	"errors"
	"log"
	"os"

	"miller/containers"
	"miller/input"
	"miller/mapping"
	"miller/output"
	"miller/runtime"
)

// ----------------------------------------------------------------
func Stream(
	filenames []string,
	inputFormatName string,
	mapperName string,
	dslString string, // xxx temp
	outputFormatName string,
) error {

	context := runtime.NewContext()

	// xxx context:
	// * NF/NR/FNR
	// * FILENAME/FILENUM

	recordReader := input.Create(inputFormatName)
	if recordReader == nil {
		return errors.New("Input format not found: " + inputFormatName)
	}

	recordMapper, err := mapping.Create(mapperName, dslString) // xxx temp
	if err != nil {
		return err
	}
	if recordMapper == nil {
		return errors.New("Mapper not found: " + mapperName)
	}

	recordWriter := output.Create(outputFormatName)
	if recordWriter == nil {
		return errors.New("Output format not found: " + outputFormatName)
	}

	inrecs := make(chan *containers.Lrec, 10)
	echan := make(chan error, 1)
	outrecs := make(chan *containers.Lrec, 1)
	donechan := make(chan bool, 1)

	// xxx need callback/arg to update context filenum/filename
	go recordReader.Read(filenames, context, inrecs, echan)
	go mapping.ChannelMapper(inrecs, context, recordMapper, outrecs)
	go output.ChannelWriter(outrecs, recordWriter, donechan, os.Stdout)

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
