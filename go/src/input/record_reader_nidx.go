package input

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"mlr/src/cli"
	"mlr/src/lib"
	"mlr/src/types"
)

type RecordReaderNIDX struct {
	readerOptions *cli.TReaderOptions
}

func NewRecordReaderNIDX(readerOptions *cli.TReaderOptions) *RecordReaderNIDX {
	return &RecordReaderNIDX{
		readerOptions: readerOptions,
	}
}

func (reader *RecordReaderNIDX) Read(
	filenames []string,
	context types.Context,
	inputChannel chan<- *types.RecordAndContext,
	errorChannel chan error,
	downstreamDoneChannel <-chan bool, // for mlr head
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
			reader.processHandle(handle, "(stdin)", &context, inputChannel, errorChannel, downstreamDoneChannel)
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
					reader.processHandle(handle, filename, &context, inputChannel, errorChannel, downstreamDoneChannel)
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
	downstreamDoneChannel <-chan bool, // for mlr head
) {
	context.UpdateForStartOfFile(filename)

	lineReader := bufio.NewReader(handle)
	eof := false

	for !eof {

		select {
		case _ = <-downstreamDoneChannel: // e.g. mlr head
			break
		default:
			break
		}

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
			if reader.readerOptions.CommentHandling == cli.PassComments {
				inputChannel <- types.NewOutputString(line, context)
				continue
			} else if reader.readerOptions.CommentHandling == cli.SkipComments {
				continue
			}
			// else comments are data
		}

		// xxx temp pending autodetect, and pending more windows-port work
		// This is how to do a chomp:
		line = strings.TrimRight(line, "\n")
		line = strings.TrimRight(line, "\r")

		record := reader.recordFromNIDXLine(line)

		context.UpdateForInputRecord()
		inputChannel <- types.NewRecordAndContext(
			record,
			context,
		)
	}
}

// ----------------------------------------------------------------
func (reader *RecordReaderNIDX) recordFromNIDXLine(
	line string,
) *types.Mlrmap {
	record := types.NewMlrmap()
	values := lib.RegexSplitString(reader.readerOptions.IFSRegex, line, -1)
	var i int = 0
	for _, value := range values {
		i++
		key := strconv.Itoa(i)
		mval := types.MlrvalPointerFromInferredTypeForDataFiles(value)
		record.PutReference(key, mval)
	}
	return record
}
