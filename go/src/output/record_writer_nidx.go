package output

import (
	"bytes"
	"io"

	"miller/src/cliutil"
	"miller/src/types"
)

type RecordWriterNIDX struct {
	writerOptions *cliutil.TWriterOptions
}

func NewRecordWriterNIDX(writerOptions *cliutil.TWriterOptions) *RecordWriterNIDX {
	return &RecordWriterNIDX{
		writerOptions: writerOptions,
	}
}

func (writer *RecordWriterNIDX) Write(
	outrec *types.Mlrmap,
	ostream io.WriteCloser,
) {
	// End of record stream: nothing special for this output format
	if outrec == nil {
		return
	}

	var buffer bytes.Buffer // 5x faster than fmt.Print() separately
	for pe := outrec.Head; pe != nil; pe = pe.Next {
		buffer.WriteString(pe.Value.String())
		if pe.Next != nil {
			buffer.WriteString(writer.writerOptions.OFS)
		}
	}
	buffer.WriteString(writer.writerOptions.ORS)
	ostream.Write(buffer.Bytes())
}
