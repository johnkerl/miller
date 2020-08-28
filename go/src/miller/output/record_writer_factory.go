package output

func Create(outputFormatName string) IRecordWriter {
	switch outputFormatName {
	case "csv":
		return NewRecordWriterCSV() // TODO: parameterize
	case "dkvp":
		return NewRecordWriterDKVP(",", "=") // TODO: parameterize
	case "json":
		return NewRecordWriterJSON() // TODO: parameterize
	case "nidx":
		return NewRecordWriterNIDX(",") // TODO: parameterize
	case "xtab":
		return NewRecordWriterXTAB() // TODO: parameterize
	default:
		return nil
	}
}
