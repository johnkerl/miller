package input

func Create(inputFormatName string) RecordReader {
	switch inputFormatName {
	case "csv":
		return NewRecordReaderCSV() // TODO: parameterize
	case "dkvp":
		return NewRecordReaderDKVP(",", "=") // TODO: parameterize
	case "json":
		return NewRecordReaderJSON()
	case "nidx":
		return NewRecordReaderNIDX() // TODO: parameterize
	default:
		return nil
	}
}
