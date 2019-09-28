package explorer

type TransferType string

var (
	TransferSend                TransferType = "send"
	TransferReceive             TransferType = "receive"
	TransferStake               TransferType = "stake"
	TransferColdStake           TransferType = "cold_stake"
	TransferPoolStake           TransferType = "pool_stake"
	TransferPoolFee             TransferType = "pool_fee"
	TransferDelegateStake       TransferType = "delegate_stake"
	TransferCommunityFund       TransferType = "community_fund"
	TransferCommunityFundPayout TransferType = "community_fund_payout"
)
