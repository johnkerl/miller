package output

import (
	"miller/types"
)

// Convention: nil outrec signifies end of record stream
type IRecordWriter interface {
	Write(outrec *types.Mlrmap)
}
