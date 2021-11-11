package output

import (
	"errors"
	"fmt"

	"mlr/internal/pkg/cli"
)

func Create(writerOptions *cli.TWriterOptions) (IRecordWriter, error) {
	switch writerOptions.OutputFileFormat {
	case "csv":
		return NewRecordWriterCSV(writerOptions)
	case "csvlite":
		return NewRecordWriterCSVLite(writerOptions)
	case "dkvp":
		return NewRecordWriterDKVP(writerOptions)
	case "json":
		return NewRecordWriterJSON(writerOptions)
	case "markdown":
		return NewRecordWriterMarkdown(writerOptions)
	case "nidx":
		return NewRecordWriterNIDX(writerOptions)
	case "pprint":
		return NewRecordWriterPPRINT(writerOptions)
	case "xtab":
		return NewRecordWriterXTAB(writerOptions)
	default:
		return nil, errors.New(fmt.Sprintf("output file format \"%s\" not found", writerOptions.OutputFileFormat))
	}
}
