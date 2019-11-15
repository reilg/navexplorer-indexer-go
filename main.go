package main

import (
	"github.com/NavExplorer/navexplorer-indexer-go/generated/dic"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/config"
	"github.com/sarulabs/dingo/v3"
	log "github.com/sirupsen/logrus"
)

var container *dic.Container

func main() {
	log.Info("Launching Indexer")
	lastBlock := setup()

	if err := container.GetIndexer().Index(lastBlock+1, true); err != nil {
		if err.Error() == "-8: Block height out of range" {
			log.Info("Persist any pending requests")
			container.GetElastic().Persist()
		} else {
			log.WithError(err).Fatal("Failed to index blocks")
		}
	}
}

func setup() uint64 {
	config.Init()
	container, _ = dic.NewContainer(dingo.App)
	log.Info("Container init")

	container.GetElastic().InstallMappings()

	height, err := container.GetRedis().Start()
	if err != nil {
		log.WithError(err).Fatal("Failed to start redis")
	}

	//height = 2000000
	container.GetSoftforkService().LoadSoftForks()

	if height < uint64(config.Get().BulkIndexSize) {
		height = 0
	} else {
		height = height - uint64(config.Get().BulkIndexSize)
	}

	if err := container.GetRewinder().RewindToHeight(height); err != nil {
		log.WithError(err).Fatal("Failed to rewind index")
	}

	return height
}
