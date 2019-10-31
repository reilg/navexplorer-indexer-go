package dao_indexer

import (
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/index"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
)

type CfundProposalIndexer struct {
	navcoind *navcoind.Navcoind
	elastic  *index.Index
}

func NewCfundProposalIndexer(navcoind *navcoind.Navcoind, elastic *index.Index) *CfundProposalIndexer {
	return &CfundProposalIndexer{navcoind, elastic}
}

func (i *CfundProposalIndexer) IndexProposalsForTransactions(txs *[]explorer.BlockTransaction) {
	for _, tx := range *txs {
		if tx.IsSpend() && tx.Version == 4 {
			i.indexProposal(tx)
		}
	}
}

func (i *CfundProposalIndexer) indexProposal(tx explorer.BlockTransaction) {
	navProposal, err := i.navcoind.GetProposal(tx.Hash)
	if err != nil {
		log.WithError(err).Errorf("Proposal not found in tx %s", tx.Hash)
		return
	}

	log.Info("Indexing proposal in tx ", tx.Hash)
	i.elastic.AddRequest(index.ProposalIndex.Get(), tx.Hash, CreateProposal(navProposal))
}
