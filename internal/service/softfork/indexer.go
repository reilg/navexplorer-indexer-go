package softfork

import (
	"fmt"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/elastic_cache"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/softfork/signal"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
)

type Indexer struct {
	elastic       *elastic_cache.Index
	blocksInCycle uint
}

func NewIndexer(elastic *elastic_cache.Index, blocksInCycle uint) *Indexer {
	return &Indexer{elastic, blocksInCycle}
}

func (i *Indexer) Index(block *explorer.Block) {
	s := signal.CreateSignal(block, &SoftForks)

	if s != nil {
		i.elastic.AddIndexRequest(elastic_cache.SignalIndex.Get(), fmt.Sprintf("%d", block.Height), s)
		i.updateSoftForks(s, block)
	}

	i.updateState(block)

	for _, softFork := range SoftForks {
		i.elastic.AddUpdateRequest(elastic_cache.SoftForkIndex.Get(), softFork.Name, softFork, softFork.MetaData.Id)
	}

}

func (i *Indexer) updateSoftForks(signal *explorer.Signal, block *explorer.Block) {
	if signal == nil || !signal.IsSignalling() {
		return
	}

	blockCycle := block.BlockCycle(i.blocksInCycle)

	for _, s := range signal.SoftForks {
		softFork := SoftForks.GetSoftFork(s)
		if softFork == nil || softFork.SignalHeight >= signal.Height {
			continue
		}

		softFork.SignalHeight = signal.Height
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

func (i *Indexer) updateState(block *explorer.Block) {
	blockCycle := explorer.GetCycleForHeight(block.Height, i.blocksInCycle)

	for idx, s := range SoftForks {
		if s.Cycles == nil {
			continue
		}

		if s.State == explorer.SoftForkStarted && s.LatestCycle().BlocksSignalling >= explorer.GetQuorum(i.blocksInCycle) {
			SoftForks[idx].LockedInHeight = uint64(i.blocksInCycle * blockCycle)
			SoftForks[idx].ActivationHeight = SoftForks[idx].LockedInHeight + uint64(i.blocksInCycle)
		}

		if s.LockedInHeight != 0 && s.ActivationHeight != 0 {
			if s.State == explorer.SoftForkStarted && block.Height >= s.LockedInHeight {
				SoftForks[idx].State = explorer.SoftForkLockedIn
			}
			if s.State == explorer.SoftForkLockedIn && block.Height >= s.ActivationHeight {
				SoftForks[idx].State = explorer.SoftForkActive
			}
		}
	}
}
