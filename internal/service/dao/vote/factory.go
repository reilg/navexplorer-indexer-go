package vote

import (
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
)

func CreateVotes(block *explorer.Block, tx *explorer.BlockTransaction) *explorer.DaoVote {
	if !tx.IsCoinbase() {
		return nil
	}

	daoVote := &explorer.DaoVote{Height: tx.Height, Address: block.StakedBy}
	for _, vout := range tx.Vout {
		if !vout.IsProposalVote() && !vout.IsPaymentRequestVote() {
			continue
		}
		vote := explorer.Vote{Hash: vout.ScriptPubKey.Hash, Vote: persuasion(vout.ScriptPubKey.Type)}
		if vout.IsProposalVote() {
			vote.Type = explorer.ProposalVote
		}
		if vout.IsPaymentRequestVote() {
			vote.Type = explorer.PaymentRequestVote
		}

		daoVote.Votes = append(daoVote.Votes, vote)
	}

	if len(daoVote.Votes) == 0 {
		return nil
	}

	return daoVote
}

func persuasion(voutType explorer.VoutType) int {
	switch voutType {

	case explorer.VoutProposalYesVote:
		return 1
	case explorer.VoutPaymentRequestYesVote:
		return 1
	case explorer.VoutProposalNoVote:
		return -1
	case explorer.VoutPaymentRequestNoVote:
		return -1
	}

	return 0
}
