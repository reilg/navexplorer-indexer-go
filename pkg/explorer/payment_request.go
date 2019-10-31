package explorer

type PaymentRequest struct {
	Version             uint32 `json:"version"`
	Hash                string `json:"hash"`
	BlockHash           string `json:"blockHash"`
	Description         string `json:"description"`
	RequestedAmount     uint64 `json:"requestedAmount"`
	Status              string `json:"status"`
	State               uint   `json:"state"`
	StateChangedOnBlock string `json:"stateChangedOnBlock,omitempty"`
	PaidOnBlock         string `json:"paidOnBlock,omitempty"`

	// Custom
	Height uint64      `json:"height"`
	Cycles CfundCycles `json:"cycles"`
}
