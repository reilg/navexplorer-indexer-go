package consensus

import (
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/elastic_cache"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	log "github.com/sirupsen/logrus"
)

var Consensus *explorer.Consensus

type Indexer struct {
	navcoin *navcoind.Navcoind
	elastic *elastic_cache.Index
	repo    *Repository
}

func NewIndexer(navcoin *navcoind.Navcoind, elastic *elastic_cache.Index, repo *Repository) *Indexer {
	return &Indexer{navcoin, elastic, repo}
}

func (i *Indexer) Index() {
	cfundStats, err := i.navcoin.CfundStats()
	if err != nil {
		log.WithError(err).Error("Failed to get CfundStats")
		return
	}

	consensus, err := i.repo.GetConsensus()
	if err != nil {
		consensus = CreateConsensus(&cfundStats)
		i.elastic.AddIndexRequest(elastic_cache.ConsensusIndex.Get(), "consensus", consensus)
	} else {
		UpdateConsensus(&cfundStats, consensus)
		i.elastic.AddUpdateRequest(elastic_cache.ConsensusIndex.Get(), "consensus", consensus, consensus.MetaData.Id)
	}
}
