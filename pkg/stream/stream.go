package stream

import (
	"bufio"
	"container/list"
	"errors"
	"io"

	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/input"
	"github.com/johnkerl/miller/pkg/output"
	"github.com/johnkerl/miller/pkg/transformers"
	"github.com/johnkerl/miller/pkg/types"
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
	initialContext := types.NewContext()

	// Instantiate the record-reader.
	// RecordsPerBatch is tracked separately from ReaderOptions since join/repl
	// may use batch size of 1.
	recordReader, err := input.Create(&options.ReaderOptions, options.ReaderOptions.RecordsPerBatch)
	if err != nil {
		return err
	}

	// Instantiate the record-writer
	recordWriter, err := output.Create(&options.WriterOptions)
	if err != nil {
		return err
	}

	// Set up the reader-to-transformer and transformer-to-writer channels.
	readerChannel := make(chan *list.List, 2) // list of *types.RecordAndContext
	writerChannel := make(chan *list.List, 1) // list of *types.RecordAndContext

	// We're done when a fatal error is registered on input (file not found,
	// etc) or when the record-writer has written all its output. We use
	// channels to communicate both of these conditions.
	inputErrorChannel := make(chan error, 1)
	doneWritingChannel := make(chan bool, 1)
	dataProcessingErrorChannel := make(chan bool, 1)

	// For mlr head, so a transformer can communicate it will disregard all
	// further input.  It writes this back upstream, and that is passed back to
	// the record-reader which then stops reading input. This is necessary to
	// get quick response from, for example, mlr head -n 10 on input files with
	// millions or billions of records.
	readerDownstreamDoneChannel := make(chan bool, 1)

	// Start the reader, transformer, and writer. Let them run until fatal input
	// error or end-of-processing happens.
	bufferedOutputStream := bufio.NewWriter(outputStream)

	go recordReader.Read(fileNames, *initialContext, readerChannel, inputErrorChannel, readerDownstreamDoneChannel)
	go transformers.ChainTransformer(readerChannel, readerDownstreamDoneChannel, recordTransformers,
		writerChannel, options)
	go output.ChannelWriter(writerChannel, recordWriter, &options.WriterOptions, doneWritingChannel,
		dataProcessingErrorChannel, bufferedOutputStream, outputIsStdout)

	var retval error
	done := false
	for !done {
		select {
		case ierr := <-inputErrorChannel:
			retval = ierr
			break
		case _ = <-dataProcessingErrorChannel:
			retval = errors.New("exiting due to data error") // details already printed
			break
		case _ = <-doneWritingChannel:
			done = true
			break
		}
	}

	bufferedOutputStream.Flush()

	return retval
}
