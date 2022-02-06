package output

import (
	"fmt"

	"github.com/johnkerl/miller/internal/pkg/cli"
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
	case "tsv":
		return NewRecordWriterTSV(writerOptions)
	case "xtab":
		return NewRecordWriterXTAB(writerOptions)
	default:
		return nil, fmt.Errorf("output file format \"%s\" not found", writerOptions.OutputFileFormat)
	}
}
