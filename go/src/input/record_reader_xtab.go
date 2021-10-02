package input

import (
	"bufio"
	"container/list"
	"errors"
	"io"
	"strings"

	"mlr/src/cli"
	"mlr/src/lib"
	"mlr/src/types"
)

type RecordReaderXTAB struct {
	readerOptions *cli.TReaderOptions
	// TODO: parameterize IRS
}

// ----------------------------------------------------------------
func NewRecordReaderXTAB(readerOptions *cli.TReaderOptions) *RecordReaderXTAB {
	return &RecordReaderXTAB{
		readerOptions: readerOptions,
	}
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

	lineReader := bufio.NewReader(handle)

	linesForRecord := list.New()

	eof := false
	for !eof {

		select {
		case _ = <-downstreamDoneChannel: // e.g. mlr head
			break
		default:
			break
		}

		//line, err := lineReader.ReadString(reader.readerOptions.IRS[0]) // xxx temp
		line, err := lineReader.ReadString('\n')
		if lib.IsEOF(err) {
			err = nil
			eof = true

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
			continue
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
		//line = strings.TrimRight(line, reader.readerOptions.IRS)

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

		kv := lib.RegexSplitString(reader.readerOptions.IPSRegex, line, 2)
		if len(kv) < 1 {
			return nil, errors.New("mlr: internal coding error in XTAB reader")
		}

		key := kv[0]
		if len(kv) == 1 {
			value := types.MLRVAL_VOID
			record.PutReference(key, value)
		} else {
			value := types.MlrvalPointerFromInferredTypeForDataFiles(kv[1])
			record.PutReference(key, value)
		}
	}

	return record, nil
}
