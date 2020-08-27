package output

import (
	// System:
	"bytes"
	"os"
	// Miller:
	"miller/containers"
)

// ostream *os.File in constructors/factory
type RecordWriterDKVP struct {
	ifs string
	ips string
	ors string
}

func NewRecordWriterDKVP(ifs string, ips string) *RecordWriterDKVP {
	return &RecordWriterDKVP{
		ifs,
		ips,
		"\n", // TODO: parameterize
	}
}

func (this *RecordWriterDKVP) Write(
	outrec *containers.Lrec,
) {
	// End of record stream: nothing special for this output format
	if outrec == nil {
		return
	}

	var buffer bytes.Buffer // 5x faster than fmt.Print() separately
	for pe := outrec.Head; pe != nil; pe = pe.Next {
		buffer.WriteString(*pe.Key)
		buffer.WriteString(this.ips)
		buffer.WriteString(*pe.Value)
		if pe.Next != nil {
			buffer.WriteString(this.ifs)
		}
	}
	buffer.WriteString(this.ors)
	os.Stdout.WriteString(buffer.String())
}
