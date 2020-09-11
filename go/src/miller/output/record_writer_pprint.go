package output

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"miller/clitypes"
	"miller/lib"
)

// ostream *os.File in constructors/factory
type RecordWriterPPRINT struct {
	records *list.List
	// For detecting schema changes: we print a newline and the new header
}

func NewRecordWriterPPRINT(writerOptions *clitypes.TWriterOptions) *RecordWriterPPRINT {
	return &RecordWriterPPRINT{
		records: list.New(),
	}
}

// ----------------------------------------------------------------
// xxx this is very naive at present -- needs copy from the C version.
func (this *RecordWriterPPRINT) Write(
	outrec *lib.Mlrmap,
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
		record := head.Value.(*lib.Mlrmap)
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
				this.writeHeterogenousList(batch)
				os.Stdout.WriteString("\n")
				// Print a newline
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
		this.writeHeterogenousList(batch)
	}
}

// ----------------------------------------------------------------
func (this *RecordWriterPPRINT) writeHeterogenousList(
	records *list.List,
) {

	// TODO: heterogeneity. keep previous header and reset if need.
	maxWidths := make(map[string]int)

	for e := records.Front(); e != nil; e = e.Next() {
		outrec := e.Value.(*lib.Mlrmap)
		for pe := outrec.Head; pe != nil; pe = pe.Next {
			width := len(pe.Value.String())
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

	onFirst := true
	for e := records.Front(); e != nil; e = e.Next() {
		outrec := e.Value.(*lib.Mlrmap)

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
			if pe.Next != nil {
				fmt.Printf("%-*s ", maxWidths[*pe.Key], pe.Value.String())
			} else {
				fmt.Println(pe.Value.String())
			}
		}
	}
}
