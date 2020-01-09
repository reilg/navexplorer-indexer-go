package indexer

import (
	"github.com/NavExplorer/navexplorer-indexer-go/internal/elastic_cache"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/address"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/block"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/dao"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/softfork"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Indexer struct {
	elastic         *elastic_cache.Index
	blockIndexer    *block.Indexer
	addressIndexer  *address.Indexer
	softForkIndexer *softfork.Indexer
	daoIndexer      *dao.Indexer
	rewinder        *Rewinder
}

type IndexOption int

var (
	SingleIndex      IndexOption = 1
	BatchIndex       IndexOption = 2
	LastBlockIndexed uint64      = 0
)

func NewIndexer(
	elastic *elastic_cache.Index,
	blockIndexer *block.Indexer,
	addressIndexer *address.Indexer,
	softForkIndexer *softfork.Indexer,
	daoIndexer *dao.Indexer,
	rewinder *Rewinder,
) *Indexer {
	return &Indexer{
		elastic,
		blockIndexer,
		addressIndexer,
		softForkIndexer,
		daoIndexer,
		rewinder,
	}
}

func (i *Indexer) BulkIndex() {
	if err := i.Index(BatchIndex); err != nil {
		if err.Error() == "-8: Block height out of range" {
			i.elastic.Persist()
		} else {
			log.WithError(err).Fatal("Failed to index blocks")
		}
	}
}

func (i *Indexer) Index(option IndexOption) error {
	if err := i.index(LastBlockIndexed+1, option); err != block.ErrOrphanBlockFound {
		// Unexpected indexing error
		return err
	}

	if err := i.rewinder.RewindToHeight(LastBlockIndexed - 9); err != nil {
		// Unable to rewind blocks
		return err
	}

	return i.Index(option)
}

func (i *Indexer) index(height uint64, option IndexOption) error {
	b, txs, err := i.blockIndexer.Index(height, int(option))
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
		if softfork.SoftForks.GetSoftFork("communityfund").State == "active" {
			i.daoIndexer.Index(b, txs)
		}
	}()

	wg.Wait()

	LastBlockIndexed = height

	if option == BatchIndex {
		i.elastic.BatchPersist(height)
	} else {
		i.elastic.Persist()
		log.Infof("Indexed height: %d", height)
	}

	return i.index(height+1, option)
}
