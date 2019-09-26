package indexer

import (
	"errors"
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/config"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/events"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/lib/elastic"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/lib/redis"
	log "github.com/sirupsen/logrus"
)

type Indexer struct {
	Navcoin *navcoind.Navcoind
	Elastic *elastic.Elastic
	Redis   *redis.Redis
	Events  *events.Events

	Debug      bool
	Ready      bool
	Reindex    bool
	Network    string
	MappingDir string
}

var (
	ErrLastBlockReadIndexed        = errors.New("Failed to read last block indexed")
	ErrLastBlockWriteIndexed       = errors.New("Failed to write last block indexed")
	ErrIndexerNotReady             = errors.New("Indexer not ready")
	ErrUnableToUpdatePreviousBlock = errors.New("Unable to update previous block")
	ErrUnableToIndexBlock          = errors.New("Unable to index block")
	ErrInvalidBlock                = errors.New("Invalid block detected")
	ErrPreviousBlockNotIndexed     = errors.New("Previous Block not indexed")
	ErrOrphanBlockFound            = errors.New("Orphan Block Found")
)

func New() *Indexer {
	i := &Indexer{
		Debug:      config.Get().Debug,
		Ready:      false,
		Reindex:    config.Get().Reindex,
		Network:    config.Get().Navcoind.Network,
		MappingDir: config.Get().BaseDir + "/mappings",
	}

	i.Navcoin, _ = navcoind.New(
		config.Get().Navcoind.Host,
		config.Get().Navcoind.Port,
		config.Get().Navcoind.User,
		config.Get().Navcoind.Password,
		config.Get().Navcoind.Ssl,
	)

	i.Elastic, _ = elastic.New(
		config.Get().ElasticSearch.Hosts,
		config.Get().ElasticSearch.Sniff,
		config.Get().ElasticSearch.HealthCheck,
		config.Get().ElasticSearch.Debug,
	)

	i.Redis = redis.New(config.Get().Redis.Host, config.Get().Redis.Password, config.Get().Redis.Db)

	i.Events = events.New()

	return i.init()
}

func (i *Indexer) init() *Indexer {
	log.Info("Initialize Indexer")

	i.InstallMappings()

	if config.Get().Reindex {
		if err := i.Purge(); err != nil {
			log.WithError(err).Fatal("Failed to reindex blocks")
		}
	} else {
		if err := i.RewindBy(10); err != nil {
			log.WithError(err).Fatal("Failed to rewind blocks")
		}
	}

	i.Ready = true

	return i
}

func (i *Indexer) Close() {
	log.Println("Closing Indexer")
}

func (i *Indexer) isReady() bool {
	return i.Ready
}

func (i *Indexer) Purge() error {
	log.Info("Purging the index")
	//if err := i.Queue.Purge(); err != nil {
	//	log.Fatal(err.Error())
	//}

	return i.setLastBlock(0)
}
