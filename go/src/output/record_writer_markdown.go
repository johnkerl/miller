package output

import (
	"bytes"
	"fmt"
	"io"

	"miller/src/cliutil"
	"miller/src/types"
)

type RecordWriterMarkdown struct {
	writerOptions *cliutil.TWriterOptions

	numHeaderLinesOutput int
	lastJoinedHeader     string
}

func NewRecordWriterMarkdown(writerOptions *cliutil.TWriterOptions) *RecordWriterMarkdown {
	return &RecordWriterMarkdown{
		writerOptions: writerOptions,

		numHeaderLinesOutput: 0,
		lastJoinedHeader:     "",
	}
}

// ----------------------------------------------------------------
func (this *RecordWriterMarkdown) Write(
	outrec *types.Mlrmap,
	ostream io.WriteCloser,
) {
	if outrec == nil { // end of record stream
		return
	}

	currentJoinedHeader := outrec.GetKeysJoined()
	if this.lastJoinedHeader != "" {
		if currentJoinedHeader != this.lastJoinedHeader {
			this.lastJoinedHeader = ""
			if this.numHeaderLinesOutput > 0 {
				fmt.Fprintf(ostream, this.writerOptions.ORS)
			}
		}
	}

	var buffer bytes.Buffer

	if this.lastJoinedHeader == "" {
		buffer.WriteString("|")
		for pe := outrec.Head; pe != nil; pe = pe.Next {
			buffer.WriteString(" ")
			buffer.WriteString(pe.Key)
			buffer.WriteString(" |")
		}
		buffer.WriteString(this.writerOptions.ORS)

		buffer.WriteString("|")
		for pe := outrec.Head; pe != nil; pe = pe.Next {
			buffer.WriteString(" --- |")
		}
		buffer.WriteString(this.writerOptions.ORS)

		this.lastJoinedHeader = currentJoinedHeader
		this.numHeaderLinesOutput++
	}

	buffer.WriteString("|")
	for pe := outrec.Head; pe != nil; pe = pe.Next {
		buffer.WriteString(" ")
		buffer.WriteString(pe.Value.String())
		buffer.WriteString(" |")
	}
	buffer.WriteString(this.writerOptions.ORS)

	fmt.Fprint(ostream, buffer.String())
}
