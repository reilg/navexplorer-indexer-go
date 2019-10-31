package explorer

type BlockTransaction struct {
	Hex             string `json:"hex"`
	Txid            string `json:"txid"`
	Hash            string `json:"hash"`
	Size            uint64 `json:"size"`
	VSize           uint64 `json:"vsize"`
	Version         uint32 `json:"version"`
	LockTime        uint32 `json:"locktime"`
	AnonDestination string `json:"anon-destination"`
	Vin             Vins   `json:"vin"`
	Vout            Vouts  `json:"vout"`
	BlockHash       string `json:"blockhash, omitempty"`
	Height          uint64 `json:"height"`
	Confirmations   uint64 `json:"confirmations, omitempty"`
	Time            int64  `json:"time, omitempty"`
	BlockTime       int64  `json:"blocktime, omitempty"`

	// Custom
	Type  BlockTransactionType `json:"type"`
	Stake uint64               `json:"stake"`
	Spend uint64               `json:"spend"`
	Fees  uint64               `json:"fees"`
}

func (tx *BlockTransaction) GetAllAddresses() []string {
	var addressMap = make(map[string]struct{})
	for _, vin := range tx.Vin {
		for _, a := range vin.Addresses {
			if _, ok := addressMap[a]; !ok {
				addressMap[a] = struct{}{}
			}
		}
	}
	for _, vout := range tx.Vout {
		for _, a := range vout.ScriptPubKey.Addresses {
			if _, ok := addressMap[a]; !ok {
				addressMap[a] = struct{}{}
			}
		}
	}

	var addresses = make([]string, 0)
	for address, _ := range addressMap {
		addresses = append(addresses, address)
	}

	return addresses
}

func (tx *BlockTransaction) IsCoinbase() bool {
	return !tx.Vin.Empty() && tx.Vin[0].IsCoinbase()
}

func (tx *BlockTransaction) IsSpend() bool {
	return tx.Type == TxSpend
}

func (tx *BlockTransaction) IsAnyStaking() bool {
	return tx.Type == TxColdStaking || tx.Type == TxStaking || tx.Type == TxPoolStaking
}

func (tx *BlockTransaction) IsStaking() bool {
	return tx.Type == TxStaking
}

func (tx *BlockTransaction) IsColdStaking() bool {
	return tx.Type == TxColdStaking
}

func (tx *BlockTransaction) IsPoolStaking() bool {
	return tx.Type == TxPoolStaking
}

func (tx *BlockTransaction) HasColdInput(address string) bool {
	for _, i := range tx.Vin {
		if i.PreviousOutput.Type == VoutColdStaking && i.Addresses[0] == address {
			return true
		}
	}
	return false
}

func (tx *BlockTransaction) HasColdStakeStake(address string) bool {
	for _, o := range tx.Vout {
		if o.ScriptPubKey.Type == VoutColdStaking && o.ScriptPubKey.Addresses[0] == address {
			return true
		}
	}
	return false
}

func (tx *BlockTransaction) HasColdStakeSpend(address string) bool {
	for _, o := range tx.Vout {
		if o.ScriptPubKey.Type == VoutColdStaking && o.ScriptPubKey.Addresses[1] == address {
			return true
		}
	}
	return false
}

func isValueInList(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
