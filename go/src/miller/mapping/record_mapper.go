package mapping

import (
	"flag"

	"miller/clitypes"
	"miller/lib"
)

type IRecordMapper interface {
	Map(
		inrecAndContext *lib.RecordAndContext,
		outputChannel chan<- *lib.RecordAndContext,
	)
}

type RecordMapperFunc func(
	inrecAndContext *lib.RecordAndContext,
	outputChannel chan<- *lib.RecordAndContext,
)

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
