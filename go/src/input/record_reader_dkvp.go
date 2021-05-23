package input

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"miller/src/cliutil"
	"miller/src/lib"
	"miller/src/types"
)

type RecordReaderDKVP struct {
	readerOptions *cliutil.TReaderOptions
	// TODO: parameterize IRS
}

func NewRecordReaderDKVP(readerOptions *cliutil.TReaderOptions) *RecordReaderDKVP {
	return &RecordReaderDKVP{
		readerOptions: readerOptions,
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
			handle, err := lib.OpenStdin(
				this.readerOptions.Prepipe,
				this.readerOptions.FileInputEncoding,
			)
			if err != nil {
				errorChannel <- err
			}
			this.processHandle(handle, "(stdin)", &context, inputChannel, errorChannel)
		} else {
			for _, filename := range filenames {
				handle, err := lib.OpenFileForRead(
					filename,
					this.readerOptions.Prepipe,
					this.readerOptions.FileInputEncoding,
				)
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

func (this *RecordReaderDKVP) processHandle(
	handle io.Reader,
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

			// xxx fix
			// xxx comment multiple places
			// xxx do this only if this.readerOptions.Prepipe != ""
			// xxx libify even that
			if strings.Contains(err.Error(), "file already closed") {
				err = nil
				eof = true
			} else {
				errorChannel <- err
			}
		} else {
			// This is how to do a chomp:
			line = strings.TrimRight(line, "\n")

			// xxx temp pending autodetect, and pending more windows-port work
			line = strings.TrimRight(line, "\r")

			record := this.recordFromDKVPLine(&line)
			context.UpdateForInputRecord()
			inputChannel <- types.NewRecordAndContext(
				record,
				context,
			)
		}
	}
}

// ----------------------------------------------------------------
func (this *RecordReaderDKVP) recordFromDKVPLine(
	line *string,
) *types.Mlrmap {
	record := types.NewMlrmap()
	pairs := lib.SplitString(*line, this.readerOptions.IFS)
	for i, pair := range pairs {
		kv := strings.SplitN(pair, this.readerOptions.IPS, 2)
		// TODO check length 0. also, check input is empty since "".split() -> [""] not []
		if len(kv) == 1 {
			// E.g the pair has no equals sign: "a" rather than "a=1" or
			// "a=".  Here we use the positional index as the key. This way
			// DKVP is a generalization of NIDX.
			key := strconv.Itoa(i + 1) // Miller userspace indices are 1-up
			value := types.MlrvalPointerFromInferredType(kv[0])
			record.PutReference(key, value)
		} else {
			key := kv[0]
			value := types.MlrvalPointerFromInferredType(kv[1])
			record.PutReference(key, value)
		}
	}
	return record
}
