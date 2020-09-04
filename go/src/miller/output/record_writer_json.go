package output

import (
	"bytes"
	"os"

	"miller/clitypes"
	"miller/lib"
)

// ostream *os.File in constructors/factory
type RecordWriterJSON struct {
	onFirst bool
}

func NewRecordWriterJSON(writerOptions *clitypes.TWriterOptions) *RecordWriterJSON {
	return &RecordWriterJSON{
		onFirst: true,
	}
}

func (this *RecordWriterJSON) Write(
	outrec *lib.Lrec,
) {
	// End of record stream
	if outrec == nil {
		// TODO: only if jlistwrap
		//os.Stdout.WriteString("\n")
		return
	}

	var buffer bytes.Buffer // 5x faster than fmt.Print() separately

	// Put a comma between records, but not before the first or after the last
	if this.onFirst {
		this.onFirst = false
	} else {
		// TODO: only if jlistwrap
		//buffer.WriteString(",\n")
	}

	// TODO: value-quoting (numeric) ...
	buffer.WriteString("{\n")
	for pe := outrec.Head; pe != nil; pe = pe.Next {
		// Write the key which is necessarily string-valued in Miller
		buffer.WriteString("  \"")
		buffer.WriteString(*pe.Key)

		// Write the value which is a mlrval
		sval, needsQuote := pe.Value.StringWithQuoteInfo()
		buffer.WriteString("\": ")
		if needsQuote {
			buffer.WriteString("\"")
		}
		buffer.WriteString(sval)
		if needsQuote {
			buffer.WriteString("\"")
		}

		if pe.Next != nil {
			buffer.WriteString(",")
		}
		buffer.WriteString("\n")
	}
	buffer.WriteString("}\n")
	os.Stdout.WriteString(buffer.String())
}
