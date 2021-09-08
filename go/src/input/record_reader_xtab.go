package input

import (
	"bufio"
	"container/list"
	"errors"
	"io"
	"regexp"
	"strings"

	"mlr/src/cliutil"
	"mlr/src/lib"
	"mlr/src/types"
)

type RecordReaderXTAB struct {
	readerOptions *cliutil.TReaderOptions
	ifsRegex      *regexp.Regexp
	// TODO: parameterize IRS

	// TODO: port from C
	// 	int    allow_repeat_ips;
	// 	int    do_auto_line_term;
	// 	int    at_eof;
}

// ----------------------------------------------------------------
func NewRecordReaderXTAB(readerOptions *cliutil.TReaderOptions) *RecordReaderXTAB {
	return &RecordReaderXTAB{
		readerOptions: readerOptions,
		// TODO: incorporate IFS
		ifsRegex: regexp.MustCompile("\\s+"),
	}
}

// ----------------------------------------------------------------
func (reader *RecordReaderXTAB) Read(
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

func (reader *RecordReaderXTAB) processHandle(
	handle io.Reader,
	filename string,
	context *types.Context,
	inputChannel chan<- *types.RecordAndContext,
	errorChannel chan error,
) {
	context.UpdateForStartOfFile(filename)

	lineReader := bufio.NewReader(handle)

	linesForRecord := list.New()

	eof := false
	for !eof {
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

		// TODO -- incorporate IFS
		kv := reader.ifsRegex.Split(line, 2)
		if len(kv) < 1 {
			return nil, errors.New("Miller: internal coding error in XTAB reader")
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
