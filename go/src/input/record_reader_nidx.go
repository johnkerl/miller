package input

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"mlr/src/cliutil"
	"mlr/src/lib"
	"mlr/src/types"
)

type RecordReaderNIDX struct {
	// TODO: use the parameterization for readerOptions.IFS/readerOptions.IPS
	readerOptions *cliutil.TReaderOptions
}

func NewRecordReaderNIDX(readerOptions *cliutil.TReaderOptions) *RecordReaderNIDX {
	return &RecordReaderNIDX{
		readerOptions: readerOptions,
	}
}

func (reader *RecordReaderNIDX) Read(
	filenames []string,
	context types.Context,
	inputChannel chan<- *types.RecordAndContext,
	errorChannel chan error,
) {
	if filenames != nil { // nil for mlr -n
		if len(filenames) == 0 { // read from stdin
			handle, err := lib.OpenStdin(
				reader.readerOptions.Prepipe,
				reader.readerOptions.PrepipeIsRaw,
				reader.readerOptions.FileInputEncoding,
			)
			if err != nil {
				errorChannel <- err
			}
			reader.processHandle(handle, "(stdin)", &context, inputChannel, errorChannel)
		} else {
			for _, filename := range filenames {
				handle, err := lib.OpenFileForRead(
					filename,
					reader.readerOptions.Prepipe,
					reader.readerOptions.PrepipeIsRaw,
					reader.readerOptions.FileInputEncoding,
				)
				if err != nil {
					errorChannel <- err
				} else {
					reader.processHandle(handle, filename, &context, inputChannel, errorChannel)
					handle.Close()
				}
			}
		}
	}
	inputChannel <- types.NewEndOfStreamMarker(&context)
}

func (reader *RecordReaderNIDX) processHandle(
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
		if lib.IsEOF(err) {
			err = nil
			eof = true
			break
		}

		if err != nil {
			errorChannel <- err
			break
		}

		// Check for comments-in-data feature
		if strings.HasPrefix(line, reader.readerOptions.CommentString) {
			if reader.readerOptions.CommentHandling == cliutil.PassComments {
				inputChannel <- types.NewOutputString(line, context)
				continue
			} else if reader.readerOptions.CommentHandling == cliutil.SkipComments {
				continue
			}
			// else comments are data
		}

		// xxx temp pending autodetect, and pending more windows-port work
		// This is how to do a chomp:
		line = strings.TrimRight(line, "\n")
		line = strings.TrimRight(line, "\r")

		record := recordFromNIDXLine(line, reader.readerOptions.IFS)

		context.UpdateForInputRecord()
		inputChannel <- types.NewRecordAndContext(
			record,
			context,
		)
	}
}

// ----------------------------------------------------------------
func recordFromNIDXLine(
	line string,
	ifs string,
) *types.Mlrmap {
	record := types.NewMlrmap()
	values := lib.SplitString(line, ifs) // TODO: repifs ...
	var i int = 0
	for _, value := range values {
		i++
		key := strconv.Itoa(i)
		mval := types.MlrvalPointerFromInferredType(value)
		record.PutReference(key, mval)
	}
	return record
}
