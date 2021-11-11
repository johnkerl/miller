package input

import (
	"container/list"
	"errors"
	"io"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/types"
)

type RecordReaderXTAB struct {
	readerOptions *cli.TReaderOptions
	// Note: XTAB uses two consecutive IFS in place of an IRS; IRS is ignored
}

// ----------------------------------------------------------------
func NewRecordReaderXTAB(readerOptions *cli.TReaderOptions) (*RecordReaderXTAB, error) {
	return &RecordReaderXTAB{
		readerOptions: readerOptions,
	}, nil
}

// ----------------------------------------------------------------
func (reader *RecordReaderXTAB) Read(
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

func (reader *RecordReaderXTAB) processHandle(
	handle io.Reader,
	filename string,
	context *types.Context,
	inputChannel chan<- *types.RecordAndContext,
	errorChannel chan error,
	downstreamDoneChannel <-chan bool, // for mlr head
) {
	context.UpdateForStartOfFile(filename)

	scanner := NewLineScanner(handle, reader.readerOptions.IFS)

	linesForRecord := list.New()

	eof := false
	for !eof {

		// See if downstream processors will be ignoring further data (e.g. mlr
		// head).  If so, stop reading. This makes 'mlr head hugefile' exit
		// quickly, as it should.
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

		if !scanner.Scan() {

			if linesForRecord.Len() > 0 {
				record, err := reader.recordFromXTABLines(linesForRecord)
				if err != nil {
					errorChannel <- err
					return
				}
				context.UpdateForInputRecord()
				inputChannel <- types.NewRecordAndContext(record, context)
				linesForRecord = list.New()
			}

			break
		}

		line := scanner.Text()

		// Check for comments-in-data feature
		if strings.HasPrefix(line, reader.readerOptions.CommentString) {
			if reader.readerOptions.CommentHandling == cli.PassComments {
				inputChannel <- types.NewOutputString(line+reader.readerOptions.IFS, context)
				continue
			} else if reader.readerOptions.CommentHandling == cli.SkipComments {
				continue
			}
			// else comments are data
		}

		if line != "" {
			linesForRecord.PushBack(line)

		} else {
			if linesForRecord.Len() > 0 {
				record, err := reader.recordFromXTABLines(linesForRecord)
				if err != nil {
					errorChannel <- err
					return
				}
				context.UpdateForInputRecord()
				inputChannel <- types.NewRecordAndContext(record, context)
				linesForRecord = list.New()
			}
		}
	}
}

// ----------------------------------------------------------------
func (reader *RecordReaderXTAB) recordFromXTABLines(
	lines *list.List,
) (*types.Mlrmap, error) {
	record := types.NewMlrmap()

	for entry := lines.Front(); entry != nil; entry = entry.Next() {
		line := entry.Value.(string)

		var kv []string
		if reader.readerOptions.IPSRegex == nil { // e.g. --no-ips-regex
			kv = strings.SplitN(line, reader.readerOptions.IPS, 2)
		} else {
			kv = lib.RegexSplitString(reader.readerOptions.IPSRegex, line, 2)
		}
		if len(kv) < 1 {
			return nil, errors.New("mlr: internal coding error in XTAB reader")
		}

		key := kv[0]
		if len(kv) == 1 {
			value := types.MLRVAL_VOID
			record.PutReference(key, value)
		} else {
			value := types.MlrvalFromInferredTypeForDataFiles(kv[1])
			record.PutReference(key, value)
		}
	}

	return record, nil
}
