package input

func Create(inputFormatName string) RecordReader {
	switch inputFormatName {
	case "dkvp":
		return NewRecordReaderDKVP(",", "=") // TODO: parameterize
	default:
		return nil
	}
}
