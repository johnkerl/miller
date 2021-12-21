package output

import (
	"bufio"
	"fmt"
	"unicode/utf8"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/colorizer"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
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
	outrec *mlrval.Mlrmap,
	bufferedOutputStream *bufio.Writer,
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
		writer.writeWithRightAlignedValues(outrec, bufferedOutputStream, outputIsStdout, maxKeyLength)
	} else {
		writer.writeWithLeftAlignedValues(outrec, bufferedOutputStream, outputIsStdout, maxKeyLength)
	}
}

func (writer *RecordWriterXTAB) writeWithLeftAlignedValues(
	outrec *mlrval.Mlrmap,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
	maxKeyLength int,
) {
	// Put a blank line between records, but not before the first or after the last
	if writer.onFirst {
		writer.onFirst = false
	} else {
		bufferedOutputStream.WriteString(writer.writerOptions.OFS)
	}

	for pe := outrec.Head; pe != nil; pe = pe.Next {
		keyLength := utf8.RuneCountInString(pe.Key)
		keyPadLength := maxKeyLength - keyLength

		bufferedOutputStream.WriteString(colorizer.MaybeColorizeKey(pe.Key, outputIsStdout))

		if writer.opslen == 1 {
			bufferedOutputStream.WriteString(writer.writerOptions.OPS) // always at least once
			for i := 0; i < keyPadLength; i += writer.opslen {
				bufferedOutputStream.WriteString(writer.writerOptions.OPS)
			}
		} else {
			bufferedOutputStream.WriteString(writer.writerOptions.OPS)
		}

		bufferedOutputStream.WriteString(colorizer.MaybeColorizeValue(pe.Value.String(), outputIsStdout))
		bufferedOutputStream.WriteString(writer.writerOptions.OFS)
	}
}

func (writer *RecordWriterXTAB) writeWithRightAlignedValues(
	outrec *mlrval.Mlrmap,
	bufferedOutputStream *bufio.Writer,
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

	// Put a blank line between records, but not before the first or after the last
	if writer.onFirst {
		writer.onFirst = false
	} else {
		bufferedOutputStream.WriteString(writer.writerOptions.OFS)
	}

	i = 0
	for pe := outrec.Head; pe != nil; pe = pe.Next {
		keyLength := utf8.RuneCountInString(pe.Key)
		keyPadLength := maxKeyLength - keyLength

		bufferedOutputStream.WriteString(colorizer.MaybeColorizeKey(pe.Key, outputIsStdout))

		if writer.opslen == 1 {
			bufferedOutputStream.WriteString(writer.writerOptions.OPS) // always at least once
			for i := 0; i < keyPadLength; i += writer.opslen {
				bufferedOutputStream.WriteString(writer.writerOptions.OPS)
			}
		} else {
			bufferedOutputStream.WriteString(writer.writerOptions.OPS)
		}

		paddedValue := fmt.Sprintf("%*s", maxValueLength, values[i])
		bufferedOutputStream.WriteString(colorizer.MaybeColorizeValue(paddedValue, outputIsStdout))
		bufferedOutputStream.WriteString(writer.writerOptions.OFS)

		i++
	}
}
