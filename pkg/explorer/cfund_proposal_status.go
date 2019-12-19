package explorer

type ProposalStatus string

var (
	ProposalPending  ProposalStatus = "pending"
	ProposalAccepted ProposalStatus = "accepted"
	ProposalRejected ProposalStatus = "rejected"
	ProposalExpired  ProposalStatus = "expired"
)

func ProposalStatusIsValid(status string) bool {
	switch true {
	case status == string(ProposalPending):
		return true
	case status == string(ProposalAccepted):
		return true
	case status == string(ProposalRejected):
		return true
	case status == string(ProposalExpired):
		return true
	}
	return false
}
