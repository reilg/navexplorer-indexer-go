package repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/index"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	"github.com/olivere/elastic/v7"
)

type SoftForkRepository struct {
	Client *elastic.Client
}

func NewSoftForkRepo(client *elastic.Client) *SoftForkRepository {
	return &SoftForkRepository{client}
}

func (r *SoftForkRepository) GetSoftFork(name string) (*explorer.SoftFork, error) {
	var softfork *explorer.SoftFork

	results, err := r.Client.Search(index.SoftForkIndex.Get()).
		Query(elastic.NewTermQuery("name", name)).
		Size(1).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	if results == nil || len(results.Hits.Hits) != 1 {
		return nil, errors.New("Invalid result found")
	}

	err = json.Unmarshal(results.Hits.Hits[0].Source, &softfork)
	if err != nil {
		return nil, err
	}

	return softfork, nil
}
