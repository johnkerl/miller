package output

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"mlr/src/cli"
	"mlr/src/colorizer"
	"mlr/src/types"
)

type RecordWriterMarkdown struct {
	writerOptions *cli.TWriterOptions
	ors           string

	numHeaderLinesOutput int
	lastJoinedHeader     string
}

func NewRecordWriterMarkdown(writerOptions *cli.TWriterOptions) *RecordWriterMarkdown {
	return &RecordWriterMarkdown{
		writerOptions: writerOptions,

		numHeaderLinesOutput: 0,
		lastJoinedHeader:     "",
	}
}

// ----------------------------------------------------------------
func (writer *RecordWriterMarkdown) Write(
	outrec *types.Mlrmap,
	ostream io.WriteCloser,
	outputIsStdout bool,
) {
	if outrec == nil { // end of record stream
		return
	}

	currentJoinedHeader := outrec.GetKeysJoined()
	if writer.lastJoinedHeader != "" {
		if currentJoinedHeader != writer.lastJoinedHeader {
			writer.lastJoinedHeader = ""
			if writer.numHeaderLinesOutput > 0 {
				fmt.Fprintf(ostream, writer.writerOptions.ORS)
			}
		}
	}

	var buffer bytes.Buffer

	if writer.lastJoinedHeader == "" {
		buffer.WriteString("|")
		for pe := outrec.Head; pe != nil; pe = pe.Next {
			buffer.WriteString(" ")
			buffer.WriteString(colorizer.MaybeColorizeKey(pe.Key, outputIsStdout))
			buffer.WriteString(" |")
		}
		buffer.WriteString(writer.writerOptions.ORS)

		buffer.WriteString("|")
		for pe := outrec.Head; pe != nil; pe = pe.Next {
			buffer.WriteString(" --- |")
		}
		buffer.WriteString(writer.writerOptions.ORS)

		writer.lastJoinedHeader = currentJoinedHeader
		writer.numHeaderLinesOutput++
	}

	buffer.WriteString("|")
	for pe := outrec.Head; pe != nil; pe = pe.Next {
		buffer.WriteString(" ")
		value := strings.ReplaceAll(pe.Value.String(), "|", "\\|")
		buffer.WriteString(colorizer.MaybeColorizeValue(value, outputIsStdout))
		buffer.WriteString(" |")
	}
	buffer.WriteString(writer.writerOptions.ORS)

	fmt.Fprint(ostream, buffer.String())
}
