package input

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"

	"miller/clitypes"
	"miller/containers"
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

func (this *RecordReaderNIDX) processHandle(
	handle *os.File,
	filename string,
	context *containers.Context,
	inrecsAndContexts chan<- *containers.LrecAndContext,
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
			lrec := lrecFromNIDXLine(&line)

			context.UpdateForInputRecord(lrec)
			inrecsAndContexts <- containers.NewLrecAndContext(
				lrec,
				context,
			)
		}
	}
}

// ----------------------------------------------------------------
func lrecFromNIDXLine(
	line *string,
) *containers.Lrec {
	lrec := containers.LrecAlloc()
	values := strings.Split(*line, " ") // TODO: repifs ...
	var i int64 = 0
	for _, value := range values {
		i++
		key := strconv.FormatInt(i, 10)
		// to do: avoid re-walk ...
		mval := lib.MlrvalFromInferredType(value)
		lrec.Put(&key, &mval)
	}
	return lrec
}
