package output

import (
	"io"

	"mlr/internal/pkg/types"
)

// IRecordWriter is the abstract interface for all record-writers.  They are
// specific to output file format, e.g. CSV, PPRINT, JSON, etc.  Convention:
// nil outrec signifies end of record stream.
type IRecordWriter interface {
	Write(outrec *types.Mlrmap, ostream io.WriteCloser, outputIsStdout bool)
}
