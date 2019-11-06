package explorer

type RawVout struct {
	Value        float64      `json:"value"`
	ValueSat     uint64       `json:"valuesat"`
	N            int          `json:"n"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
}
type Vout struct {
	RawVout
}

func (o *Vout) IsColdStaking() bool {
	return o.ScriptPubKey.Type == VoutColdStaking
}

func (o *Vout) IsProposalVote() bool {
	return o.ScriptPubKey.Type == VoutProposalYesVote || o.ScriptPubKey.Type == VoutProposalNoVote
}

func (o *Vout) IsColdStakingAddress(address string) bool {
	return len(o.ScriptPubKey.Addresses) == 2 && o.ScriptPubKey.Addresses[0] == address
}

func (o *Vout) IsColdSpendingAddress(address string) bool {
	return len(o.ScriptPubKey.Addresses) == 2 && o.ScriptPubKey.Addresses[0] == address
}
