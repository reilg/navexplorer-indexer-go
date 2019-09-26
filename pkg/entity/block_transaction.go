package entity

type BlockTransaction struct {
	Hex             string `json:"hex"`
	Txid            string `json:"txid"`
	Hash            string `json:"hash"`
	Size            uint64 `json:"size"`
	VSize           uint64 `json:"vsize"`
	Version         uint32 `json:"version"`
	LockTime        uint32 `json:"locktime"`
	AnonDestination string `json:"anon-destination"`
	Vin             []Vin  `json:"vin"`
	Vout            []Vout `json:"vout"`
	BlockHash       string `json:"blockhash, omitempty"`
	Height          uint64 `json:"height"`
	Confirmations   uint64 `json:"confirmations, omitempty"`
	Time            int64  `json:"time, omitempty"`
	BlockTime       int64  `json:"blocktime, omitempty"`

	Type  string `json:"type"`
	Stake uint64 `json:"stake"`
	Spend uint64 `json:"spend"`
	Fees  uint64 `json:"fees"`
}

func (tx *BlockTransaction) HasInputs() bool {
	return len(tx.Vin) > 0
}

func (tx *BlockTransaction) HasInputWithAddress(a string) bool {
	for _, vin := range tx.Vin {
		if vin.Address == a {
			return true
		}
	}

	return false
}

func (tx *BlockTransaction) GetInputAmount() uint64 {
	var amount uint64 = 0
	for _, i := range tx.Vin {
		amount += i.ValueSat
	}
	return amount
}

func (tx *BlockTransaction) GetOutputsWithAddresses() []Vout {
	var vouts = make([]Vout, 0)
	for _, o := range tx.Vout {
		if len(o.ScriptPubKey.Addresses) != 0 {
			vouts = append(vouts, o)
		}
	}
	return vouts
}

func (tx *BlockTransaction) HasOutputOfType(txType VoutType) bool {
	for _, t := range tx.Vout {
		if t.ScriptPubKey.Type == string(txType) {
			return true
		}
	}

	return false
}

func (tx *BlockTransaction) GetOutputAmount() uint64 {
	var amount uint64 = 0
	for _, o := range tx.Vout {
		amount += o.ValueSat
	}
	return amount
}

func (tx *BlockTransaction) IsCoinbase() bool {
	return tx.HasInputs() && tx.Vin[0].HasCoinbase()
}

func (tx *BlockTransaction) IsSpend() bool {
	return tx.Type == string(TxSpend)
}

func (tx *BlockTransaction) IsStaking() bool {
	return tx.Type == string(TxColdStaking) || tx.Type == string(TxStaking) || tx.Type == string(TxPoolStaking)
}

//for (Output output : transaction.getOutputs()) {
//if (output.getType().equals(OutputType.PUBKEY)) {
//return output.getAmount();
//}
//}
//}

type ScriptSig struct {
	Asm string `json:"asm,omitempty"`
	Hex string `json:"hex,omitempty"`
}

type Vin struct {
	Coinbase  string     `json:"coinbase,omitempty"`
	Txid      *string    `json:"txid,omitempty"`
	Vout      *int       `json:"vout,omitempty"`
	ScriptSig *ScriptSig `json:"scriptSig,omitempty"`
	Value     float64    `json:"value,omitempty"`
	ValueSat  uint64     `json:"valuesat,omitempty"`
	Address   string     `json:"address,omitempty"`
	Sequence  uint32     `json:"sequence"`
}

func (v *Vin) HasCoinbase() bool {
	return v.Coinbase != ""
}

type ScriptPubKey struct {
	Asm       string   `json:"asm"`
	Hex       string   `json:"hex"`
	ReqSigs   int      `json:"reqSigs, omitempty"`
	Type      string   `json:"type"`
	Addresses []string `json:"addresses, omitempty"`
}

type Vout struct {
	Value        float64      `json:"value"`
	ValueSat     uint64       `json:"valuesat"`
	N            int          `json:"n"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
}

type VoutType string
type BlockTransactionType string

var (
	TxCoinbase    BlockTransactionType = "coinbase"
	TxStaking     BlockTransactionType = "staking"
	TxColdStaking BlockTransactionType = "cold_staking"
	TxPoolStaking BlockTransactionType = "pool_staking"
	TxSpend       BlockTransactionType = "spend"

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
