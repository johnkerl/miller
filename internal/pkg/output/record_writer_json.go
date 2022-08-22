package output

import (
	"bufio"
	"fmt"
	"os"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
)

// ----------------------------------------------------------------
type RecordWriterJSON struct {
	// Parameters:
	writerOptions  *cli.TWriterOptions
	jsonFormatting mlrval.TJSONFormatting
	jvQuoteAll     bool

	// State:
	onFirst bool
}

// ----------------------------------------------------------------
func NewRecordWriterJSON(writerOptions *cli.TWriterOptions) (*RecordWriterJSON, error) {
	var jsonFormatting mlrval.TJSONFormatting = mlrval.JSON_SINGLE_LINE
	if writerOptions.JSONOutputMultiline {
		jsonFormatting = mlrval.JSON_MULTILINE
	}
	return &RecordWriterJSON{
		writerOptions:  writerOptions,
		jsonFormatting: jsonFormatting,
		jvQuoteAll:     writerOptions.JVQuoteAll,
		onFirst:        true,
	}, nil
}

// ----------------------------------------------------------------
func (writer *RecordWriterJSON) Write(
	outrec *mlrval.Mlrmap,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) {
	if outrec != nil && writer.jvQuoteAll {
		outrec.StringifyValuesRecursively()
	}

	if writer.writerOptions.WrapJSONOutputInOuterList {
		writer.writeWithListWrap(outrec, bufferedOutputStream, outputIsStdout)
	} else {
		writer.writeWithoutListWrap(outrec, bufferedOutputStream, outputIsStdout)
	}
}

// ----------------------------------------------------------------
func (writer *RecordWriterJSON) writeWithListWrap(
	outrec *mlrval.Mlrmap,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) {
	if outrec != nil { // Not end of record stream
		if writer.onFirst {
			bufferedOutputStream.WriteString("[\n")
		}

		// The Mlrmap MarshalJSON doesn't include the final newline, so that we
		// can place it neatly with commas here (if the user requested them).
		s, err := outrec.MarshalJSON(writer.jsonFormatting, outputIsStdout)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if !writer.onFirst {
			bufferedOutputStream.WriteString(",\n")
		}

		bufferedOutputStream.WriteString(s)

		writer.onFirst = false

	} else { // End of record stream
		if writer.onFirst { // zero records in the entire output stream
			bufferedOutputStream.WriteString("[")
		}
		bufferedOutputStream.WriteString("\n]\n")
	}
}

// ----------------------------------------------------------------
func (writer *RecordWriterJSON) writeWithoutListWrap(
	outrec *mlrval.Mlrmap,
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
