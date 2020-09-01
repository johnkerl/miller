package input

import (
	"bufio"
	"io"
	"os"
	"strings"

	"miller/clitypes"
	"miller/containers"
	"miller/lib"
)

type RecordReaderDKVP struct {
	ifs string
	ips string
}

func NewRecordReaderDKVP(readerOptions *clitypes.TReaderOptions) *RecordReaderDKVP {
	return &RecordReaderDKVP{
		ifs: readerOptions.IFS,
		ips: readerOptions.IPS,
	}
}

func (this *RecordReaderDKVP) Read(
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

func (this *RecordReaderDKVP) processHandle(
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
			lrec := lrecFromDKVPLine(&line, &this.ifs, &this.ips)
			context.UpdateForInputRecord(lrec)
			inrecsAndContexts <- containers.NewLrecAndContext(
				lrec,
				context,
			)
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
