package main

import (
	"bufio"
	"container/list"
	"errors"
	"fmt"
	"os"

	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/input"
	"github.com/johnkerl/miller/pkg/output"
	"github.com/johnkerl/miller/pkg/transformers"
	"github.com/johnkerl/miller/pkg/types"
)

func convert_csv_to_json(fileNames []string) error {
	options := &cli.TOptions{
		ReaderOptions: cli.TReaderOptions{
			InputFileFormat: "csv",
			IFS:             ",",
			IRS:             "\n",
			RecordsPerBatch: 1,
		},
		WriterOptions: cli.TWriterOptions{
			OutputFileFormat: "json",
		},
	}
	outputStream := os.Stdout
	outputIsStdout := true

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

	cat, err := transformers.NewTransformerCat(
		false, // doCounters bool,
		"",    // counterFieldName string,
		nil,   // groupByFieldNames []string,
		false, // doFileName bool,
		false, // doFileNum bool,
	)
	if err != nil {
		return err
	}
	recordTransformers := []transformers.IRecordTransformer{cat}

	// Set up the reader-to-transformer and transformer-to-writer channels.
	readerChannel := make(chan *list.List, 2) // list of *types.RecordAndContext
	writerChannel := make(chan *list.List, 1) // list of *types.RecordAndContext

	// We're done when a fatal error is registered on input (file not found,
	// etc) or when the record-writer has written all its output. We use
	// channels to communicate both of these conditions.
	inputErrorChannel := make(chan error, 1)
	doneWritingChannel := make(chan bool, 1)
	dataProcessingErrorChannel := make(chan bool, 1)

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

func main() {
	err := convert_csv_to_json(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}
