package input

import (
	"regexp"

	"github.com/johnkerl/miller/v6/pkg/cli"
)

// NewRecordReaderBox reads Unicode box-drawing tables, e.g.:
//
// ┌─────┬─────┬────┬─────────────────────┬─────────────────────┐
// │ a   │ b   │ i  │ x                   │ y                   │
// ├─────┼─────┼────┼─────────────────────┼─────────────────────┤
// │ pan │ pan │ 1  │ 0.3467901443380824  │ 0.7268028627434533  │
// ├─────┼─────┼────┼─────────────────────┼─────────────────────┤
// │ eks │ pan │ 2  │ 0.7586799647899636  │ 0.5221511083334797  │
// └─────┴─────┴────┴─────────────────────┴─────────────────────┘
//
// This is structurally the same as barred-PPRINT/Markdown input: border/divider
// lines are recognized and skipped wherever they occur, so a divider after
// every row (not just the header) already works with the shared
// RecordReaderPprintBarredOrMarkdown parsing logic -- only the separator
// regex and IFS differ.
func NewRecordReaderBox(
	readerOptions *cli.TReaderOptions,
	recordsPerBatch int64,
) (IRecordReader, error) {

	readerOptions.IFS = "│" // U+2502 BOX DRAWINGS LIGHT VERTICAL
	readerOptions.AllowRepeatIFS = false

	reader := &RecordReaderPprintBarredOrMarkdown{
		readerOptions:     readerOptions,
		recordsPerBatch:   recordsPerBatch,
		separatorMatcher:  regexp.MustCompile(`^[┌┬┐├┼┤└┴┘─]+$`),
		fieldSplitter:     newFieldSplitter(readerOptions),
		formatDisplayName: "box",
	}
	if reader.readerOptions.UseImplicitHeader {
		reader.recordBatchGetter = getRecordBatchImplicitPprintHeader
	} else {
		reader.recordBatchGetter = getRecordBatchExplicitPprintHeader
	}
	return reader, nil
}
