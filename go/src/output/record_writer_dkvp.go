package output

import (
	"bytes"
	"io"

	"miller/src/cliutil"
	"miller/src/types"
)

type RecordWriterDKVP struct {
	writerOptions *cliutil.TWriterOptions
}

func NewRecordWriterDKVP(writerOptions *cliutil.TWriterOptions) *RecordWriterDKVP {
	return &RecordWriterDKVP{
		writerOptions: writerOptions,
	}
}

func (writer *RecordWriterDKVP) Write(
	outrec *types.Mlrmap,
	ostream io.WriteCloser,
) {
	// End of record stream: nothing special for this output format
	if outrec == nil {
		return
	}

	if outrec.FieldCount == 0 {
		ostream.Write([]byte(writer.writerOptions.ORS))
		return
	}

	var buffer bytes.Buffer // 5x faster than fmt.Print() separately
	for pe := outrec.Head; pe != nil; pe = pe.Next {
		buffer.WriteString(pe.Key)
		buffer.WriteString(writer.writerOptions.OPS)
		buffer.WriteString(pe.Value.String())
		if pe.Next != nil {
			buffer.WriteString(writer.writerOptions.OFS)
		}
	}
	buffer.WriteString(writer.writerOptions.ORS)
	ostream.Write(buffer.Bytes())
}
