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
	context lib.Context,
	inrecsAndContexts chan<- *lib.RecordAndContext,
	echan chan error,
) {
	if len(filenames) == 0 { // read from stdin
		handle := os.Stdin
		this.processHandle(handle, "(stdin)", &context, inrecsAndContexts, echan)
	} else {
		for _, filename := range filenames {
			handle, err := os.Open(filename)
			if err != nil {
				echan <- err
			} else {
				this.processHandle(handle, filename, &context, inrecsAndContexts, echan)
				handle.Close()
			}
		}
	}
	inrecsAndContexts <- lib.NewRecordAndContext(
		nil, // signals end of input record stream
		&context,
	)
}

func (this *RecordReaderJSON) processHandle(
	handle *os.File,
	filename string,
	context *lib.Context,
	inrecsAndContexts chan<- *lib.RecordAndContext,
	echan chan error,
) {
	context.UpdateForStartOfFile(filename)
	decoder := json.NewDecoder(handle)

	for {
		mlrval, eof, err := lib.MlrvalDecodeFromJSON(decoder)
		if eof {
			break
		}
		if err != nil {
			echan <- err
			return
		}

		// Find out what we got.
		// * Map is an input record: deliver it.
		// * Array is OK if it's array of input record: deliver them.
		// * Non-collection types are valid but unmillerable JSON.

		if mlrval.IsMap() {
			record := mlrval.GetMap()
			if record == nil {
				echan <- errors.New("Internal coding error detected in JSON record-reader")
				return
			}
			context.UpdateForInputRecord(record)
			inrecsAndContexts <- lib.NewRecordAndContext(
				record,
				context,
			)
		} else {
		}
	}
}
