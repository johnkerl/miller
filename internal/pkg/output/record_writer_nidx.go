package output

import (
	"bytes"
	"io"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/types"
)

type RecordWriterNIDX struct {
	writerOptions *cli.TWriterOptions
	ofs           string
	ors           string
}

func NewRecordWriterNIDX(writerOptions *cli.TWriterOptions) (*RecordWriterNIDX, error) {
	return &RecordWriterNIDX{
		writerOptions: writerOptions,
	}, nil
}

func (writer *RecordWriterNIDX) Write(
	outrec *types.Mlrmap,
	ostream io.WriteCloser,
	outputIsStdout bool,
) {
	// End of record stream: nothing special for this output format
	if outrec == nil {
		return
	}

	var buffer bytes.Buffer // stdio is non-buffered in Go, so buffer for ~5x speed increase
	for pe := outrec.Head; pe != nil; pe = pe.Next {
		buffer.WriteString(pe.Value.String())
		if pe.Next != nil {
			buffer.WriteString(writer.writerOptions.OFS)
		}
	}
	buffer.WriteString(writer.writerOptions.ORS)
	ostream.Write(buffer.Bytes())
}
