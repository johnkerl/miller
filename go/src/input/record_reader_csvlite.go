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
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"mlr/src/cliutil"
	"mlr/src/lib"
	"mlr/src/types"
)

// ----------------------------------------------------------------
type RecordReaderCSVLite struct {
	readerOptions     *cliutil.TReaderOptions
	emptyStringMlrval types.Mlrval
}

// ----------------------------------------------------------------
func NewRecordReaderCSVLite(readerOptions *cliutil.TReaderOptions) *RecordReaderCSVLite {
	return &RecordReaderCSVLite{
		readerOptions:     readerOptions,
		emptyStringMlrval: types.MlrvalFromString(""),
	}
}

// ----------------------------------------------------------------
func NewRecordReaderPPRINT(readerOptions *cliutil.TReaderOptions) *RecordReaderCSVLite {
	return &RecordReaderCSVLite{
		readerOptions:     readerOptions,
		emptyStringMlrval: types.MlrvalFromString(""),
	}
}

// ----------------------------------------------------------------
func (reader *RecordReaderCSVLite) Read(
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
			if reader.readerOptions.UseImplicitCSVHeader {
				reader.processHandleImplicitCSVHeader(
					handle,
					"(stdin)",
					&context,
					inputChannel,
					errorChannel,
				)
			} else {
				reader.processHandleExplicitCSVHeader(
					handle,
					"(stdin)",
					&context,
					inputChannel,
					errorChannel,
				)
			}
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
					if reader.readerOptions.UseImplicitCSVHeader {
						reader.processHandleImplicitCSVHeader(
							handle,
							filename,
							&context,
							inputChannel,
							errorChannel,
						)
					} else {
						reader.processHandleExplicitCSVHeader(
							handle,
							filename,
							&context,
							inputChannel,
							errorChannel,
						)
					}
					handle.Close()
				}
			}
		}
	}
	inputChannel <- types.NewEndOfStreamMarker(&context)
}

// ----------------------------------------------------------------
func (reader *RecordReaderCSVLite) processHandleExplicitCSVHeader(
	handle io.Reader,
	filename string,
	context *types.Context,
	inputChannel chan<- *types.RecordAndContext,
	errorChannel chan error,
) {
	var inputLineNumber int = 0
	var headerStrings []string = nil

	context.UpdateForStartOfFile(filename)

	lineReader := bufio.NewReader(handle)
	eof := false
	for !eof {
		line, err := lineReader.ReadString(reader.readerOptions.IRS[0]) // xxx temp
		if lib.IsEOF(err) {
			err = nil
			eof = true
		} else if err != nil {
			errorChannel <- err
		} else {
			inputLineNumber++

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

			// This is how to do a chomp:
			line = strings.TrimRight(line, reader.readerOptions.IRS)
			// xxx temp pending autodetect, and pending more windows-port work
			line = strings.TrimRight(line, "\r")

			if line == "" {
				// Reset to new schema
				headerStrings = nil
				continue
			}

			fields := lib.SplitString(line, reader.readerOptions.IFS)
			if reader.readerOptions.AllowRepeatIFS {
				fields = reader.stripEmpties(fields)
			}
			if headerStrings == nil {
				headerStrings = fields
				// Get data lines on subsequent loop iterations
			} else {
				if !reader.readerOptions.AllowRaggedCSVInput && len(headerStrings) != len(fields) {
					err := errors.New(
						fmt.Sprintf(
							"Miller: CSV header/data length mismatch %d != %d "+
								"at filename %s line  %d.\n",
							len(headerStrings), len(fields), filename, inputLineNumber,
						),
					)
					errorChannel <- err
					return
				}

				record := types.NewMlrmap()
				if !reader.readerOptions.AllowRaggedCSVInput {
					for i, field := range fields {
						value := types.MlrvalPointerFromInferredType(field)
						record.PutCopy(headerStrings[i], value)
					}
				} else {
					nh := len(headerStrings)
					nd := len(fields)
					n := lib.IntMin2(nh, nd)
					var i int
					for i = 0; i < n; i++ {
						value := types.MlrvalPointerFromInferredType(fields[i])
						record.PutCopy(headerStrings[i], value)
					}
					if nh < nd {
						// if header shorter than data: use 1-up itoa keys
						for i = nh; i < nd; i++ {
							key := strconv.Itoa(i + 1)
							value := types.MlrvalPointerFromInferredType(fields[i])
							record.PutCopy(key, value)
						}
					}
					if nh > nd {
						// if header longer than data: use "" values
						for i = nd; i < nh; i++ {
							record.PutCopy(headerStrings[i], &reader.emptyStringMlrval)
						}
					}
				}

				context.UpdateForInputRecord()
				inputChannel <- types.NewRecordAndContext(
					record,
					context,
				)
			}
		}
	}
}

// ----------------------------------------------------------------
func (reader *RecordReaderCSVLite) processHandleImplicitCSVHeader(
	handle io.Reader,
	filename string,
	context *types.Context,
	inputChannel chan<- *types.RecordAndContext,
	errorChannel chan error,
) {
	var inputLineNumber int = 0
	var headerStrings []string = nil

	context.UpdateForStartOfFile(filename)

	lineReader := bufio.NewReader(handle)
	eof := false
	for !eof {
		line, err := lineReader.ReadString(reader.readerOptions.IRS[0]) // xxx temp
		if lib.IsEOF(err) {
			err = nil
			eof = true
		} else if err != nil {
			errorChannel <- err
		} else {
			inputLineNumber++

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

			// This is how to do a chomp:
			line = strings.TrimRight(line, reader.readerOptions.IRS)

			// xxx temp pending autodetect, and pending more windows-port work
			line = strings.TrimRight(line, "\r")

			if line == "" {
				// Reset to new schema
				headerStrings = nil
				continue
			}

			fields := lib.SplitString(line, reader.readerOptions.IFS)
			if reader.readerOptions.AllowRepeatIFS {
				fields = reader.stripEmpties(fields)
			}
			if headerStrings == nil {
				n := len(fields)
				headerStrings = make([]string, n)
				for i := 0; i < n; i++ {
					headerStrings[i] = strconv.Itoa(i + 1)
				}
			} else {
				if !reader.readerOptions.AllowRaggedCSVInput && len(headerStrings) != len(fields) {
					err := errors.New(
						fmt.Sprintf(
							"Miller: CSV header/data length mismatch %d != %d "+
								"at filename %s line  %d.\n",
							len(headerStrings), len(fields), filename, inputLineNumber,
						),
					)
					errorChannel <- err
					return
				}
			}

			record := types.NewMlrmap()
			if !reader.readerOptions.AllowRaggedCSVInput {
				for i, field := range fields {
					value := types.MlrvalPointerFromInferredType(field)
					record.PutCopy(headerStrings[i], value)
				}
			} else {
				nh := len(headerStrings)
				nd := len(fields)
				n := lib.IntMin2(nh, nd)
				var i int
				for i = 0; i < n; i++ {
					value := types.MlrvalPointerFromInferredType(fields[i])
					record.PutCopy(headerStrings[i], value)
				}
				if nh < nd {
					// if header shorter than data: use 1-up itoa keys
					key := strconv.Itoa(i + 1)
					value := types.MlrvalPointerFromInferredType(fields[i])
					record.PutCopy(key, value)
				}
				if nh > nd {
					// if header longer than data: use "" values
					for i = nd; i < nh; i++ {
						record.PutCopy(headerStrings[i], &reader.emptyStringMlrval)
					}
				}
			}

			context.UpdateForInputRecord()
			inputChannel <- types.NewRecordAndContext(
				record,
				context,
			)

		}
	}
}

// ----------------------------------------------------------------
// For CSV, we have "a,,c" -> ["a", "", "c"]. But for PPRINT, "a  b" -> ["a", "b"].
// One way to do this is split on single spaces, then strip empty-string slots.
func (reader *RecordReaderCSVLite) stripEmpties(input []string) []string {
	output := make([]string, 0, len(input))
	for _, e := range input {
		if e != "" {
			output = append(output, e)
		}
	}
	return output
}
