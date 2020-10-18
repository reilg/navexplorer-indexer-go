package proposal

import (
	"github.com/NavExplorer/navexplorer-indexer-go/v2/pkg/explorer"
)

var Proposals proposals

type proposals []*explorer.Proposal

func (p *proposals) GetByHash(hash string) *explorer.Proposal {
	for _, p := range Proposals {
		if p.Hash == hash {
			return p
		}
	}

	return nil
}

func (p *proposals) Delete(hash string) {
	for i := range Proposals {
		if Proposals[i].Hash == hash {
			Proposals[i] = Proposals[len(Proposals)-1]                                     // Copy last element to index i.
			Proposals[len(Proposals)-1] = nil                                              // Erase last element (write zero value).
			Proposals = append([]*explorer.Proposal(nil), Proposals[:len(Proposals)-1]...) // Truncate slice.
			break
		}
	}
}
