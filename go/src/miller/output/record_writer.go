package output

import (
	"miller/lib"
)

// Convention: nil outrec signifies end of record stream
type IRecordWriter interface {
	Write(outrec *types.Mlrmap)
}
