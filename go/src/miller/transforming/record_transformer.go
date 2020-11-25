package transforming

import (
	"flag"

	"miller/clitypes"
	"miller/types"
)

type IRecordTransformer interface {
	Map(
		inrecAndContext *types.RecordAndContext,
		outputChannel chan<- *types.RecordAndContext,
	)
}

type RecordTransformerFunc func(
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
) IRecordTransformer

type TransformerSetup struct {
	Verb         string
	ParseCLIFunc MapperParseCLIFunc
	IgnoresInput bool
}
