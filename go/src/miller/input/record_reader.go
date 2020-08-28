package input

import (
	"miller/containers"
)

// reader *bufio.Reader in constructors/factory
type IRecordReader interface {
	Read(filenames []string, inrecs chan<- *containers.Lrec, echan chan error)
}
