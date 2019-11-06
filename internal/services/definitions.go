package services

import (
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/config"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/events"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/index"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/indexer/address_indexer"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/indexer/block_indexer"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/indexer/dao_indexer"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/indexer/signal_indexer"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/indexer/softfork_indexer"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/redis"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/repository"
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
		Build: func() (*index.Index, error) {
			elastic, err := index.New()
			if err != nil {
				log.WithError(err).Fatal("Failed toStart ES")
			}

			return elastic, nil
		},
	},
	{
		Name: "softfork.repo",
		Build: func(elastic *index.Index) (*repository.SoftForkRepository, error) {
			return repository.NewSoftForkRepo(elastic.Client), nil
		},
	},
	{
		Name: "address.indexer",
		Build: func(elastic *index.Index) (*address_indexer.Indexer, error) {
			return address_indexer.New(elastic), nil
		},
	},
	{
		Name: "block.indexer",
		Build: func(elastic *index.Index, cache *redis.Redis, navcoin *navcoind.Navcoind) (*block_indexer.Indexer, error) {
			return block_indexer.New(elastic, cache, navcoin), nil
		},
	},
	{
		Name: "dao.proposal.indexer",
		Build: func(navcoin *navcoind.Navcoind, elastic *index.Index) (*dao_indexer.ProposalIndexer, error) {
			return dao_indexer.NewProposalIndexer(navcoin, elastic), nil
		},
	},
	{
		Name: "dao.payment-request.indexer",
		Build: func(navcoin *navcoind.Navcoind, elastic *index.Index) (*dao_indexer.PaymentRequestIndexer, error) {
			return dao_indexer.NewPaymentRequestIndexer(navcoin, elastic), nil
		},
	},
	{
		Name: "softfork.indexer",
		Build: func(elastic *index.Index, navcoin *navcoind.Navcoind, repo *repository.SoftForkRepository) (*softfork_indexer.Indexer, error) {
			return softfork_indexer.New(elastic, navcoin, repo), nil
		},
	},
	{
		Name: "signal.indexer",
		Build: func(elastic *index.Index) (*signal_indexer.Indexer, error) {
			return signal_indexer.New(elastic), nil
		},
	},
	{
		Name: "event.subscriber",
		Build: func(
			elastic *index.Index,
			addressIndexer *address_indexer.Indexer,
			proposalIndexer *dao_indexer.ProposalIndexer,
			paymentRequestIndexer *dao_indexer.PaymentRequestIndexer,
			signalIndexer *signal_indexer.Indexer,
			softForkIndexer *softfork_indexer.Indexer,
		) (*events.Subscriber, error) {
			return events.NewSubscriber(elastic, addressIndexer, proposalIndexer, paymentRequestIndexer, signalIndexer, softForkIndexer), nil
		},
	},
}
