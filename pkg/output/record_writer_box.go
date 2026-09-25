package output

import (
	"bufio"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/colorizer"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
)

// NewRecordWriterBox produces Unicode box-drawing tables, e.g.:
//
// ┌─────┬─────┬────┬─────────────────────┬─────────────────────┐
// │ a   │ b   │ i  │ x                   │ y                   │
// ├─────┼─────┼────┼─────────────────────┼─────────────────────┤
// │ pan │ pan │ 1  │ 0.3467901443380824  │ 0.7268028627434533  │
// ├─────┼─────┼────┼─────────────────────┼─────────────────────┤
// │ eks │ pan │ 2  │ 0.7586799647899636  │ 0.5221511083334797  │
// └─────┴─────┴────┴─────────────────────┴─────────────────────┘
//
// This is the same as `--opprint --barred-unicode` except a divider is drawn
// after every row, not just after the header. It reuses RecordWriterPPRINT's
// same-schema batching/width computation (see writeHeterogenousList in
// record_writer_pprint.go), which is exactly what --barred-unicode needs too.
func NewRecordWriterBox(writerOptions *cli.TWriterOptions) (*RecordWriterPPRINT, error) {
	return NewRecordWriterPPRINT(writerOptions)
}

// escapeBoxValue escapes a literal box-drawing vertical bar in a cell value
// so it can't be mistaken for a column boundary. As with Markdown's analogous
// "|" escaping, this is output-only: box-format input does not unescape it.
func escapeBoxValue(value string) string {
	return strings.ReplaceAll(value, "│", "\\│")
}

func (writer *RecordWriterPPRINT) writeHeterogenousListBoxed(
	records []*mlrval.Mlrmap,
	maxWidths map[string]int,
	rightAlignedHeaders map[string]bool,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) {

	bc := newBarredChars(writer.writerOptions.OFS, true)

	horizontalBars := make(map[string]string)
	for key, width := range maxWidths {
		horizontalBars[key] = strings.Repeat(bc.horizontalBar, width)
	}

	writeHorizontalRule := func(start, middle, end string, outrec *mlrval.Mlrmap) {
		bufferedOutputStream.WriteString(start)
		for pe := outrec.Head; pe != nil; pe = pe.Next {
			bufferedOutputStream.WriteString(horizontalBars[pe.Key])
			if pe.Next != nil {
				bufferedOutputStream.WriteString(middle)
			} else {
				bufferedOutputStream.WriteString(end)
				bufferedOutputStream.WriteString(writer.writerOptions.ORS)
			}
		}
	}

	onFirst := true
	for i, outrec := range records {

		if onFirst {
			// Unlike --barred/--barred-unicode, the top border is always
			// printed (even with --headerless-output) so the box stays
			// visually complete; only the header text line and its
			// divider are suppressed.
			writeHorizontalRule(bc.firstRowHorizontalStart, bc.firstRowHorizontalMiddle, bc.firstRowHorizontalEnd, outrec)

			if !writer.writerOptions.HeaderlessOutput {
				bufferedOutputStream.WriteString(bc.verticalStart)
				for pe := outrec.Head; pe != nil; pe = pe.Next {
					if !writer.headerIsRightAligned(pe.Key, rightAlignedHeaders) { // left-align
						bufferedOutputStream.WriteString(colorizer.MaybeColorizeKey(pe.Key, outputIsStdout))
						writer.writePadding(pe.Key, maxWidths[pe.Key], bufferedOutputStream)
					} else { // right-align
						writer.writePadding(pe.Key, maxWidths[pe.Key], bufferedOutputStream)
						bufferedOutputStream.WriteString(colorizer.MaybeColorizeKey(pe.Key, outputIsStdout))
					}
					if pe.Next != nil {
						bufferedOutputStream.WriteString(bc.verticalMiddle)
					} else {
						bufferedOutputStream.WriteString(bc.verticalEnd)
						bufferedOutputStream.WriteString(writer.writerOptions.ORS)
					}
				}

				writeHorizontalRule(bc.horizontalStart, bc.horizontalMiddle, bc.horizontalEnd, outrec)
			}
		}
		onFirst = false

		// Data line
		bufferedOutputStream.WriteString(bc.verticalStart)
		for pe := outrec.Head; pe != nil; pe = pe.Next {
			s := escapeBoxValue(pe.Value.String())
			if !writer.cellIsRightAligned(pe.Value) { // left-align
				bufferedOutputStream.WriteString(colorizer.MaybeColorizeValue(s, outputIsStdout))
				writer.writePadding(s, maxWidths[pe.Key], bufferedOutputStream)
			} else { // right-align
				writer.writePadding(s, maxWidths[pe.Key], bufferedOutputStream)
				bufferedOutputStream.WriteString(colorizer.MaybeColorizeValue(s, outputIsStdout))
			}
			if pe.Next != nil {
				bufferedOutputStream.WriteString(bc.verticalMiddle)
			} else {
				bufferedOutputStream.WriteString(bc.verticalEnd)
				bufferedOutputStream.WriteString(writer.writerOptions.ORS)
			}
		}

		// Divider after every row: the bottom border on the last row, a
		// mid-divider otherwise.
		if i == len(records)-1 {
			writeHorizontalRule(bc.lastRowHorizontalStart, bc.lastRowHorizontalMiddle, bc.lastRowHorizontalEnd, outrec)
		} else {
			writeHorizontalRule(bc.horizontalStart, bc.horizontalMiddle, bc.horizontalEnd, outrec)
		}

		if writer.writerOptions.FlushOnEveryRecord {
			// bufio.Writer errors are sticky; the final Flush in pkg/stream is checked
			_ = bufferedOutputStream.Flush()
		}
	}
}
