package output

import (
	"containers"
)

// ostream *os.File in constructors/factory
type RecordWriter interface {
	Write(outrec *containers.Lrec)
}
