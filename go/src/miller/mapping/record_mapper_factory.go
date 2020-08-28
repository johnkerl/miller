package mapping

func Create(mapperName string) IRecordMapper {
	switch mapperName {
	case "cat":
		return NewMapperCat()
	case "nothing":
		return NewMapperNothing()
	case "tac":
		return NewMapperTac()
	default:
		return nil
	}
}
