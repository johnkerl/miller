package output

import (
	"bufio"
	"fmt"
	"os"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/types"
)

// ----------------------------------------------------------------
type RecordWriterJSON struct {
	// Parameters:
	writerOptions  *cli.TWriterOptions
	jsonFormatting types.TJSONFormatting

	// State:
	onFirst bool
}

// ----------------------------------------------------------------
func NewRecordWriterJSON(writerOptions *cli.TWriterOptions) (*RecordWriterJSON, error) {
	var jsonFormatting types.TJSONFormatting = types.JSON_SINGLE_LINE
	if writerOptions.JSONOutputMultiline {
		jsonFormatting = types.JSON_MULTILINE
	}
	return &RecordWriterJSON{
		writerOptions:  writerOptions,
		jsonFormatting: jsonFormatting,
		onFirst:        true,
	}, nil
}

// ----------------------------------------------------------------
func (writer *RecordWriterJSON) Write(
	outrec *types.Mlrmap,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) {
	if writer.writerOptions.WrapJSONOutputInOuterList {
		writer.writeWithListWrap(outrec, bufferedOutputStream, outputIsStdout)
	} else {
		writer.writeWithoutListWrap(outrec, bufferedOutputStream, outputIsStdout)
	}
}

// ----------------------------------------------------------------
func (writer *RecordWriterJSON) writeWithListWrap(
	outrec *types.Mlrmap,
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
	outrec *types.Mlrmap,
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
