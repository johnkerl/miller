package input

import (
	"miller/containers"
	"miller/runtime"
)

type IRecordReader interface {
	Read(
		filenames []string,
		context *runtime.Context,
		inrecs chan<- *containers.Lrec,
		echan chan error,
	)
}
