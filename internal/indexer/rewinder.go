package indexer

import (
	"github.com/NavExplorer/navexplorer-indexer-go/internal/elastic_cache"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/redis"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/address"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/block"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/dao"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/softfork"
	log "github.com/sirupsen/logrus"
)

type Rewinder struct {
	redis            *redis.Redis
	elastic          *elastic_cache.Index
	blockRewinder    *block.Rewinder
	addressRewinder  *address.Rewinder
	softforkRewinder *softfork.Rewinder
	daoRewinder      *dao.Rewinder
}

func NewRewinder(
	redis *redis.Redis,
	elastic *elastic_cache.Index,
	blockRewinder *block.Rewinder,
	addressRewinder *address.Rewinder,
	softforkRewinder *softfork.Rewinder,
	daoRewinder *dao.Rewinder,
) *Rewinder {
	return &Rewinder{
		redis,
		elastic,
		blockRewinder,
		addressRewinder,
		softforkRewinder,
		daoRewinder,
	}
}

func (r *Rewinder) RewindToHeight(height uint64) error {
	log.Infof("Rewinding to height: %d", height)

	if err := r.addressRewinder.Rewind(height); err != nil {
		return err
	}
	if err := r.blockRewinder.Rewind(height); err != nil {
		return err
	}
	if err := r.softforkRewinder.Rewind(height); err != nil {
		return err
	}
	if err := r.daoRewinder.Rewind(height); err != nil {
		return err
	}

	if err := r.redis.SetLastBlock(height); err != nil {
		log.WithError(err).Fatal("Failed to set last block indexed")
		return err
	}

	return nil
}
