package input

import (
	"miller/cliutil"
)

func Create(readerOptions *cliutil.TReaderOptions) IRecordReader {
	switch readerOptions.InputFileFormat {
	case "csv":
		return NewRecordReaderCSV(readerOptions)
	case "csvlite":
		return NewRecordReaderCSVLite(readerOptions)
	case "dkvp":
		return NewRecordReaderDKVP(readerOptions)
	case "json":
		return NewRecordReaderJSON(readerOptions)
	case "nidx":
		return NewRecordReaderNIDX(readerOptions)
	case "xtab":
		return NewRecordReaderXTAB(readerOptions)
	default:
		return nil
	}
}
