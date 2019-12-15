package explorer

type Proposal struct {
	MetaData MetaData `json:"-"`

	Version             uint32         `json:"version"`
	Hash                string         `json:"hash"`
	BlockHash           string         `json:"blockHash"`
	Description         string         `json:"description"`
	RequestedAmount     float64        `json:"requestedAmount"`
	NotPaidYet          float64        `json:"notPaidYet"`
	NotRequestedYet     float64        `json:"notRequestedYet"`
	UserPaidFee         float64        `json:"userPaidFee"`
	PaymentAddress      string         `json:"paymentAddress"`
	ProposalDuration    uint64         `json:"proposalDuration"`
	ExpiresOn           uint64         `json:"expiresOn"`
	Status              ProposalStatus `json:"status"`
	State               uint           `json:"state"`
	StateChangedOnBlock string         `json:"stateChangedOnBlock,omitempty"`

	// Custom
	Height         uint64 `json:"height"`
	UpdatedOnBlock uint64 `json:"updatedOnBlock"`
}

type ProposalStatus string

var (
	PROPOSAL_PENDING  ProposalStatus = "pending"
	PROPOSAL_ACCEPTED ProposalStatus = "accepted"
	PROPOSAL_REJECTED ProposalStatus = "rejected"
	PROPOSAL_EXPIRED  ProposalStatus = "expired"
)
