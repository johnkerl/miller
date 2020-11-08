package output

import (
	"encoding/csv"
	"os"
	"strings"

	"miller/clitypes"
	"miller/types"
)

// ostream *os.File in constructors/factory
type RecordWriterCSV struct {
	csvWriter          *csv.Writer
	doHeaderlessOutput bool
	// For reporting schema changes: we print a newline and the new header
	lastJoinedHeader *string
	// Only write one blank line for schema changes / blank input lines
	justWroteEmptyLine bool
}

func NewRecordWriterCSV(writerOptions *clitypes.TWriterOptions) *RecordWriterCSV {
	csvWriter := csv.NewWriter(os.Stdout)
	// xxx temp
	csvWriter.Comma = rune(writerOptions.OFS[0])

	return &RecordWriterCSV{
		csvWriter:          csvWriter,
		doHeaderlessOutput: writerOptions.HeaderlessCSVOutput,
		lastJoinedHeader:   nil,
		justWroteEmptyLine: false,
	}
}

func (this *RecordWriterCSV) Write(
	outrec *types.Mlrmap,
) {
	// End of record stream: nothing special for this output format
	if outrec == nil {
		return
	}

	if outrec.FieldCount == 0 {
		if !this.justWroteEmptyLine {
			os.Stdout.WriteString("\n")
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
				os.Stdout.WriteString("\n")
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
			fields[i] = *pe.Key
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
