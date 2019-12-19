package explorer

type BlockTransactionType string

var (
	TxCoinbase    BlockTransactionType = "coinbase"
	TxStaking     BlockTransactionType = "staking"
	TxColdStaking BlockTransactionType = "cold_staking"
	TxPoolStaking BlockTransactionType = "pool_staking"
	TxSpend       BlockTransactionType = "spend"
)
