package output

import (
	"bytes"
	"io"
	"unicode/utf8"

	"mlr/src/cli"
	"mlr/src/colorizer"
	"mlr/src/types"
)

type RecordWriterXTAB struct {
	writerOptions *cli.TWriterOptions
	onFirst       bool
}

func NewRecordWriterXTAB(writerOptions *cli.TWriterOptions) *RecordWriterXTAB {
	return &RecordWriterXTAB{
		writerOptions: writerOptions,
		onFirst:       true,
	}
}

func (writer *RecordWriterXTAB) Write(
	outrec *types.Mlrmap,
	ostream io.WriteCloser,
	outputIsStdout bool,
) {
	// End of record stream: nothing special for this output format
	if outrec == nil {
		return
	}

	maxKeyLength := 1
	for pe := outrec.Head; pe != nil; pe = pe.Next {
		keyLength := utf8.RuneCountInString(pe.Key)
		if keyLength > maxKeyLength {
			maxKeyLength = keyLength
		}
	}

	var buffer bytes.Buffer // 5x faster than fmt.Print() separately

	// Put a blank line between records, but not before the first or after the last
	if writer.onFirst {
		writer.onFirst = false
	} else {
		buffer.WriteString("\n") // TODO: ORS
	}

	for pe := outrec.Head; pe != nil; pe = pe.Next {
		keyLength := utf8.RuneCountInString(pe.Key)
		padLength := maxKeyLength - keyLength

		buffer.WriteString(colorizer.MaybeColorizeKey(pe.Key, outputIsStdout))
		buffer.WriteString(" ")
		for i := 0; i < padLength; i++ {
			buffer.WriteString(writer.writerOptions.OPS)
		}
		buffer.WriteString(colorizer.MaybeColorizeValue(pe.Value.String(), outputIsStdout))
		buffer.WriteString("\n") // TODO: ORS
	}
	ostream.Write(buffer.Bytes())
}
