package output

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"miller/src/cliutil"
	"miller/src/types"
)

// ----------------------------------------------------------------
type RecordWriterJSON struct {
	// Parameters:
	writerOptions  *cliutil.TWriterOptions
	jsonFormatting types.TJSONFormatting

	// State:
	onFirst bool
}

// ----------------------------------------------------------------
func NewRecordWriterJSON(writerOptions *cliutil.TWriterOptions) *RecordWriterJSON {
	var jsonFormatting types.TJSONFormatting = types.JSON_SINGLE_LINE
	if writerOptions.JSONOutputMultiline {
		jsonFormatting = types.JSON_MULTILINE
	}
	return &RecordWriterJSON{
		writerOptions:  writerOptions,
		jsonFormatting: jsonFormatting,
		onFirst:        true,
	}
}

// ----------------------------------------------------------------
func (this *RecordWriterJSON) Write(
	outrec *types.Mlrmap,
	ostream io.WriteCloser,
) {
	if this.writerOptions.WrapJSONOutputInOuterList {
		this.writeWithListWrap(outrec, ostream)
	} else {
		this.writeWithoutListWrap(outrec, ostream)
	}
}

// ----------------------------------------------------------------
func (this *RecordWriterJSON) writeWithListWrap(
	outrec *types.Mlrmap,
	ostream io.WriteCloser,
) {
	var buffer bytes.Buffer // 5x faster than fmt.Print() separately

	if outrec != nil { // Not end of record stream
		if this.onFirst {
			buffer.WriteString("[\n")
		}

		// The Mlrmap MarshalJSON doesn't include the final newline, so that we
		// can place it neatly with commas here (if the user requested them).
		bytes, err := outrec.MarshalJSON(this.jsonFormatting)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if !this.onFirst {
			buffer.WriteString(",\n")
		}

		buffer.Write(bytes)

		this.onFirst = false

	} else { // End of record stream
		if this.onFirst { // zero records in the entire output stream
			buffer.Write([]byte("["))
		}
		buffer.Write([]byte("\n]\n"))
	}
	ostream.Write(buffer.Bytes())
}

// ----------------------------------------------------------------
func (this *RecordWriterJSON) writeWithoutListWrap(
	outrec *types.Mlrmap,
	ostream io.WriteCloser,
) {
	if outrec == nil {
		// End of record stream
		return
	}

	// The Mlrmap MarshalJSON doesn't include the final newline, so that we
	// can place it neatly with commas here (if the user requested them).
	bytes, err := outrec.MarshalJSON(this.jsonFormatting)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ostream.Write(bytes)
	ostream.Write([]byte("\n"))
}
