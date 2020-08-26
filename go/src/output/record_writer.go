package output

import (
	"containers"
)

// ostream *os.File in constructors/factory
type RecordWriter interface {
	Writer(outrecs <-chan *containers.Lrec)
}
