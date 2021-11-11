package stream

import (
	"fmt"
	"io"
	"os"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/input"
	"mlr/internal/pkg/output"
	"mlr/internal/pkg/transformers"
	"mlr/internal/pkg/types"
)

// Since Go is concurrent, the context struct (AWK-like variables such as
// FILENAME, NF, NF, FNR, etc.) needs to be duplicated and passed through the
// channels along with each record.
//
// * Record-readers update FILENAME, FILENUM, NF, NR, FNR within context structs.
//
// * Record-transformers can read these from the context structs.
//
// * Record-writers don't need them (OPS et al. are already in the
//   writer-options struct). However, we have chained transformers using the
//   'then' command-line syntax. This means a given transformer might be piping
//   its output to a record-writer, or another transformer. So, the
//   record-and-context pair goes to the record-writers even though they don't
//   need the contexts.

// Stream is the high-level sketch of Miller. It coordinates instantiating
// format-specific record-reader and record-writer objects, using flags from
// the command line; setting up I/O channels; running the record stream from
// the record-reader object, through the specified chain of transformers
// (verbs), to the record-writer object.
func Stream(
	// fileNames argument is separate from options.FileNames for in-place mode,
	// which sends along only one file name per call to Stream():
	fileNames []string,
	options *cli.TOptions,
	recordTransformers []transformers.IRecordTransformer,
	outputStream io.WriteCloser,
	outputIsStdout bool,
) error {

	// Since Go is concurrent, the context struct needs to be duplicated and
	// passed through the channels along with each record.
	initialContext := types.NewContext(
		options.ReaderOptions.IPS,
		options.ReaderOptions.IFS,
		options.ReaderOptions.IRS,
		options.WriterOptions.OPS,
		options.WriterOptions.OFS,
		options.WriterOptions.ORS,
		options.WriterOptions.FLATSEP,
	)

	// Instantiate the record-reader
	recordReader, err := input.Create(&options.ReaderOptions)
	if err != nil {
		return err
	}

	// Instantiate the record-writer
	recordWriter, err := output.Create(&options.WriterOptions)
	if err != nil {
		return err
	}

	// Set up the reader-to-transformer and transformer-to-writer channels.
	inputChannel := make(chan *types.RecordAndContext, 10)
	outputChannel := make(chan *types.RecordAndContext, 1)

	// We're done when a fatal error is registered on input (file not found,
	// etc) or when the record-writer has written all its output. We use
	// channels to communicate both of these conditions.
	errorChannel := make(chan error, 1)
	doneWritingChannel := make(chan bool, 1)

	// For mlr head, so a transformer can communicate it will disregard all
	// further input.  It writes this back upstream, and that is passed back to
	// the record-reader which then stops reading input. This is necessary to
	// get quick response from, for example, mlr head -n 10 on input files with
	// millions or billions of records.
	readerDownstreamDoneChannel := make(chan bool, 1)

	// Start the reader, transformer, and writer. Let them run until fatal input
	// error or end-of-processing happens.

	go recordReader.Read(fileNames, *initialContext, inputChannel, errorChannel, readerDownstreamDoneChannel)
	go transformers.ChainTransformer(inputChannel, readerDownstreamDoneChannel, recordTransformers, outputChannel,
		options)
	go output.ChannelWriter(outputChannel, recordWriter, doneWritingChannel, outputStream, outputIsStdout)

	done := false
	for !done {
		select {
		case err := <-errorChannel:
			fmt.Fprintln(os.Stderr, "mlr", ": ", err)
			os.Exit(1)
		case _ = <-doneWritingChannel:
			done = true
			break
		}
	}

	return nil
}
