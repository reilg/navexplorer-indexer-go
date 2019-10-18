package softfork_indexer

import (
	"context"
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/config"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/index"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
)

var SoftForks explorer.SoftForks

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

func (i *Indexer) UpdateForSignal(signal *explorer.Signal, block *explorer.Block) {
	if !signal.IsSignalling() {
		return
	}

	size := config.Get().SoftForkBlockCycle
	blockCycle := block.BlockCycle(size)
	cycleIndex := block.CycleIndex(size)

	for _, s := range signal.SoftForks {
		softFork := SoftForks.GetSoftFork(s)
		if softFork == nil || softFork.SignalHeight >= signal.Height {
			continue
		}

		cycle := softFork.GetCycle(blockCycle)
		if cycle == nil {
			softFork.Cycles = append(softFork.Cycles, explorer.SoftForkCycle{Cycle: blockCycle, BlocksSignalling: 0})
			cycle = softFork.GetCycle(blockCycle)
		}
		cycle.BlocksSignalling++

		if cycle.BlocksSignalling >= uint(float64(size)*0.75) {
			softFork.LockedInHeight = uint64(size * blockCycle)
			softFork.ActivationHeight = softFork.LockedInHeight + uint64(size)
			if cycleIndex == size {
				softFork.State = explorer.SoftForkLockedIn
			}
		}

		softFork.SignalHeight = signal.Height
		if softFork.State == explorer.SoftForkDefined {
			softFork.State = explorer.SoftForkStarted
		}

		i.elastic.GetBulkRequest(block.Height).Add(
			elastic.NewBulkIndexRequest().Index(index.SoftForkIndex.Get()).Id(softFork.Name).Doc(softFork),
		)
	}
}
