// This file contains the interface for file-format-specific record-readers, as
// well as a collection of utility functions.

package input

import (
	"container/list"

	"github.com/johnkerl/miller/pkg/types"
)

// Since Go is concurrent, the context struct (AWK-like variables such as
// FILENAME, NF, NF, FNR, etc.) needs to be duplicated and passed through the
// channels along with each record. Hence the initial context, which readers
// update on each new file/record, and the channel of types.RecordAndContext
// rather than channel of mlrval.Mlrmap.

type IRecordReader interface {
	Read(
		filenames []string,
		initialContext types.Context,
		readerChannel chan<- *list.List, // list of *types.RecordAndContext
		errorChannel chan error,
		downstreamDoneChannel <-chan bool, // for mlr head
	)
}
