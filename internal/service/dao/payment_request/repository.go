package payment_request

import (
	"context"
	"encoding/json"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/elastic_cache"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
)

type Repository struct {
	Client *elastic.Client
}

func NewRepo(client *elastic.Client) *Repository {
	return &Repository{client}
}

func (r *Repository) GetPaymentRequests(status string) ([]*explorer.PaymentRequest, error) {
	var paymentRequests []*explorer.PaymentRequest

	results, err := r.Client.Search(elastic_cache.PaymentRequestIndex.Get()).
		Query(elastic.NewTermsQuery("status", status)).
		Size(9999).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	if results == nil {
		return nil, elastic_cache.ErrResultsNotFound
	}

	for _, hit := range results.Hits.Hits {
		var paymentRequest *explorer.PaymentRequest
		if err := json.Unmarshal(hit.Source, &paymentRequest); err != nil {
			log.WithError(err).Fatal("Failed to unmarshall proposal")
		}
		paymentRequest.MetaData = explorer.NewMetaData(hit.Id, hit.Index)
		paymentRequests = append(paymentRequests, paymentRequest)
	}

	return paymentRequests, nil
}
