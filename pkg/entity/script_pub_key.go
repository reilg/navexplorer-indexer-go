package entity

type ScriptPubKey struct {
	Asm       string   `json:"asm"`
	Hex       string   `json:"hex"`
	ReqSigs   int      `json:"reqSigs, omitempty"`
	Type      string   `json:"type"`
	Addresses []string `json:"addresses, omitempty"`
}
