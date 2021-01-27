package transforming

import (
	"flag"
	"os"

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
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError // TODO: remove
	readerOptions *clitypes.TReaderOptions,
	writerOptions *clitypes.TWriterOptions,
) IRecordTransformer

type TransformerUsageFunc func(
	ostream *os.File,
	doExit bool,
	exitCode int,
)

type TransformerSetup struct {
	Verb         string
	ParseCLIFunc TransformerParseCLIFunc
	UsageFunc    TransformerUsageFunc
	IgnoresInput bool
}
