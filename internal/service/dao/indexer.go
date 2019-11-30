package dao

import (
	"fmt"
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/elastic_cache"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	log "github.com/sirupsen/logrus"
)

type Indexer struct {
	navcoin       *navcoind.Navcoind
	elastic       *elastic_cache.Index
	blocksInCycle uint
	quorum        uint
}

func NewIndexer(navcoin *navcoind.Navcoind, elastic *elastic_cache.Index, blocksInCycle uint, quorum uint) *Indexer {
	return &Indexer{navcoin, elastic, blocksInCycle, quorum}
}

func (i *Indexer) Index(block *explorer.Block, txs []*explorer.BlockTransaction) {
	i.indexProposals(txs)
	i.indexPaymentRequests(txs)

	for _, tx := range txs {
		if tx.IsCoinbase() {
			i.indexDaoVote(block, tx)
		}
	}
}

func (i *Indexer) indexProposals(txs []*explorer.BlockTransaction) {
	for _, tx := range txs {
		if !tx.IsSpend() && tx.Version != 4 {
			continue
		}

		if navP, err := i.navcoin.GetProposal(tx.Hash); err == nil {
			p := CreateProposal(navP, tx.Height)
			log.Infof("Index Proposal: %s", p.Hash)
			i.elastic.AddIndexRequest(elastic_cache.ProposalIndex.Get(), p.Hash, p)
		}
	}
}

func (i *Indexer) indexPaymentRequests(txs []*explorer.BlockTransaction) {
	for _, tx := range txs {
		if tx.Version != 5 {
			continue
		}

		if navP, err := i.navcoin.GetPaymentRequest(tx.Hash); err == nil {
			p := CreatePaymentRequest(navP, tx.Height)
			log.Infof("Index PaymentRequest: %s", p.Hash)
			i.elastic.AddIndexRequest(elastic_cache.PaymentRequestIndex.Get(), p.Hash, p)
		}
	}
}

func (i *Indexer) indexDaoVote(block *explorer.Block, coinbase *explorer.BlockTransaction) *explorer.DaoVote {
	if vote := CreateVotes(block, coinbase); vote != nil {
		i.elastic.AddIndexRequest(
			elastic_cache.DaoVoteIndex.Get(),
			fmt.Sprintf("%d-%s", vote.Height, vote.Address),
			vote,
		)
		return vote
	}

	return nil
}
