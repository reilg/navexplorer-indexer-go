package vote

import (
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	log "github.com/sirupsen/logrus"
)

func CreateVotes(block *explorer.Block, tx *explorer.BlockTransaction, header *navcoind.BlockHeader) *explorer.DaoVotes {
	if !tx.IsCoinbase() {
		return nil
	}

	daoVote := &explorer.DaoVotes{Height: tx.Height, Address: block.StakedBy}

	for _, cfundVote := range header.CfundVotes {
		vote := explorer.Vote{Type: explorer.ProposalVote, Hash: cfundVote.Hash, Vote: cfundVote.Vote}
		log.Infof("Adding cfund proposal vote for %s - %d", cfundVote.Hash, cfundVote.Vote)
		daoVote.Votes = append(daoVote.Votes, vote)
	}

	for _, cfundVote := range header.CfundRequestVotes {
		vote := explorer.Vote{Type: explorer.PaymentRequestVote, Hash: cfundVote.Hash, Vote: cfundVote.Vote}
		log.Infof("Adding cfund payment request vote for %s - %d", cfundVote.Hash, cfundVote.Vote)
		daoVote.Votes = append(daoVote.Votes, vote)
	}

	for _, cfundVote := range header.DaoSupport {
		vote := explorer.Vote{Type: explorer.DaoSupport, Hash: cfundVote.Hash, Vote: cfundVote.Vote}
		log.Infof("Adding dao support for %s - %d", cfundVote.Hash, cfundVote.Vote)
		daoVote.Votes = append(daoVote.Votes, vote)
	}

	for _, cfundVote := range header.DaoVotes {
		vote := explorer.Vote{Type: explorer.DaoVote, Hash: cfundVote.Hash, Vote: cfundVote.Vote}
		log.Infof("Adding dao vote for %s - %d", cfundVote.Hash, cfundVote.Vote)
		daoVote.Votes = append(daoVote.Votes, vote)
	}

	if len(daoVote.Votes) == 0 {
		return nil
	}

	return daoVote
}
