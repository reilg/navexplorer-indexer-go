package dao_indexer

import (
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func CreateProposal(proposal navcoind.Proposal) explorer.Proposal {
	return explorer.Proposal{
		Version:             proposal.Version,
		Hash:                proposal.Hash,
		BlockHash:           proposal.BlockHash,
		Description:         proposal.Description,
		RequestedAmount:     convertStringToUint(proposal.RequestedAmount),
		NotPaidYet:          convertStringToUint(proposal.NotPaidYet),
		UserPaidFee:         convertStringToUint(proposal.UserPaidFee),
		PaymentAddress:      proposal.PaymentAddress,
		ProposalDuration:    proposal.ProposalDuration,
		ExpiresOn:           proposal.ExpiresON,
		VotesYes:            proposal.VotesYes,
		VotesNo:             proposal.VotesNo,
		VotingCycle:         proposal.VotingCycle,
		Status:              proposal.Status,
		State:               proposal.State,
		StateChangedOnBlock: proposal.StateChangedOnBlock,
	}
}

func convertStringToUint(input string) uint64 {
	output, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		log.WithError(err).Errorf("Unable to convert %s to uint64", input)
		return 0
	}

	return output
}
