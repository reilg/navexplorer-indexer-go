package softfork

import (
	"context"
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/elastic_cache"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	log "github.com/sirupsen/logrus"
)

var SoftForks explorer.SoftForks

type Service struct {
	navcoin *navcoind.Navcoind
	elastic *elastic_cache.Index
	repo    *Repository
}

func New(navcoin *navcoind.Navcoind, elastic *elastic_cache.Index, repo *Repository) *Service {
	return &Service{navcoin, elastic, repo}
}

func (i *Service) LoadSoftForks() {
	log.Info("Load SoftForks")

	info, err := i.navcoin.GetBlockchainInfo()
	if err != nil {
		log.WithError(err).Fatal("Failed to get blockchaininfo")
	}

	SoftForks, err = i.repo.GetSoftForks()
	if err != nil {
		log.WithError(err).Fatal("Failed to get softforks from repo")
	}

	for name, bip9fork := range info.Bip9SoftForks {
		if SoftForks.GetSoftFork(name) == nil {
			softFork := &explorer.SoftFork{Name: name, SignalBit: bip9fork.Bit, State: explorer.SoftForkDefined}
			resp, err := i.elastic.Client.Index().Index(elastic_cache.SoftForkIndex.Get()).BodyJson(softFork).Do(context.Background())
			if err != nil {
				log.WithError(err).Fatal("Failed to save new softfork")
			}

			log.Info("Saving new softfork ", name)
			softFork.Id = resp.Id
			SoftForks = append(SoftForks, softFork)
		}
	}

	for _, sf := range SoftForks {
		log.WithFields(log.Fields{"Name": sf.Name, "id": sf.Id}).Info("SoftFork")
	}
}
