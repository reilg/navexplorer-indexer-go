package explorer

type VoutType string

var (
	VoutNonstandard           VoutType = "nonstandard"
	VoutNulldata              VoutType = "nulldata"
	VoutPubkeyhash            VoutType = "pubkeyhash"
	VoutPubkey                VoutType = "pubkey"
	VoutScripthash            VoutType = "scripthash"
	VoutColdStaking           VoutType = "cold_staking"
	VoutCfundContribution     VoutType = "cfund_contribution"
	VoutProposalNoVote        VoutType = "proposal_no_vote"
	VoutProposalYesVote       VoutType = "proposal_yes_vote"
	VoutPaymentRequestNoVote  VoutType = "payment_request_no_vote"
	VoutPaymentRequestYesVote VoutType = "payment_request_yes_vote"
	VoutPoolStaking           VoutType = "pool_staking"
)
