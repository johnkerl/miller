package input

import (
	"miller/containers"
	"miller/runtime"
)

// reader *bufio.Reader in constructors/factory
type IRecordReader interface {
	Read(
		filenames []string,
		context *runtime.Context,
		inrecs chan<- *containers.Lrec,
		echan chan error,
	)
}
