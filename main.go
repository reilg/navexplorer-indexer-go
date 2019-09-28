package main

import (
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/config"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/events"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/index"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/indexer/address_indexer"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/indexer/block_indexer"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/indexer/signal_indexer"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/indexer/softfork_indexer"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/redis"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	"github.com/gookit/event"
	log "github.com/sirupsen/logrus"
)

func main() {
	config.Init()

	elastic, _ := index.New()
	if err := elastic.Init(); err != nil {
		log.WithError(err).Fatal("Failed to initialize Elastic Search")
	}

	cache := redis.New()
	if err := cache.Init(); err != nil {
		log.WithError(err).Fatal("Failed to initialize Redis")
	}

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

	softfork_indexer.New(elastic, navcoin).Init()

	go eventSubscription(elastic)

	blockIndexer := block_indexer.New(elastic, cache, navcoin)
	defer blockIndexer.Close()

	if err := blockIndexer.IndexBlocks(); err != nil {
		log.WithError(err).Fatal("Failed to index blocks")
	}
}

func eventSubscription(elastic *index.Index) {
	elastic, _ = index.New()
	addressIndexer := address_indexer.New(elastic)
	event.On(string(events.EventBlockIndexed), event.ListenerFunc(func(e event.Event) error {
		go addressIndexer.IndexAddressesForTransactions(e.Get("txs").(*[]explorer.BlockTransaction))
		return nil
	}), event.Normal)

	elastic, _ = index.New()
	signalIndexer := signal_indexer.New(elastic)
	event.On(string(events.EventBlockIndexed), event.ListenerFunc(func(e event.Event) error {
		go signalIndexer.IndexSignal(e.Get("block").(*explorer.Block))
		return nil
	}), event.Normal)
}
