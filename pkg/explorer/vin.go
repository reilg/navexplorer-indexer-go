package explorer

type Vin struct {
	Coinbase  string     `json:"coinbase,omitempty"`
	Txid      *string    `json:"txid,omitempty"`
	Vout      *int       `json:"vout,omitempty"`
	ScriptSig *ScriptSig `json:"scriptSig,omitempty"`
	Value     float64    `json:"value,omitempty"`
	ValueSat  uint64     `json:"valuesat,omitempty"`
	Address   string     `json:"address,omitempty"`
	Sequence  uint32     `json:"sequence"`

	// Custom
	Height uint64 `json:"height"`
	Type   string `json:"type"`
}

func (v *Vin) IsCoinbase() bool {
	return v.Coinbase != ""
}
