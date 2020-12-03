package transforming

import (
	"flag"

	"miller/clitypes"
	"miller/types"
)

type IRecordTransformer interface {
	Transform(
		inrecAndContext *types.RecordAndContext,
		outputChannel chan<- *types.RecordAndContext,
	)
}

type RecordTransformerFunc func(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
)

type TransformerParseCLIFunc func(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	readerOptions *clitypes.TReaderOptions,
	writerOptions *clitypes.TWriterOptions,
) IRecordTransformer

type TransformerSetup struct {
	Verb         string
	ParseCLIFunc TransformerParseCLIFunc
	IgnoresInput bool
}
