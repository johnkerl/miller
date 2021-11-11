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
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
type RecordReaderCSVLite struct {
	readerOptions *cli.TReaderOptions
}

// ----------------------------------------------------------------
func NewRecordReaderCSVLite(readerOptions *cli.TReaderOptions) (*RecordReaderCSVLite, error) {
	return &RecordReaderCSVLite{
		readerOptions: readerOptions,
	}, nil
}

// ----------------------------------------------------------------
func NewRecordReaderPPRINT(readerOptions *cli.TReaderOptions) (*RecordReaderCSVLite, error) {
	return &RecordReaderCSVLite{
		readerOptions: readerOptions,
	}, nil
}

// ----------------------------------------------------------------
func (reader *RecordReaderCSVLite) Read(
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
			if reader.readerOptions.UseImplicitCSVHeader {
				reader.processHandleImplicitCSVHeader(
					handle,
					"(stdin)",
					&context,
					inputChannel,
					errorChannel,
					downstreamDoneChannel,
				)
			} else {
				reader.processHandleExplicitCSVHeader(
					handle,
					"(stdin)",
					&context,
					inputChannel,
					errorChannel,
					downstreamDoneChannel,
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
							downstreamDoneChannel,
						)
					} else {
						reader.processHandleExplicitCSVHeader(
							handle,
							filename,
							&context,
							inputChannel,
							errorChannel,
							downstreamDoneChannel,
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
	downstreamDoneChannel <-chan bool, // for mlr head
) {
	var inputLineNumber int = 0
	var headerStrings []string = nil

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

		line := scanner.Text()

		inputLineNumber++

		// Strip CSV BOM
		if inputLineNumber == 1 {
			if strings.HasPrefix(line, CSV_BOM) {
				line = strings.Replace(line, CSV_BOM, "", 1)
			}
		}

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

		if line == "" {
			// Reset to new schema
			headerStrings = nil
			continue
		}

		var fields []string
		if reader.readerOptions.IFSRegex == nil { // e.g. --no-ifs-regex
			fields = lib.SplitString(line, reader.readerOptions.IFS)
		} else {
			fields = lib.RegexSplitString(reader.readerOptions.IFSRegex, line, -1)
		}

		if reader.readerOptions.AllowRepeatIFS {
			fields = lib.StripEmpties(fields) // left/right trim
		}
		if headerStrings == nil {
			headerStrings = fields
			// Get data lines on subsequent loop iterations
		} else {
			if !reader.readerOptions.AllowRaggedCSVInput && len(headerStrings) != len(fields) {
				err := errors.New(
					fmt.Sprintf(
						"mlr: CSV header/data length mismatch %d != %d "+
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
					value := types.MlrvalFromInferredTypeForDataFiles(field)
					record.PutCopy(headerStrings[i], value)
				}
			} else {
				nh := len(headerStrings)
				nd := len(fields)
				n := lib.IntMin2(nh, nd)
				var i int
				for i = 0; i < n; i++ {
					value := types.MlrvalFromInferredTypeForDataFiles(fields[i])
					record.PutCopy(headerStrings[i], value)
				}
				if nh < nd {
					// if header shorter than data: use 1-up itoa keys
					for i = nh; i < nd; i++ {
						key := strconv.Itoa(i + 1)
						value := types.MlrvalFromInferredTypeForDataFiles(fields[i])
						record.PutCopy(key, value)
					}
				}
				if nh > nd {
					// if header longer than data: use "" values
					for i = nd; i < nh; i++ {
						record.PutCopy(headerStrings[i], types.MLRVAL_VOID)
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
func (reader *RecordReaderCSVLite) processHandleImplicitCSVHeader(
	handle io.Reader,
	filename string,
	context *types.Context,
	inputChannel chan<- *types.RecordAndContext,
	errorChannel chan error,
	downstreamDoneChannel <-chan bool, // for mlr head
) {
	var inputLineNumber int = 0
	var headerStrings []string = nil

	context.UpdateForStartOfFile(filename)

	scanner := NewLineScanner(handle, reader.readerOptions.IRS)
	for scanner.Scan() {

		// See if downstream processors will be ignoring further data (e.g. mlr
		// head).  If so, stop reading. This makes 'mlr head hugefile' exit
		// quickly, as it should.

		// TODO: extract a helper function
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

		inputLineNumber++

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

		// This is how to do a chomp:
		line = strings.TrimRight(line, reader.readerOptions.IRS)

		// xxx temp pending autodetect, and pending more windows-port work
		line = strings.TrimRight(line, "\r")

		if line == "" {
			// Reset to new schema
			headerStrings = nil
			continue
		}

		var fields []string
		if reader.readerOptions.IFSRegex == nil { // e.g. --no-ifs-regex
			fields = lib.SplitString(line, reader.readerOptions.IFS)
		} else {
			fields = lib.RegexSplitString(reader.readerOptions.IFSRegex, line, -1)
		}

		if reader.readerOptions.AllowRepeatIFS {
			fields = lib.StripEmpties(fields) // left/right trim
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
						"mlr: CSV header/data length mismatch %d != %d "+
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
				value := types.MlrvalFromInferredTypeForDataFiles(field)
				record.PutCopy(headerStrings[i], value)
			}
		} else {
			nh := len(headerStrings)
			nd := len(fields)
			n := lib.IntMin2(nh, nd)
			var i int
			for i = 0; i < n; i++ {
				value := types.MlrvalFromInferredTypeForDataFiles(fields[i])
				record.PutCopy(headerStrings[i], value)
			}
			if nh < nd {
				// if header shorter than data: use 1-up itoa keys
				key := strconv.Itoa(i + 1)
				value := types.MlrvalFromInferredTypeForDataFiles(fields[i])
				record.PutCopy(key, value)
			}
			if nh > nd {
				// if header longer than data: use "" values
				for i = nd; i < nh; i++ {
					record.PutCopy(headerStrings[i], types.MLRVAL_VOID)
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
