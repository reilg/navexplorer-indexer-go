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

func (s *Service) LoadVotingPaymentRequests() {
	log.Info("Load Pending Payment Requests")
	paymentRequests, err := s.repo.GetPaymentRequests("pending")
	if err != nil {
		log.WithError(err).Fatal("Failed to load pending proposals")
	}

	PaymentRequests = paymentRequests
}

func getPaymentRequestByHash(hash string) *explorer.PaymentRequest {
	for idx, _ := range PaymentRequests {
		if PaymentRequests[idx].Hash == hash {
			return PaymentRequests[idx]
		}
	}

	return nil
}
