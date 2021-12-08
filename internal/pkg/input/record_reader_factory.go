package input

import (
	"errors"
	"fmt"

	"github.com/johnkerl/miller/internal/pkg/cli"
)

func Create(readerOptions *cli.TReaderOptions, recordsPerBatch int) (IRecordReader, error) {
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
	case "pprint":
		return NewRecordReaderPPRINT(readerOptions, recordsPerBatch)
	case "xtab":
		return NewRecordReaderXTAB(readerOptions, recordsPerBatch)
	case "gen":
		return NewPseudoReaderGen(readerOptions, recordsPerBatch)
	default:
		return nil, errors.New(fmt.Sprintf("input file format \"%s\" not found", readerOptions.InputFileFormat))
	}
}
