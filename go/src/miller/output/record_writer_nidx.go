package output

import (
	"bytes"
	"os"

	"miller/clitypes"
	"miller/types"
)

// ostream *os.File in constructors/factory
type RecordWriterNIDX struct {
	ofs string
	ors string
}

func NewRecordWriterNIDX(writerOptions *clitypes.TWriterOptions) *RecordWriterNIDX {
	return &RecordWriterNIDX{
		ofs: writerOptions.OFS,
		ors: writerOptions.ORS,
	}
}

func (this *RecordWriterNIDX) Write(
	outrec *types.Mlrmap,
) {
	// End of record stream: nothing special for this output format
	if outrec == nil {
		return
	}

	var buffer bytes.Buffer // 5x faster than fmt.Print() separately
	for pe := outrec.Head; pe != nil; pe = pe.Next {
		buffer.WriteString(pe.Value.String())
		if pe.Next != nil {
			buffer.WriteString(this.ofs)
		}
	}
	buffer.WriteString(this.ors)
	os.Stdout.WriteString(buffer.String())
}
