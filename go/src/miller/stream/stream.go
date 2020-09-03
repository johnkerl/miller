package stream

import (
	"errors"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/containers"
	"miller/input"
	"miller/mapping"
	"miller/output"
)

// Since Go is concurrent, the context struct (AWK-like variables such as
// FILENAME, NF, NF, FNR, etc.) needs to be duplicated and passed through the
// channels along with each record.
//
// * Record-readers update FILENAME, FILENUM, NF, NR, FNR within context structs.
//
// * Record-mappers can read these from the context structs.
//
// * Record-writes don't need them (OPS et al. are already in the
//   writer-options struct). However, we have chained mappers using the 'then'
//   command-line syntax. This means a given mapper might be piping its output
//   to a record-writer, or another mapper. So, the lrec-and-context pair goes
//   to the record-writers even though they don't need the contexts.

func Stream(
	options clitypes.TOptions,
	recordMappers []mapping.IRecordMapper,
) error {

	// Since Go is concurrent, the context struct needs to be duplicated and
	// passed through the channels along with each record.
	initialContext := containers.NewContext()

	// Instantiate the record-reader
	recordReader := input.Create(&options.ReaderOptions)
	if recordReader == nil {
		return errors.New("Input format not found: " + options.ReaderOptions.InputFileFormat)
	}

	// Instantiate the record-writer
	recordWriter := output.Create(&options.WriterOptions)
	if recordWriter == nil {
		return errors.New("Output format not found: " + options.WriterOptions.OutputFileFormat)
	}

	// Set up the reader-to-mapper and mapper-to-writer channels.
	inrecs := make(chan *containers.LrecAndContext, 10)
	outrecs := make(chan *containers.LrecAndContext, 1)

	// We're done when a fatal error is registered on input (file not found,
	// etc) or when the record-writer has written all its output. We use
	// channels to communicate both of these conditions.
	echan := make(chan error, 1)
	donechan := make(chan bool, 1)

	// Start the reader, mapper, and writer. Let them run until fatal input
	// error or end-of-processing happens.

	go recordReader.Read(options.FileNames, *initialContext, inrecs, echan)
	go mapping.ChainMapper(inrecs, recordMappers, outrecs)
	go output.ChannelWriter(outrecs, recordWriter, donechan, os.Stdout)

	done := false
	for !done {
		select {
		case err := <-echan:
			fmt.Fprintln(os.Stderr, os.Args[0], ": ", err)
			os.Exit(1)
		case _ = <-donechan:
			done = true
			break
		}
	}

	return nil
}
