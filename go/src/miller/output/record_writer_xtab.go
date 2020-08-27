package output

import (
	// System:
	"bytes"
	"os"
	// Miller:
	"miller/containers"
)

// ostream *os.File in constructors/factory
type RecordWriterXTAB struct {
	onFirst bool
}

func NewRecordWriterXTAB() *RecordWriterXTAB {
	return &RecordWriterXTAB{
		true,
	}
}

func (this *RecordWriterXTAB) Write(
	outrec *containers.Lrec,
) {
	// End of record stream: nothing special for this output format
	if outrec == nil {
		return
	}

	maxKeyLength := 1
	for pe := outrec.Head; pe != nil; pe = pe.Next {
		keyLength := len(*pe.Key)
		if keyLength > maxKeyLength {
			maxKeyLength = keyLength
		}
	}

	var buffer bytes.Buffer // 5x faster than fmt.Print() separately

	// Put a blank line between records, but not before the first or after the last
	if this.onFirst {
		this.onFirst = false
	} else {
		buffer.WriteString("\n")
	}

	for pe := outrec.Head; pe != nil; pe = pe.Next {
		keyLength := len(*pe.Key)
		padLength := maxKeyLength - keyLength

		buffer.WriteString(*pe.Key)
		buffer.WriteString(" ")
		for i := 0; i < padLength; i++ {
			buffer.WriteString(" ")
		}
		buffer.WriteString(*pe.Value)
		buffer.WriteString("\n")
	}
	os.Stdout.WriteString(buffer.String())
}
