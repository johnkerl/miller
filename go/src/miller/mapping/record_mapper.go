package mapping

import (
	"flag"

	"miller/clitypes"
	"miller/types"
)

type IRecordMapper interface {
	Map(
		inrecAndContext *types.RecordAndContext,
		outputChannel chan<- *types.RecordAndContext,
	)
}

type RecordMapperFunc func(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
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
