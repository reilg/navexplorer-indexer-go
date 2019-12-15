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

func (s *Service) LoadVotingProposals(block *explorer.Block, blockCycle *explorer.BlockCycle) {
	log.Info("Load Voting Proposals")

	excludeOlderThan := block.Height - (uint64(blockCycle.Size * 2))
	if excludeOlderThan < 0 {
		excludeOlderThan = 0
	}

	proposals, err := s.repo.GetPossibleVotingProposals(excludeOlderThan)
	if err != nil {
		log.WithError(err).Fatal("Failed to load pending proposals")
	}

	Proposals = proposals
}
