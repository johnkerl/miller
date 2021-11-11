package output

import (
	"bytes"
	"fmt"
	"io"
	"unicode/utf8"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/colorizer"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
// Note: If OPS is single-character then we can do alignment of the form
//
//   ab  123
//   def 4567
//
// On the other hand, if it's multi-character, we won't be able to align
// neatly in all cases. Yet we do allow multi-character OPS, just without
// repetition: if someone wants to use OPS ": " and format data as
//
//   ab: 123
//   def: 4567
//
// then they can do that.
// ----------------------------------------------------------------

type RecordWriterXTAB struct {
	writerOptions *cli.TWriterOptions
	// Note: XTAB uses two consecutive OFS in place of an ORS; ORS is ignored
	opslen  int
	onFirst bool
}

func NewRecordWriterXTAB(writerOptions *cli.TWriterOptions) (*RecordWriterXTAB, error) {
	return &RecordWriterXTAB{
		writerOptions: writerOptions,
		opslen:        utf8.RuneCountInString(writerOptions.OPS),
		onFirst:       true,
	}, nil
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
		buffer.WriteString(writer.writerOptions.OFS)
	}

	for pe := outrec.Head; pe != nil; pe = pe.Next {
		keyLength := utf8.RuneCountInString(pe.Key)
		keyPadLength := maxKeyLength - keyLength

		buffer.WriteString(colorizer.MaybeColorizeKey(pe.Key, outputIsStdout))

		if writer.opslen == 1 {
			buffer.WriteString(writer.writerOptions.OPS) // always at least once
			for i := 0; i < keyPadLength; i += writer.opslen {
				buffer.WriteString(writer.writerOptions.OPS)
			}
		} else {
			buffer.WriteString(writer.writerOptions.OPS)
		}

		buffer.WriteString(colorizer.MaybeColorizeValue(pe.Value.String(), outputIsStdout))
		buffer.WriteString(writer.writerOptions.OFS)
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
		buffer.WriteString(writer.writerOptions.OFS)
	}

	i = 0
	for pe := outrec.Head; pe != nil; pe = pe.Next {
		keyLength := utf8.RuneCountInString(pe.Key)
		keyPadLength := maxKeyLength - keyLength

		buffer.WriteString(colorizer.MaybeColorizeKey(pe.Key, outputIsStdout))

		if writer.opslen == 1 {
			buffer.WriteString(writer.writerOptions.OPS) // always at least once
			for i := 0; i < keyPadLength; i += writer.opslen {
				buffer.WriteString(writer.writerOptions.OPS)
			}
		} else {
			buffer.WriteString(writer.writerOptions.OPS)
		}

		paddedValue := fmt.Sprintf("%*s", maxValueLength, values[i])
		buffer.WriteString(colorizer.MaybeColorizeValue(paddedValue, outputIsStdout))
		buffer.WriteString(writer.writerOptions.OFS)

		i++
	}
	ostream.Write(buffer.Bytes())
}
