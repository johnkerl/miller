package input

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"

	"miller/clitypes"
	"miller/types"
)

// ----------------------------------------------------------------
type RecordReaderCSV struct {
	ifs string
	// TODO: parameterize for ASV.
	//irs string
	useImplicitHeader bool
}

// ----------------------------------------------------------------
func NewRecordReaderCSV(readerOptions *clitypes.TReaderOptions) *RecordReaderCSV {
	return &RecordReaderCSV{
		ifs:               readerOptions.IFS,
		useImplicitHeader: readerOptions.UseImplicitCSVHeader,
	}
}

// ----------------------------------------------------------------
func (this *RecordReaderCSV) Read(
	filenames []string,
	context types.Context,
	inputChannel chan<- *types.RecordAndContext,
	errorChannel chan error,
) {
	if filenames != nil { // nil for mlr -n
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
	}
	inputChannel <- types.NewRecordAndContext(
		nil, // signals end of input record stream
		&context,
	)
}

// ----------------------------------------------------------------
func (this *RecordReaderCSV) processHandle(
	handle *os.File,
	filename string,
	context *types.Context,
	inputChannel chan<- *types.RecordAndContext,
	errorChannel chan error,
) {
	context.UpdateForStartOfFile(filename)
	needHeader := !this.useImplicitHeader
	var header []string = nil

	csvReader := csv.NewReader(handle)
	// xxx temp
	csvReader.Comma = rune(this.ifs[0])

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

		record := types.NewMlrmap()

		if header == nil { // implicit CSV header
			n := len(csvRecord)
			header = make([]string, n)
			for i := 0; i < n; i++ {
				header[i] = strconv.Itoa(i + 1)
			}
		}

		// TODO: check for length mismatches
		n := len(header)
		for i := 0; i < n; i++ {
			key := header[i]
			value := types.MlrvalFromInferredType(csvRecord[i])
			record.PutReference(&key, &value)
		}
		context.UpdateForInputRecord(record)

		inputChannel <- types.NewRecordAndContext(
			record,
			context,
		)
	}
}
