package explorer

type PaymentRequestVote struct {
	MetaData MetaData `json:"-"`

	Height  uint64 `json:"height"`
	Address string `json:"address"`
	Votes   []Vote `json:"votes"`
}
