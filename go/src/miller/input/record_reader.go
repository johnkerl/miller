package input

import (
	"miller/lib"
)

// Since Go is concurrent, the context struct (AWK-like variables such as
// FILENAME, NF, NF, FNR, etc.) needs to be duplicated and passed through the
// channels along with each record. Hence the initial context, which readers
// update on each new file/record, and the channel of lib.RecordAndContext
// rather than channel of lib.Mlrmap.

type IRecordReader interface {
	Read(
		filenames []string,
		initialContext lib.Context,
		inputChannel chan<- *lib.RecordAndContext,
		errorChannel chan error,
	)
}
