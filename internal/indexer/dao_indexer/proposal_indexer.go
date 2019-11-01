package dao_indexer

import (
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/index"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	log "github.com/sirupsen/logrus"
)

type ProposalIndexer struct {
	navcoind *navcoind.Navcoind
	elastic  *index.Index
}

func NewProposalIndexer(navcoind *navcoind.Navcoind, elastic *index.Index) *ProposalIndexer {
	return &ProposalIndexer{navcoind, elastic}
}

func (i *ProposalIndexer) IndexProposals(txs *[]explorer.BlockTransaction) {
	for _, tx := range *txs {
		if tx.IsSpend() && tx.Version == 4 {
			i.indexProposal(tx)
		}
	}
}

func (i *ProposalIndexer) indexProposal(tx explorer.BlockTransaction) {
	navProposal, err := i.navcoind.GetProposal(tx.Hash)
	if err != nil {
		log.WithError(err).Errorf("Proposal not found in tx %s", tx.Hash)
		return
	}

	log.Info("Indexing proposal in tx ", tx.Hash)
	i.elastic.AddIndexRequest(index.ProposalIndex.Get(), tx.Hash, CreateProposal(navProposal, tx.Height))
}
