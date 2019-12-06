package explorer

type PaymentRequest struct {
	MetaData MetaData `json:"-"`

	Version             uint32  `json:"version"`
	Hash                string  `json:"hash"`
	BlockHash           string  `json:"blockHash"`
	ProposalHash        string  `json:"proposalHash"`
	Description         string  `json:"description"`
	RequestedAmount     float64 `json:"requestedAmount"`
	Status              string  `json:"status"`
	State               uint    `json:"state"`
	StateChangedOnBlock string  `json:"stateChangedOnBlock,omitempty"`
	PaidOnBlock         string  `json:"paidOnBlock,omitempty"`

	// Custom
	Height uint64      `json:"height"`
	Cycles CfundCycles `json:"cycles"`
}

func (p *PaymentRequest) LatestCycle() *CfundCycle {
	if len(p.Cycles) == 0 {
		return nil
	}

	return &(p.Cycles)[len(p.Cycles)-1]
}

func (p *PaymentRequest) GetCycle(cycle uint) *CfundCycle {
	for i, c := range p.Cycles {
		if c.VotingCycle == cycle {
			return &p.Cycles[i]
		}
	}

	return nil
}
