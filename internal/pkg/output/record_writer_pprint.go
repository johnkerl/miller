package output

import (
	"bytes"
	"container/list"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/colorizer"
	"mlr/internal/pkg/types"
)

type RecordWriterPPRINT struct {
	writerOptions *cli.TWriterOptions
	// Input:
	records *list.List

	// State:
	lastJoinedHeader *string
	batch            *list.List
}

func NewRecordWriterPPRINT(writerOptions *cli.TWriterOptions) (*RecordWriterPPRINT, error) {
	return &RecordWriterPPRINT{
		writerOptions: writerOptions,
		records:       list.New(),

		lastJoinedHeader: nil,
		batch:            list.New(),
	}, nil
}

// ----------------------------------------------------------------
func (writer *RecordWriterPPRINT) Write(
	outrec *types.Mlrmap,
	ostream io.WriteCloser,
	outputIsStdout bool,
) {
	// Group records by have-same-schema or not. Pretty-print each
	// homoegeneous sublist, or "batch".
	//
	// No output until end of a homogeneous batch of records, since we need to
	// find out max width down each column.

	if outrec != nil { // Not end of record stream

		if writer.lastJoinedHeader == nil {
			// First output record:
			// * New batch
			// * No old batch to print
			writer.batch.PushBack(outrec)
			temp := strings.Join(outrec.GetKeys(), ",")
			writer.lastJoinedHeader = &temp
		} else {
			// May or may not continue the same homogeneous batch
			joinedHeader := strings.Join(outrec.GetKeys(), ",")
			if *writer.lastJoinedHeader != joinedHeader {
				// Print and free old batch
				nonEmpty := writer.writeHeterogenousList(
					writer.batch,
					writer.writerOptions.BarredPprintOutput,
					ostream,
					outputIsStdout,
				)
				if nonEmpty {
					// Print a newline
					ostream.Write([]byte(writer.writerOptions.ORS))
				}
				// Start a new batch
				writer.batch = list.New()
				writer.batch.PushBack(outrec)
				writer.lastJoinedHeader = &joinedHeader
			} else {
				// Continue the batch
				writer.batch.PushBack(outrec)
			}
		}

	} else { // End of record stream

		if writer.batch.Front() != nil {
			writer.writeHeterogenousList(writer.batch, writer.writerOptions.BarredPprintOutput, ostream, outputIsStdout)
		}
	}
}

// ----------------------------------------------------------------
// Returns false if there was nothing but empty record(s), e.g. 'mlr gap -n 10'.
func (writer *RecordWriterPPRINT) writeHeterogenousList(
	records *list.List,
	barred bool,
	ostream io.WriteCloser,
	outputIsStdout bool,
) bool {
	maxWidths := make(map[string]int)
	var maxNR int = 0

	for e := records.Front(); e != nil; e = e.Next() {
		outrec := e.Value.(*types.Mlrmap)
		nr := outrec.FieldCount
		if maxNR < nr {
			maxNR = nr
		}
		for pe := outrec.Head; pe != nil; pe = pe.Next {
			width := utf8.RuneCountInString(pe.Value.String())
			if width == 0 {
				width = 1 // We'll rewrite "" to "-" below
			}
			oldMaxWidth := maxWidths[pe.Key]
			if width > oldMaxWidth {
				maxWidths[pe.Key] = width
			}
		}
	}

	if maxNR == 0 {
		return false
	} else {
		// Column name may be longer/shorter than all data values in the column
		for key, oldMaxWidth := range maxWidths {
			width := utf8.RuneCountInString(key)
			if width > oldMaxWidth {
				maxWidths[key] = width
			}
		}
		if barred {
			writer.writeHeterogenousListBarred(records, maxWidths, ostream, outputIsStdout)
		} else {
			writer.writeHeterogenousListNonBarred(records, maxWidths, ostream, outputIsStdout)
		}
		return true
	}
}

// ----------------------------------------------------------------
// Example:
//
// a   b   i  x                    y
// pan pan 1  0.3467901443380824   0.7268028627434533
// eks pan 2  -0.7586799647899636  0.5221511083334797
// wye wye 3  0.20460330576630303  0.33831852551664776
// eks wye 4  -0.38139939387114097 0.13418874328430463
// wye pan 5  0.5732889198020006   0.8636244699032729

func (writer *RecordWriterPPRINT) writeHeterogenousListNonBarred(
	records *list.List,
	maxWidths map[string]int,
	ostream io.WriteCloser,
	outputIsStdout bool,
) {

	onFirst := true
	for e := records.Front(); e != nil; e = e.Next() {
		outrec := e.Value.(*types.Mlrmap)

		// Print header line
		if onFirst && !writer.writerOptions.HeaderlessCSVOutput {
			var buffer bytes.Buffer // faster than fmt.Print() separately
			for pe := outrec.Head; pe != nil; pe = pe.Next {
				if !writer.writerOptions.RightAlignedPPRINTOutput { // left-align
					if pe.Next != nil {
						// Header line, left-align, not last column
						buffer.WriteString(colorizer.MaybeColorizeKey(pe.Key, outputIsStdout))
						writer.writePadding(pe.Key, maxWidths[pe.Key], &buffer)
						buffer.WriteString(writer.writerOptions.OFS)
					} else {
						// Header line, left-align, last column
						buffer.WriteString(colorizer.MaybeColorizeKey(pe.Key, outputIsStdout))
						buffer.WriteString(writer.writerOptions.ORS)
					}
				} else { // right-align
					writer.writePadding(pe.Key, maxWidths[pe.Key], &buffer)
					buffer.WriteString(colorizer.MaybeColorizeKey(pe.Key, outputIsStdout))
					if pe.Next != nil {
						// Header line, right-align, not last column
						buffer.WriteString(writer.writerOptions.OFS)
					} else {
						// Header line, right-align, last column
						buffer.WriteString(writer.writerOptions.ORS)
					}
				}

			}
			ostream.Write(buffer.Bytes())
		}
		onFirst = false

		// Print data lines
		var buffer bytes.Buffer // faster than fmt.Print() separately
		for pe := outrec.Head; pe != nil; pe = pe.Next {
			s := pe.Value.String()
			if s == "" {
				s = "-"
			}
			if !writer.writerOptions.RightAlignedPPRINTOutput { // left-align
				if pe.Next != nil {
					// Data line, left-align, not last column
					buffer.WriteString(colorizer.MaybeColorizeValue(s, outputIsStdout))
					writer.writePadding(s, maxWidths[pe.Key], &buffer)
					buffer.WriteString(writer.writerOptions.OFS)
				} else {
					// Data line, left-align, last column
					buffer.WriteString(colorizer.MaybeColorizeValue(s, outputIsStdout))
					buffer.WriteString(writer.writerOptions.ORS)
				}
			} else { // right-align
				writer.writePadding(s, maxWidths[pe.Key], &buffer)
				buffer.WriteString(colorizer.MaybeColorizeValue(s, outputIsStdout))
				if pe.Next != nil {
					// Data line, right-align, not last column
					buffer.WriteString(writer.writerOptions.OFS)
				} else {
					// Data line, right-align, last column
					buffer.WriteString(writer.writerOptions.ORS)
				}
			}
		}
		ostream.Write(buffer.Bytes())
	}
}

// ----------------------------------------------------------------
// Example:
//
// +-----+-----+----+----------------------+---------------------+
// | a   | b   | i  | x                    | y                   |
// +-----+-----+----+----------------------+---------------------+
// | pan | pan | 1  | 0.3467901443380824   | 0.7268028627434533  |
// | eks | pan | 2  | -0.7586799647899636  | 0.5221511083334797  |
// | wye | wye | 3  | 0.20460330576630303  | 0.33831852551664776 |
// | eks | wye | 4  | -0.38139939387114097 | 0.13418874328430463 |
// | wye | pan | 5  | 0.5732889198020006   | 0.8636244699032729  |
// +-----+-----+----+----------------------+---------------------+

// TODO: for better performance, uuse string-buffer as in DKVP for this and all
// record-writers

func (writer *RecordWriterPPRINT) writeHeterogenousListBarred(
	records *list.List,
	maxWidths map[string]int,
	ostream io.WriteCloser,
	outputIsStdout bool,
) {

	horizontalBars := make(map[string]string)
	for key, width := range maxWidths {
		horizontalBars[key] = strings.Repeat("-", width)
	}
	ofs := writer.writerOptions.OFS
	horizontalStart := "+-"
	horizontalMiddle := "-+-"
	horizontalEnd := "-+"
	verticalStart := "|" + ofs
	verticalMiddle := ofs + "|" + ofs
	verticalEnd := ofs + "|"

	onFirst := true
	for e := records.Front(); e != nil; e = e.Next() {
		outrec := e.Value.(*types.Mlrmap)

		// Print header line
		if onFirst && !writer.writerOptions.HeaderlessCSVOutput {
			var buffer bytes.Buffer // faster than fmt.Print() separately

			buffer.WriteString(horizontalStart)
			for pe := outrec.Head; pe != nil; pe = pe.Next {
				buffer.WriteString(horizontalBars[pe.Key])
				if pe.Next != nil {
					buffer.WriteString(horizontalMiddle)
				} else {
					buffer.WriteString(horizontalEnd)
					buffer.WriteString(writer.writerOptions.ORS)
				}
			}

			buffer.WriteString(verticalStart)
			for pe := outrec.Head; pe != nil; pe = pe.Next {
				if !writer.writerOptions.RightAlignedPPRINTOutput { // left-align
					buffer.WriteString(colorizer.MaybeColorizeKey(pe.Key, outputIsStdout))
					writer.writePadding(pe.Key, maxWidths[pe.Key], &buffer)
				} else { // right-align
					writer.writePadding(pe.Key, maxWidths[pe.Key], &buffer)
					buffer.WriteString(colorizer.MaybeColorizeKey(pe.Key, outputIsStdout))
				}
				if pe.Next != nil {
					buffer.WriteString(verticalMiddle)
				} else {
					buffer.WriteString(verticalEnd)
					buffer.WriteString(writer.writerOptions.ORS)
				}
			}

			buffer.WriteString(horizontalStart)
			for pe := outrec.Head; pe != nil; pe = pe.Next {
				buffer.WriteString(horizontalBars[pe.Key])
				if pe.Next != nil {
					buffer.WriteString(horizontalMiddle)
				} else {
					buffer.WriteString(horizontalEnd)
					buffer.WriteString(writer.writerOptions.ORS)
				}
			}

			ostream.Write(buffer.Bytes())
		}
		onFirst = false

		// Print data lines
		var buffer bytes.Buffer // faster than fmt.Print() separately
		buffer.WriteString(verticalStart)
		for pe := outrec.Head; pe != nil; pe = pe.Next {
			s := pe.Value.String()
			if !writer.writerOptions.RightAlignedPPRINTOutput { // left-align
				buffer.WriteString(colorizer.MaybeColorizeValue(s, outputIsStdout))
				writer.writePadding(s, maxWidths[pe.Key], &buffer)
			} else { // right-align
				writer.writePadding(s, maxWidths[pe.Key], &buffer)
				buffer.WriteString(colorizer.MaybeColorizeValue(s, outputIsStdout))
			}
			if pe.Next != nil {
				buffer.WriteString(fmt.Sprint(verticalMiddle))
			} else {
				buffer.WriteString(verticalEnd)
				buffer.WriteString(writer.writerOptions.ORS)
			}
		}

		if e.Next() == nil {
			buffer.WriteString(horizontalStart)
			for pe := outrec.Head; pe != nil; pe = pe.Next {
				buffer.WriteString(horizontalBars[pe.Key])
				if pe.Next != nil {
					buffer.WriteString(horizontalMiddle)
				} else {
					buffer.WriteString(horizontalEnd)
					buffer.WriteString(writer.writerOptions.ORS)
				}
			}
		}

		ostream.Write(buffer.Bytes())
	}
}

func (writer *RecordWriterPPRINT) writePadding(
	text string,
	fieldWidth int,
	buffer *bytes.Buffer,
) {
	textWidth := utf8.RuneCountInString(text)
	padWidth := fieldWidth - textWidth
	ofs := writer.writerOptions.OFS
	for i := 0; i < padWidth; i++ {
		buffer.WriteString(ofs)
	}
}
