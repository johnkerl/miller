package input

import (
	// System:
	"bufio"
	"io"
	"strconv"
	"strings"
	// Miller:
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

	istream, err := lib.Argf(filenames) // can't stay -- each CSV file has its own header, etc
	if err != nil {
		echan <- err
		return
	}
	lineReader := bufio.NewReader(istream)

	if len(filenames) > 0 {
		context.UpdateForStartOfFile(filenames[0]) // xxx temp
	}

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

	inrecs <- nil // signals end of input record stream
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
		copy := value // copy
		lrec.Put(&key, &copy)
	}
	return lrec
}
