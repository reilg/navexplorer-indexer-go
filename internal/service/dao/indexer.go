package dao

import (
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/dao/consensus"
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
	consensusIndexer      *consensus.Indexer
	blocksInCycle         uint
	quorum                uint
}

func NewIndexer(
	proposalIndexer *proposal.Indexer,
	paymentRequestIndexer *payment_request.Indexer,
	voteIndexer *vote.Indexer,
	consensusIndexer *consensus.Indexer,
	blocksInCycle uint,
	quorum uint,
) *Indexer {
	return &Indexer{
		proposalIndexer,
		paymentRequestIndexer,
		voteIndexer,
		consensusIndexer,
		blocksInCycle,
		quorum,
	}
}

func (i *Indexer) Index(block *explorer.Block, txs []*explorer.BlockTransaction) {
	blockCycle := block.BlockCycle(i.blocksInCycle, i.quorum)

	i.proposalIndexer.Index(txs)
	i.paymentRequestIndexer.Index(txs)
	i.voteIndexer.IndexVotes(txs, block)

	if blockCycle.IsEnd() {
		log.WithFields(log.Fields{"Quorum": blockCycle.Quorum, "height": block.Height}).Info("Dao - End of voting cycle")
		i.proposalIndexer.Update(blockCycle, block)
		i.paymentRequestIndexer.Update(blockCycle, block)
		i.consensusIndexer.Index()
	}
}
