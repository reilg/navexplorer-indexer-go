package explorer

type PaymentRequestVote struct {
	Height  uint64 `json:"height"`
	Address string `json:"address"`
	Votes   []Vote `json:"votes"`
}
