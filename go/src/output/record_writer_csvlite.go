package output

import (
	"bytes"
	"io"
	"strings"

	"mlr/src/cliutil"
	"mlr/src/colorizer"
	"mlr/src/types"
)

type RecordWriterCSVLite struct {
	writerOptions *cliutil.TWriterOptions
	// For reporting schema changes: we print a newline and the new header
	lastJoinedHeader *string
	// Only write one blank line for schema changes / blank input lines
	justWroteEmptyLine bool
}

func NewRecordWriterCSVLite(writerOptions *cliutil.TWriterOptions) *RecordWriterCSVLite {

	return &RecordWriterCSVLite{
		writerOptions:      writerOptions,
		lastJoinedHeader:   nil,
		justWroteEmptyLine: false,
	}
}

func (writer *RecordWriterCSVLite) Write(
	outrec *types.Mlrmap,
	ostream io.WriteCloser,
	outputIsStdout bool,
) {
	// End of record stream: nothing special for this output format
	if outrec == nil {
		return
	}

	if outrec.FieldCount == 0 {
		if !writer.justWroteEmptyLine {
			ostream.Write([]byte(writer.writerOptions.ORS))
		}
		joinedHeader := ""
		writer.lastJoinedHeader = &joinedHeader
		writer.justWroteEmptyLine = true
		return
	}

	needToPrintHeader := false
	joinedHeader := strings.Join(outrec.GetKeys(), ",")
	if writer.lastJoinedHeader == nil || *writer.lastJoinedHeader != joinedHeader {
		if writer.lastJoinedHeader != nil {
			if !writer.justWroteEmptyLine {
				ostream.Write([]byte(writer.writerOptions.ORS))
			}
			writer.justWroteEmptyLine = true
		}
		writer.lastJoinedHeader = &joinedHeader
		needToPrintHeader = true
	}

	if needToPrintHeader && !writer.writerOptions.HeaderlessCSVOutput {
		var buffer bytes.Buffer // faster than fmt.Print() separately
		for pe := outrec.Head; pe != nil; pe = pe.Next {
			buffer.WriteString(colorizer.MaybeColorizeKey(pe.Key, outputIsStdout))

			if pe.Next != nil {
				buffer.WriteString(writer.writerOptions.OFS)
			}
		}

		buffer.WriteString(writer.writerOptions.ORS)
		ostream.Write(buffer.Bytes())
	}

	var buffer bytes.Buffer // faster than fmt.Print() separately
	for pe := outrec.Head; pe != nil; pe = pe.Next {
		buffer.WriteString(colorizer.MaybeColorizeValue(pe.Value.String(), outputIsStdout))
		if pe.Next != nil {
			buffer.WriteString(writer.writerOptions.OFS)
		}
	}
	buffer.WriteString(writer.writerOptions.ORS)
	ostream.Write(buffer.Bytes())

	writer.justWroteEmptyLine = false
}
