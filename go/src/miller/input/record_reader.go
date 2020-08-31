package input

import (
	"miller/runtime"
)

type IRecordReader interface {
	Read(
		filenames []string,
		initialContext runtime.Context,
		inrecsAndContexts chan<- *runtime.LrecAndContext,
		echan chan error,
	)
}
