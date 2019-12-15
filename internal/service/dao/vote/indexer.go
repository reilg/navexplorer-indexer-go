package vote

import (
	"fmt"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/elastic_cache"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
)

type Indexer struct {
	elastic           *elastic_cache.Index
	maxProposalCycles uint
}

func NewIndexer(elastic *elastic_cache.Index, maxProposalCycles uint) *Indexer {
	return &Indexer{elastic, maxProposalCycles}
}

func (i *Indexer) IndexVotes(txs []*explorer.BlockTransaction, block *explorer.Block) {
	for _, tx := range txs {
		if !tx.IsCoinbase() {
			continue
		}

		if v := CreateVotes(block, tx); v != nil {
			i.elastic.AddIndexRequest(elastic_cache.DaoVoteIndex.Get(), fmt.Sprintf("%d-%s", v.Height, v.Address), v)
			return
		}
	}
}
