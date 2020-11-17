package output

import (
	"miller/clitypes"
)

func Create(writerOptions *clitypes.TWriterOptions) IRecordWriter {
	switch writerOptions.OutputFileFormat {
	case "csv":
		return NewRecordWriterCSV(writerOptions)
	case "csvlite":
		return NewRecordWriterCSVLite(writerOptions)
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
