package input

import (
	"regexp"

	"github.com/johnkerl/miller/pkg/cli"
)

func NewRecordReaderMarkdown(
	readerOptions *cli.TReaderOptions,
	recordsPerBatch int64,
) (IRecordReader, error) {

	readerOptions.IFS = "|"
	readerOptions.AllowRepeatIFS = false

	reader := &RecordReaderPprintBarred{
		readerOptions:    readerOptions,
		recordsPerBatch:  recordsPerBatch,
		//separatorMatcher: regexp.MustCompile(`^|[ -|]*|$`),
		//separatorMatcher: regexp.MustCompile(`^|[ -|]+|$`),
		//separatorMatcher: regexp.MustCompile(`^X$`),
		//separatorMatcher: regexp.MustCompile(`^\|$`),
		separatorMatcher: regexp.MustCompile(`^\|[ -\|]+\|$`),
		fieldSplitter:    newFieldSplitter(readerOptions),
	}
	if reader.readerOptions.UseImplicitHeader {
		reader.recordBatchGetter = getRecordBatchImplicitPprintHeader
	} else {
		reader.recordBatchGetter = getRecordBatchExplicitPprintHeader
	}
	return reader, nil

}
