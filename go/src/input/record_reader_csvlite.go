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

	"miller/src/cliutil"
	"miller/src/lib"
	"miller/src/types"
)

// ----------------------------------------------------------------
type RecordReaderCSVLite struct {
	readerOptions     *cliutil.TReaderOptions
	emptyStringMlrval types.Mlrval

	// TODO: from reader options
	allowRepeatIFS bool

	// TODO: port from C
	//	comment_handling_t comment_handling;
	//	char*  comment_string;
}

// ----------------------------------------------------------------
func NewRecordReaderCSVLite(readerOptions *cliutil.TReaderOptions) *RecordReaderCSVLite {
	return &RecordReaderCSVLite{
		readerOptions:  readerOptions,
		allowRepeatIFS: false,

		// TODO: port from C
		//	pstate->comment_handling        = comment_handling;
		//	pstate->comment_string          = comment_string;
		emptyStringMlrval: types.MlrvalFromString(""),
	}
}

// ----------------------------------------------------------------
func NewRecordReaderPPRINT(readerOptions *cliutil.TReaderOptions) *RecordReaderCSVLite {
	return &RecordReaderCSVLite{
		readerOptions:  readerOptions,
		allowRepeatIFS: true,

		// TODO: port from C
		//	pstate->comment_handling        = comment_handling;
		//	pstate->comment_string          = comment_string;
		emptyStringMlrval: types.MlrvalFromString(""),
	}
}

// ----------------------------------------------------------------
func (this *RecordReaderCSVLite) Read(
	filenames []string,
	context types.Context,
	inputChannel chan<- *types.RecordAndContext,
	errorChannel chan error,
) {
	if filenames != nil { // nil for mlr -n
		if len(filenames) == 0 { // read from stdin
			handle, err := lib.OpenStdin(
				this.readerOptions.Prepipe,
				this.readerOptions.FileInputEncoding,
			)
			if err != nil {
				errorChannel <- err
			}
			if this.readerOptions.UseImplicitCSVHeader {
				this.processHandleImplicitCSVHeader(
					handle,
					"(stdin)",
					&context,
					inputChannel,
					errorChannel,
				)
			} else {
				this.processHandleExplicitCSVHeader(
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
					this.readerOptions.Prepipe,
					this.readerOptions.FileInputEncoding,
				)
				if err != nil {
					errorChannel <- err
				} else {
					if this.readerOptions.UseImplicitCSVHeader {
						this.processHandleImplicitCSVHeader(
							handle,
							filename,
							&context,
							inputChannel,
							errorChannel,
						)
					} else {
						this.processHandleExplicitCSVHeader(
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
func (this *RecordReaderCSVLite) processHandleExplicitCSVHeader(
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
		line, err := lineReader.ReadString(this.readerOptions.IRS[0]) // xxx temp
		if lib.IsEOF(err) {
			err = nil
			eof = true
		} else if err != nil {
			errorChannel <- err
		} else {
			inputLineNumber++
			// This is how to do a chomp:
			line = strings.TrimRight(line, this.readerOptions.IRS)

			// xxx temp pending autodetect, and pending more windows-port work
			line = strings.TrimRight(line, "\r")

			if line == "" {
				// Reset to new schema
				headerStrings = nil
				continue
			}

			fields := lib.SplitString(line, this.readerOptions.IFS)
			if this.allowRepeatIFS {
				fields = this.stripEmpties(fields)
			}
			if headerStrings == nil {
				headerStrings = fields
				// Get data lines on subsequent loop iterations
			} else {
				if !this.readerOptions.AllowRaggedCSVInput && len(headerStrings) != len(fields) {
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
				if !this.readerOptions.AllowRaggedCSVInput {
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
							record.PutCopy(headerStrings[i], &this.emptyStringMlrval)
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
func (this *RecordReaderCSVLite) processHandleImplicitCSVHeader(
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
		line, err := lineReader.ReadString(this.readerOptions.IRS[0]) // xxx temp
		if lib.IsEOF(err) {
			err = nil
			eof = true
		} else if err != nil {
			errorChannel <- err
		} else {
			inputLineNumber++
			// This is how to do a chomp:
			line = strings.TrimRight(line, this.readerOptions.IRS)

			// xxx temp pending autodetect, and pending more windows-port work
			line = strings.TrimRight(line, "\r")

			if line == "" {
				// Reset to new schema
				headerStrings = nil
				continue
			}

			fields := lib.SplitString(line, this.readerOptions.IFS)
			if this.allowRepeatIFS {
				fields = this.stripEmpties(fields)
			}
			if headerStrings == nil {
				n := len(fields)
				headerStrings = make([]string, n)
				for i := 0; i < n; i++ {
					headerStrings[i] = strconv.Itoa(i + 1)
				}
			} else {
				if !this.readerOptions.AllowRaggedCSVInput && len(headerStrings) != len(fields) {
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
			if !this.readerOptions.AllowRaggedCSVInput {
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
						record.PutCopy(headerStrings[i], &this.emptyStringMlrval)
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
func (this *RecordReaderCSVLite) stripEmpties(input []string) []string {
	output := make([]string, 0, len(input))
	for _, e := range input {
		if e != "" {
			output = append(output, e)
		}
	}
	return output
}
