package input

import (
	"encoding/csv"
	"io"
	"os"

	"miller/containers"
	"miller/lib"
	"miller/runtime"
)

type RecordReaderCSV struct {
	// TODO: parameterize
	//ifs string
	//irs string
	needHeader bool
	header     []string
}

func NewRecordReaderCSV( /*ifs string, ips string*/ ) *RecordReaderCSV {
	return &RecordReaderCSV{
		true,
		nil,
		//ifs,
		//irs,
	}
}

func (this *RecordReaderCSV) Read(
	filenames []string,
	context *runtime.Context,
	inrecs chan<- *containers.Lrec,
	echan chan error,
) {
	// TODO: loop over filenames
	// TODO: handle empty filenames array as read-from-stdin
	filename := filenames[0]
	context.UpdateForStartOfFile(filename)

	handle, err := os.Open(filename)
	if err != nil {
		echan <- err
	}

	csvReader := csv.NewReader(handle)

	for {
		if this.needHeader {
			// TODO: make this a helper function
			record, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				echan <- err
				return
			}
			this.header = record

			this.needHeader = false
		}

		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			echan <- err
			return
		}

		lrec := containers.LrecAlloc()

		// TODO: check for length mismatches
		n := len(this.header)
		for i := 0; i < n; i++ {
			key := this.header[i]
			value := lib.MlrvalFromInferredType(record[i])
			// to do: avoid re-walk ...
			lrec.Put(&key, &value)
		}

		inrecs <- lrec
	}

	inrecs <- nil // signals end of input record stream
}
