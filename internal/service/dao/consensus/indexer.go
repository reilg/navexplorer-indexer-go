package consensus

import (
	"context"
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

func (i *Indexer) Index() error {
	cfundStats, err := i.navcoin.CfundStats()
	if err != nil {
		log.WithError(err).Error("Failed to get CfundStats")
		return err
	}

	consensus, _ := i.repo.GetConsensus()
	if consensus == nil {
		consensus = new(explorer.Consensus)
		UpdateConsensus(&cfundStats, consensus)
		_, err := i.elastic.Client.Index().Index(elastic_cache.ConsensusIndex.Get()).BodyJson(consensus).Do(context.Background())
		if err != nil {
			log.WithError(err).Fatalf("Failed to persist consensus")
			return err
		}
	} else {
		UpdateConsensus(&cfundStats, consensus)
		log.Info("Index Update Cfund Consensus")
		i.elastic.AddUpdateRequest(elastic_cache.ConsensusIndex.Get(), "consensus", consensus, consensus.MetaData.Id)
	}

	Consensus = consensus

	return nil
}
