package dao_indexer

import (
	"fmt"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/index"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
)

type ProposalVoteIndexer struct {
	elastic *index.Index
}

func NewProposalVoteIndexer(elastic *index.Index) *ProposalVoteIndexer {
	return &ProposalVoteIndexer{elastic}
}

func (i *ProposalVoteIndexer) IndexVotes(block *explorer.Block, tx *explorer.BlockTransaction) {
	if !tx.IsCoinbase() {
		return
	}

	for _, vout := range tx.Vout {
		if !vout.IsProposalVote() {
			continue
		}

		vote := explorer.ProposalVote{
			Height:   tx.Height,
			Address:  block.StakedBy,
			Proposal: vout.ScriptPubKey.Hash,
			Vote:     vout.ScriptPubKey.Type == explorer.VoutProposalYesVote,
		}

		i.elastic.AddIndexRequest(index.ProposalVoteIndex.Get(), fmt.Sprintf("%d-%s", tx.Height, vout.ScriptPubKey.Hash), vote)
	}
}
