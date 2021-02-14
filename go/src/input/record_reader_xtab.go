package input

import (
	"bufio"
	"container/list"
	"errors"
	"io"
	"os"
	"regexp"
	"strings"

	"miller/src/cliutil"
	"miller/src/types"
)

type RecordReaderXTAB struct {
	ifsRegex *regexp.Regexp
	ifs      string
	ips      string
	irs      string
	// TODO: parameterize IRS

	// TODO: port from C
	// 	int    allow_repeat_ips;
	// 	int    do_auto_line_term;
	// 	int    at_eof;
	// 	comment_handling_t comment_handling;
	// 	char*  comment_string;
}

// ----------------------------------------------------------------
func NewRecordReaderXTAB(readerOptions *cliutil.TReaderOptions) *RecordReaderXTAB {
	return &RecordReaderXTAB{
		// TODO: incorporate IFS
		ifsRegex: regexp.MustCompile("\\s+"),
		ifs:      readerOptions.IFS,
		ips:      readerOptions.IPS,
		irs:      "\n",
	}
}

// ----------------------------------------------------------------
func (this *RecordReaderXTAB) Read(
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
	inputChannel <- types.NewEndOfStreamMarker(&context)
}

func (this *RecordReaderXTAB) processHandle(
	handle *os.File,
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
		line, err := lineReader.ReadString(this.irs[0]) // xxx temp
		if err == io.EOF {
			err = nil
			eof = true

			if linesForRecord.Len() > 0 {
				record, err := this.recordFromXTABLines(linesForRecord)
				if err != nil {
					errorChannel <- err
					return
				}
				context.UpdateForInputRecord()
				inputChannel <- types.NewRecordAndContext(record, context)
				linesForRecord = list.New()
			}
		} else if err != nil {
			errorChannel <- err
		} else {
			// This is how to do a chomp:
			line = strings.TrimRight(line, this.irs)

			if line == "" {
				if linesForRecord.Len() > 0 {
					record, err := this.recordFromXTABLines(linesForRecord)
					if err != nil {
						errorChannel <- err
						return
					}
					context.UpdateForInputRecord()
					inputChannel <- types.NewRecordAndContext(record, context)
					linesForRecord = list.New()
				}
			} else {
				linesForRecord.PushBack(line)
			}
		}
	}
}

// ----------------------------------------------------------------
func (this *RecordReaderXTAB) recordFromXTABLines(
	lines *list.List,
) (*types.Mlrmap, error) {
	record := types.NewMlrmap()

	for entry := lines.Front(); entry != nil; entry = entry.Next() {
		line := entry.Value.(string)

		// TODO -- incorporate IFS
		kv := this.ifsRegex.Split(line, 2)
		if len(kv) < 1 {
			return nil, errors.New("Miller: internal coding error in XTAB reader")
		}

		key := kv[0]
		if len(kv) == 1 {
			value := types.MlrvalFromVoid()
			record.PutReference(key, &value)
		} else {
			value := types.MlrvalFromInferredType(kv[1])
			record.PutReference(key, &value)
		}
	}

	return record, nil
}
