package input

import (
	"miller/lib"
)

// Since Go is concurrent, the context struct (AWK-like variables such as
// FILENAME, NF, NF, FNR, etc.) needs to be duplicated and passed through the
// channels along with each record. Hence the initial context, which readers
// update on each new file/record, and the channel of lib.LrecAndContext
// rather than channel of lib.Lrec.

type IRecordReader interface {
	Read(
		filenames []string,
		initialContext lib.Context,
		inrecsAndContexts chan<- *lib.LrecAndContext,
		echan chan error,
	)
}
