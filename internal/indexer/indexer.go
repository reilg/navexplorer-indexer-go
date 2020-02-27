package indexer

import (
	"github.com/NavExplorer/navexplorer-indexer-go/internal/elastic_cache"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/indexer/IndexOption"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/address"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/block"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/dao"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/softfork"
	"github.com/getsentry/raven-go"
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

var (
	LastBlockIndexed uint64 = 0
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
	if err := i.Index(IndexOption.BatchIndex); err != nil {
		if err.Error() == "-8: Block height out of range" {
			i.elastic.Persist()
		} else {
			log.WithError(err).Fatal("Failed to index blocks")
		}
	}
}

func (i *Indexer) Index(option IndexOption.IndexOption) error {
	err := i.index(LastBlockIndexed+1, option)
	if err != block.ErrOrphanBlockFound {
		raven.CaptureError(err, nil)
		log.WithError(err).Error("Failed to index block")
		// Unexpected indexing error
		return err
	}

	if err := i.rewinder.RewindToHeight(LastBlockIndexed - 9); err != nil {
		raven.CaptureError(err, nil)
		// Unable to rewind blocks
		return err
	}

	return i.Index(option)
}

func (i *Indexer) index(height uint64, option IndexOption.IndexOption) error {
	b, txs, err := i.blockIndexer.Index(height, option)
	if err != nil {
		raven.CaptureError(err, nil)
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

	if option == IndexOption.BatchIndex {
		i.elastic.BatchPersist(height)
	} else {
		i.elastic.Persist()
		log.Infof("Indexed height: %d", height)
	}

	return i.index(height+1, option)
}
