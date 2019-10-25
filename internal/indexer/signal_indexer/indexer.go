package signal_indexer

import (
	"fmt"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/events"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/index"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/indexer/softfork_indexer"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	"github.com/gookit/event"
)

type Indexer struct {
	elastic *index.Index
}

func New(elastic *index.Index) *Indexer {
	return &Indexer{elastic}
}

func (i *Indexer) IndexSignal(block *explorer.Block) {
	signal := explorer.Signal{
		Address:   block.MetaData.StakedBy,
		Height:    block.Height,
		SoftForks: make([]string, 0),
	}
	for _, sf := range softfork_indexer.SoftForks {
		if sf.IsOpen() && block.Version>>sf.SignalBit&1 == 1 {
			signal.SoftForks = append(signal.SoftForks, sf.Name)
		}
	}

	if len(signal.SoftForks) != 0 {
		id := fmt.Sprintf("%d-%s", block.Height, block.MetaData.StakedBy)
		i.elastic.AddRequest(index.SignalIndex.Get(), id, signal)
	}

	event.MustFire(string(events.EventSignalIndexed), event.M{"signal": &signal, "block": block})
}
