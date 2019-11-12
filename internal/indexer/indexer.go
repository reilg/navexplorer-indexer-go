package indexer

import (
	"github.com/NavExplorer/navexplorer-indexer-go/internal/elastic_cache"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/redis"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/address"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/block"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/dao"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/softfork"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Indexer struct {
	redis           *redis.Redis
	elastic         *elastic_cache.Index
	blockIndexer    *block.Indexer
	addressIndexer  *address.Indexer
	softForkIndexer *softfork.Indexer
	daoIndexer      *dao.Indexer
}

func NewIndexer(
	redis *redis.Redis,
	elastic *elastic_cache.Index,
	blockIndexer *block.Indexer,
	addressIndexer *address.Indexer,
	softForkIndexer *softfork.Indexer,
	daoIndexer *dao.Indexer,
) *Indexer {
	return &Indexer{
		redis,
		elastic,
		blockIndexer,
		addressIndexer,
		softForkIndexer,
		daoIndexer,
	}
}

func (i *Indexer) Index(height uint64) error {
	b, txs, err := i.blockIndexer.Index(height)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		i.addressIndexer.Index(txs)
	}()

	go func() {
		defer wg.Done()
		i.softForkIndexer.Index(b)
	}()

	go func() {
		defer wg.Done()
		i.daoIndexer.Index(b, txs)
	}()

	wg.Wait()

	if err := i.redis.SetLastBlock(height); err != nil {
		log.WithError(err).Fatal("Failed to set last block indexed")
		return err
	}

	i.elastic.PersistRequests(height)

	return i.Index(height + 1)
}
