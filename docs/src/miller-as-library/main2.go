// This is an example of using Miller as a library.
package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"

	"github.com/johnkerl/miller/pkg/bifs"
	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/input"
	"github.com/johnkerl/miller/pkg/output"
	"github.com/johnkerl/miller/pkg/types"
)

// Put your record-processing logic here.
func custom_record_processor(irac *types.RecordAndContext) (*types.RecordAndContext, error) {
	irec := irac.Record

	v := irec.Get("i")
	if v == nil {
		return nil, fmt.Errorf("did not find key \"i\" at filename %s record number %d",
			irac.Context.FILENAME, irac.Context.FNR,
		)
	}
	v2 := bifs.BIF_times(v, v)
	irec.PutReference("i2", v2)

	return irac, nil
}

// Put your various options here.
func custom_options() *cli.TOptions {
	return &cli.TOptions{
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
}

// This function you don't need to modify.
func convert_csv_to_json(
	fileNames []string,
	options *cli.TOptions,
	record_processor func (irac *types.RecordAndContext) (*types.RecordAndContext, error),
) error {
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

	// Set up the channels for the record-reader.
	readerChannel := make(chan *list.List, 2) // list of *types.RecordAndContext
	inputErrorChannel := make(chan error, 1)
	// Not needed in this example
	readerDownstreamDoneChannel := make(chan bool, 1)

	// Instantiate the record-writer
	recordWriter, err := output.Create(&options.WriterOptions)
	if err != nil {
		return err
	}
	bufferedOutputStream := bufio.NewWriter(outputStream)

	// Start the record-reader.
	go recordReader.Read(
		fileNames, *initialContext, readerChannel, inputErrorChannel, readerDownstreamDoneChannel)

	// Loop through the record stream.
	var retval error
	done := false
	for !done {
		select {

		case ierr := <-inputErrorChannel:
			retval = ierr
			break

		case iracs := <-readerChannel:
			// Handle the record batch
			for e := iracs.Front(); e != nil; e = e.Next() {
				irac := e.Value.(*types.RecordAndContext)
				if irac.Record != nil {
					orac, err := record_processor(irac)
					if err != nil {
						retval = err
						done = true
						break
					}
					recordWriter.Write(orac.Record, bufferedOutputStream, outputIsStdout)
				}
				if irac.OutputString != "" {
					fmt.Fprintln(bufferedOutputStream, irac.OutputString)
				}
				if irac.EndOfStream {
					done = true
				}
			}
			break

		}
	}

	bufferedOutputStream.Flush()

	return retval
}

func main() {
	options := custom_options()
	err := convert_csv_to_json(os.Args[1:], options, custom_record_processor)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}
