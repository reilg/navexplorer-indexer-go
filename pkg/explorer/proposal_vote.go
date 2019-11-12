package explorer

type ProposalVote struct {
	Height  uint64 `json:"height"`
	Address string `json:"address"`
	Votes   []Vote `json:"votes"`
}

type Vote struct {
	Hash string `json:"hash"`
	Vote bool   `json:"vote"`
}
