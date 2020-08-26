package input

import (
	"containers"
)

// reader *bufio.Reader in constructors/factory
type RecordReader interface {
	Read(inrecs chan<- *containers.Lrec, echan chan error)
}
