package input

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"

	"miller/containers"
	"miller/lib"
	"miller/runtime"
)

type RecordReaderNIDX struct {
	// TODO: parameterize
	//ifs string
	//ips string
}

func NewRecordReaderNIDX() *RecordReaderNIDX {
	return &RecordReaderNIDX{
		//ifs,
		//ips,
	}
}

func (this *RecordReaderNIDX) Read(
	filenames []string,
	context *runtime.Context,
	inrecs chan<- *containers.Lrec,
	echan chan error,
) {
	if len(filenames) == 0 { // read from stdin
		handle := os.Stdin
		this.processHandle(handle, "(stdin)", context, inrecs, echan)
	} else {
		for _, filename := range filenames {
			handle, err := os.Open(filename)
			if err != nil {
				echan <- err
			} else {
				this.processHandle(handle, filename, context, inrecs, echan)
				handle.Close()
			}
		}
	}
	inrecs <- nil // signals end of input record stream
}

func (this *RecordReaderNIDX) processHandle(
	handle *os.File,
	filename string,
	context *runtime.Context,
	inrecs chan<- *containers.Lrec,
	echan chan error,
) {
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
			inrecs <- lrec
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
