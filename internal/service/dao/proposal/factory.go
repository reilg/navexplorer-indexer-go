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
		UserPaidFee:         convertStringToFloat(proposal.UserPaidFee),
		PaymentAddress:      proposal.PaymentAddress,
		ProposalDuration:    proposal.ProposalDuration,
		ExpiresOn:           proposal.ExpiresOn,
		Status:              "pending",
		State:               0,
		StateChangedOnBlock: proposal.StateChangedOnBlock,
		Height:              height,
	}
}

func CreateProposalVotes(block *explorer.Block, tx *explorer.BlockTransaction) *explorer.ProposalVote {
	if !tx.IsCoinbase() {
		return nil
	}

	proposalVote := &explorer.ProposalVote{Height: tx.Height, Address: block.StakedBy}
	for _, vout := range tx.Vout {
		if !vout.IsProposalVote() {
			continue
		}

		vote := explorer.Vote{Hash: vout.ScriptPubKey.Hash, Vote: vout.ScriptPubKey.Type == explorer.VoutProposalYesVote}
		proposalVote.Votes = append(proposalVote.Votes, vote)
	}

	return proposalVote
}

func convertStringToFloat(input string) float64 {
	output, err := strconv.ParseFloat(input, 64)
	if err != nil {
		log.WithError(err).Errorf("Unable to convert %s to uint64", input)
		return 0
	}

	return output
}
