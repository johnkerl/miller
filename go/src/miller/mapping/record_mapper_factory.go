package mapping

func Create(mapperName string, dslString string) (IRecordMapper, error) {
	switch mapperName {
	case "cat":
		return NewMapperCat()
	case "check":
		return NewMapperNothing()
	case "nothing":
		return NewMapperNothing()
	case "put":
		return NewMapperPut(dslString)
	case "tac":
		return NewMapperTac()
	default:
		return nil, nil
	}
}
