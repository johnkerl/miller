package output

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/colorizer"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
)

type RecordWriterTSV struct {
	writerOptions *cli.TWriterOptions
	// For reporting schema changes: we print a newline and the new header
	lastJoinedHeader *string
	// Only write one blank line for schema changes / blank input lines
	justWroteEmptyLine bool
}

func NewRecordWriterTSV(writerOptions *cli.TWriterOptions) (*RecordWriterTSV, error) {
	if writerOptions.OFS != "\t" {
		return nil, fmt.Errorf("for TSV, OFS cannot be altered")
	}
	if writerOptions.ORS != "\n" && writerOptions.ORS != "\r\n" {
		return nil, fmt.Errorf("for CSV, ORS cannot be altered")
	}
	return &RecordWriterTSV{
		writerOptions:      writerOptions,
		lastJoinedHeader:   nil,
		justWroteEmptyLine: false,
	}, nil
}

func (writer *RecordWriterTSV) Write(
	outrec *mlrval.Mlrmap,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) {
	// End of record stream: nothing special for this output format
	if outrec == nil {
		return
	}

	if outrec.IsEmpty() {
		if !writer.justWroteEmptyLine {
			bufferedOutputStream.WriteString(writer.writerOptions.ORS)
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
				bufferedOutputStream.WriteString(writer.writerOptions.ORS)
			}
			writer.justWroteEmptyLine = true
		}
		writer.lastJoinedHeader = &joinedHeader
		needToPrintHeader = true
	}

	if needToPrintHeader && !writer.writerOptions.HeaderlessCSVOutput {
		for pe := outrec.Head; pe != nil; pe = pe.Next {
			bufferedOutputStream.WriteString(
				colorizer.MaybeColorizeKey(
					lib.TSVEncodeField(
						pe.Key,
					),
					outputIsStdout,
				),
			)

			if pe.Next != nil {
				bufferedOutputStream.WriteString(writer.writerOptions.OFS)
			}
		}

		bufferedOutputStream.WriteString(writer.writerOptions.ORS)
	}

	for pe := outrec.Head; pe != nil; pe = pe.Next {
		bufferedOutputStream.WriteString(
			colorizer.MaybeColorizeValue(
				lib.TSVEncodeField(
					pe.Value.String(),
				),
				outputIsStdout,
			),
		)
		if pe.Next != nil {
			bufferedOutputStream.WriteString(writer.writerOptions.OFS)
		}
	}
	bufferedOutputStream.WriteString(writer.writerOptions.ORS)

	writer.justWroteEmptyLine = false
}
