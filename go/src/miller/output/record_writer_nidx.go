package output

import (
	"bytes"
	"os"

	"miller/containers"
)

// ostream *os.File in constructors/factory
type RecordWriterNIDX struct {
	ifs string
	ors string
}

func NewRecordWriterNIDX(ifs string) *RecordWriterNIDX {
	return &RecordWriterNIDX{
		ifs,
		"\n", // TODO: parameterize
	}
}

func (this *RecordWriterNIDX) Write(
	outrec *containers.Lrec,
) {
	// End of record stream: nothing special for this output format
	if outrec == nil {
		return
	}

	var buffer bytes.Buffer // 5x faster than fmt.Print() separately
	for pe := outrec.Head; pe != nil; pe = pe.Next {
		buffer.WriteString(*pe.Value)
		if pe.Next != nil {
			buffer.WriteString(this.ifs)
		}
	}
	buffer.WriteString(this.ors)
	os.Stdout.WriteString(buffer.String())
}
