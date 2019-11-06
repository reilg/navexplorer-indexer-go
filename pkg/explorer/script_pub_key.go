package explorer

type ScriptPubKey struct {
	Asm       string   `json:"asm"`
	Hex       string   `json:"hex"`
	ReqSigs   int      `json:"reqSigs,omitempty"`
	Type      VoutType `json:"type"`
	Addresses []string `json:"addresses,omitempty"`
	Hash      string   `json:"hash,omitempty"`
}
