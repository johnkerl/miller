package input

import (
	"bufio"
	"io"
	"os"
	"strings"

	"miller/containers"
	"miller/lib"
	"miller/runtime"
)

type RecordReaderDKVP struct {
	ifs string
	ips string
}

func NewRecordReaderDKVP(ifs string, ips string) *RecordReaderDKVP {
	return &RecordReaderDKVP{
		ifs,
		ips,
	}
}

func (this *RecordReaderDKVP) Read(
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

func (this *RecordReaderDKVP) processHandle(
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
			lrec := lrecFromDKVPLine(&line, &this.ifs, &this.ips)
			inrecs <- lrec
		}
	}
}

// ----------------------------------------------------------------
func lrecFromDKVPLine(
	line *string,
	ifs *string,
	ips *string,
) *containers.Lrec {
	lrec := containers.LrecAlloc()
	pairs := strings.Split(*line, *ifs)
	for _, pair := range pairs {
		kv := strings.SplitN(pair, *ips, 2)
		// xxx range-check
		key := kv[0]
		value := lib.MlrvalFromInferredType(kv[1])
		// to do: avoid re-walk ...
		lrec.Put(&key, &value)
	}
	return lrec
}
