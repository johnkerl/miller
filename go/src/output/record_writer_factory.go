package output

import (
	"miller/src/cliutil"
)

func Create(writerOptions *cliutil.TWriterOptions) IRecordWriter {
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
		return nil
	}
}
