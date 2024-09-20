package output

import (
	"bufio"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/colorizer"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

type RecordWriterCSVLite struct {
	writerOptions *cli.TWriterOptions
	// For reporting schema changes: we print a newline and the new header
	lastJoinedHeader *string
	// Only write one blank line for schema changes / blank input lines
	justWroteEmptyLine bool
}

func NewRecordWriterCSVLite(writerOptions *cli.TWriterOptions) (*RecordWriterCSVLite, error) {
	return &RecordWriterCSVLite{
		writerOptions:      writerOptions,
		lastJoinedHeader:   nil,
		justWroteEmptyLine: false,
	}, nil
}

func (writer *RecordWriterCSVLite) Write(
	outrec *mlrval.Mlrmap,
	_ *types.Context,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) error {
	if outrec == nil {
		// End of record stream: nothing special for this output format
		return nil
	}

	if outrec.IsEmpty() {
		if !writer.justWroteEmptyLine {
			bufferedOutputStream.WriteString(writer.writerOptions.ORS)
		}
		joinedHeader := ""
		writer.lastJoinedHeader = &joinedHeader
		writer.justWroteEmptyLine = true
		return nil
	}

	needToPrintHeader := false
	joinedHeader := strings.Join(outrec.GetKeys(), ",")
	if writer.lastJoinedHeader == nil || *writer.lastJoinedHeader != joinedHeader {
		if writer.lastJoinedHeader != nil {
			if !writer.justWroteEmptyLine {
				bufferedOutputStream.WriteString(writer.writerOptions.ORS)
			}
			writer.justWroteEmptyLine = true
		}
		writer.lastJoinedHeader = &joinedHeader
		needToPrintHeader = true
	}

	if needToPrintHeader && !writer.writerOptions.HeaderlessOutput {
		for pe := outrec.Head; pe != nil; pe = pe.Next {
			bufferedOutputStream.WriteString(colorizer.MaybeColorizeKey(pe.Key, outputIsStdout))

			if pe.Next != nil {
				bufferedOutputStream.WriteString(writer.writerOptions.OFS)
			}
		}

		bufferedOutputStream.WriteString(writer.writerOptions.ORS)
	}

	for pe := outrec.Head; pe != nil; pe = pe.Next {
		bufferedOutputStream.WriteString(colorizer.MaybeColorizeValue(pe.Value.String(), outputIsStdout))
		if pe.Next != nil {
			bufferedOutputStream.WriteString(writer.writerOptions.OFS)
		}
	}
	bufferedOutputStream.WriteString(writer.writerOptions.ORS)

	writer.justWroteEmptyLine = false

	return nil
}
