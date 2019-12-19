package explorer

type MetaData struct {
	Id    string
	Index string
}

func NewMetaData(id string, index string) MetaData {
	return MetaData{id, index}
}
