package softfork_indexer

import (
	"context"
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/index"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	log "github.com/sirupsen/logrus"
)

var SoftForks []explorer.SoftFork

type Indexer struct {
	elastic *index.Index
	navcoin *navcoind.Navcoind
}

func New(elastic *index.Index, navcoin *navcoind.Navcoind) *Indexer {
	return &Indexer{elastic: elastic, navcoin: navcoin}
}

func (i *Indexer) Init() *Indexer {
	info, err := i.navcoin.GetBlockchainInfo()
	if err != nil {
		log.WithError(err).Error("Failed to get blockchaininfo")
	}

	for name, fork := range info.Bip9SoftForks {
		var softFork *explorer.SoftFork
		if err := i.elastic.GetById(index.SoftForkIndex.Get(), name, &softFork); err != nil {
			softFork = &explorer.SoftFork{
				Name:      name,
				SignalBit: uint(fork.Bit),
				State:     explorer.SoftForkDefined,
			}
			_, err := i.elastic.Client.
				Index().
				Index(index.SoftForkIndex.Get()).
				Id(name).
				BodyJson(softFork).
				Do(context.Background())

			if err != nil {
				log.WithError(err).Fatal("Could not index new soft fork")
			}
		}

		SoftForks = append(SoftForks, *softFork)
	}

	return i
}

func (i *Indexer) Update(softFork *explorer.SoftFork) error {
	_, err := i.elastic.Client.
		Index().
		Index(index.SoftForkIndex.Get()).
		Id(softFork.Name).
		BodyJson(softFork).
		Do(context.Background())

	return err
}
