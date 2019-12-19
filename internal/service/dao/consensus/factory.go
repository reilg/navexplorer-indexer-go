package consensus

import (
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
)

func CreateConsensus(cfundStats *navcoind.CFundStats) *explorer.Consensus {
	return &explorer.Consensus{
		BlocksPerVotingCycle:                cfundStats.Consensus.BlocksPerVotingCycle,
		MinSumVotesPerVotingCycle:           cfundStats.Consensus.MinSumVotesPerVotingCycle,
		MaxCountVotingCycleProposals:        cfundStats.Consensus.MaxCountVotingCycleProposals,
		MaxCountVotingCyclePaymentRequests:  cfundStats.Consensus.MaxCountVotingCyclePaymentRequests,
		VotesAcceptProposalPercentage:       cfundStats.Consensus.VotesAcceptProposalPercentage,
		VotesRejectProposalPercentage:       cfundStats.Consensus.VotesRejectProposalPercentage,
		VotesAcceptPaymentRequestPercentage: cfundStats.Consensus.VotesAcceptPaymentRequestPercentage,
		VotesRejectPaymentRequestPercentage: cfundStats.Consensus.VotesRejectPaymentRequestPercentage,
		ProposalMinimalFee:                  cfundStats.Consensus.ProposalMinimalFee,
	}
}

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
