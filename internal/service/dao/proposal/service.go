package proposal

import (
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
)

type Proposals []*explorer.Proposal

type Service struct {
	repo *Repository
}

func New(repo *Repository) *Service {
	return &Service{repo}
}

func (s *Service) LoadVotingProposals() {

}
