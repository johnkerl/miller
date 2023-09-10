package output

import (
	"bufio"
	"strings"

	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/colorizer"
	"github.com/johnkerl/miller/pkg/mlrval"
)

type RecordWriterMarkdown struct {
	writerOptions *cli.TWriterOptions
	ors           string

	numHeaderLinesOutput int
	lastJoinedHeader     string
}

func NewRecordWriterMarkdown(writerOptions *cli.TWriterOptions) (*RecordWriterMarkdown, error) {
	return &RecordWriterMarkdown{
		writerOptions: writerOptions,

		numHeaderLinesOutput: 0,
		lastJoinedHeader:     "",
	}, nil
}

// ----------------------------------------------------------------
func (writer *RecordWriterMarkdown) Write(
	outrec *mlrval.Mlrmap,
	bufferedOutputStream *bufio.Writer,
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
				bufferedOutputStream.WriteString(writer.writerOptions.ORS)
			}
		}
	}

	if writer.lastJoinedHeader == "" {
		bufferedOutputStream.WriteString("|")
		for pe := outrec.Head; pe != nil; pe = pe.Next {
			bufferedOutputStream.WriteString(" ")
			bufferedOutputStream.WriteString(colorizer.MaybeColorizeKey(pe.Key, outputIsStdout))
			bufferedOutputStream.WriteString(" |")
		}
		bufferedOutputStream.WriteString(writer.writerOptions.ORS)

		bufferedOutputStream.WriteString("|")
		for pe := outrec.Head; pe != nil; pe = pe.Next {
			bufferedOutputStream.WriteString(" --- |")
		}
		bufferedOutputStream.WriteString(writer.writerOptions.ORS)

		writer.lastJoinedHeader = currentJoinedHeader
		writer.numHeaderLinesOutput++
	}

	bufferedOutputStream.WriteString("|")
	for pe := outrec.Head; pe != nil; pe = pe.Next {
		bufferedOutputStream.WriteString(" ")
		value := strings.ReplaceAll(pe.Value.String(), "|", "\\|")
		bufferedOutputStream.WriteString(colorizer.MaybeColorizeValue(value, outputIsStdout))
		bufferedOutputStream.WriteString(" |")
	}
	bufferedOutputStream.WriteString(writer.writerOptions.ORS)
}
