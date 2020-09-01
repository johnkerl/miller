package input

import (
	"miller/cli"
)

func Create(readerOptions *cli.TReaderOptions) IRecordReader {
	switch readerOptions.InputFileFormat {
	case "csv":
		return NewRecordReaderCSV(readerOptions)
	case "dkvp":
		return NewRecordReaderDKVP(readerOptions)
	case "json":
		return NewRecordReaderJSON(readerOptions)
	case "nidx":
		return NewRecordReaderNIDX(readerOptions)
	default:
		return nil
	}
}
