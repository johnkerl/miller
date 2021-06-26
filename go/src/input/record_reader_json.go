package input

import (
	"errors"
	"fmt"
	"io"
	"os"

	"encoding/json"

	"miller/src/cliutil"
	"miller/src/lib"
	"miller/src/types"
)

type RecordReaderJSON struct {
	readerOptions *cliutil.TReaderOptions
}

func NewRecordReaderJSON(readerOptions *cliutil.TReaderOptions) *RecordReaderJSON {
	return &RecordReaderJSON{
		readerOptions: readerOptions,
	}
}

func (reader *RecordReaderJSON) Read(
	filenames []string,
	context types.Context,
	inputChannel chan<- *types.RecordAndContext,
	warningChannel chan error,
	fatalErrorChannel chan error,
) {
	if filenames != nil { // nil for mlr -n
		if len(filenames) == 0 { // read from stdin
			handle, err := lib.OpenStdin(
				reader.readerOptions.Prepipe,
				reader.readerOptions.PrepipeIsRaw,
				reader.readerOptions.FileInputEncoding,
			)
			if err != nil {
				fatalErrorChannel <- err
			}
			reader.processHandle(handle, "(stdin)", &context, inputChannel, warningChannel, fatalErrorChannel)
		} else {
			for _, filename := range filenames {
				handle, err := lib.OpenFileForRead(
					filename,
					reader.readerOptions.Prepipe,
					reader.readerOptions.PrepipeIsRaw,
					reader.readerOptions.FileInputEncoding,
				)
				if err != nil {
					if reader.readerOptions.KeepGoing {
						fmt.Fprint(os.Stderr, err)
					} else {
						fatalErrorChannel <- err
						return
					}
				} else {
					reader.processHandle(handle, filename, &context, inputChannel, warningChannel, fatalErrorChannel)
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
	warningChannel chan error,
	fatalErrorChannel chan error,
) {
	context.UpdateForStartOfFile(filename)
	decoder := json.NewDecoder(handle)

	for {
		mlrval, eof, err := types.MlrvalDecodeFromJSON(decoder)
		if eof {
			break
		}
		if err != nil {
			if reader.readerOptions.KeepGoing {
				warningChannel <- err
				continue
			} else {
				fatalErrorChannel <- err
				return
			}
		}

		// Find out what we got.
		// * Map is an input record: deliver it.
		// * Array is OK if it's array of input record: deliver them.
		// * Non-collection types are valid but unmillerable JSON.

		if mlrval.IsMap() {
			// TODO: make a helper method
			record := mlrval.GetMap()
			if record == nil {
				fatalErrorChannel <- errors.New("Internal coding error detected in JSON record-reader")
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
				fatalErrorChannel <- errors.New("Internal coding error detected in JSON record-reader")
				return
			}
			for _, mlrval := range records {
				if !mlrval.IsMap() {
					// TODO: more context
					err := errors.New(
						fmt.Sprintf(
							"Valid but unmillerable JSON. Expected map (JSON object); got %s.",
							mlrval.GetTypeName(),
						),
					)
					if reader.readerOptions.KeepGoing {
						fmt.Fprint(os.Stderr, err)
						continue
					} else {
						fatalErrorChannel <- err
						return
					}
				}
				record := mlrval.GetMap()
				if record == nil {
					fatalErrorChannel <- errors.New("Internal coding error detected in JSON record-reader")
					return
				}
				context.UpdateForInputRecord()
				inputChannel <- types.NewRecordAndContext(
					record,
					context,
				)

			}

		} else {
			err := errors.New(
				fmt.Sprintf(
					"Valid but unmillerable JSON. Expected map (JSON object); got %s.",
					mlrval.GetTypeName(),
				),
			)
			if reader.readerOptions.KeepGoing {
				fmt.Fprint(os.Stderr, err)
			} else {
				fatalErrorChannel <- err
				return
			}
		}
	}
}
