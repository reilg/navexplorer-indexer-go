package dao

import (
	"github.com/NavExplorer/navexplorer-indexer-go/internal/elastic_cache"
	log "github.com/sirupsen/logrus"
)

type Rewinder struct {
	elastic *elastic_cache.Index
}

func NewRewinder(elastic *elastic_cache.Index) *Rewinder {
	return &Rewinder{elastic}
}

func (r *Rewinder) Rewind(height uint64) error {
	log.Infof("Rewinding DAO index to height: %d", height)

	return r.elastic.DeleteHeightGT(height,
		elastic_cache.ProposalIndex.Get(),
		elastic_cache.PaymentRequestIndex.Get(),
		elastic_cache.ProposalVoteIndex.Get(),
		elastic_cache.PaymentRequestVoteIndex.Get(),
	)
}
