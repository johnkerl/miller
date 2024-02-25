package input

// Multi-file cases:
//
// a,a        a,b        c          d
// -- FILE1:  -- FILE1:  -- FILE1:  -- FILE1:
// a,b,c      a,b,c      a,b,c      a,b,c
// 1,2,3      1,2,3      1,2,3      1,2,3
// 4,5,6      4,5,6      4,5,6      4,5,6
// -- FILE2:  -- FILE2:
// a,b,c      d,e,f,g    a,b,c      d,e,f
// 7,8,9      3,4,5,6    7,8,9      3,4,5
// --OUTPUT:  --OUTPUT:  --OUTPUT:  --OUTPUT:
// a,b,c      a,b,c      a,b,c      a,b,c
// 1,2,3      1,2,3      1,2,3      1,2,3
// 4,5,6      4,5,6      4,5,6      4,5,6
// 7,8,9                 7,8,9
//            d,e,f,g               d,e,f
//            3,4,5,6               3,4,5

import (
	"container/list"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/lib"
	"github.com/johnkerl/miller/pkg/mlrval"
	"github.com/johnkerl/miller/pkg/types"
)

// recordBatchGetterCSV points to either an explicit-CSV-header or
// implicit-CSV-header record-batch getter.
type recordBatchGetterCSV func(
	reader *RecordReaderCSVLite,
	linesChannel <-chan *list.List,
	filename string,
	context *types.Context,
	errorChannel chan error,
) (
	recordsAndContexts *list.List,
	eof bool,
)

type RecordReaderCSVLite struct {
	readerOptions   *cli.TReaderOptions
	recordsPerBatch int64 // distinct from readerOptions.RecordsPerBatch for join/repl

	fieldSplitter     iFieldSplitter
	recordBatchGetter recordBatchGetterCSV

	inputLineNumber int64
	headerStrings   []string

	useVoidRep bool
	voidRep    string // For pprint output, empty strings are mapped to "-"; this is for reading them back in
}

func NewRecordReaderCSVLite(
	readerOptions *cli.TReaderOptions,
	recordsPerBatch int64,
) (*RecordReaderCSVLite, error) {
	reader := &RecordReaderCSVLite{
		readerOptions:   readerOptions,
		recordsPerBatch: recordsPerBatch,
		fieldSplitter:   newFieldSplitter(readerOptions),

		useVoidRep: false,
		voidRep:    "",
	}
	if reader.readerOptions.UseImplicitHeader {
		reader.recordBatchGetter = getRecordBatchImplicitCSVHeader
	} else {
		reader.recordBatchGetter = getRecordBatchExplicitCSVHeader
	}
	return reader, nil
}

func (reader *RecordReaderCSVLite) Read(
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

func (reader *RecordReaderCSVLite) processHandle(
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

func getRecordBatchExplicitCSVHeader(
	reader *RecordReaderCSVLite,
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

		// Strip CSV BOM
		if reader.inputLineNumber == 1 {
			if strings.HasPrefix(line, CSV_BOM) {
				line = strings.Replace(line, CSV_BOM, "", 1)
			}
		}

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

		fields := reader.fieldSplitter.Split(line)

		if reader.headerStrings == nil {
			reader.headerStrings = fields
			// Get data lines on subsequent loop iterations
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

			record := mlrval.NewMlrmapAsRecord()
			if !reader.readerOptions.AllowRaggedCSVInput {
				for i, field := range fields {
					if reader.useVoidRep && field == reader.voidRep {
						field = ""
					}
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
					if reader.useVoidRep && field == reader.voidRep {
						field = ""
					}
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

func getRecordBatchImplicitCSVHeader(
	reader *RecordReaderCSVLite,
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

		// This is how to do a chomp:
		line = strings.TrimRight(line, reader.readerOptions.IRS)

		line = strings.TrimRight(line, "\r")

		if line == "" {
			// Reset to new schema
			reader.headerStrings = nil
			continue
		}

		fields := reader.fieldSplitter.Split(line)

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
				if reader.useVoidRep && field == reader.voidRep {
					field = ""
				}
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
				if reader.useVoidRep && field == reader.voidRep {
					field = ""
				}
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
					_, err := record.PutReferenceMaybeDedupe(reader.headerStrings[i], mlrval.VOID.Copy(), dedupeFieldNames)
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
