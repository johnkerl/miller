package input

import (
	"errors"
	"fmt"
	"io"

	"encoding/json"

	"mlr/src/cli"
	"mlr/src/lib"
	"mlr/src/types"
)

type RecordReaderJSON struct {
	readerOptions *cli.TReaderOptions
}

func NewRecordReaderJSON(readerOptions *cli.TReaderOptions) *RecordReaderJSON {
	return &RecordReaderJSON{
		readerOptions: readerOptions,
	}
}

func (reader *RecordReaderJSON) Read(
	filenames []string,
	context types.Context,
	inputChannel chan<- *types.RecordAndContext,
	errorChannel chan error,
) {
	if filenames != nil { // nil for mlr -n
		if len(filenames) == 0 { // read from stdin
			handle, err := lib.OpenStdin(
				reader.readerOptions.Prepipe,
				reader.readerOptions.PrepipeIsRaw,
				reader.readerOptions.FileInputEncoding,
			)
			if err != nil {
				errorChannel <- err
			}
			reader.processHandle(handle, "(stdin)", &context, inputChannel, errorChannel)
		} else {
			for _, filename := range filenames {
				handle, err := lib.OpenFileForRead(
					filename,
					reader.readerOptions.Prepipe,
					reader.readerOptions.PrepipeIsRaw,
					reader.readerOptions.FileInputEncoding,
				)
				if err != nil {
					errorChannel <- err
				} else {
					reader.processHandle(handle, filename, &context, inputChannel, errorChannel)
					handle.Close()
				}
			}
		}
	}
	inputChannel <- types.NewEndOfStreamMarker(&context)
}

func (reader *RecordReaderJSON) processHandle(
	handle io.Reader,
	filename string,
	context *types.Context,
	inputChannel chan<- *types.RecordAndContext,
	errorChannel chan error,
) {
	context.UpdateForStartOfFile(filename)
	decoder := json.NewDecoder(handle)

	for {
		mlrval, eof, err := types.MlrvalDecodeFromJSON(decoder)
		if eof {
			break
		}
		if err != nil {
			errorChannel <- err
			return
		}

		// Find out what we got.
		// * Map is an input record: deliver it.
		// * Array is OK if it's array of input record: deliver them.
		// * Non-collection types are valid but unmillerable JSON.

		if mlrval.IsMap() {
			// TODO: make a helper method
			record := mlrval.GetMap()
			if record == nil {
				errorChannel <- errors.New("Internal coding error detected in JSON record-reader")
				return
			}
			context.UpdateForInputRecord()
			inputChannel <- types.NewRecordAndContext(
				record,
				context,
			)
		} else if mlrval.IsArray() {
			records := mlrval.GetArray()
			if records == nil {
				errorChannel <- errors.New("Internal coding error detected in JSON record-reader")
				return
			}
			for _, mlrval := range records {
				if !mlrval.IsMap() {
					// TODO: more context
					errorChannel <- errors.New(
						fmt.Sprintf(
							"Valid but unmillerable JSON. Expected map (JSON object); got %s.",
							mlrval.GetTypeName(),
						),
					)
					return
				}
				record := mlrval.GetMap()
				if record == nil {
					errorChannel <- errors.New("Internal coding error detected in JSON record-reader")
					return
				}
				context.UpdateForInputRecord()
				inputChannel <- types.NewRecordAndContext(
					record,
					context,
				)

			}

		} else {
			errorChannel <- errors.New(
				fmt.Sprintf(
					"Valid but unmillerable JSON. Expected map (JSON object); got %s.",
					mlrval.GetTypeName(),
				),
			)
			return
		}
	}
}
