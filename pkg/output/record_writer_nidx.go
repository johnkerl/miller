package output

import (
	"bufio"

	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/mlrval"
)

type RecordWriterNIDX struct {
	writerOptions *cli.TWriterOptions
}

func NewRecordWriterNIDX(writerOptions *cli.TWriterOptions) (*RecordWriterNIDX, error) {
	return &RecordWriterNIDX{
		writerOptions: writerOptions,
	}, nil
}

func (writer *RecordWriterNIDX) Write(
	outrec *mlrval.Mlrmap,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) error {
	// End of record stream: nothing special for this output format
	if outrec == nil {
		return nil
	}

	for pe := outrec.Head; pe != nil; pe = pe.Next {
		bufferedOutputStream.WriteString(pe.Value.String())
		if pe.Next != nil {
			bufferedOutputStream.WriteString(writer.writerOptions.OFS)
		}
	}
	bufferedOutputStream.WriteString(writer.writerOptions.ORS)

	return nil
}
