package output

import (
	"encoding/csv"
	"os"

	"miller/clitypes"
	"miller/lib"
)

// ostream *os.File in constructors/factory
type RecordWriterCSV struct {
	onFirst   bool
	csvWriter *csv.Writer
}

func NewRecordWriterCSV(writerOptions *clitypes.TWriterOptions) *RecordWriterCSV {
	return &RecordWriterCSV{
		onFirst:   true,
		csvWriter: csv.NewWriter(os.Stdout),
	}
}

func (this *RecordWriterCSV) Write(
	outrec *lib.Mlrmap,
) {
	// End of record stream: nothing special for this output format
	if outrec == nil {
		return
	}

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
		fields[i] = pe.Value.String()
		i++
	}
	this.csvWriter.Write(fields)
	this.csvWriter.Flush()
}
