package mapping

import (
	"os"

	"miller/clitypes"
	"miller/containers"
	"miller/runtime"
)

type IRecordMapper interface {
	Map(
		inrec *containers.Lrec,
		context *runtime.Context,
		outrecs chan<- *containers.Lrec,
	)
}

type MapperParseCLIFunc func(
	pargi *int,
	argc int,
	args []string,
	readerOptions *clitypes.TReaderOptions,
	writerOptions *clitypes.TWriterOptions,
) IRecordMapper

type MapperUsageFunc func(
	ostream *os.File,
	argv0 string,
	verb string,
)

type MapperSetup struct {
	Verb         string
	ParseCLIFunc MapperParseCLIFunc
	UsageFunc    MapperUsageFunc
	IgnoresInput bool
}
