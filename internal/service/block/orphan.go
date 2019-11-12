package block

import (
	"errors"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
)

type OrphanService struct {
}

func NewOrphanService() *OrphanService {
	return &OrphanService{}
}

var (
	ErrOrphanBlockFound = errors.New("Orphan block_indexer found")
)

func (o *OrphanService) IsOrphanBlock(block *explorer.Block) (bool, error) {
	if block.Height == 1 {
		return false, nil
	}
	//
	//previousBlockJson, err := i.elastic.Client.Get().
	//	Index(elastic_cache.BlockTransactionIndex.Get()).
	//	Id(block.Previousblockhash).
	//	Do(context.Background())
	//if err != nil {
	//	return false, err
	//}
	//
	//var previousBlock navcoind.Block
	//if err := json.Unmarshal(previousBlockJson.Source, &previousBlock); err != nil {
	//	return false, err
	//}
	//
	//orphan := previousBlock.Hash != block.Previousblockhash
	//if orphan == true {
	//	log.Printf("INFO: Orphan block_indexer found: %s", block.Hash)
	//}
	//
	//return orphan, nil

	return false, nil
}
