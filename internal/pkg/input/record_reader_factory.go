package input

import (
	"errors"
	"fmt"

	"mlr/internal/pkg/cli"
)

func Create(readerOptions *cli.TReaderOptions) (IRecordReader, error) {
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
	case "pprint":
		return NewRecordReaderPPRINT(readerOptions)
	case "xtab":
		return NewRecordReaderXTAB(readerOptions)
	case "gen":
		return NewPseudoReaderGen(readerOptions)
	default:
		return nil, errors.New(fmt.Sprintf("input file format \"%s\" not found", readerOptions.InputFileFormat))
	}
}
