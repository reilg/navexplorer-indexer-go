package explorer

type TransferType string

var (
	TransferSend                TransferType = "send"
	TransferReceive             TransferType = "receive"
	TransferStake               TransferType = "stake"
	TransferDelegateStake       TransferType = "delegate_stake"
	TransferPoolStake           TransferType = "pool_stake"
	TransferPoolFee             TransferType = "pool_fee"
	TransferCommunityFund       TransferType = "community_fund"
	TransferCommunityFundPayout TransferType = "community_fund_payout"
)
