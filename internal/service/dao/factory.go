package dao

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

func CreatePaymentRequest(paymentRequest navcoind.PaymentRequest, height uint64) *explorer.PaymentRequest {
	return &explorer.PaymentRequest{
		Version:             paymentRequest.Version,
		Hash:                paymentRequest.Hash,
		BlockHash:           paymentRequest.BlockHash,
		Description:         paymentRequest.Description,
		RequestedAmount:     convertStringToFloat(paymentRequest.RequestedAmount),
		Status:              "pending",
		State:               0,
		StateChangedOnBlock: paymentRequest.StateChangedOnBlock,
		Height:              height,
	}
}

func CreateVotes(block *explorer.Block, tx *explorer.BlockTransaction) *explorer.DaoVote {
	if !tx.IsCoinbase() {
		return nil
	}

	daoVote := &explorer.DaoVote{Height: tx.Height, Address: block.StakedBy}
	for _, vout := range tx.Vout {
		if vout.IsProposalVote() || !vout.IsPaymentRequestVote() {
			vote := explorer.Vote{Hash: vout.ScriptPubKey.Hash, Vote: persuasion(vout.ScriptPubKey.Type)}
			if vout.IsProposalVote() {
				vote.Type = explorer.ProposalVote
			}
			if vout.IsPaymentRequestVote() {
				vote.Type = explorer.PaymentRequestVote
			}

			daoVote.Votes = append(daoVote.Votes, vote)
		}
	}

	return daoVote
}

func persuasion(voutType explorer.VoutType) int {
	switch voutType {

	case explorer.VoutProposalYesVote:
	case explorer.VoutPaymentRequestYesVote:
		return 1

	case explorer.VoutProposalNoVote:
	case explorer.VoutPaymentRequestNoVote:
		return 1
	}

	return 0
}

func convertStringToFloat(input string) float64 {
	output, err := strconv.ParseFloat(input, 64)
	if err != nil {
		log.WithError(err).Errorf("Unable to convert %s to uint64", input)
		return 0
	}

	return output
}
