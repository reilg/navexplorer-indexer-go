package explorer

type ProposalVote struct {
	Height   uint64 `json:"height"`
	Address  string `json:"address"`
	Proposal string `json:"proposal"`
	Vote     bool   `json:"vote"`
}
