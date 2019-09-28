package signal_indexer

import (
	"context"
	"fmt"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/index"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/indexer/softfork_indexer"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Indexer struct {
	elastic *index.Index
}

func New(elastic *index.Index) *Indexer {
	return &Indexer{elastic}
}

func (i *Indexer) IndexSignal(block *explorer.Block) {
	signal := explorer.Signal{
		Address: block.MetaData.StakedBy,
		Height:  block.Height,
		Signals: make([]string, 0),
	}
	for _, sf := range softfork_indexer.SoftForks {
		if sf.IsOpen() && block.Version>>sf.SignalBit&1 == 1 {
			signal.Signals = append(signal.Signals, sf.Name)
		}
	}

	if len(signal.Signals) == 0 {
		return
	}

	log.WithFields(log.Fields{"height": signal.Height}).Info("Signals Indexed: ", strings.Join(signal.Signals, ","))

	_, err := i.elastic.Client.Index().
		Index(index.SignalIndex.Get()).
		Id(fmt.Sprintf("%d-%s", block.Height, block.MetaData.StakedBy)).
		BodyJson(signal).
		Do(context.Background())

	if err != nil {
		log.WithError(err).Fatal("Failed to index signal at height: ", block.Height)
	}

}

//public void indexBlock(Block block) {
//List<SoftFork> softForks = softForkRepository.findAll();
//
//softForks.forEach(softFork -> {
//boolean signalling = (block.getVersion() >> softFork.getSignalBit() & 1) == 1;
//
//block.getSignals().add(new Signal(softFork.getName(), signalling));
//});
//
//block.setBlockCycle(
//((Double) Math.ceil(block.getHeight().intValue() / blocksInCycle)).intValue() + 1
//);
//blockRepository.save(block);
//
//updateSoftForks(block);
//}
