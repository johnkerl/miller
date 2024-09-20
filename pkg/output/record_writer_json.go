package output

import (
	"bufio"
	"fmt"
	"os"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

// ----------------------------------------------------------------
type RecordWriterJSON struct {
	// Parameters:
	writerOptions  *cli.TWriterOptions
	jsonFormatting mlrval.TJSONFormatting
	jvQuoteAll     bool

	// State:
	wroteAnyRecords bool
}

// ----------------------------------------------------------------
func NewRecordWriterJSON(writerOptions *cli.TWriterOptions) (*RecordWriterJSON, error) {
	var jsonFormatting mlrval.TJSONFormatting = mlrval.JSON_SINGLE_LINE
	if writerOptions.JSONOutputMultiline {
		jsonFormatting = mlrval.JSON_MULTILINE
	}
	return &RecordWriterJSON{
		writerOptions:   writerOptions,
		jsonFormatting:  jsonFormatting,
		jvQuoteAll:      writerOptions.JVQuoteAll,
		wroteAnyRecords: false,
	}, nil
}

// ----------------------------------------------------------------
func (writer *RecordWriterJSON) Write(
	outrec *mlrval.Mlrmap,
	context *types.Context,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) error {
	if outrec != nil && writer.jvQuoteAll {
		outrec.StringifyValuesRecursively()
	}

	if writer.writerOptions.WrapJSONOutputInOuterList {
		writer.writeWithListWrap(outrec, context, bufferedOutputStream, outputIsStdout)
	} else {
		writer.writeWithoutListWrap(outrec, context, bufferedOutputStream, outputIsStdout)
	}
	return nil
}

// ----------------------------------------------------------------
func (writer *RecordWriterJSON) writeWithListWrap(
	outrec *mlrval.Mlrmap,
	context *types.Context,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) {
	if outrec != nil { // Not end of record stream
		if !writer.wroteAnyRecords {
			bufferedOutputStream.WriteString("[\n")
		}

		// The Mlrmap MarshalJSON doesn't include the final newline, so that we
		// can place it neatly with commas here (if the user requested them).
		s, err := outrec.MarshalJSON(writer.jsonFormatting, outputIsStdout)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if writer.wroteAnyRecords {
			bufferedOutputStream.WriteString(",\n")
		}

		bufferedOutputStream.WriteString(s)

		writer.wroteAnyRecords = true

	} else { // End of record stream

		if !writer.wroteAnyRecords {
			if context.JSONHadBrackets {
				bufferedOutputStream.WriteString("[")
				bufferedOutputStream.WriteString("\n]\n")
			}
		} else {
			bufferedOutputStream.WriteString("\n]\n")
		}

	}
}

// ----------------------------------------------------------------
func (writer *RecordWriterJSON) writeWithoutListWrap(
	outrec *mlrval.Mlrmap,
	_ *types.Context,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) {
	if outrec == nil {
		// End of record stream
		return
	}

	// The Mlrmap MarshalJSON doesn't include the final newline, so that we
	// can place it neatly with commas here (if the user requested them).
	s, err := outrec.MarshalJSON(writer.jsonFormatting, outputIsStdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bufferedOutputStream.WriteString(s)
	bufferedOutputStream.WriteString("\n")
}
