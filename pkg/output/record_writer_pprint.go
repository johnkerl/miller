package output

import (
	"bufio"
	"container/list"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/colorizer"
	"github.com/johnkerl/miller/pkg/mlrval"
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
	outrec *mlrval.Mlrmap,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) error {
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
					bufferedOutputStream,
					outputIsStdout,
				)
				if nonEmpty {
					// Print a newline
					bufferedOutputStream.WriteString(writer.writerOptions.ORS)
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
			writer.writeHeterogenousList(writer.batch, writer.writerOptions.BarredPprintOutput,
				bufferedOutputStream, outputIsStdout)
		}
	}

	return nil
}

// ----------------------------------------------------------------
// Returns false if there was nothing but empty record(s), e.g. 'mlr gap -n 10'.
func (writer *RecordWriterPPRINT) writeHeterogenousList(
	records *list.List,
	barred bool,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) bool {
	maxWidths := make(map[string]int)
	var maxNR int64 = 0

	for e := records.Front(); e != nil; e = e.Next() {
		outrec := e.Value.(*mlrval.Mlrmap)
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
			writer.writeHeterogenousListBarred(records, maxWidths, bufferedOutputStream, outputIsStdout)
		} else {
			writer.writeHeterogenousListNonBarred(records, maxWidths, bufferedOutputStream, outputIsStdout)
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
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) {

	onFirst := true
	for e := records.Front(); e != nil; e = e.Next() {
		outrec := e.Value.(*mlrval.Mlrmap)

		// Print header line
		if onFirst && !writer.writerOptions.HeaderlessOutput {
			for pe := outrec.Head; pe != nil; pe = pe.Next {
				if !writer.writerOptions.RightAlignedPPRINTOutput { // left-align
					if pe.Next != nil {
						// Header line, left-align, not last column
						bufferedOutputStream.WriteString(colorizer.MaybeColorizeKey(pe.Key, outputIsStdout))
						writer.writePadding(pe.Key, maxWidths[pe.Key], bufferedOutputStream)
						bufferedOutputStream.WriteString(writer.writerOptions.OFS)
					} else {
						// Header line, left-align, last column
						bufferedOutputStream.WriteString(colorizer.MaybeColorizeKey(pe.Key, outputIsStdout))
						bufferedOutputStream.WriteString(writer.writerOptions.ORS)
					}
				} else { // right-align
					writer.writePadding(pe.Key, maxWidths[pe.Key], bufferedOutputStream)
					bufferedOutputStream.WriteString(colorizer.MaybeColorizeKey(pe.Key, outputIsStdout))
					if pe.Next != nil {
						// Header line, right-align, not last column
						bufferedOutputStream.WriteString(writer.writerOptions.OFS)
					} else {
						// Header line, right-align, last column
						bufferedOutputStream.WriteString(writer.writerOptions.ORS)
					}
				}

			}
		}
		onFirst = false

		// Print data lines
		for pe := outrec.Head; pe != nil; pe = pe.Next {
			s := pe.Value.String()
			if s == "" {
				s = "-"
			}
			if !writer.writerOptions.RightAlignedPPRINTOutput { // left-align
				if pe.Next != nil {
					// Data line, left-align, not last column
					bufferedOutputStream.WriteString(colorizer.MaybeColorizeValue(s, outputIsStdout))
					writer.writePadding(s, maxWidths[pe.Key], bufferedOutputStream)
					bufferedOutputStream.WriteString(writer.writerOptions.OFS)
				} else {
					// Data line, left-align, last column
					bufferedOutputStream.WriteString(colorizer.MaybeColorizeValue(s, outputIsStdout))
					bufferedOutputStream.WriteString(writer.writerOptions.ORS)
				}
			} else { // right-align
				writer.writePadding(s, maxWidths[pe.Key], bufferedOutputStream)
				bufferedOutputStream.WriteString(colorizer.MaybeColorizeValue(s, outputIsStdout))
				if pe.Next != nil {
					// Data line, right-align, not last column
					bufferedOutputStream.WriteString(writer.writerOptions.OFS)
				} else {
					// Data line, right-align, last column
					bufferedOutputStream.WriteString(writer.writerOptions.ORS)
				}
			}
		}

		if writer.writerOptions.FlushOnEveryRecord {
			bufferedOutputStream.Flush()
		}
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

func (writer *RecordWriterPPRINT) writeHeterogenousListBarred(
	records *list.List,
	maxWidths map[string]int,
	bufferedOutputStream *bufio.Writer,
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
		outrec := e.Value.(*mlrval.Mlrmap)

		// Print header line
		if onFirst && !writer.writerOptions.HeaderlessOutput {
			bufferedOutputStream.WriteString(horizontalStart)
			for pe := outrec.Head; pe != nil; pe = pe.Next {
				bufferedOutputStream.WriteString(horizontalBars[pe.Key])
				if pe.Next != nil {
					bufferedOutputStream.WriteString(horizontalMiddle)
				} else {
					bufferedOutputStream.WriteString(horizontalEnd)
					bufferedOutputStream.WriteString(writer.writerOptions.ORS)
				}
			}

			bufferedOutputStream.WriteString(verticalStart)
			for pe := outrec.Head; pe != nil; pe = pe.Next {
				if !writer.writerOptions.RightAlignedPPRINTOutput { // left-align
					bufferedOutputStream.WriteString(colorizer.MaybeColorizeKey(pe.Key, outputIsStdout))
					writer.writePadding(pe.Key, maxWidths[pe.Key], bufferedOutputStream)
				} else { // right-align
					writer.writePadding(pe.Key, maxWidths[pe.Key], bufferedOutputStream)
					bufferedOutputStream.WriteString(colorizer.MaybeColorizeKey(pe.Key, outputIsStdout))
				}
				if pe.Next != nil {
					bufferedOutputStream.WriteString(verticalMiddle)
				} else {
					bufferedOutputStream.WriteString(verticalEnd)
					bufferedOutputStream.WriteString(writer.writerOptions.ORS)
				}
			}

			bufferedOutputStream.WriteString(horizontalStart)
			for pe := outrec.Head; pe != nil; pe = pe.Next {
				bufferedOutputStream.WriteString(horizontalBars[pe.Key])
				if pe.Next != nil {
					bufferedOutputStream.WriteString(horizontalMiddle)
				} else {
					bufferedOutputStream.WriteString(horizontalEnd)
					bufferedOutputStream.WriteString(writer.writerOptions.ORS)
				}
			}
		}
		onFirst = false

		// Print data lines
		bufferedOutputStream.WriteString(verticalStart)
		for pe := outrec.Head; pe != nil; pe = pe.Next {
			s := pe.Value.String()
			if !writer.writerOptions.RightAlignedPPRINTOutput { // left-align
				bufferedOutputStream.WriteString(colorizer.MaybeColorizeValue(s, outputIsStdout))
				writer.writePadding(s, maxWidths[pe.Key], bufferedOutputStream)
			} else { // right-align
				writer.writePadding(s, maxWidths[pe.Key], bufferedOutputStream)
				bufferedOutputStream.WriteString(colorizer.MaybeColorizeValue(s, outputIsStdout))
			}
			if pe.Next != nil {
				bufferedOutputStream.WriteString(fmt.Sprint(verticalMiddle))
			} else {
				bufferedOutputStream.WriteString(verticalEnd)
				bufferedOutputStream.WriteString(writer.writerOptions.ORS)
			}
		}

		if e.Next() == nil {
			bufferedOutputStream.WriteString(horizontalStart)
			for pe := outrec.Head; pe != nil; pe = pe.Next {
				bufferedOutputStream.WriteString(horizontalBars[pe.Key])
				if pe.Next != nil {
					bufferedOutputStream.WriteString(horizontalMiddle)
				} else {
					bufferedOutputStream.WriteString(horizontalEnd)
					bufferedOutputStream.WriteString(writer.writerOptions.ORS)
				}
			}
		}

		if writer.writerOptions.FlushOnEveryRecord {
			bufferedOutputStream.Flush()
		}
	}
}

func (writer *RecordWriterPPRINT) writePadding(
	text string,
	fieldWidth int,
	bufferedOutputStream *bufio.Writer,
) {
	textWidth := utf8.RuneCountInString(text)
	padWidth := fieldWidth - textWidth
	ofs := writer.writerOptions.OFS
	for i := 0; i < padWidth; i++ {
		bufferedOutputStream.WriteString(ofs)
	}
}
