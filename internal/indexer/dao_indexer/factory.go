package dao_indexer

import (
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func CreateProposal(proposal navcoind.Proposal, height uint64) explorer.Proposal {
	return explorer.Proposal{
		Version:             proposal.Version,
		Hash:                proposal.Hash,
		BlockHash:           proposal.BlockHash,
		Description:         proposal.Description,
		RequestedAmount:     convertStringToUint(proposal.RequestedAmount),
		NotPaidYet:          convertStringToUint(proposal.RequestedAmount),
		UserPaidFee:         convertStringToUint(proposal.UserPaidFee),
		PaymentAddress:      proposal.PaymentAddress,
		ProposalDuration:    proposal.ProposalDuration,
		ExpiresOn:           proposal.ExpiresON,
		Status:              "pending",
		State:               0,
		StateChangedOnBlock: proposal.StateChangedOnBlock,
		Height:              height,
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
