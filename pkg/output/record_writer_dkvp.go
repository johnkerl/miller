package output

import (
	"bufio"

	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/colorizer"
	"github.com/johnkerl/miller/pkg/mlrval"
)

type RecordWriterDKVP struct {
	writerOptions *cli.TWriterOptions
}

func NewRecordWriterDKVP(writerOptions *cli.TWriterOptions) (*RecordWriterDKVP, error) {
	return &RecordWriterDKVP{
		writerOptions: writerOptions,
	}, nil
}

func (writer *RecordWriterDKVP) Write(
	outrec *mlrval.Mlrmap,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) error {
	// End of record stream: nothing special for this output format
	if outrec == nil {
		return nil
	}

	if outrec.IsEmpty() {
		bufferedOutputStream.WriteString(writer.writerOptions.ORS)
		return nil
	}

	for pe := outrec.Head; pe != nil; pe = pe.Next {
		bufferedOutputStream.WriteString(colorizer.MaybeColorizeKey(pe.Key, outputIsStdout))
		bufferedOutputStream.WriteString(writer.writerOptions.OPS)
		bufferedOutputStream.WriteString(colorizer.MaybeColorizeValue(pe.Value.String(), outputIsStdout))
		if pe.Next != nil {
			bufferedOutputStream.WriteString(writer.writerOptions.OFS)
		}
	}
	bufferedOutputStream.WriteString(writer.writerOptions.ORS)

	return nil
}
