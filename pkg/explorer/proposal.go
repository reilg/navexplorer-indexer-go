package explorer

type Proposal struct {
	MetaData MetaData `json:"-"`

	Version             uint32  `json:"version"`
	Hash                string  `json:"hash"`
	BlockHash           string  `json:"blockHash"`
	Description         string  `json:"description"`
	RequestedAmount     float64 `json:"requestedAmount"`
	NotPaidYet          float64 `json:"notPaidYet"`
	NotRequestedYet     float64 `json:"notRequestedYet"`
	UserPaidFee         float64 `json:"userPaidFee"`
	PaymentAddress      string  `json:"paymentAddress"`
	ProposalDuration    uint64  `json:"proposalDuration"`
	ExpiresOn           uint64  `json:"expiresOn"`
	Status              string  `json:"status"`
	State               uint    `json:"state"`
	StateChangedOnBlock string  `json:"stateChangedOnBlock,omitempty"`

	// Custom
	Height uint64      `json:"height"`
	Cycles CfundCycles `json:"cycles"`
}

func (p *Proposal) LatestCycle() *CfundCycle {
	if len(p.Cycles) == 0 {
		return nil
	}

	return &(p.Cycles)[len(p.Cycles)-1]
}

type CfundCycles []CfundCycle

type CfundCycle struct {
	VotingCycle uint `json:"votingCycle"`
	VotesYes    uint `json:"votesYes"`
	VotesNo     uint `json:"votesNo"`
}

func (cfc *CfundCycle) Votes() uint {
	return cfc.VotesYes + cfc.VotesNo
}

func (p *Proposal) GetCycle(cycle uint) *CfundCycle {
	for i, c := range p.Cycles {
		if c.VotingCycle == cycle {
			return &p.Cycles[i]
		}
	}

	return nil
}
