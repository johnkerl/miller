package output

import (
	"encoding/csv"
	"io"
	"os"
	"strings"

	"miller/cliutil"
	"miller/types"
)

type RecordWriterCSV struct {
	csvWriter          *csv.Writer
	doHeaderlessOutput bool
	// For reporting schema changes: we print a newline and the new header
	lastJoinedHeader *string
	// Only write one blank line for schema changes / blank input lines
	justWroteEmptyLine bool

	ofs string
}

func NewRecordWriterCSV(writerOptions *cliutil.TWriterOptions) *RecordWriterCSV {
	return &RecordWriterCSV{
		csvWriter:          nil, // will be set on first Write() wherein we have the ostream
		doHeaderlessOutput: writerOptions.HeaderlessCSVOutput,
		lastJoinedHeader:   nil,
		justWroteEmptyLine: false,
		ofs:                writerOptions.OFS,
	}
}

func (this *RecordWriterCSV) Write(
	outrec *types.Mlrmap,
	ostream io.WriteCloser,
) {
	// End of record stream: nothing special for this output format
	if outrec == nil {
		return
	}

	if this.csvWriter == nil {
		this.csvWriter = csv.NewWriter(os.Stdout)
		this.csvWriter.Comma = rune(this.ofs[0]) // xxx temp
	}

	if outrec.FieldCount == 0 {
		if !this.justWroteEmptyLine {
			ostream.Write([]byte("\n"))
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
				ostream.Write([]byte("\n"))
			}
			this.justWroteEmptyLine = true
		}
		this.lastJoinedHeader = &joinedHeader
		needToPrintHeader = true
	}

	if needToPrintHeader && !this.doHeaderlessOutput {
		fields := make([]string, outrec.FieldCount)
		i := 0
		for pe := outrec.Head; pe != nil; pe = pe.Next {
			fields[i] = pe.Key
			i++
		}
		this.csvWriter.Write(fields)
	}

	fields := make([]string, outrec.FieldCount)
	i := 0
	for pe := outrec.Head; pe != nil; pe = pe.Next {
		fields[i] = pe.Value.String()
		i++
	}
	this.csvWriter.Write(fields)
	this.csvWriter.Flush()
	this.justWroteEmptyLine = false
}
