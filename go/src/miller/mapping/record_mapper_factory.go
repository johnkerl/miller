package mapping

func Create(mapperName string) RecordMapper {
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
