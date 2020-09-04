package output

import (
	"bytes"
	"fmt"
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
	outrec *lib.Mlrmap,
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
		buffer.WriteString("\": ")

		// Write the value which is a mlrval
		valueBytes, err := pe.Value.MarshalJSON()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		_, err = buffer.Write(valueBytes)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if pe.Next != nil {
			buffer.WriteString(",")
		}
		buffer.WriteString("\n")
	}
	buffer.WriteString("}\n")
	os.Stdout.WriteString(buffer.String())
}
