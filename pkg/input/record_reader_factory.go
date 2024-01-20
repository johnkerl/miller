package input

import (
	"fmt"

	"github.com/johnkerl/miller/pkg/cli"
)

func Create(readerOptions *cli.TReaderOptions, recordsPerBatch int64) (IRecordReader, error) {
	switch readerOptions.InputFileFormat {
	case "csv":
		return NewRecordReaderCSV(readerOptions, recordsPerBatch)
	case "csvlite":
		return NewRecordReaderCSVLite(readerOptions, recordsPerBatch)
	case "dkvp":
		return NewRecordReaderDKVP(readerOptions, recordsPerBatch)
	case "json":
		return NewRecordReaderJSON(readerOptions, recordsPerBatch)
	case "nidx":
		return NewRecordReaderNIDX(readerOptions, recordsPerBatch)
	case "markdown":
		return NewRecordReaderMarkdown(readerOptions, recordsPerBatch)
	case "pprint":
		return NewRecordReaderPPRINT(readerOptions, recordsPerBatch)
	case "tsv":
		return NewRecordReaderTSV(readerOptions, recordsPerBatch)
	case "xtab":
		return NewRecordReaderXTAB(readerOptions, recordsPerBatch)
	case "gen":
		return NewPseudoReaderGen(readerOptions, recordsPerBatch)
	default:
		return nil, fmt.Errorf("input file format \"%s\" not found", readerOptions.InputFileFormat)
	}
}
