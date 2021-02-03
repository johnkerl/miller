package output

import (
	"bytes"
	"io"
	"strings"

	"miller/cliutil"
	"miller/types"
)

type RecordWriterCSVLite struct {
	ofs string
	ors string

	doHeaderlessOutput bool
	// For reporting schema changes: we print a newline and the new header
	lastJoinedHeader *string
	// Only write one blank line for schema changes / blank input lines
	justWroteEmptyLine bool
}

func NewRecordWriterCSVLite(writerOptions *cliutil.TWriterOptions) *RecordWriterCSVLite {

	return &RecordWriterCSVLite{
		ofs: writerOptions.OFS,
		ors: writerOptions.ORS,

		doHeaderlessOutput: writerOptions.HeaderlessCSVOutput,
		lastJoinedHeader:   nil,
		justWroteEmptyLine: false,
	}
}

func (this *RecordWriterCSVLite) Write(
	outrec *types.Mlrmap,
	ostream io.WriteCloser,
) {
	// End of record stream: nothing special for this output format
	if outrec == nil {
		return
	}

	if outrec.FieldCount == 0 {
		if !this.justWroteEmptyLine {
			ostream.Write([]byte(this.ors))
		}
		joinedHeader := ""
		this.lastJoinedHeader = &joinedHeader
		this.justWroteEmptyLine = true
		return
	}

	needToPrintHeader := false
	joinedHeader := strings.Join(outrec.GetKeys(), ",")
	if this.lastJoinedHeader == nil || *this.lastJoinedHeader != joinedHeader {
		if this.lastJoinedHeader != nil {
			if !this.justWroteEmptyLine {
				ostream.Write([]byte(this.ors))
			}
			this.justWroteEmptyLine = true
		}
		this.lastJoinedHeader = &joinedHeader
		needToPrintHeader = true
	}

	if needToPrintHeader && !this.doHeaderlessOutput {
		var buffer bytes.Buffer // faster than fmt.Print() separately
		for pe := outrec.Head; pe != nil; pe = pe.Next {
			buffer.WriteString(pe.Key)

			if pe.Next != nil {
				buffer.WriteString(this.ofs)
			}
		}

		buffer.WriteString(this.ors)
		ostream.Write(buffer.Bytes())
	}

	var buffer bytes.Buffer // faster than fmt.Print() separately
	for pe := outrec.Head; pe != nil; pe = pe.Next {
		buffer.WriteString(pe.Value.String())
		if pe.Next != nil {
			buffer.WriteString(this.ofs)
		}
	}
	buffer.WriteString(this.ors)
	ostream.Write(buffer.Bytes())

	this.justWroteEmptyLine = false
}
