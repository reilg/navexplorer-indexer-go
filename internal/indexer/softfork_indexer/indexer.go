package softfork_indexer

import (
	"context"
	"fmt"
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/config"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/index"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer/repository"
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

func (i *Indexer) reset() {
	for idx, s := range SoftForks {
		SoftForks[idx] = explorer.SoftFork{Name: s.Name, SignalBit: s.SignalBit, State: explorer.SoftForkDefined}
	}
}

func (i *Indexer) Reindex(height uint64) {
	i.reset()

	size := config.Get().SoftForkBlockCycle
	start := uint64(1)
	end := uint64(size)

	for {
		if height == 0 || end >= height {
			break
		}

		signalRepository := repository.New(i.elastic.Client)
		signals := signalRepository.GetSignals(start, end)
		cycle := explorer.GetCycleForHeight(start, size)

		log.WithFields(log.Fields{"Start": start, "End": end, "height": height, "signals": len(*signals)}).
			Info(fmt.Sprintf("Cycle ", cycle))

		for _, s := range *signals {
			for _, sf := range s.SoftForks {
				softFork := SoftForks.GetSoftFork(sf)
				if softFork.IsOpen() {
					softFork.Dirty = true
					softFork.SignalHeight = end
					softFork.State = explorer.SoftForkStarted
					cycleIndex := explorer.GetCycleForHeight(s.Height, size)
					var cycle *explorer.SoftForkCycle
					if cycle = softFork.GetCycle(cycleIndex); cycle == nil {
						softFork.Cycles = append(softFork.Cycles, explorer.SoftForkCycle{Cycle: cycleIndex, BlocksSignalling: 0})
						cycle = softFork.GetCycle(cycleIndex)
					}
					cycle.BlocksSignalling++
				}
			}
		}

		for idx, s := range SoftForks {
			if s.State == explorer.SoftForkStarted && s.LatestCycle().BlocksSignalling >= explorer.GetQuorum(size) {
				SoftForks[idx].State = explorer.SoftForkLockedIn
				SoftForks[idx].LockedInHeight = end
				SoftForks[idx].ActivationHeight = end + uint64(size)
				SoftForks[idx].Dirty = true
			}
			if s.State == explorer.SoftForkLockedIn && height >= SoftForks[idx].ActivationHeight {
				SoftForks[idx].State = explorer.SoftForkActive
				SoftForks[idx].Dirty = true
			}
		}

		start, end = func(start uint64, end uint64, height uint64) (uint64, uint64) {
			start += uint64(config.Get().SoftForkBlockCycle)
			end += uint64(config.Get().SoftForkBlockCycle)
			if end > height {
				end = height
			}
			return start, end
		}(start, end, height)
	}
}

func (i *Indexer) Update(signal *explorer.Signal, block *explorer.Block) {
	i.UpdateForSignal(signal, block)
	i.UpdateSoftForkState(&SoftForks, block.Height)
	i.persist(block.Height)
}

func (i *Indexer) UpdateForSignal(signal *explorer.Signal, block *explorer.Block) {
	if signal == nil || !signal.IsSignalling() {
		return
	}

	size := config.Get().SoftForkBlockCycle
	blockCycle := block.BlockCycle(size)

	for _, s := range signal.SoftForks {
		softFork := SoftForks.GetSoftFork(s)
		if softFork == nil || softFork.SignalHeight >= signal.Height {
			continue
		}

		softFork.SignalHeight = signal.Height
		softFork.Dirty = true
		if softFork.State == explorer.SoftForkDefined {
			softFork.State = explorer.SoftForkStarted
		}

		var cycle *explorer.SoftForkCycle
		if cycle = softFork.GetCycle(blockCycle.Cycle); cycle == nil {
			softFork.Cycles = append(softFork.Cycles, explorer.SoftForkCycle{Cycle: blockCycle.Cycle, BlocksSignalling: 0})
			cycle = softFork.GetCycle(blockCycle.Cycle)
		}

		cycle.BlocksSignalling++
	}
}

func (i *Indexer) UpdateSoftForkState(softForks *explorer.SoftForks, height uint64) {
	size := config.Get().SoftForkBlockCycle
	blockCycle := explorer.GetCycleForHeight(height, size)

	for idx, s := range SoftForks {
		if s.Cycles == nil {
			continue
		}
		if s.State == explorer.SoftForkStarted && s.LatestCycle().BlocksSignalling >= explorer.GetQuorum(size) {
			SoftForks[idx].LockedInHeight = uint64(size * blockCycle)
			SoftForks[idx].ActivationHeight = SoftForks[idx].LockedInHeight + uint64(size)
			SoftForks[idx].Dirty = true
		}
		if s.LockedInHeight != 0 && s.ActivationHeight != 0 {
			if s.State == explorer.SoftForkStarted && height >= s.LockedInHeight {
				SoftForks[idx].State = explorer.SoftForkLockedIn
				SoftForks[idx].Dirty = true
			}
			if s.State == explorer.SoftForkLockedIn && height >= s.ActivationHeight {
				SoftForks[idx].State = explorer.SoftForkActive
				SoftForks[idx].Dirty = true
			}
		}
	}
}

func (i *Indexer) persist(height uint64) {
	for _, s := range SoftForks {
		if s.Dirty == false {
			continue
		}

		i.elastic.AddRequest(index.SoftForkIndex.Get(), s.Name, s)
	}
}
