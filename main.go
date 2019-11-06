package main

import (
	"github.com/NavExplorer/navexplorer-indexer-go/generated/dic"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/config"
	"github.com/sarulabs/dingo/v3"
	log "github.com/sirupsen/logrus"
)

var container *dic.Container

func main() {
	setup()

	if err := container.GetBlockIndexer().IndexBlocks(); err != nil {
		log.WithError(err).Fatal("Failed to index blocks")
	}
}

func setup() {
	config.Init()
	container, _ = dic.NewContainer(dingo.App)

	height, err := container.GetRedis().Start()
	if err != nil {
		log.WithError(err).Fatal("Failed to rewind redis")
	}

	container.GetElastic().RewindTo(height)
	container.GetSoftforkIndexer().Init().RewindTo(height)
}
