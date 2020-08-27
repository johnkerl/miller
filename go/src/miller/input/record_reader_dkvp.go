package input

import (
	// System:
	"bufio"
	"io"
	"strings"
	// Miller:
	"miller/containers"
	"miller/lib"
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
	inrecs chan<- *containers.Lrec,
	echan chan error,
) {

	istream, err := lib.Argf(filenames) // can't stay -- each CSV file has its own header, etc
	if err != nil {
		echan <- err
		return
	}
	lineReader := bufio.NewReader(istream)

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

	inrecs <- nil // signals end of input record stream
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
		value := kv[1]
		// to do: avoid re-walk ...
		lrec.Put(&key, &value)
	}
	return lrec
}
