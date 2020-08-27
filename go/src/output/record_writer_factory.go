package output
func Create(outputFormatName string) RecordWriter {
	switch outputFormatName {
	case "dkvp": return NewRecordWriterDKVP(",", "=") // TODO: parameterize
	case "nidx": return NewRecordWriterNIDX(",") // TODO: parameterize
	default: return nil
	}
}
