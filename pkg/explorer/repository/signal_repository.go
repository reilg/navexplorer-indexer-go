package repository

import (
	"context"
	"encoding/json"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/index"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	"github.com/olivere/elastic/v7"
)

type SignalRepository struct {
	Client *elastic.Client
}

func New(client *elastic.Client) SignalRepository {
	return SignalRepository{client}
}

func (r *SignalRepository) GetSignals(start uint64, end uint64) *[]explorer.Signal {
	signals := make([]explorer.Signal, 0)

	results, err := r.Client.Search(index.SignalIndex.Get()).
		Sort("height", true).
		Query(elastic.NewRangeQuery("height").Gte(start).Lte(end)).
		Size(int(end-start) + 1).
		Do(context.Background())

	if err == nil && results != nil && results.Hits != nil {
		for _, hit := range results.Hits.Hits {
			var signal explorer.Signal
			err := json.Unmarshal(hit.Source, &signal)
			if err == nil {
				signals = append(signals, signal)
			}
		}
	}

	return &signals
}
