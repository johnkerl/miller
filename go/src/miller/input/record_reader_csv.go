package input

import (
	"encoding/csv"
	"io"
	"os"

	"miller/clitypes"
	"miller/containers"
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
	context containers.Context,
	inrecsAndContexts chan<- *containers.LrecAndContext,
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
	inrecsAndContexts <- containers.NewLrecAndContext(
		nil, // signals end of input record stream
		&context,
	)
}

func (this *RecordReaderCSV) processHandle(
	handle *os.File,
	filename string,
	context *containers.Context,
	inrecsAndContexts chan<- *containers.LrecAndContext,
	echan chan error,
) {
	context.UpdateForStartOfFile(filename)
	needHeader := true
	var header []string = nil

	csvReader := csv.NewReader(handle)

	for {
		if needHeader {
			// TODO: make this a helper function
			record, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				echan <- err
				return
			}
			header = record

			needHeader = false
		}

		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			echan <- err
			return
		}

		lrec := containers.NewLrec()

		// TODO: check for length mismatches
		n := len(header)
		for i := 0; i < n; i++ {
			key := header[i]
			value := lib.MlrvalFromInferredType(record[i])
			// to do: avoid re-walk ...
			lrec.Put(&key, &value)
		}
		context.UpdateForInputRecord(lrec)

		inrecsAndContexts <- containers.NewLrecAndContext(
			lrec,
			context,
		)
	}
}
