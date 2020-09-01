package mapping

import (
	"flag"

	"miller/clitypes"
	"miller/containers"
)

type IRecordMapper interface {
	Map(
		inrecAndContext *containers.LrecAndContext,
		outrecsAndContexts chan<- *containers.LrecAndContext,
	)
}

type MapperParseCLIFunc func(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	readerOptions *clitypes.TReaderOptions,
	writerOptions *clitypes.TWriterOptions,
) IRecordMapper

type MapperSetup struct {
	Verb         string
	ParseCLIFunc MapperParseCLIFunc
	IgnoresInput bool
}
