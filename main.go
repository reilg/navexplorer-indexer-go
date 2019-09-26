package main

import (
	"github.com/NavExplorer/navexplorer-indexer-go/internal/config"
	index "github.com/NavExplorer/navexplorer-indexer-go/internal/indexer"
)

func main() {
	config.Init()
	indexer := index.New()
	defer indexer.Close()

	//go indexer.SubscribeBlockIndexed()

	indexer.IndexBlocks()
}
