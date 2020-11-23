package input

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"

	"miller/clitypes"
	"miller/lib"
	"miller/types"
)

type RecordReaderDKVP struct {
	ifs string
	ips string
	// TODO: parameterize IRS
}

func NewRecordReaderDKVP(readerOptions *clitypes.TReaderOptions) *RecordReaderDKVP {
	return &RecordReaderDKVP{
		ifs: readerOptions.IFS,
		ips: readerOptions.IPS,
	}
}

func (this *RecordReaderDKVP) Read(
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
	inputChannel <- types.NewRecordAndContext(
		nil, // signals end of input record stream
		&context,
	)
}

func (this *RecordReaderDKVP) processHandle(
	handle *os.File,
	filename string,
	context *types.Context,
	inputChannel chan<- *types.RecordAndContext,
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
			inputChannel <- types.NewRecordAndContext(
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
) *types.Mlrmap {
	record := types.NewMlrmap()
	pairs := lib.SplitString(*line, *ifs)
	for i, pair := range pairs {
		kv := strings.SplitN(pair, *ips, 2)
		// TODO check length 0. also, check input is empty since "".split() -> [""] not []
		if len(kv) == 1 {
			// E.g the pair has no equals sign: "a" rather than "a=1" or
			// "a=".  Here we use the positional index as the key. This way
			// DKVP is a generalization of NIDX.
			key := strconv.Itoa(i + 1) // Miller userspace indices are 1-up
			value := types.MlrvalFromInferredType(kv[0])
			record.PutReference(&key, &value)
		} else {
			key := kv[0]
			value := types.MlrvalFromInferredType(kv[1])
			record.PutReference(&key, &value)
		}
	}
	return record
}
