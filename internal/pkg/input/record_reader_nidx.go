package input

import (
	"io"
	"strconv"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/types"
)

type RecordReaderNIDX struct {
	readerOptions *cli.TReaderOptions
}

func NewRecordReaderNIDX(readerOptions *cli.TReaderOptions) (*RecordReaderNIDX, error) {
	return &RecordReaderNIDX{
		readerOptions: readerOptions,
	}, nil
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

	scanner := NewLineScanner(handle, reader.readerOptions.IRS)

	for scanner.Scan() {

		// See if downstream processors will be ignoring further data (e.g. mlr
		// head).  If so, stop reading. This makes 'mlr head hugefile' exit
		// quickly, as it should.
		eof := false
		select {
		case _ = <-downstreamDoneChannel:
			eof = true
			break
		default:
			break
		}
		if eof {
			break
		}

		// TODO: IRS
		line := scanner.Text()

		// Check for comments-in-data feature
		if strings.HasPrefix(line, reader.readerOptions.CommentString) {
			if reader.readerOptions.CommentHandling == cli.PassComments {
				inputChannel <- types.NewOutputString(line+"\n", context)
				continue
			} else if reader.readerOptions.CommentHandling == cli.SkipComments {
				continue
			}
			// else comments are data
		}

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

	var values []string
	if reader.readerOptions.IFSRegex == nil { // e.g. --no-ifs-regex
		values = lib.SplitString(line, reader.readerOptions.IFS)
	} else {
		values = lib.RegexSplitString(reader.readerOptions.IFSRegex, line, -1)
	}

	if reader.readerOptions.AllowRepeatIFS {
		values = lib.StripEmpties(values) // left/right trim
	}

	var i int = 0
	for _, value := range values {
		i++
		key := strconv.Itoa(i)
		mval := types.MlrvalFromInferredTypeForDataFiles(value)
		record.PutReference(key, mval)
	}
	return record
}
