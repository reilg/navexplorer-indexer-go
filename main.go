package main

import (
	"github.com/NavExplorer/navexplorer-indexer-go/internal/app/config"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/app/index"
	"log"
)

func main() {
	config.Init()

	indexer := index.New()
	defer indexer.Close()

	log.Println("Indexing all blocks")
	indexer.IndexAllBlocks()
}
