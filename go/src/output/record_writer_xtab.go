package output

import (
	"bytes"
	"fmt"
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

	if writer.writerOptions.RightAlignedXTABOutput {
		writer.writeWithRightAlignedValues(outrec, ostream, outputIsStdout, maxKeyLength)
	} else {
		writer.writeWithLeftAlignedValues(outrec, ostream, outputIsStdout, maxKeyLength)
	}

}

func (writer *RecordWriterXTAB) writeWithLeftAlignedValues(
	outrec *types.Mlrmap,
	ostream io.WriteCloser,
	outputIsStdout bool,
	maxKeyLength int,
) {

	var buffer bytes.Buffer // stdio is non-buffered in Go, so buffer for ~5x speed increase

	// Put a blank line between records, but not before the first or after the last
	if writer.onFirst {
		writer.onFirst = false
	} else {
		buffer.WriteString("\n") // TODO: ORS
	}

	for pe := outrec.Head; pe != nil; pe = pe.Next {
		keyLength := utf8.RuneCountInString(pe.Key)
		keyPadLength := maxKeyLength - keyLength

		buffer.WriteString(colorizer.MaybeColorizeKey(pe.Key, outputIsStdout))
		buffer.WriteString(" ")
		for i := 0; i < keyPadLength; i++ {
			buffer.WriteString(writer.writerOptions.OPS)
		}
		buffer.WriteString(colorizer.MaybeColorizeValue(pe.Value.String(), outputIsStdout))
		buffer.WriteString("\n") // TODO: ORS
	}
	ostream.Write(buffer.Bytes())

}

func (writer *RecordWriterXTAB) writeWithRightAlignedValues(
	outrec *types.Mlrmap,
	ostream io.WriteCloser,
	outputIsStdout bool,
	maxKeyLength int,
) {

	values := make([]string, outrec.FieldCount)

	maxValueLength := 0
	i := 0
	for pe := outrec.Head; pe != nil; pe = pe.Next {
		value := pe.Value.String()
		values[i] = value
		valueLength := utf8.RuneCountInString(value)
		if valueLength > maxValueLength {
			maxValueLength = valueLength
		}
		i++
	}

	var buffer bytes.Buffer // stdio is non-buffered in Go, so buffer for ~5x speed increase

	// Put a blank line between records, but not before the first or after the last
	if writer.onFirst {
		writer.onFirst = false
	} else {
		buffer.WriteString("\n") // TODO: ORS
	}

	i = 0
	for pe := outrec.Head; pe != nil; pe = pe.Next {
		keyLength := utf8.RuneCountInString(pe.Key)
		keyPadLength := maxKeyLength - keyLength

		buffer.WriteString(colorizer.MaybeColorizeKey(pe.Key, outputIsStdout))
		buffer.WriteString(" ")
		for i := 0; i < keyPadLength; i++ {
			buffer.WriteString(writer.writerOptions.OPS)
		}
		paddedValue := fmt.Sprintf("%*s", maxValueLength, values[i])
		buffer.WriteString(colorizer.MaybeColorizeValue(paddedValue, outputIsStdout))
		buffer.WriteString("\n") // TODO: ORS

		i++
	}
	ostream.Write(buffer.Bytes())

}
