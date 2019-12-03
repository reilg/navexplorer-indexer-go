package proposal

import (
	"context"
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/elastic_cache"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	log "github.com/sirupsen/logrus"
)

type Indexer struct {
	navcoin           *navcoind.Navcoind
	elastic           *elastic_cache.Index
	maxProposalCycles uint
}

func NewIndexer(navcoin *navcoind.Navcoind, elastic *elastic_cache.Index, maxProposalCycles uint) *Indexer {
	return &Indexer{navcoin, elastic, maxProposalCycles}
}

func (i *Indexer) Index(txs []*explorer.BlockTransaction) {
	for _, tx := range txs {
		if !tx.IsSpend() && tx.Version != 4 {
			continue
		}

		if navP, err := i.navcoin.GetProposal(tx.Hash); err == nil {
			proposal := CreateProposal(navP, tx.Height)
			log.Infof("Index Proposal: %s", proposal.Hash)

			index := elastic_cache.ProposalIndex.Get()
			resp, err := i.elastic.Client.Index().Index(index).BodyJson(proposal).Do(context.Background())
			if err != nil {
				log.WithError(err).Fatal("Failed to save new proposal")
			}

			proposal.MetaData = explorer.NewMetaData(resp.Id, resp.Index)
			Proposals = append(Proposals, proposal)
		}
	}
}

func (i *Indexer) ApplyVote(vote explorer.Vote, blockCycle explorer.BlockCycle) {
	proposal := getProposalByHash(vote.Hash)
	if proposal == nil {
		log.Fatalf("Proposal not found: %s", vote.Hash)
		return
	}

	if vote.Vote == 1 {
		proposal.GetCycle(blockCycle.Cycle).VotesYes++
	}
	if vote.Vote == -1 {
		proposal.GetCycle(blockCycle.Cycle).VotesNo++
	}

	i.elastic.AddUpdateRequest(elastic_cache.ProposalIndex.Get(), proposal.Hash, proposal, proposal.MetaData.Id)
}

func (i *Indexer) UpdateState(blockCycle explorer.BlockCycle) {
	for _, proposal := range Proposals {
		if proposal.Status == "pending" && proposal.LatestCycle().Votes() >= blockCycle.Quorum {
			log.WithFields(log.Fields{"Votes": proposal.LatestCycle().Votes()}).Info("Proposal has met Quorum")
		}

		if len(proposal.Cycles) == int(i.maxProposalCycles)+1 && proposal.Status == "pending" {
			proposal.Status = "expired"
			i.elastic.AddUpdateRequest(elastic_cache.ProposalIndex.Get(), proposal.Hash, proposal, proposal.MetaData.Id)
			Proposals.Delete(proposal.Hash)
			log.Infof("Proposal Expired: %s", proposal.Hash)
		}
	}
}
