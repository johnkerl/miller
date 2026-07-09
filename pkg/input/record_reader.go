// This file contains the interface for file-format-specific record-readers, as
// well as a collection of utility functions.

package input

import "github.com/johnkerl/miller/v6/pkg/types"

// Since Go is concurrent, the context struct (AWK-like variables such as
// FILENAME, NF, NF, FNR, etc.) needs to be duplicated and passed through the
// channels along with each record. Hence the initial context, which readers
// update on each new file/record, and the channel of types.RecordAndContext
// rather than channel of mlrval.Mlrmap.

type IRecordReader interface {
	Read(
		filenames []string,
		initialContext types.Context,
		readerChannel chan<- []*types.RecordAndContext, // list of *types.RecordAndContext
		errorChannel chan error,
		downstreamDoneChannel <-chan bool, // for mlr head
	)
}

// hasNonEmptyField returns true if any of the split-out fields is non-empty.
// A blank input line splits to zero fields or to a single empty field; a line
// consisting only of field separators splits to all-empty fields. Used by the
// CSV/TSV readers to support skipping trivial records at read time, when the
// skip-trivial-records verb is present in the then-chain. See issue #1535.
func hasNonEmptyField(fields []string) bool {
	for _, field := range fields {
		if field != "" {
			return true
		}
	}
	return false
}
