package main

import (
	"fmt"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/config"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/elastic"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/events"
	index "github.com/NavExplorer/navexplorer-indexer-go/internal/indexer"
	"github.com/gookit/event"
	log "github.com/sirupsen/logrus"
)

func main() {
	config.Init()

	es, _ := elastic.New(
		config.Get().ElasticSearch.Hosts,
		config.Get().ElasticSearch.Sniff,
		config.Get().ElasticSearch.HealthCheck,
		config.Get().ElasticSearch.Debug,
	)

	addressIndexer := index.NewAddressIndexer(config.Get().Navcoind.Network, es)

	event.On(string(events.EventBlockIndexed), event.ListenerFunc(func(e event.Event) error {
		go addressIndexer.IndexAddressesForBlock(fmt.Sprintf("%v", e.Get("hash")))
		return nil
	}), event.Normal)

	indexer := index.New()
	defer indexer.Close()

	if err := indexer.IndexBlocks(); err != nil {
		log.WithError(err).Fatal("Failed to index blocks")
	}
}
