package output

import (
	"fmt"
	"os"

	"miller/clitypes"
	"miller/types"
)

// ----------------------------------------------------------------
type RecordWriterJSON struct {
	// Parameters:
	wrapJSONOutputInOuterList bool
	jsonFormatting            types.TJSONFormatting

	// State:
	onFirst bool
}

// ----------------------------------------------------------------
func NewRecordWriterJSON(writerOptions *clitypes.TWriterOptions) *RecordWriterJSON {
	var jsonFormatting types.TJSONFormatting = types.JSON_SINGLE_LINE
	if writerOptions.JSONOutputMultiline {
		jsonFormatting = types.JSON_MULTILINE
	}
	return &RecordWriterJSON{
		wrapJSONOutputInOuterList: writerOptions.WrapJSONOutputInOuterList,
		jsonFormatting:            jsonFormatting,
		onFirst:                   true,
	}
}

// ----------------------------------------------------------------
func (this *RecordWriterJSON) Write(
	outrec *types.Mlrmap,
) {
	if this.wrapJSONOutputInOuterList {
		this.writeWithListWrap(outrec)
	} else {
		this.writeWithoutListWrap(outrec)
	}
}

// ----------------------------------------------------------------
func (this *RecordWriterJSON) writeWithListWrap(
	outrec *types.Mlrmap,
) {
	if outrec != nil { // Not end of record stream
		if this.onFirst {
			os.Stdout.WriteString("[\n")
		}

		// The Mlrmap MarshalJSON doesn't include the final newline, so that we
		// can place it neatly with commas here (if the user requested them).
		bytes, err := outrec.MarshalJSON(this.jsonFormatting)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if !this.onFirst {
			os.Stdout.WriteString(",\n")
		}

		os.Stdout.Write(bytes)

		this.onFirst = false

	} else { // End of record stream
		if this.onFirst { // zero records in the entire output stream
			os.Stdout.WriteString("[")
		}
		os.Stdout.WriteString("\n]\n")
	}
}

// ----------------------------------------------------------------
func (this *RecordWriterJSON) writeWithoutListWrap(
	outrec *types.Mlrmap,
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

	os.Stdout.Write(bytes)
	os.Stdout.WriteString("\n")
}
