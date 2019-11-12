package dao

import (
	"fmt"
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/elastic_cache"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/dao/payment_request"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/dao/proposal"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
)

type Indexer struct {
	navcoin *navcoind.Navcoind
	elastic *elastic_cache.Index
}

func NewIndexer(navcoin *navcoind.Navcoind, elastic *elastic_cache.Index) *Indexer {
	return &Indexer{navcoin, elastic}
}

func (i *Indexer) Index(block *explorer.Block, txs []*explorer.BlockTransaction) {
	i.indexProposals(txs)
	i.indexPaymentRequests(txs)

	for _, tx := range txs {
		if tx.IsCoinbase() {
			i.indexProposalVotes(block, tx)
			i.indexPaymentRequestVotes(block, tx)
		}
	}
}

func (i *Indexer) indexProposals(txs []*explorer.BlockTransaction) {
	for _, tx := range txs {
		if !tx.IsSpend() && tx.Version != 4 {
			continue
		}

		if navP, err := i.navcoin.GetProposal(tx.Hash); err == nil {
			p := proposal.CreateProposal(navP, tx.Height)
			i.elastic.AddIndexRequest(elastic_cache.ProposalIndex.Get(), p.Hash, p)
		}
	}
}

func (i *Indexer) indexProposalVotes(block *explorer.Block, coinbase *explorer.BlockTransaction) {
	if vote := proposal.CreateProposalVotes(block, coinbase); vote != nil {
		if len(vote.Votes) > 0 {
			i.elastic.AddIndexRequest(
				elastic_cache.ProposalVoteIndex.Get(),
				fmt.Sprintf("%d-%s", vote.Height, vote.Address),
				vote,
			)
		}
	}
}

func (i *Indexer) indexPaymentRequests(txs []*explorer.BlockTransaction) {
	for _, tx := range txs {
		if tx.Version != 5 {
			continue
		}

		if navP, err := i.navcoin.GetPaymentRequest(tx.Hash); err == nil {
			p := payment_request.CreatePaymentRequest(navP, tx.Height)
			i.elastic.AddIndexRequest(elastic_cache.PaymentRequestIndex.Get(), p.Hash, p)
		}
	}
}

func (i *Indexer) indexPaymentRequestVotes(block *explorer.Block, coinbase *explorer.BlockTransaction) {
	if vote := payment_request.CreatePaymentRequestVotes(block, coinbase); vote != nil {
		if len(vote.Votes) > 0 {
			i.elastic.AddIndexRequest(
				elastic_cache.PaymentRequestVoteIndex.Get(),
				fmt.Sprintf("%d-%s", vote.Height, vote.Address),
				vote,
			)
		}
	}
}
