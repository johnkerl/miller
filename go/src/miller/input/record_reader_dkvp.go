package input

import (
	"bufio"
	"io"
	"os"
	"strings"

	"miller/clitypes"
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
	context lib.Context,
	inputChannel chan<- *lib.RecordAndContext,
	errorChannel chan error,
) {
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
	inputChannel <- lib.NewRecordAndContext(
		nil, // signals end of input record stream
		&context,
	)
}

func (this *RecordReaderDKVP) processHandle(
	handle *os.File,
	filename string,
	context *lib.Context,
	inputChannel chan<- *lib.RecordAndContext,
	errorChannel chan error,
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
			errorChannel <- err
		} else {
			// This is how to do a chomp:
			line = strings.TrimRight(line, "\n")
			record := recordFromDKVPLine(&line, &this.ifs, &this.ips)
			context.UpdateForInputRecord(record)
			inputChannel <- lib.NewRecordAndContext(
				record,
				context,
			)
		}
	}
}

// ----------------------------------------------------------------
func recordFromDKVPLine(
	line *string,
	ifs *string,
	ips *string,
) *lib.Mlrmap {
	record := lib.NewMlrmap()
	pairs := strings.Split(*line, *ifs)
	for _, pair := range pairs {
		kv := strings.SplitN(pair, *ips, 2)
		key := kv[0]
		// xxx check length 0. also, check input is empty since "".split() -> [""] not []
		if len(kv) == 1 {
			value := lib.MlrvalFromVoid()
			record.PutReference(&key, &value)
		} else {
			value := lib.MlrvalFromInferredType(kv[1])
			record.PutReference(&key, &value)
		}
	}
	return record
}
