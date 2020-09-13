package input

import (
	"errors"
	"os"

	"encoding/json"

	"miller/clitypes"
	"miller/lib"
)

type RecordReaderJSON struct {
}

func NewRecordReaderJSON(readerOptions *clitypes.TReaderOptions) *RecordReaderJSON {
	return &RecordReaderJSON{}
}

func (this *RecordReaderJSON) Read(
	filenames []string,
	context types.Context,
	inputChannel chan<- *lib.RecordAndContext,
	errorChannel chan error,
) {
	if len(filenames) == 0 { // read from stdin
		handle := os.Stdin
		this.processHandle(handle, "(stdin)", &context, inputChannel, errorChannel)
	} else {
		for _, filename := range filenames {
			handle, err := os.Open(filename)
			if err != nil {
				errorChannel <- err
			} else {
				this.processHandle(handle, filename, &context, inputChannel, errorChannel)
				handle.Close()
			}
		}
	}
	inputChannel <- lib.NewRecordAndContext(
		nil, // signals end of input record stream
		&context,
	)
}

func (this *RecordReaderJSON) processHandle(
	handle *os.File,
	filename string,
	context *types.Context,
	inputChannel chan<- *lib.RecordAndContext,
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
			context.UpdateForInputRecord(record)
			inputChannel <- lib.NewRecordAndContext(
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
					errorChannel <- errors.New("Valid but unmillerable JSON")
					return
				}
				record := mlrval.GetMap()
				if record == nil {
					errorChannel <- errors.New("Internal coding error detected in JSON record-reader")
					return
				}
				context.UpdateForInputRecord(record)
				inputChannel <- lib.NewRecordAndContext(
					record,
					context,
				)

			}

		} else {
			// TODO: more context
			errorChannel <- errors.New("Valid but unmillerable JSON")
			return
		}
	}
}
