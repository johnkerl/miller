package output

import (
	"miller/containers"
)

// Convention: nil outrec signifies end of record stream
type RecordWriter interface {
	Write(outrec *containers.Lrec)
}
