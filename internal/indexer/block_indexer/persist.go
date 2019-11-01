package block_indexer

import (
	"github.com/NavExplorer/navexplorer-indexer-go/internal/events"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/index"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	"github.com/gookit/event"
	log "github.com/sirupsen/logrus"
)

func (i *Indexer) persist(txs *[]explorer.BlockTransaction, block *explorer.Block) error {
	i.elastic.AddIndexRequest(index.BlockIndex.Get(), block.Hash, block)
	for _, tx := range *txs {
		i.elastic.AddIndexRequest(index.BlockTransactionIndex.Get(), tx.Hash, tx)
	}

	if err := i.cache.SetLastBlock(block.Height); err != nil {
		log.WithError(err).Fatal("Failed to set last block_indexer indexed")
		return err
	}

	event.MustFire(string(events.EventBlockIndexed), event.M{"block": block, "txs": txs})

	return nil
}
