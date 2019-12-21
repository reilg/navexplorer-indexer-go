package consensus

import (
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
)

func UpdateConsensus(cfundStats *navcoind.CFundStats, c *explorer.Consensus) {
	c.BlocksPerVotingCycle = cfundStats.Consensus.BlocksPerVotingCycle
	c.MinSumVotesPerVotingCycle = cfundStats.Consensus.MinSumVotesPerVotingCycle
	c.MaxCountVotingCycleProposals = cfundStats.Consensus.MaxCountVotingCycleProposals
	c.MaxCountVotingCyclePaymentRequests = cfundStats.Consensus.MaxCountVotingCyclePaymentRequests
	c.VotesAcceptProposalPercentage = cfundStats.Consensus.VotesAcceptProposalPercentage
	c.VotesRejectProposalPercentage = cfundStats.Consensus.VotesRejectProposalPercentage
	c.VotesAcceptPaymentRequestPercentage = cfundStats.Consensus.VotesAcceptPaymentRequestPercentage
	c.VotesRejectPaymentRequestPercentage = cfundStats.Consensus.VotesRejectPaymentRequestPercentage
	c.ProposalMinimalFee = cfundStats.Consensus.ProposalMinimalFee
}
