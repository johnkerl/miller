package output

import (
	"bufio"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/colorizer"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

type RecordWriterDCF struct {
	writerOptions *cli.TWriterOptions
}

func NewRecordWriterDCF(writerOptions *cli.TWriterOptions) (*RecordWriterDCF, error) {
	return &RecordWriterDCF{
		writerOptions: writerOptions,
	}, nil
}

func (writer *RecordWriterDCF) Write(
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
		valueStr := dcfValueString(pe.Value)
		keyStr := colorizer.MaybeColorizeKey(pe.Key, outputIsStdout)
		valueStr = colorizer.MaybeColorizeValue(valueStr, outputIsStdout)
		writeDCFField(bufferedOutputStream, keyStr, valueStr)
	}
	bufferedOutputStream.WriteString("\n")
	return nil
}

// dcfValueString returns the string form of a field value for DCF output.
// Arrays are joined with ", " to match list fields (Depends, etc.).
func dcfValueString(mv *mlrval.Mlrval) string {
	if mv == nil {
		return ""
	}
	if mv.IsArray() {
		arr := mv.GetArray()
		if arr == nil || len(arr) == 0 {
			return ""
		}
		parts := make([]string, 0, len(arr))
		for _, el := range arr {
			if el != nil {
				parts = append(parts, el.String())
			}
		}
		return strings.Join(parts, ", ")
	}
	return mv.String()
}

// writeDCFField writes one "Key: value" line, folding on newlines per DCF:
// continuation lines start with a single space.
func writeDCFField(b *bufio.Writer, key, value string) {
	lines := strings.Split(value, "\n")
	for i, line := range lines {
		if i == 0 {
			b.WriteString(key)
			b.WriteString(": ")
			b.WriteString(line)
			b.WriteString("\n")
		} else {
			b.WriteString(" ")
			b.WriteString(line)
			b.WriteString("\n")
		}
	}
}
