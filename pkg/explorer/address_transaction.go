package explorer

type AddressTransaction struct {
	Hash   string       `json:"hash"`
	Txid   string       `json:"txid"`
	Height uint64       `json:"height"`
	Time   int64        `json:"time, omitempty"`
	Type   TransferType `json:"type"`
	Input  uint64       `json:"input"`
	Output uint64       `json:"output"`
	Total  int64        `json:"total"`
	Cold   bool         `json:"cold"`
}
