package explorer

type Proposal struct {
	MetaData MetaData `json:"-"`

	Version             uint32  `json:"version"`
	Hash                string  `json:"hash"`
	BlockHash           string  `json:"blockHash"`
	Description         string  `json:"description"`
	RequestedAmount     float64 `json:"requestedAmount"`
	NotPaidYet          float64 `json:"notPaidYet"`
	UserPaidFee         float64 `json:"userPaidFee"`
	PaymentAddress      string  `json:"paymentAddress"`
	ProposalDuration    uint64  `json:"proposalDuration"`
	ExpiresOn           uint64  `json:"expiresOn"`
	Status              string  `json:"status"`
	State               uint    `json:"state"`
	StateChangedOnBlock string  `json:"stateChangedOnBlock,omitempty"`

	// Custom
	Height uint64      `json:"height"`
	Cycles CfundCycles `json:"cycles"`
}

type CfundCycles []CfundCycle

type CfundCycle struct {
	VotingCycle uint `json:"votingCycle"`
	VotesYes    uint `json:"votesYes"`
	VotesNo     uint `json:"votesNo"`
}
