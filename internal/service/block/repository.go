package block

import (
	"context"
	"encoding/json"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/elastic_cache"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	"github.com/olivere/elastic/v7"
)

type Repository struct {
	Client *elastic.Client
}

func NewRepo(Client *elastic.Client) *Repository {
	return &Repository{Client}
}

func (r *Repository) GetHeight() (uint64, error) {
	results, err := r.Client.
		Search(elastic_cache.AddressIndex.Get()).
		Sort("height", false).
		Size(1).
		Do(context.Background())
	if err != nil || results == nil {
		return 0, err
	}

	var block *explorer.Block
	if len(results.Hits.Hits) == 0 {
		return 0, nil
	}

	if err = json.Unmarshal(results.Hits.Hits[0].Source, &block); err != nil {
		return 0, err
	}
	return block.Height, nil
}
