package input

import (
	"container/list"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/lib"
	"github.com/johnkerl/miller/pkg/mlrval"
	"github.com/johnkerl/miller/pkg/types"
)

func NewRecordReaderPPRINT(
	readerOptions *cli.TReaderOptions,
	recordsPerBatch int64,
) (IRecordReader, error) {
	if readerOptions.BarredPprintInput {
		// Implemented in this file

		readerOptions.IFS = "|"
		readerOptions.AllowRepeatIFS = false

		reader := &RecordReaderPprintBarredOrMarkdown{
			readerOptions:    readerOptions,
			recordsPerBatch:  recordsPerBatch,
			separatorMatcher: regexp.MustCompile(`^\+[-+]*\+$`),
			fieldSplitter:    newFieldSplitter(readerOptions),
		}
		if reader.readerOptions.UseImplicitHeader {
			reader.recordBatchGetter = getRecordBatchImplicitPprintHeader
		} else {
			reader.recordBatchGetter = getRecordBatchExplicitPprintHeader
		}
		return reader, nil

	} else {
		// Use the CSVLite record-reader, which is implemented in another file,
		// with multiple spaces instead of commas
		reader := &RecordReaderCSVLite{
			readerOptions:   readerOptions,
			recordsPerBatch: recordsPerBatch,
			fieldSplitter:   newFieldSplitter(readerOptions),

			useVoidRep: true,
			voidRep:    "-",
		}
		if reader.readerOptions.UseImplicitHeader {
			reader.recordBatchGetter = getRecordBatchImplicitCSVHeader
		} else {
			reader.recordBatchGetter = getRecordBatchExplicitCSVHeader
		}
		return reader, nil
	}
}

type RecordReaderPprintBarredOrMarkdown struct {
	readerOptions   *cli.TReaderOptions
	recordsPerBatch int64 // distinct from readerOptions.RecordsPerBatch for join/repl

	separatorMatcher  *regexp.Regexp
	fieldSplitter     iFieldSplitter
	recordBatchGetter recordBatchGetterPprint

	inputLineNumber int64
	headerStrings   []string
}

// recordBatchGetterPprint points to either an explicit-PPRINT-header or
// implicit-PPRINT-header record-batch getter.
type recordBatchGetterPprint func(
	reader *RecordReaderPprintBarredOrMarkdown,
	linesChannel <-chan *list.List,
	filename string,
	context *types.Context,
	errorChannel chan error,
) (
	recordsAndContexts *list.List,
	eof bool,
)

func (reader *RecordReaderPprintBarredOrMarkdown) Read(
	filenames []string,
	context types.Context,
	readerChannel chan<- *list.List, // list of *types.RecordAndContext
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
				return
			}
			reader.processHandle(
				handle,
				"(stdin)",
				&context,
				readerChannel,
				errorChannel,
				downstreamDoneChannel,
			)
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
					return
				}
				reader.processHandle(
					handle,
					filename,
					&context,
					readerChannel,
					errorChannel,
					downstreamDoneChannel,
				)
				handle.Close()
			}
		}
	}
	readerChannel <- types.NewEndOfStreamMarkerList(&context)
}

func (reader *RecordReaderPprintBarredOrMarkdown) processHandle(
	handle io.Reader,
	filename string,
	context *types.Context,
	readerChannel chan<- *list.List, // list of *types.RecordAndContext
	errorChannel chan error,
	downstreamDoneChannel <-chan bool, // for mlr head
) {
	context.UpdateForStartOfFile(filename)
	reader.inputLineNumber = 0
	reader.headerStrings = nil

	recordsPerBatch := reader.recordsPerBatch
	lineReader := NewLineReader(handle, reader.readerOptions.IRS)
	linesChannel := make(chan *list.List, recordsPerBatch)
	go channelizedLineReader(lineReader, linesChannel, downstreamDoneChannel, recordsPerBatch)

	for {
		recordsAndContexts, eof := reader.recordBatchGetter(reader, linesChannel, filename, context, errorChannel)
		if recordsAndContexts.Len() > 0 {
			readerChannel <- recordsAndContexts
		}
		if eof {
			break
		}
	}
}

func getRecordBatchExplicitPprintHeader(
	reader *RecordReaderPprintBarredOrMarkdown,
	linesChannel <-chan *list.List,
	filename string,
	context *types.Context,
	errorChannel chan error,
) (
	recordsAndContexts *list.List,
	eof bool,
) {
	recordsAndContexts = list.New()
	dedupeFieldNames := reader.readerOptions.DedupeFieldNames

	lines, more := <-linesChannel
	if !more {
		return recordsAndContexts, true
	}

	for e := lines.Front(); e != nil; e = e.Next() {
		line := e.Value.(string)

		reader.inputLineNumber++

		// Check for comments-in-data feature
		// TODO: function-pointer this away
		if reader.readerOptions.CommentHandling != cli.CommentsAreData {
			if strings.HasPrefix(line, reader.readerOptions.CommentString) {
				if reader.readerOptions.CommentHandling == cli.PassComments {
					recordsAndContexts.PushBack(types.NewOutputString(line+"\n", context))
					continue
				} else if reader.readerOptions.CommentHandling == cli.SkipComments {
					continue
				}
				// else comments are data
			}
		}

		if line == "" {
			// Reset to new schema
			reader.headerStrings = nil
			continue
		}

		// Example input:
		// +-----+-----+----+---------------------+---------------------+
		// | a   | b   | i  | x                   | y                   |
		// +-----+-----+----+---------------------+---------------------+
		// | pan | pan | 1  | 0.3467901443380824  | 0.7268028627434533  |
		// | eks | pan | 2  | 0.7586799647899636  | 0.5221511083334797  |
		// +-----+-----+----+---------------------+---------------------+

		// Skip lines like
		// +-----+-----+----+---------------------+---------------------+
		if reader.separatorMatcher.MatchString(line) {
			continue
		}

		// Skip the leading and trailing pipes
		paddedFields := reader.fieldSplitter.Split(line)
		npad := len(paddedFields)
		if npad < 2 {
			continue
		}
		fields := make([]string, npad-2)
		for i, _ := range paddedFields {
			if i == 0 || i == npad-1 {
				continue
			}
			fields[i-1] = strings.TrimSpace(paddedFields[i])
		}

		if reader.headerStrings == nil {
			reader.headerStrings = fields
			// Get data lines on subsequent loop iterations
		} else {
			if !reader.readerOptions.AllowRaggedCSVInput && len(reader.headerStrings) != len(fields) {
				err := fmt.Errorf(
					"mlr: PPRINT-barred header/data length mismatch %d != %d "+
						"at filename %s line  %d.\n",
					len(reader.headerStrings), len(fields), filename, reader.inputLineNumber,
				)
				errorChannel <- err
				return
			}

			record := mlrval.NewMlrmapAsRecord()
			if !reader.readerOptions.AllowRaggedCSVInput {
				for i, field := range fields {
					value := mlrval.FromDeferredType(field)
					_, err := record.PutReferenceMaybeDedupe(reader.headerStrings[i], value, dedupeFieldNames)
					if err != nil {
						errorChannel <- err
						return
					}
				}
			} else {
				nh := int64(len(reader.headerStrings))
				nd := int64(len(fields))
				n := lib.IntMin2(nh, nd)
				var i int64
				for i = 0; i < n; i++ {
					field := fields[i]
					value := mlrval.FromDeferredType(field)
					_, err := record.PutReferenceMaybeDedupe(reader.headerStrings[i], value, dedupeFieldNames)
					if err != nil {
						errorChannel <- err
						return
					}
				}
				if nh < nd {
					// if header shorter than data: use 1-up itoa keys
					for i = nh; i < nd; i++ {
						key := strconv.FormatInt(i+1, 10)
						value := mlrval.FromDeferredType(fields[i])
						_, err := record.PutReferenceMaybeDedupe(key, value, dedupeFieldNames)
						if err != nil {
							errorChannel <- err
							return
						}
					}
				}
				if nh > nd {
					// if header longer than data: use "" values
					for i = nd; i < nh; i++ {
						record.PutCopy(reader.headerStrings[i], mlrval.VOID)
					}
				}
			}

			context.UpdateForInputRecord()
			recordsAndContexts.PushBack(types.NewRecordAndContext(record, context))

		}
	}

	return recordsAndContexts, false
}

func getRecordBatchImplicitPprintHeader(
	reader *RecordReaderPprintBarredOrMarkdown,
	linesChannel <-chan *list.List,
	filename string,
	context *types.Context,
	errorChannel chan error,
) (
	recordsAndContexts *list.List,
	eof bool,
) {
	recordsAndContexts = list.New()
	dedupeFieldNames := reader.readerOptions.DedupeFieldNames

	lines, more := <-linesChannel
	if !more {
		return recordsAndContexts, true
	}

	for e := lines.Front(); e != nil; e = e.Next() {
		line := e.Value.(string)

		reader.inputLineNumber++

		// Check for comments-in-data feature
		// TODO: function-pointer this away
		if reader.readerOptions.CommentHandling != cli.CommentsAreData {
			if strings.HasPrefix(line, reader.readerOptions.CommentString) {
				if reader.readerOptions.CommentHandling == cli.PassComments {
					recordsAndContexts.PushBack(types.NewOutputString(line+"\n", context))
					continue
				} else if reader.readerOptions.CommentHandling == cli.SkipComments {
					continue
				}
				// else comments are data
			}
		}

		if line == "" {
			// Reset to new schema
			reader.headerStrings = nil
			continue
		}

		// Example input:
		// +-----+-----+----+---------------------+---------------------+
		// | a   | b   | i  | x                   | y                   |
		// +-----+-----+----+---------------------+---------------------+
		// | pan | pan | 1  | 0.3467901443380824  | 0.7268028627434533  |
		// | eks | pan | 2  | 0.7586799647899636  | 0.5221511083334797  |
		// +-----+-----+----+---------------------+---------------------+

		// Skip lines like
		// +-----+-----+----+---------------------+---------------------+
		if reader.separatorMatcher.MatchString(line) {
			continue
		}

		// Skip the leading and trailing pipes
		paddedFields := reader.fieldSplitter.Split(line)
		npad := len(paddedFields)
		fields := make([]string, npad-2)
		for i, _ := range paddedFields {
			if i == 0 || i == npad-1 {
				continue
			}
			fields[i-1] = strings.TrimSpace(paddedFields[i])
		}

		if reader.headerStrings == nil {
			n := len(fields)
			reader.headerStrings = make([]string, n)
			for i := 0; i < n; i++ {
				reader.headerStrings[i] = strconv.Itoa(i + 1)
			}
		} else {
			if !reader.readerOptions.AllowRaggedCSVInput && len(reader.headerStrings) != len(fields) {
				err := fmt.Errorf(
					"mlr: CSV header/data length mismatch %d != %d "+
						"at filename %s line  %d.\n",
					len(reader.headerStrings), len(fields), filename, reader.inputLineNumber,
				)
				errorChannel <- err
				return
			}
		}

		record := mlrval.NewMlrmapAsRecord()
		if !reader.readerOptions.AllowRaggedCSVInput {
			for i, field := range fields {
				value := mlrval.FromDeferredType(field)
				_, err := record.PutReferenceMaybeDedupe(reader.headerStrings[i], value, dedupeFieldNames)
				if err != nil {
					errorChannel <- err
					return
				}
			}
		} else {
			nh := int64(len(reader.headerStrings))
			nd := int64(len(fields))
			n := lib.IntMin2(nh, nd)
			var i int64
			for i = 0; i < n; i++ {
				field := fields[i]
				value := mlrval.FromDeferredType(field)
				_, err := record.PutReferenceMaybeDedupe(reader.headerStrings[i], value, dedupeFieldNames)
				if err != nil {
					errorChannel <- err
					return
				}
			}
			if nh < nd {
				// if header shorter than data: use 1-up itoa keys
				key := strconv.FormatInt(i+1, 10)
				value := mlrval.FromDeferredType(fields[i])
				_, err := record.PutReferenceMaybeDedupe(key, value, dedupeFieldNames)
				if err != nil {
					errorChannel <- err
					return
				}
			}
			if nh > nd {
				// if header longer than data: use "" values
				for i = nd; i < nh; i++ {
					_, err := record.PutReferenceMaybeDedupe(
						reader.headerStrings[i],
						mlrval.VOID.Copy(),
						dedupeFieldNames,
					)
					if err != nil {
						errorChannel <- err
						return
					}
				}
			}
		}

		context.UpdateForInputRecord()
		recordsAndContexts.PushBack(types.NewRecordAndContext(record, context))
	}

	return recordsAndContexts, false
}
