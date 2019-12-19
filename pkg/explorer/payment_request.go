package explorer

type PaymentRequest struct {
	MetaData MetaData `json:"-"`

	Version             uint32               `json:"version"`
	Hash                string               `json:"hash"`
	BlockHash           string               `json:"blockHash"`
	ProposalHash        string               `json:"proposalHash"`
	Description         string               `json:"description"`
	RequestedAmount     float64              `json:"requestedAmount"`
	Status              PaymentRequestStatus `json:"status"`
	State               uint                 `json:"state"`
	StateChangedOnBlock string               `json:"stateChangedOnBlock,omitempty"`
	PaidOnBlock         string               `json:"paidOnBlock,omitempty"`

	// Custom
	Height         uint64 `json:"height"`
	UpdatedOnBlock uint64 `json:"height"`
}

type PaymentRequestStatus string

var (
	PAYMENT_REQUEST_PENDING  PaymentRequestStatus = "pending"
	PAYMENT_REQUEST_ACCEPTED PaymentRequestStatus = "accepted"
	PAYMENT_REQUEST_REJECTED PaymentRequestStatus = "rejected"
	PAYMENT_REQUEST_EXPIRED  PaymentRequestStatus = "expired"
)

func PaymentRequestStatusIsValid(status string) bool {
	switch true {
	case status == string(PAYMENT_REQUEST_PENDING):
		return true
	case status == string(PAYMENT_REQUEST_ACCEPTED):
		return true
	case status == string(PAYMENT_REQUEST_REJECTED):
		return true
	case status == string(PAYMENT_REQUEST_EXPIRED):
		return true
	}
	return false
}
