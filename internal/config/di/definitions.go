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
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/dao/payment_request"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/dao/proposal"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/dao/vote"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/softfork"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/softfork/signal"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/zeromq"
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
				config.Get().Network,
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
		Name: "address.repo",
		Build: func(elastic *elastic_cache.Index) (*address.Repository, error) {
			return address.NewRepo(elastic.Client), nil
		},
	},
	{
		Name: "block.repo",
		Build: func(elastic *elastic_cache.Index) (*block.Repository, error) {
			return block.NewRepo(elastic.Client), nil
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
			rewinder *indexer.Rewinder,
		) (*indexer.Indexer, error) {
			return indexer.NewIndexer(redis, elastic, blockIndexer, addressIndexer, softForkIndexer, daoIndexer, rewinder), nil
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
		Build: func(elastic *elastic_cache.Index, repo *address.Repository) (*address.Indexer, error) {
			return address.NewIndexer(elastic, repo), nil
		},
	},
	{
		Name: "address.rewinder",
		Build: func(elastic *elastic_cache.Index, repo *address.Repository) (*address.Rewinder, error) {
			return address.NewRewinder(elastic, repo), nil
		},
	},
	{
		Name: "block.orphan.service",
		Build: func(navcoin *navcoind.Navcoind) (*block.OrphanService, error) {
			return block.NewOrphanService(navcoin), nil
		},
	},
	{
		Name: "softfork.indexer",
		Build: func(navcoin *navcoind.Navcoind, elastic *elastic_cache.Index, repo *softfork.Repository) (*softfork.Indexer, error) {
			return softfork.NewIndexer(elastic, config.Get().SoftForkBlockCycle, config.Get().SoftForkQuorum), nil
		},
	},
	{
		Name: "softfork.rewinder",
		Build: func(elastic *elastic_cache.Index, signalRepo *signal.Repository) (*softfork.Rewinder, error) {
			return softfork.NewRewinder(elastic, signalRepo, config.Get().SoftForkBlockCycle, config.Get().SoftForkQuorum), nil
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
		Build: func(proposalIndexer *proposal.Indexer, paymentRequestIndexer *payment_request.Indexer, voteIndexer *vote.Indexer) (*dao.Indexer, error) {
			return dao.NewIndexer(
				proposalIndexer,
				paymentRequestIndexer,
				voteIndexer,
				config.Get().DaoCfundConsensus.BlocksPerVotingCycle,
				config.Get().DaoCfundConsensus.Quorum,
			), nil
		},
	},
	{
		Name: "dao.proposal.Indexer",
		Build: func(navcoin *navcoind.Navcoind, elastic *elastic_cache.Index) (*proposal.Indexer, error) {
			return proposal.NewIndexer(
				navcoin,
				elastic,
				uint64(config.Get().ReindexSize),
			), nil
		},
	},
	{
		Name: "dao.payment_request.Indexer",
		Build: func(navcoin *navcoind.Navcoind, elastic *elastic_cache.Index) (*payment_request.Indexer, error) {
			return payment_request.NewIndexer(
				navcoin,
				elastic,
			), nil
		},
	},
	{
		Name: "dao.vote.Indexer",
		Build: func(elastic *elastic_cache.Index) (*vote.Indexer, error) {
			return vote.NewIndexer(
				elastic,
				config.Get().DaoCfundConsensus.MaxCountVotingCycleProposals,
			), nil
		},
	},
	{
		Name: "dao.Rewinder",
		Build: func(elastic *elastic_cache.Index) (*dao.Rewinder, error) {
			return dao.NewRewinder(elastic), nil
		},
	},
	{
		Name: "dao.proposal.Service",
		Build: func(repo *proposal.Repository) (*proposal.Service, error) {
			return proposal.NewService(repo), nil
		},
	},
	{
		Name: "dao.payment_request.Service",
		Build: func(repo *payment_request.Repository) (*payment_request.Service, error) {
			return payment_request.NewService(repo), nil
		},
	},
	{
		Name: "dao.payment_request.repo",
		Build: func(elastic *elastic_cache.Index) (*payment_request.Repository, error) {
			return payment_request.NewRepo(elastic.Client), nil
		},
	},
	{
		Name: "dao.proposal.repo",
		Build: func(elastic *elastic_cache.Index) (*proposal.Repository, error) {
			return proposal.NewRepo(elastic.Client), nil
		},
	},
	{
		Name: "subscriber",
		Build: func(indexer *indexer.Indexer) (*zeromq.Subscriber, error) {
			return zeromq.New(config.Get().ZeroMq.Address, indexer), nil
		},
	},
}
