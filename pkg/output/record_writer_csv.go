package output

import (
	"bufio"
	"fmt"
	"strings"

	csv "github.com/johnkerl/miller/pkg/go-csv"

	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/mlrval"
)

type RecordWriterCSV struct {
	writerOptions *cli.TWriterOptions
	ofs0          byte // Go's CSV library only lets its 'Comma' be a single character
	csvWriter     *csv.Writer
	// For reporting schema changes: we print a newline and the new header
	lastJoinedHeader *string
	// Only write one blank line for schema changes / blank input lines
	justWroteEmptyLine bool
	// For double-quote around all fields
	quoteAll bool
}

func NewRecordWriterCSV(writerOptions *cli.TWriterOptions) (*RecordWriterCSV, error) {
	if len(writerOptions.OFS) != 1 {
		return nil, fmt.Errorf("for CSV, OFS can only be a single character")
	}
	if writerOptions.ORS != "\n" && writerOptions.ORS != "\r\n" {
		return nil, fmt.Errorf("for CSV, ORS cannot be altered")
	}
	return &RecordWriterCSV{
		writerOptions:      writerOptions,
		csvWriter:          nil, // will be set on first Write() wherein we have the output stream
		lastJoinedHeader:   nil,
		justWroteEmptyLine: false,
		quoteAll:           writerOptions.CSVQuoteAll,
	}, nil
}

func (writer *RecordWriterCSV) Write(
	outrec *mlrval.Mlrmap,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) {
	// End of record stream: nothing special for this output format
	if outrec == nil {
		return
	}

	if writer.csvWriter == nil {
		writer.csvWriter = csv.NewWriter(bufferedOutputStream)
		writer.csvWriter.Comma = rune(writer.writerOptions.OFS[0]) // xxx temp
	}

	if outrec.IsEmpty() {
		if !writer.justWroteEmptyLine {
			bufferedOutputStream.WriteString("\n")
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
				bufferedOutputStream.WriteString("\n")
			}
			writer.justWroteEmptyLine = true
		}
		writer.lastJoinedHeader = &joinedHeader
		needToPrintHeader = true
	}

	if needToPrintHeader && !writer.writerOptions.HeaderlessOutput {
		fields := make([]string, outrec.FieldCount)
		i := 0
		for pe := outrec.Head; pe != nil; pe = pe.Next {
			fields[i] = pe.Key
			i++
		}
		//////writer.csvWriter.Write(fields)
		writer.WriteCSVRecordMaybeColorized(fields, bufferedOutputStream, outputIsStdout, true, writer.quoteAll)
	}

	fields := make([]string, outrec.FieldCount)
	i := 0
	for pe := outrec.Head; pe != nil; pe = pe.Next {
		fields[i] = pe.Value.String()
		i++
	}
	writer.WriteCSVRecordMaybeColorized(fields, bufferedOutputStream, outputIsStdout, false, writer.quoteAll)
	writer.justWroteEmptyLine = false
}
