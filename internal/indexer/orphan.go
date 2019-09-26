package indexer

import (
	"context"
	"encoding/json"
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/entity"
	"log"
)

func (i *Indexer) isOrphanBlock(block entity.Block) (bool, error) {
	if block.Height == 1 {
		return false, nil
	}

	previousBlockJson, err := i.Elastic.Client.Get().
		Index(BlockTransactionIndex.Get(i.Network)).
		Id(block.Previousblockhash).
		Do(context.Background())
	if err != nil {
		return false, err
	}

	var previousBlock navcoind.Block
	if err := json.Unmarshal(previousBlockJson.Source, &previousBlock); err != nil {
		return false, err
	}

	orphan := previousBlock.Hash != block.Previousblockhash
	if orphan == true {
		log.Printf("INFO: Orphan block found: %s", block.Hash)
	}

	return orphan, nil
}
