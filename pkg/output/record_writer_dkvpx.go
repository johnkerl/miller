package output

import (
	"bufio"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/colorizer"
	"github.com/johnkerl/miller/v6/pkg/dkvpx"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

type RecordWriterDKVPX struct {
	writerOptions *cli.TWriterOptions
}

func NewRecordWriterDKVPX(writerOptions *cli.TWriterOptions) (*RecordWriterDKVPX, error) {
	return &RecordWriterDKVPX{
		writerOptions: writerOptions,
	}, nil
}

func (writer *RecordWriterDKVPX) Write(
	outrec *mlrval.Mlrmap,
	_ *types.Context,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) error {
	if outrec == nil {
		return nil
	}

	if outrec.IsEmpty() {
		bufferedOutputStream.WriteString(writer.writerOptions.ORS)
		return nil
	}

	first := true
	for pe := outrec.Head; pe != nil; pe = pe.Next {
		if !first {
			bufferedOutputStream.WriteString(writer.writerOptions.OFS)
		}
		first = false

		keyStr := dkvpx.FormatField(pe.Key)
		valStr := dkvpx.FormatField(pe.Value.String())

		bufferedOutputStream.WriteString(colorizer.MaybeColorizeKey(keyStr, outputIsStdout))
		bufferedOutputStream.WriteString(writer.writerOptions.OPS)
		bufferedOutputStream.WriteString(colorizer.MaybeColorizeValue(valStr, outputIsStdout))
	}
	bufferedOutputStream.WriteString(writer.writerOptions.ORS)

	return nil
}
