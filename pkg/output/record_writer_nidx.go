package output

import (
	"bufio"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
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
	_ *types.Context,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) error {
	if outrec == nil {
		// End of record stream: nothing special for this output format
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
