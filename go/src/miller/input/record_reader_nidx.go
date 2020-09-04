package input

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"

	"miller/clitypes"
	"miller/lib"
)

type RecordReaderNIDX struct {
	// TODO: use the parameterization
	ifs string
	ips string
}

func NewRecordReaderNIDX(readerOptions *clitypes.TReaderOptions) *RecordReaderNIDX {
	return &RecordReaderNIDX{
		ifs: readerOptions.IFS,
		ips: readerOptions.IPS,
	}
}

func (this *RecordReaderNIDX) Read(
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

func (this *RecordReaderNIDX) processHandle(
	handle *os.File,
	filename string,
	context *lib.Context,
	inrecsAndContexts chan<- *lib.RecordAndContext,
	echan chan error,
) {
	context.UpdateForStartOfFile(filename)

	lineReader := bufio.NewReader(handle)
	eof := false

	for !eof {
		line, err := lineReader.ReadString('\n') // TODO: auto-detect
		if err == io.EOF {
			err = nil
			eof = true
		} else if err != nil {
			echan <- err
		} else {
			// This is how to do a chomp:
			line = strings.TrimRight(line, "\n")
			record := recordFromNIDXLine(&line)

			context.UpdateForInputRecord(record)
			inrecsAndContexts <- lib.NewRecordAndContext(
				record,
				context,
			)
		}
	}
}

// ----------------------------------------------------------------
func recordFromNIDXLine(
	line *string,
) *lib.Mlrmap {
	record := lib.NewMlrmap()
	values := strings.Split(*line, " ") // TODO: repifs ...
	var i int64 = 0
	for _, value := range values {
		i++
		key := strconv.FormatInt(i, 10)
		// to do: avoid re-walk ...
		mval := lib.MlrvalFromInferredType(value)
		record.Put(&key, &mval)
	}
	return record
}
