package proposal

import "github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"

var Proposals proposals

type proposals []*explorer.Proposal

func (p *proposals) Delete(hash string) {
	for i, _ := range Proposals {
		if Proposals[i].Hash == hash {
			Proposals[i] = Proposals[len(Proposals)-1] // Copy last element to index i.
			Proposals[len(Proposals)-1] = nil          // Erase last element (write zero value).
			Proposals = Proposals[:len(Proposals)-1]   // Truncate slice.
			break
		}
	}
}
