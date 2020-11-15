package output

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"miller/clitypes"
	"miller/types"
)

// ostream *os.File in constructors/factory
type RecordWriterPPRINT struct {
	records *list.List
	// For detecting schema changes: we print a newline and the new header.
	barred bool
}

func NewRecordWriterPPRINT(writerOptions *clitypes.TWriterOptions) *RecordWriterPPRINT {
	return &RecordWriterPPRINT{
		records: list.New(),
		barred:  writerOptions.BarredPprintOutput,
	}
}

// ----------------------------------------------------------------
// xxx this is very naive at present -- needs copy from the C version.
func (this *RecordWriterPPRINT) Write(
	outrec *types.Mlrmap,
) {
	// No output until end of record stream, since we need to find out max
	// width down each column.
	if outrec != nil {
		this.records.PushBack(outrec)
		return
	}

	// Group records by have-same-schema or not. Pretty-print each
	// homoegeneous sublist, or "batch".

	var lastJoinedHeader *string = nil
	batch := list.New()
	for {
		head := this.records.Front()
		if head == nil {
			break
		}
		record := head.Value.(*types.Mlrmap)
		this.records.Remove(head)

		if lastJoinedHeader == nil {
			// First output record
			// New batch
			// No old batch to print
			batch.PushBack(record)
			temp := strings.Join(record.GetKeys(), ",")
			lastJoinedHeader = &temp
		} else {
			joinedHeader := strings.Join(record.GetKeys(), ",")
			if *lastJoinedHeader != joinedHeader {
				// Print and free old batch
				this.writeHeterogenousList(batch, this.barred)
				// Print a newline
				os.Stdout.WriteString("\n")
				// Start a new batch
				batch = list.New()
				batch.PushBack(record)
				lastJoinedHeader = &joinedHeader
			} else {
				// Continue the batch
				batch.PushBack(record)
			}
		}
	}
	if batch.Front() != nil {
		this.writeHeterogenousList(batch, this.barred)
	}
}

// ----------------------------------------------------------------
func (this *RecordWriterPPRINT) writeHeterogenousList(
	records *list.List,
	barred bool,
) {
	maxWidths := make(map[string]int)

	for e := records.Front(); e != nil; e = e.Next() {
		outrec := e.Value.(*types.Mlrmap)
		for pe := outrec.Head; pe != nil; pe = pe.Next {
			width := len(pe.Value.String())
			if width == 0 {
				width = 1 // We'll rewrite "" to "-" below
			}
			oldMaxWidth := maxWidths[*pe.Key]
			if width > oldMaxWidth {
				maxWidths[*pe.Key] = width
			}
		}
	}

	// Column name may be longer/shorter than all data values in the column
	for key, oldMaxWidth := range maxWidths {
		width := len(key)
		if width > oldMaxWidth {
			maxWidths[key] = width
		}
	}
	if barred {
		this.writeHeterogenousListBarred(records, maxWidths)
	} else {
		this.writeHeterogenousListNonBarred(records, maxWidths)
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

func (this *RecordWriterPPRINT) writeHeterogenousListNonBarred(
	records *list.List,
	maxWidths map[string]int,
) {

	onFirst := true
	for e := records.Front(); e != nil; e = e.Next() {
		outrec := e.Value.(*types.Mlrmap)

		// Print header line
		if onFirst {
			for pe := outrec.Head; pe != nil; pe = pe.Next {
				if pe.Next != nil {
					fmt.Printf("%-*s ", maxWidths[*pe.Key], *pe.Key)
				} else {
					fmt.Println(*pe.Key)
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
			if pe.Next != nil {
				fmt.Printf("%-*s ", maxWidths[*pe.Key], s)
			} else {
				fmt.Println(s)
			}
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

// TODO: for better performance, uuse string-buffer as in DKVP for this and all
// record-writers

func (this *RecordWriterPPRINT) writeHeterogenousListBarred(
	records *list.List,
	maxWidths map[string]int,
) {

	horizontalBars := make(map[string]string)
	for key, width := range maxWidths {
		horizontalBars[key] = strings.Repeat("-", width)
	}
	horizontalStart := "+-"
	horizontalMiddle := "-+-"
	horizontalEnd := "-+"
	verticalStart := "| "
	verticalMiddle := " | "
	verticalEnd := " |"

	onFirst := true
	for e := records.Front(); e != nil; e = e.Next() {
		outrec := e.Value.(*types.Mlrmap)

		// Print header line
		if onFirst {

			fmt.Print(horizontalStart)
			for pe := outrec.Head; pe != nil; pe = pe.Next {
				fmt.Print(horizontalBars[*pe.Key])
				if pe.Next != nil {
					fmt.Print(horizontalMiddle)
				} else {
					fmt.Println(horizontalEnd)
				}
			}

			fmt.Print(verticalStart)
			for pe := outrec.Head; pe != nil; pe = pe.Next {
				fmt.Printf("%-*s", maxWidths[*pe.Key], *pe.Key)
				if pe.Next != nil {
					fmt.Print(verticalMiddle)
				} else {
					fmt.Println(verticalEnd)
				}
			}

			fmt.Print(horizontalStart)
			for pe := outrec.Head; pe != nil; pe = pe.Next {
				fmt.Print(horizontalBars[*pe.Key])
				if pe.Next != nil {
					fmt.Print(horizontalMiddle)
				} else {
					fmt.Println(horizontalEnd)
				}
			}

		}
		onFirst = false

		// Print data lines
		fmt.Print(verticalStart)
		for pe := outrec.Head; pe != nil; pe = pe.Next {
			s := pe.Value.String()
			fmt.Printf("%-*s", maxWidths[*pe.Key], s)
			if pe.Next != nil {
				fmt.Print(verticalMiddle)
			} else {
				fmt.Println(verticalEnd)
			}
		}

		if e.Next() == nil {
			fmt.Print(horizontalStart)
			for pe := outrec.Head; pe != nil; pe = pe.Next {
				fmt.Print(horizontalBars[*pe.Key])
				if pe.Next != nil {
					fmt.Print(horizontalMiddle)
				} else {
					fmt.Println(horizontalEnd)
				}
			}
		}

	}
}
