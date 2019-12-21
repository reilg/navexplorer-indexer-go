package explorer

type RawPaymentRequest struct {
	MetaData MetaData `json:"-"`

	Version             uint32               `json:"version"`
	Hash                string               `json:"hash"`
	BlockHash           string               `json:"blockHash"`
	ProposalHash        string               `json:"proposalHash,omitempty"`
	Description         string               `json:"description"`
	RequestedAmount     float64              `json:"requestedAmount"`
	Status              PaymentRequestStatus `json:"status"`
	State               uint                 `json:"state"`
	StateChangedOnBlock string               `json:"stateChangedOnBlock,omitempty"`
	PaidOnBlock         string               `json:"paidOnBlock,omitempty"`
}

type PaymentRequest struct {
	RawPaymentRequest
	Height         uint64 `json:"height"`
	UpdatedOnBlock uint64 `json:"height"`
}
