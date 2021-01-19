package input

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"miller/clitypes"
	"miller/lib"
	"miller/types"
)

// ----------------------------------------------------------------
type RecordReaderCSV struct {
	ifs                 string
	useImplicitHeader   bool
	allowRaggedCSVInput bool
	emptyStringMlrval   types.Mlrval
}

// ----------------------------------------------------------------
func NewRecordReaderCSV(readerOptions *clitypes.TReaderOptions) *RecordReaderCSV {
	return &RecordReaderCSV{
		ifs:                 readerOptions.IFS,
		useImplicitHeader:   readerOptions.UseImplicitCSVHeader,
		allowRaggedCSVInput: readerOptions.AllowRaggedCSVInput,
		emptyStringMlrval:   types.MlrvalFromString(""),
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
	inputChannel <- types.NewEndOfStreamMarker(&context)
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
	var rowNumber int64 = 0

	csvReader := csv.NewReader(handle)
	csvReader.Comma = rune(this.ifs[0]) // xxx temp

	for {
		if needHeader {
			// TODO: make this a helper function
			csvRecord, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil && csvRecord == nil {
				// See https://golang.org/pkg/encoding/csv.
				// We handle field-count ourselves.
				errorChannel <- err
				return
			}
			header = csvRecord
			rowNumber++

			needHeader = false
		}

		csvRecord, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil && csvRecord == nil {
			// See https://golang.org/pkg/encoding/csv.
			// We handle field-count ourselves.
			errorChannel <- err
			return
		}
		rowNumber++

		if header == nil { // implicit CSV header
			n := len(csvRecord)
			header = make([]string, n)
			for i := 0; i < n; i++ {
				header[i] = strconv.Itoa(i + 1)
			}
		}

		record := types.NewMlrmap()

		nh := len(header)
		nd := len(csvRecord)

		if nh == nd {
			for i := 0; i < nh; i++ {
				key := header[i]
				value := types.MlrvalFromInferredType(csvRecord[i])
				record.PutReference(key, &value)
			}

		} else {
			if !this.allowRaggedCSVInput {
				err := errors.New(
					fmt.Sprintf(
						"Miller: CSV header/data length mismatch %d != %d "+
							"at filename %s row %d.\n",
						nh, nd, filename, rowNumber,
					),
				)
				errorChannel <- err
				return
			} else {
				i := 0
				n := lib.IntMin2(nh, nd)
				for i = 0; i < n; i++ {
					key := header[i]
					value := types.MlrvalFromInferredType(csvRecord[i])
					record.PutReference(key, &value)
				}
				if nh < nd {
					// if header shorter than data: use 1-up itoa keys
					key := strconv.Itoa(i + 1)
					value := types.MlrvalFromInferredType(csvRecord[i])
					record.PutCopy(key, &value)
				}
				if nh > nd {
					// if header longer than data: use "" values
					for i = nd; i < nh; i++ {
						record.PutCopy(header[i], &this.emptyStringMlrval)
					}
				}
			}
		}

		context.UpdateForInputRecord()

		inputChannel <- types.NewRecordAndContext(
			record,
			context,
		)
	}
}
