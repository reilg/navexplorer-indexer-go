package vote

import (
	"fmt"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/elastic_cache"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/dao/payment_request"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/service/dao/proposal"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
)

type Indexer struct {
	elastic *elastic_cache.Index
}

func NewIndexer(elastic *elastic_cache.Index) *Indexer {
	return &Indexer{elastic}
}

func (i *Indexer) InitCycles(blockCycle explorer.BlockCycle) {
	for _, p := range proposal.Proposals {
		var cycle *explorer.CfundCycle
		if cycle = p.GetCycle(blockCycle.Cycle); cycle == nil {
			p.Cycles = append(p.Cycles, explorer.CfundCycle{VotingCycle: blockCycle.Cycle})
		}
	}
	for _, p := range payment_request.PaymentRequests {
		var cycle *explorer.CfundCycle
		if cycle = p.GetCycle(blockCycle.Cycle); cycle == nil {
			p.Cycles = append(p.Cycles, explorer.CfundCycle{VotingCycle: blockCycle.Cycle})
		}
	}
}

func (i *Indexer) IndexVotes(txs []*explorer.BlockTransaction, block *explorer.Block) *explorer.DaoVote {
	for _, tx := range txs {
		if !tx.IsCoinbase() {
			continue
		}

		if v := CreateVotes(block, tx); v != nil {
			i.elastic.AddIndexRequest(elastic_cache.DaoVoteIndex.Get(), fmt.Sprintf("%d-%s", v.Height, v.Address), v)
			return v
		}
	}
	return nil
}
