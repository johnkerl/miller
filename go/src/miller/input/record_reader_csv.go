package input

import (
	"encoding/csv"
	"io"
	"os"

	"miller/clitypes"
	"miller/lib"
)

type RecordReaderCSV struct {
	// TODO: parameterize
	//ifs string
	//irs string
}

func NewRecordReaderCSV(readerOptions *clitypes.TReaderOptions) *RecordReaderCSV {
	return &RecordReaderCSV{
		//ifs,
		//irs,
	}
}

func (this *RecordReaderCSV) Read(
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

func (this *RecordReaderCSV) processHandle(
	handle *os.File,
	filename string,
	context *types.Context,
	inputChannel chan<- *lib.RecordAndContext,
	errorChannel chan error,
) {
	context.UpdateForStartOfFile(filename)
	needHeader := true
	var header []string = nil

	csvReader := csv.NewReader(handle)

	for {
		if needHeader {
			// TODO: make this a helper function
			csvRecord, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				errorChannel <- err
				return
			}
			header = csvRecord

			needHeader = false
		}

		csvRecord, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errorChannel <- err
			return
		}

		record := lib.NewMlrmap()

		// TODO: check for length mismatches
		n := len(header)
		for i := 0; i < n; i++ {
			key := header[i]
			value := types.MlrvalFromInferredType(csvRecord[i])
			record.PutReference(&key, &value)
		}
		context.UpdateForInputRecord(record)

		inputChannel <- lib.NewRecordAndContext(
			record,
			context,
		)
	}
}
