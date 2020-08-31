package output

import (
	"bytes"
	"os"

	"miller/containers"
)

// ostream *os.File in constructors/factory
type RecordWriterJSON struct {
	onFirst bool
}

func NewRecordWriterJSON() *RecordWriterJSON {
	return &RecordWriterJSON{
		true,
	}
}

func (this *RecordWriterJSON) Write(
	outrec *containers.Lrec,
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
		buffer.WriteString("  \"")
		buffer.WriteString(*pe.Key)
		buffer.WriteString("\": \"")
		buffer.WriteString(pe.Value.String())
		buffer.WriteString("\"")
		if pe.Next != nil {
			buffer.WriteString(",")
		}
		buffer.WriteString("\n")
	}
	buffer.WriteString("}\n")
	os.Stdout.WriteString(buffer.String())
}
