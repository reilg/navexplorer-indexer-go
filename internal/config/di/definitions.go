package di

import (
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/config"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/elastic_cache"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/indexer"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/redis"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/address"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/block"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/dao"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/softfork"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/softfork/signal"
	"github.com/sarulabs/dingo/v3"
	log "github.com/sirupsen/logrus"
)

var Definitions = []dingo.Def{
	{
		Name: "navcoin",
		Build: func() (*navcoind.Navcoind, error) {
			c := config.Get().Navcoind
			navcoin, err := navcoind.New(c.Host, c.Port, c.User, c.Password, c.Ssl)
			if err != nil {
				log.WithError(err).Fatal("Failed to initialize Navcoind")
			}
			return navcoin, nil
		},
	},
	{
		Name: "redis",
		Build: func() (*redis.Redis, error) {
			return redis.NewRedis(
				config.Get().Redis.Host,
				config.Get().Redis.Password,
				config.Get().Redis.Db,
				config.Get().ReindexSize), nil
		},
	},
	{
		Name: "elastic",
		Build: func() (*elastic_cache.Index, error) {
			elastic, err := elastic_cache.New()
			if err != nil {
				log.WithError(err).Fatal("Failed toStart ES")
			}

			return elastic, nil
		},
	},
	{
		Name: "signal.repo",
		Build: func(elastic *elastic_cache.Index) (*signal.Repository, error) {
			return signal.NewRepo(elastic.Client), nil
		},
	},
	{
		Name: "softfork.repo",
		Build: func(elastic *elastic_cache.Index) (*softfork.Repository, error) {
			return softfork.NewRepo(elastic.Client), nil
		},
	},
	{
		Name: "indexer",
		Build: func(
			redis *redis.Redis,
			elastic *elastic_cache.Index,
			blockIndexer *block.Indexer,
			addressIndexer *address.Indexer,
			softForkIndexer *softfork.Indexer,
			daoIndexer *dao.Indexer,
		) (*indexer.Indexer, error) {
			return indexer.NewIndexer(redis, elastic, blockIndexer, addressIndexer, softForkIndexer, daoIndexer), nil
		},
	},
	{
		Name: "rewinder",
		Build: func(
			redis *redis.Redis,
			elastic *elastic_cache.Index,
			addressRewinder *address.Rewinder,
			blockRewinder *block.Rewinder,
			softforkRewinder *softfork.Rewinder,
			daoRewinder *dao.Rewinder,
		) (*indexer.Rewinder, error) {
			return indexer.NewRewinder(redis, elastic, blockRewinder, addressRewinder, softforkRewinder, daoRewinder), nil
		},
	},
	{
		Name: "block.indexer",
		Build: func(navcoin *navcoind.Navcoind, elastic *elastic_cache.Index, orphanedService *block.OrphanService) (*block.Indexer, error) {
			return block.NewIndexer(navcoin, elastic, orphanedService), nil
		},
	},
	{
		Name: "block.rewinder",
		Build: func(elastic *elastic_cache.Index) (*block.Rewinder, error) {
			return block.NewRewinder(elastic), nil
		},
	},
	{
		Name: "address.indexer",
		Build: func(elastic *elastic_cache.Index) (*address.Indexer, error) {
			return address.NewIndexer(elastic), nil
		},
	},
	{
		Name: "address.rewinder",
		Build: func(elastic *elastic_cache.Index) (*address.Rewinder, error) {
			return address.NewRewinder(elastic), nil
		},
	},
	{
		Name: "block.orphan.service",
		Build: func() (*block.OrphanService, error) {
			return block.NewOrphanService(), nil
		},
	},
	{
		Name: "softfork.indexer",
		Build: func(navcoin *navcoind.Navcoind, elastic *elastic_cache.Index, repo *softfork.Repository) (*softfork.Indexer, error) {
			return softfork.NewIndexer(elastic, config.Get().SoftForkBlockCycle), nil
		},
	},
	{
		Name: "softfork.rewinder",
		Build: func(elastic *elastic_cache.Index, signalRepo *signal.Repository) (*softfork.Rewinder, error) {
			return softfork.NewRewinder(elastic, signalRepo, config.Get().SoftForkBlockCycle), nil
		},
	},
	{
		Name: "softfork.service",
		Build: func(navcoin *navcoind.Navcoind, elastic *elastic_cache.Index, repo *softfork.Repository) (*softfork.Service, error) {
			return softfork.New(navcoin, elastic, repo), nil
		},
	},
	{
		Name: "dao.Indexer",
		Build: func(navcoin *navcoind.Navcoind, elastic *elastic_cache.Index) (*dao.Indexer, error) {
			return dao.NewIndexer(navcoin, elastic), nil
		},
	},
	{
		Name: "dao.Rewinder",
		Build: func(elastic *elastic_cache.Index) (*dao.Rewinder, error) {
			return dao.NewRewinder(elastic), nil
		},
	},
}
