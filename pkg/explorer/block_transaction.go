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

	MetaData struct {
		Type  string `json:"type"`
		Stake uint64 `json:"stake"`
		Spend uint64 `json:"spend"`
		Fees  uint64 `json:"fees"`
	}
}

func (blockTransaction *BlockTransaction) GetAllAddresses() []string {
	var addressMap = make(map[string]struct{})
	for _, vin := range blockTransaction.Vin {
		if vin.Address != "" {
			if _, ok := addressMap[vin.Address]; !ok {
				addressMap[vin.Address] = struct{}{}
			}
		}
	}
	for _, vout := range blockTransaction.Vout {
		if len(vout.ScriptPubKey.Addresses) != 0 {
			for _, address := range vout.ScriptPubKey.Addresses {
				if _, ok := addressMap[address]; !ok {
					addressMap[address] = struct{}{}
				}
			}
		}
	}

	var addresses = make([]string, 0)
	for address, _ := range addressMap {
		addresses = append(addresses, address)
	}

	return addresses
}

func (blockTransaction *BlockTransaction) IsCoinbase() bool {
	return !blockTransaction.Vin.Empty() && blockTransaction.Vin[0].IsCoinbase()
}

func (blockTransaction *BlockTransaction) IsSpend() bool {
	return blockTransaction.MetaData.Type == string(TxSpend)
}

func (blockTransaction *BlockTransaction) IsAnyStaking() bool {
	return blockTransaction.MetaData.Type == string(TxColdStaking) ||
		blockTransaction.MetaData.Type == string(TxStaking) ||
		blockTransaction.MetaData.Type == string(TxPoolStaking)
}

func (blockTransaction *BlockTransaction) IsStaking() bool {
	return blockTransaction.MetaData.Type == string(TxStaking)
}

func (blockTransaction *BlockTransaction) IsColdStaking() bool {
	return blockTransaction.MetaData.Type == string(TxColdStaking)
}

func (blockTransaction *BlockTransaction) IsPoolStaking() bool {
	return blockTransaction.MetaData.Type == string(TxPoolStaking)
}

func isValueInList(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
