package address

import (
	"fmt"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/elastic_cache"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
)

type Indexer struct {
	elastic *elastic_cache.Index
}

func NewIndexer(elastic *elastic_cache.Index) *Indexer {
	return &Indexer{elastic}
}

func (i *Indexer) Index(txs []*explorer.BlockTransaction) {
	for _, addressTx := range CreateAddressTransactions(txs) {
		i.elastic.AddIndexRequest(
			elastic_cache.AddressTransactionIndex.Get(),
			fmt.Sprintf("%s-%s-%t", addressTx.Hash, addressTx.Txid, addressTx.Cold),
			addressTx,
		)
	}
}
