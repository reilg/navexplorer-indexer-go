package main

import (
	"github.com/NavExplorer/navexplorer-indexer-go/internal/cli"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Launching Indexer")

	if err := cli.Execute(); err != nil {
		log.WithError(err).Fatal(err)
	}
}
