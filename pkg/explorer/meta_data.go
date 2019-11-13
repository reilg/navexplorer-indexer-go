package explorer

type MetaData struct {
	Id    string
	Index string
	Dirty bool
}

func NewMetaData(id string, index string) MetaData {
	return MetaData{id, index, false}
}
