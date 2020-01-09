package payment_request

import (
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

func (s *Service) LoadVotingPaymentRequests(block *explorer.Block, blockCycle *explorer.BlockCycle) {
	log.Info("Load Voting Payment Requests")

	excludeOlderThan := block.Height - (uint64(blockCycle.Size * 2))
	if excludeOlderThan < 0 {
		excludeOlderThan = 0
	}

	paymentRequests, err := s.repo.GetPossibleVotingRequests(excludeOlderThan)
	if err != nil {
		log.WithError(err).Error("Failed to load pending proposals")
	}

	PaymentRequests = paymentRequests
}
