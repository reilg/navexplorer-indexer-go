package daemon

import (
	"github.com/NavExplorer/navexplorer-indexer-go/generated/dic"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/config"
	"github.com/sarulabs/dingo/v3"
	log "github.com/sirupsen/logrus"
)

var container *dic.Container

func Execute() {
	config.Init()

	container, _ = dic.NewContainer(dingo.App)

	container.GetElastic().InstallMappings()
	container.GetSoftforkService().LoadSoftForks()

	if err := container.GetRewinder().RewindToHeight(getHeight()); err != nil {
		log.WithError(err).Fatal("Failed to rewind index")
	}

	container.GetIndexer().BulkIndex()
	container.GetSubscriber().Subscribe()
}

func getHeight() uint64 {
	if height, err := container.GetRedis().Start(); err != nil {
		log.WithError(err).Fatal("Failed to start redis")
	} else {
		if height >= uint64(config.Get().BulkIndexSize) {
			return height - uint64(config.Get().BulkIndexSize)
		}
	}

	return 0
}
