package dao_indexer

import (
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/index"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	log "github.com/sirupsen/logrus"
)

type PaymentRequestIndexer struct {
	navcoind *navcoind.Navcoind
	elastic  *index.Index
}

func NewPaymentRequestIndexer(navcoind *navcoind.Navcoind, elastic *index.Index) *PaymentRequestIndexer {
	return &PaymentRequestIndexer{navcoind, elastic}
}

func (i *PaymentRequestIndexer) IndexPaymentRequests(txs *[]explorer.BlockTransaction) {
	for _, tx := range *txs {
		if tx.Version == 5 {
			i.indexPaymentRequest(tx)
		}
	}
}

func (i *PaymentRequestIndexer) indexPaymentRequest(tx explorer.BlockTransaction) {
	navPaymentRequest, err := i.navcoind.GetPaymentRequest(tx.Hash)
	if err != nil {
		log.WithError(err).Errorf("Payment Request not found in tx %s", tx.Hash)
		return
	}

	log.Info("Indexing proposal in tx ", tx.Hash)
	i.elastic.AddRequest(index.PaymentRequestIndex.Get(), tx.Hash, CreatePaymentRequest(navPaymentRequest, tx.Height))
}
