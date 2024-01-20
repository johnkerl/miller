package output

import (
	"bufio"

	"github.com/johnkerl/miller/pkg/mlrval"
)

// IRecordWriter is the abstract interface for all record-writers.  They are
// specific to output file format, e.g. CSV, PPRINT, JSON, etc.
//
// Convention: nil outrec signifies end of record stream.
//
// The ChannelWriter will call bufferedOutputStream.Flush() after every record
// if the --fflush flag (writerOptions.FlushOnEveryRecord) is present, so each
// writer does not have to -- unless the writer retains records e.g. for PPRINT
// format.
type IRecordWriter interface {
	Write(
		outrec *mlrval.Mlrmap,
		bufferedOutputStream *bufio.Writer,
		outputIsStdout bool,
	) error
}
