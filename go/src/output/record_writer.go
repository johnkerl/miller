package output

import (
	"io"

	"miller/src/types"
)

// Convention: nil outrec signifies end of record stream
type IRecordWriter interface {
	Write(outrec *types.Mlrmap, ostream io.WriteCloser)
}
