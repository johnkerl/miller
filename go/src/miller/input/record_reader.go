package input

import (
	"miller/containers"
)

type IRecordReader interface {
	Read(
		filenames []string,
		initialContext containers.Context,
		inrecsAndContexts chan<- *containers.LrecAndContext,
		echan chan error,
	)
}
