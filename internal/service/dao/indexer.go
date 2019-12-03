package dao

import (
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/dao/payment_request"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/dao/proposal"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/dao/vote"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	log "github.com/sirupsen/logrus"
)

type Indexer struct {
	proposalIndexer       *proposal.Indexer
	paymentRequestIndexer *payment_request.Indexer
	voteIndexer           *vote.Indexer
	blocksInCycle         uint
	quorum                uint
}

func NewIndexer(
	proposalIndexer *proposal.Indexer,
	paymentRequestIndexer *payment_request.Indexer,
	voteIndexer *vote.Indexer,
	blocksInCycle uint,
	quorum uint,
) *Indexer {
	return &Indexer{
		proposalIndexer,
		paymentRequestIndexer,
		voteIndexer,
		blocksInCycle,
		quorum,
	}
}

func (i *Indexer) Index(block *explorer.Block, txs []*explorer.BlockTransaction) {
	i.proposalIndexer.Index(txs)
	i.paymentRequestIndexer.Index(txs)

	blockCycle := block.BlockCycle(i.blocksInCycle, i.quorum)
	i.voteIndexer.InitCycles(blockCycle)

	if daoVote := i.voteIndexer.IndexVotes(txs, block); daoVote != nil {
		i.applyVotes(daoVote, blockCycle)
	}

	if blockCycle.Index == blockCycle.Size {
		log.WithFields(log.Fields{"Quorum": blockCycle.Quorum}).Info("Dao - End of voting cycle")
		i.proposalIndexer.UpdateState(blockCycle)
	}
}

func (i *Indexer) applyVotes(daoVote *explorer.DaoVote, blockCycle explorer.BlockCycle) {
	for _, v := range daoVote.Votes {
		if v.Type == explorer.ProposalVote {
			i.proposalIndexer.ApplyVote(v, blockCycle)
		}

		if v.Type == explorer.PaymentRequestVote {
			i.paymentRequestIndexer.ApplyVote(v, blockCycle)

		}
	}
}
