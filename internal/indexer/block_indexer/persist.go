package block_indexer

import (
	"context"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/events"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/index"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	"github.com/gookit/event"
	log "github.com/sirupsen/logrus"
)

func (i *Indexer) persist(txs *[]explorer.BlockTransaction, block *explorer.Block) error {
	if err := i.persistBlockTransactions(txs, block); err != nil {
		log.WithError(err).Error("Failed to persist block_indexer transactions")
		return err
	}

	if err := i.persistBlock(*block); err != nil {
		log.WithError(err).Error("Failed to persist block_indexer")
		return err
	}

	if err := i.cache.SetLastBlock(block.Height); err != nil {
		log.WithError(err).Error("Failed to set last block_indexer indexed")
		return err
	}

	i.elastic.Flush([]string{
		index.BlockIndex.Get(),
		index.BlockTransactionIndex.Get(),
	}...)

	log.WithFields(log.Fields{"height": block.Height}).Info("Block Indexed")

	go event.MustFire(string(events.EventBlockIndexed), event.M{"block": block, "txs": txs})

	return nil
}

func (i *Indexer) persistBlockTransactions(txs *[]explorer.BlockTransaction, block *explorer.Block) error {
	for _, tx := range *txs {
		_, err := i.elastic.Client.
			Index().
			Index(index.BlockTransactionIndex.Get()).
			Id(tx.Hash).
			BodyJson(tx).
			Do(context.Background())
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Indexer) persistBlock(block explorer.Block) error {
	_, err := i.elastic.Client.
		Index().
		Index(index.BlockIndex.Get()).
		Id(block.Hash).
		BodyJson(block).
		Do(context.Background())
	return err
}
