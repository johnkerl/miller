package output

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/types"
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
	ostream io.WriteCloser,
	outputIsStdout bool,
) {
	if writer.writerOptions.WrapJSONOutputInOuterList {
		writer.writeWithListWrap(outrec, ostream, outputIsStdout)
	} else {
		writer.writeWithoutListWrap(outrec, ostream, outputIsStdout)
	}
}

// ----------------------------------------------------------------
func (writer *RecordWriterJSON) writeWithListWrap(
	outrec *types.Mlrmap,
	ostream io.WriteCloser,
	outputIsStdout bool,
) {
	var buffer bytes.Buffer // stdio is non-buffered in Go, so buffer for ~5x speed increase

	if outrec != nil { // Not end of record stream
		if writer.onFirst {
			buffer.WriteString("[\n")
		}

		// The Mlrmap MarshalJSON doesn't include the final newline, so that we
		// can place it neatly with commas here (if the user requested them).
		s, err := outrec.MarshalJSON(writer.jsonFormatting, outputIsStdout)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if !writer.onFirst {
			buffer.WriteString(",\n")
		}

		buffer.WriteString(s)

		writer.onFirst = false

	} else { // End of record stream
		if writer.onFirst { // zero records in the entire output stream
			buffer.Write([]byte("["))
		}
		buffer.Write([]byte("\n]\n"))
	}
	ostream.Write(buffer.Bytes())
}

// ----------------------------------------------------------------
func (writer *RecordWriterJSON) writeWithoutListWrap(
	outrec *types.Mlrmap,
	ostream io.WriteCloser,
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

	ostream.Write([]byte(s))
	ostream.Write([]byte("\n"))
}
