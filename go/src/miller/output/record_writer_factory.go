package output

import (
	"miller/cli"
)

func Create(writerOptions *cli.TWriterOptions) IRecordWriter {
	switch writerOptions.OutputFileFormat {
	case "csv":
		return NewRecordWriterCSV(writerOptions)
	case "dkvp":
		return NewRecordWriterDKVP(writerOptions)
	case "json":
		return NewRecordWriterJSON(writerOptions)
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
