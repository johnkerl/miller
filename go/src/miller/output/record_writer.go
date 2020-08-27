package output

import (
	"miller/containers"
)

// ostream *os.File in constructors/factory
type RecordWriter interface {
	Write(outrec *containers.Lrec)
}
