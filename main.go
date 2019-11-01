package main

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
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	"github.com/gookit/event"
	log "github.com/sirupsen/logrus"
)

func main() {
	config.Init()

	redis, lastBlock := initRedis()
	elastic := initElasticCache(lastBlock)
	navcoin := initNavcoin()

	softForkRepo := repository.NewSoftForkRepo(elastic.Client)
	softforkIndexer := softfork_indexer.New(elastic, navcoin, &softForkRepo).Init()
	softforkIndexer.RewindTo(lastBlock)

	go eventSubscription(navcoin, elastic)

	blockIndexer := block_indexer.New(elastic, redis, navcoin)
	if err := blockIndexer.IndexBlocks(); err != nil {
		log.WithError(err).Fatal("Failed to index blocks")
	}
}

func eventSubscription(navcoin *navcoind.Navcoind, elastic *index.Index) {
	event.On(string(events.EventBlockIndexed), event.ListenerFunc(func(e event.Event) error {
		block := e.Get("block").(*explorer.Block)
		txs := e.Get("txs").(*[]explorer.BlockTransaction)

		address_indexer.New(elastic).IndexAddressesForTransactions(txs)
		dao_indexer.NewProposalIndexer(navcoin, elastic).IndexProposals(txs)
		dao_indexer.NewPaymentRequestIndexer(navcoin, elastic).IndexPaymentRequests(txs)
		signal_indexer.New(elastic).IndexSignal(block)
		return nil
	}), event.Normal)

	softForkRepo := repository.NewSoftForkRepo(elastic.Client)
	event.On(string(events.EventSignalIndexed), event.ListenerFunc(func(e event.Event) error {
		signal := e.Get("signal").(*explorer.Signal)
		block := e.Get("block").(*explorer.Block)

		softfork_indexer.New(elastic, navcoin, &softForkRepo).Update(signal, block)
		elastic.PersistRequests(signal.Height)
		return nil
	}), event.Normal)
}

func initElasticCache(lastBlock uint64) *index.Index {
	elastic, err := index.New()
	if err != nil {
		log.WithError(err).Fatal("Failed toStart ES")
	}
	if err := elastic.Init(); err != nil {
		log.WithError(err).Fatal("Failed to initialize ES")
	}

	elastic.RewindTo(lastBlock)

	return elastic
}

func initNavcoin() *navcoind.Navcoind {
	navcoin, err := navcoind.New(
		config.Get().Navcoind.Host,
		config.Get().Navcoind.Port,
		config.Get().Navcoind.User,
		config.Get().Navcoind.Password,
		config.Get().Navcoind.Ssl,
	)

	if err != nil {
		log.WithError(err).Fatal("Failed to initialize Navcoind")
	}

	return navcoin
}

func initRedis() (*redis.Redis, uint64) {
	cache := redis.New()
	lastBlock, err := cache.Init()
	if err != nil {
		log.WithError(err).Fatal("Failed to initialize Redis")
	}

	return cache, lastBlock
}
