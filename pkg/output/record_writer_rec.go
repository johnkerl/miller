package output

import (
	"bufio"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/colorizer"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

// RecordWriterREC writes GNU recutils (.rec) format: "FieldName: Value"
// lines, one record per blank-line-separated stanza. Unlike DCF, recutils
// has no hardcoded list-valued field names -- array/map values are handled
// by Miller's generic auto-flatten mechanism upstream of this writer.
type RecordWriterREC struct {
	writerOptions *cli.TWriterOptions
}

func NewRecordWriterREC(writerOptions *cli.TWriterOptions) (*RecordWriterREC, error) {
	return &RecordWriterREC{
		writerOptions: writerOptions,
	}, nil
}

func (writer *RecordWriterREC) Write(
	outrec *mlrval.Mlrmap,
	_ *types.Context,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) error {
	if outrec == nil {
		return nil
	}
	if outrec.IsEmpty() {
		bufferedOutputStream.WriteString("\n")
		return nil
	}

	for pe := outrec.Head; pe != nil; pe = pe.Next {
		valueStr := pe.Value.String()
		keyStr := colorizer.MaybeColorizeKey(pe.Key, outputIsStdout)
		valueStr = colorizer.MaybeColorizeValue(valueStr, outputIsStdout)
		writeRECField(bufferedOutputStream, keyStr, valueStr)
	}
	bufferedOutputStream.WriteString("\n")
	return nil
}

// writeRECField writes one "Key: value" line, folding embedded newlines in
// the value into recutils continuation lines: each continuation line is
// prefixed with "+ ".
func writeRECField(b *bufio.Writer, key, value string) {
	lines := strings.Split(value, "\n")
	for i, line := range lines {
		if i == 0 {
			b.WriteString(key)
			b.WriteString(": ")
			b.WriteString(line)
			b.WriteString("\n")
		} else {
			b.WriteString("+ ")
			b.WriteString(line)
			b.WriteString("\n")
		}
	}
}
