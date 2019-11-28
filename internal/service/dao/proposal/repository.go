package proposal

import "github.com/olivere/elastic/v7"

type Repository struct {
	Client *elastic.Client
}

func NewRepo(client *elastic.Client) *Repository {
	return &Repository{client}
}

func (r *Repository) GetVotingProposals() []*Proposals {

	return nil
}
