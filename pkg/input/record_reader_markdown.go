package input

import (
	"regexp"

	"github.com/johnkerl/miller/v6/pkg/cli"
)

func NewRecordReaderMarkdown(
	readerOptions *cli.TReaderOptions,
	recordsPerBatch int64,
) (IRecordReader, error) {

	readerOptions.IFS = "|"
	readerOptions.AllowRepeatIFS = false

	reader := &RecordReaderPprintBarredOrMarkdown{
		readerOptions:     readerOptions,
		recordsPerBatch:   recordsPerBatch,
		separatorMatcher:  regexp.MustCompile(`^\|[-\| ]+\|$`),
		fieldSplitter:     newFieldSplitter(readerOptions),
		formatDisplayName: "Markdown",
	}
	if reader.readerOptions.UseImplicitHeader {
		reader.recordBatchGetter = getRecordBatchImplicitPprintHeader
	} else {
		reader.recordBatchGetter = getRecordBatchExplicitPprintHeader
	}
	return reader, nil

}
