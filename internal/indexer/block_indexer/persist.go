package block_indexer

import (
	"github.com/NavExplorer/navexplorer-indexer-go/internal/events"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/index"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	"github.com/gookit/event"
	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
)

func (i *Indexer) persist(txs *[]explorer.BlockTransaction, block *explorer.Block) error {
	i.elastic.GetBulkRequest(block.Height).Add(elastic.NewBulkIndexRequest().
		Index(index.BlockIndex.Get()).
		Id(block.Hash).
		Doc(block))

	for _, tx := range *txs {
		i.elastic.GetBulkRequest(block.Height).Add(elastic.NewBulkIndexRequest().
			Index(index.BlockTransactionIndex.Get()).
			Id(tx.Hash).
			Doc(tx))
	}

	if err := i.cache.SetLastBlock(block.Height); err != nil {
		log.WithError(err).Fatal("Failed to set last block_indexer indexed")
		return err
	}

	event.MustFire(string(events.EventBlockIndexed), event.M{"block": block, "txs": txs})

	return nil
}
