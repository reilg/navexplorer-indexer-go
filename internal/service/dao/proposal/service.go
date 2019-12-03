package proposal

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

func (s *Service) LoadVotingProposals() {
	log.Info("Load Pending Proposals")
	proposals, err := s.repo.GetProposals("pending")
	if err != nil {
		log.WithError(err).Fatal("Failed to load pending proposals")
	}

	Proposals = proposals
}

func getProposalByHash(hash string) *explorer.Proposal {
	for idx, _ := range Proposals {
		if Proposals[idx].Hash == hash {
			return Proposals[idx]
		}
	}

	return nil
}
