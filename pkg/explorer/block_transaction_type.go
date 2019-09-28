package explorer

type BlockTransactionType string

var (
	TxCoinbase     BlockTransactionType = "coinbase"
	TxStaking      BlockTransactionType = "staking"
	TxColdStaking  BlockTransactionType = "cold_staking"
	TxPoolStaking  BlockTransactionType = "pool_staking"
	TxStakeDeposit BlockTransactionType = "staking_deposit"
	TxSpend        BlockTransactionType = "spend"
)
