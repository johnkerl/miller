package output

import (
	"io"

	"mlr/src/types"
)

// Convention: nil outrec signifies end of record stream
type IRecordWriter interface {
	Write(outrec *types.Mlrmap, ostream io.WriteCloser, outputIsStdout bool)
}
