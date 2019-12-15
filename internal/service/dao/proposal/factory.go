package proposal

import (
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func CreateProposal(proposal navcoind.Proposal, height uint64) *explorer.Proposal {
	return &explorer.Proposal{
		Version:             proposal.Version,
		Hash:                proposal.Hash,
		BlockHash:           proposal.BlockHash,
		Description:         proposal.Description,
		RequestedAmount:     convertStringToFloat(proposal.RequestedAmount),
		NotPaidYet:          convertStringToFloat(proposal.RequestedAmount),
		NotRequestedYet:     convertStringToFloat(proposal.RequestedAmount),
		UserPaidFee:         convertStringToFloat(proposal.UserPaidFee),
		PaymentAddress:      proposal.PaymentAddress,
		ProposalDuration:    proposal.ProposalDuration,
		ExpiresOn:           proposal.ExpiresOn,
		Status:              "pending",
		State:               proposal.State,
		StateChangedOnBlock: proposal.StateChangedOnBlock,
		Height:              height,
		UpdatedOnBlock:      height,
	}
}

func UpdateProposal(proposal navcoind.Proposal, height uint64, p *explorer.Proposal) {
	if p.NotPaidYet != convertStringToFloat(proposal.RequestedAmount) {
		p.NotPaidYet = convertStringToFloat(proposal.RequestedAmount)
		p.UpdatedOnBlock = height
	}

	if p.NotRequestedYet != convertStringToFloat(proposal.RequestedAmount) {
		p.NotRequestedYet = convertStringToFloat(proposal.RequestedAmount)
		p.UpdatedOnBlock = height
	}

	if p.Status != explorer.ProposalStatus(proposal.Status) {
		p.Status = explorer.ProposalStatus(proposal.Status)
		p.UpdatedOnBlock = height
	}

	if p.State != proposal.State {
		p.State = proposal.State
		p.UpdatedOnBlock = height
	}

	if p.StateChangedOnBlock != proposal.StateChangedOnBlock {
		p.StateChangedOnBlock = proposal.StateChangedOnBlock
		p.UpdatedOnBlock = height
	}
}

func convertStringToFloat(input string) float64 {
	output, err := strconv.ParseFloat(input, 64)
	if err != nil {
		log.WithError(err).Errorf("Unable to convert %s to uint64", input)
		return 0
	}

	return output
}
