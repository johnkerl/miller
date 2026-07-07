package output

import (
	"bufio"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/colorizer"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

type RecordWriterMarkdown struct {
	writerOptions *cli.TWriterOptions

	// Streaming-mode state (when MarkdownAlignedOutput is false)
	numHeaderLinesOutput int
	lastJoinedHeader     string

	// Aligned-mode state (when MarkdownAlignedOutput is true)
	batch             []*mlrval.Mlrmap
	batchJoinedHeader *string
	numBatchesOutput  int
}

func NewRecordWriterMarkdown(writerOptions *cli.TWriterOptions) (*RecordWriterMarkdown, error) {
	return &RecordWriterMarkdown{
		writerOptions: writerOptions,

		numHeaderLinesOutput: 0,
		lastJoinedHeader:     "",

		batch:             []*mlrval.Mlrmap{},
		batchJoinedHeader: nil,
		numBatchesOutput:  0,
	}, nil
}

func (writer *RecordWriterMarkdown) Write(
	outrec *mlrval.Mlrmap,
	_ *types.Context,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) error {
	if writer.writerOptions.MarkdownAlignedOutput {
		return writer.writeAligned(outrec, bufferedOutputStream, outputIsStdout)
	}
	return writer.writeStreaming(outrec, bufferedOutputStream, outputIsStdout)
}

func (writer *RecordWriterMarkdown) writeStreaming(
	outrec *mlrval.Mlrmap,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) error {
	if outrec == nil { // end of record stream
		return nil
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
			// In streaming mode the header-separator line must be emitted
			// before subsequent records are seen, so with
			// --right-align-numeric the alignment marker is chosen from this
			// first record's value.
			if writer.writerOptions.RightAlignNumericOutput && pe.Value.IsNumeric() {
				bufferedOutputStream.WriteString(" ---: |")
			} else {
				bufferedOutputStream.WriteString(" --- |")
			}
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

	return nil
}

// writeAligned accumulates records into a same-schema batch and flushes when
// the schema changes or the stream ends. We need the whole batch in hand
// before emitting any row, since column widths depend on every value.
func (writer *RecordWriterMarkdown) writeAligned(
	outrec *mlrval.Mlrmap,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) error {
	if outrec == nil {
		if len(writer.batch) > 0 {
			writer.flushBatch(bufferedOutputStream, outputIsStdout)
		}
		return nil
	}

	joinedHeader := outrec.GetKeysJoined()
	if writer.batchJoinedHeader == nil {
		writer.batch = append(writer.batch, outrec)
		writer.batchJoinedHeader = &joinedHeader
	} else if *writer.batchJoinedHeader != joinedHeader {
		writer.flushBatch(bufferedOutputStream, outputIsStdout)
		writer.batch = []*mlrval.Mlrmap{outrec}
		writer.batchJoinedHeader = &joinedHeader
	} else {
		writer.batch = append(writer.batch, outrec)
	}
	return nil
}

func (writer *RecordWriterMarkdown) flushBatch(
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) {
	if writer.numBatchesOutput > 0 {
		bufferedOutputStream.WriteString(writer.writerOptions.ORS)
	}

	first := writer.batch[0]

	// With --right-align-numeric, a column gets a right-alignment marker
	// (`---:`) when every value in the batch's column is numeric.
	columnRightAligned := make(map[string]bool)
	if writer.writerOptions.RightAlignNumericOutput {
		for pe := first.Head; pe != nil; pe = pe.Next {
			columnRightAligned[pe.Key] = true
		}
		for _, rec := range writer.batch {
			for pe := rec.Head; pe != nil; pe = pe.Next {
				if !pe.Value.IsNumeric() {
					columnRightAligned[pe.Key] = false
				}
			}
		}
	}

	// Floor of 3 so "---" never overflows the column -- or 4 for "---:".
	maxWidths := make(map[string]int)
	for pe := first.Head; pe != nil; pe = pe.Next {
		minWidth := 3
		if columnRightAligned[pe.Key] {
			minWidth = 4
		}
		maxWidths[pe.Key] = max(lib.DisplayWidth(pe.Key), minWidth)
	}
	for _, rec := range writer.batch {
		for pe := rec.Head; pe != nil; pe = pe.Next {
			value := strings.ReplaceAll(pe.Value.String(), "|", "\\|")
			width := lib.DisplayWidth(value)
			if width > maxWidths[pe.Key] {
				maxWidths[pe.Key] = width
			}
		}
	}

	// Header. Right-aligned columns get right-justified header text so that
	// header and data share the same alignment in the raw markdown, matching
	// how Markdown viewers render the `---:` marker.
	bufferedOutputStream.WriteString("|")
	for pe := first.Head; pe != nil; pe = pe.Next {
		bufferedOutputStream.WriteString(" ")
		if columnRightAligned[pe.Key] {
			writePadding(bufferedOutputStream, maxWidths[pe.Key]-lib.DisplayWidth(pe.Key))
			bufferedOutputStream.WriteString(colorizer.MaybeColorizeKey(pe.Key, outputIsStdout))
		} else {
			bufferedOutputStream.WriteString(colorizer.MaybeColorizeKey(pe.Key, outputIsStdout))
			writePadding(bufferedOutputStream, maxWidths[pe.Key]-lib.DisplayWidth(pe.Key))
		}
		bufferedOutputStream.WriteString(" |")
	}
	bufferedOutputStream.WriteString(writer.writerOptions.ORS)

	// Separator
	bufferedOutputStream.WriteString("|")
	for pe := first.Head; pe != nil; pe = pe.Next {
		if columnRightAligned[pe.Key] {
			writePadding(bufferedOutputStream, maxWidths[pe.Key]-4)
			bufferedOutputStream.WriteString(" ---:")
		} else {
			bufferedOutputStream.WriteString(" ---")
			writePadding(bufferedOutputStream, maxWidths[pe.Key]-3)
		}
		bufferedOutputStream.WriteString(" |")
	}
	bufferedOutputStream.WriteString(writer.writerOptions.ORS)

	// Data
	for _, rec := range writer.batch {
		bufferedOutputStream.WriteString("|")
		for pe := rec.Head; pe != nil; pe = pe.Next {
			value := strings.ReplaceAll(pe.Value.String(), "|", "\\|")
			bufferedOutputStream.WriteString(" ")
			if columnRightAligned[pe.Key] {
				writePadding(bufferedOutputStream, maxWidths[pe.Key]-lib.DisplayWidth(value))
				bufferedOutputStream.WriteString(colorizer.MaybeColorizeValue(value, outputIsStdout))
			} else {
				bufferedOutputStream.WriteString(colorizer.MaybeColorizeValue(value, outputIsStdout))
				writePadding(bufferedOutputStream, maxWidths[pe.Key]-lib.DisplayWidth(value))
			}
			bufferedOutputStream.WriteString(" |")
		}
		bufferedOutputStream.WriteString(writer.writerOptions.ORS)
	}

	writer.batch = nil
	writer.batchJoinedHeader = nil
	writer.numBatchesOutput++
}

func writePadding(bufferedOutputStream *bufio.Writer, n int) {
	if n <= 0 {
		return
	}
	bufferedOutputStream.WriteString(strings.Repeat(" ", n))
}
