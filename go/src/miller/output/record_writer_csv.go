package output

import (
	// System:
	"encoding/csv"
	"os"
	// Miller:
	"miller/containers"
)

// ostream *os.File in constructors/factory
type RecordWriterCSV struct {
	onFirst bool
	csvWriter *csv.Writer
}

func NewRecordWriterCSV() *RecordWriterCSV {
	return &RecordWriterCSV{
		true,
		csv.NewWriter(os.Stdout),
	}
}

func (this *RecordWriterCSV) Write(
	outrec *containers.Lrec,
) {
	// TODO: heterogeneity. keep previous header and reset if need.

	if this.onFirst {
		fields := make([]string, outrec.FieldCount)
		i := 0
		for pe := outrec.Head; pe != nil; pe = pe.Next {
			fields[i] = *pe.Key
			i++
		}
		this.csvWriter.Write(fields)

		this.onFirst = false
	}

	fields := make([]string, outrec.FieldCount)
	i := 0
	for pe := outrec.Head; pe != nil; pe = pe.Next {
		fields[i] = *pe.Value
		i++
	}
	this.csvWriter.Write(fields)
	this.csvWriter.Flush()
}
