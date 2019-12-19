package explorer

type PaymentRequestStatus string

var (
	PaymentRequestPending  PaymentRequestStatus = "pending"
	PaymentRequestAccepted PaymentRequestStatus = "accepted"
	PaymentRequestRejected PaymentRequestStatus = "rejected"
	PaymentRequestExpired  PaymentRequestStatus = "expired"
)

func PaymentRequestStatusIsValid(status string) bool {
	switch true {
	case status == string(PaymentRequestPending):
		return true
	case status == string(PaymentRequestAccepted):
		return true
	case status == string(PaymentRequestRejected):
		return true
	case status == string(PaymentRequestExpired):
		return true
	}
	return false
}
