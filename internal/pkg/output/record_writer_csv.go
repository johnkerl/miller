package output

import (
	"encoding/csv"
	"errors"
	"io"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/colorizer"
	"mlr/internal/pkg/types"
)

type RecordWriterCSV struct {
	writerOptions *cli.TWriterOptions
	ofs0          byte // Go's CSV library only lets its 'Comma' be a single character
	csvWriter     *csv.Writer
	// For reporting schema changes: we print a newline and the new header
	lastJoinedHeader *string
	// Only write one blank line for schema changes / blank input lines
	justWroteEmptyLine bool
}

func NewRecordWriterCSV(writerOptions *cli.TWriterOptions) (*RecordWriterCSV, error) {
	if len(writerOptions.OFS) != 1 {
		return nil, errors.New("CSV OFS can only be a single character")
	}
	if writerOptions.ORS != "\n" {
		return nil, errors.New("CSV ORS cannot be altered")
	}
	return &RecordWriterCSV{
		writerOptions:      writerOptions,
		csvWriter:          nil, // will be set on first Write() wherein we have the ostream
		lastJoinedHeader:   nil,
		justWroteEmptyLine: false,
	}, nil
}

func (writer *RecordWriterCSV) Write(
	outrec *types.Mlrmap,
	ostream io.WriteCloser,
	outputIsStdout bool,
) {
	// End of record stream: nothing special for this output format
	if outrec == nil {
		return
	}

	if writer.csvWriter == nil {
		writer.csvWriter = csv.NewWriter(ostream)
		writer.csvWriter.Comma = rune(writer.writerOptions.OFS[0]) // xxx temp
	}

	if outrec.IsEmpty() {
		if !writer.justWroteEmptyLine {
			ostream.Write([]byte("\n"))
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
				ostream.Write([]byte("\n"))
			}
			writer.justWroteEmptyLine = true
		}
		writer.lastJoinedHeader = &joinedHeader
		needToPrintHeader = true
	}

	if needToPrintHeader && !writer.writerOptions.HeaderlessCSVOutput {
		fields := make([]string, outrec.FieldCount)
		i := 0
		for pe := outrec.Head; pe != nil; pe = pe.Next {
			fields[i] = colorizer.MaybeColorizeKey(pe.Key, outputIsStdout)
			i++
		}
		writer.csvWriter.Write(fields)
	}

	fields := make([]string, outrec.FieldCount)
	i := 0
	for pe := outrec.Head; pe != nil; pe = pe.Next {
		fields[i] = colorizer.MaybeColorizeValue(pe.Value.String(), outputIsStdout)
		i++
	}
	writer.csvWriter.Write(fields)
	writer.csvWriter.Flush()
	writer.justWroteEmptyLine = false
}
