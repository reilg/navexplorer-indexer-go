package explorer

type Vouts []Vout

func (vouts *Vouts) HasOutputOfType(txType VoutType) bool {
	for _, vout := range *vouts {
		if vout.ScriptPubKey.Type == string(txType) {
			return true
		}
	}

	return false
}

func (vouts *Vouts) GetAmount() uint64 {
	var amount uint64 = 0
	for _, o := range *vouts {
		amount += o.ValueSat
	}

	return amount
}

func (vouts *Vouts) GetAmountByAddress(address string) (value float64, valuesat uint64) {
	for _, o := range *vouts {
		if isValueInList(address, o.ScriptPubKey.Addresses) {
			value += o.Value
			valuesat += o.ValueSat
		}
	}
	return
}

func (vouts *Vouts) FilterWithAddresses() Vouts {
	var filtered = make(Vouts, 0)
	for _, o := range *vouts {
		if len(o.ScriptPubKey.Addresses) != 0 {
			filtered = append(filtered, o)
		}
	}
	return filtered
}
