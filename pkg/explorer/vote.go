package explorer

type DaoVote struct {
	MetaData MetaData `json:"-"`

	Height  uint64 `json:"height"`
	Address string `json:"address"`
	Votes   []Vote `json:"votes"`
}

type Vote struct {
	Type VoteType `json:"type"`
	Hash string   `json:"hash"`
	Vote int      `json:"vote"`
}

type VoteType string

var (
	ProposalVote       VoteType = "Proposal"
	PaymentRequestVote VoteType = "PaymentRequest"
)
