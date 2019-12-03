package payment_request

import (
	"context"
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/elastic_cache"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	log "github.com/sirupsen/logrus"
)

type Indexer struct {
	navcoin *navcoind.Navcoind
	elastic *elastic_cache.Index
}

func NewIndexer(navcoin *navcoind.Navcoind, elastic *elastic_cache.Index) *Indexer {
	return &Indexer{navcoin, elastic}
}

func (i *Indexer) Index(txs []*explorer.BlockTransaction) {
	for _, tx := range txs {
		if tx.Version != 5 {
			continue
		}

		if navP, err := i.navcoin.GetPaymentRequest(tx.Hash); err == nil {
			paymentRequest := CreatePaymentRequest(navP, tx.Height)
			log.Infof("Index PaymentRequest: %s", paymentRequest.Hash)

			index := elastic_cache.PaymentRequestIndex.Get()
			resp, err := i.elastic.Client.Index().Index(index).BodyJson(paymentRequest).Do(context.Background())
			if err != nil {
				log.WithError(err).Fatal("Failed to save new payment request")
			}

			paymentRequest.MetaData = explorer.NewMetaData(resp.Id, resp.Index)
			PaymentRequests = append(PaymentRequests, paymentRequest)
		}
	}
}

func (i *Indexer) ApplyVote(vote explorer.Vote, blockCycle explorer.BlockCycle) {
	paymentRequest := getPaymentRequestByHash(vote.Hash)
	if paymentRequest == nil {
		log.Fatalf("Payment Request not found: %s", vote.Hash)
		return
	}

	if vote.Vote == 1 {
		paymentRequest.GetCycle(blockCycle.Cycle).VotesYes++
	}
	if vote.Vote == -1 {
		paymentRequest.GetCycle(blockCycle.Cycle).VotesNo++
	}

	i.elastic.AddUpdateRequest(elastic_cache.PaymentRequestIndex.Get(), paymentRequest.Hash, paymentRequest, paymentRequest.MetaData.Id)
}

func (i *Indexer) updateState(blockCycle explorer.BlockCycle) {

}
