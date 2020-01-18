package daemon

import (
	"github.com/NavExplorer/navexplorer-indexer-go/generated/dic"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/config"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/indexer"
	"github.com/sarulabs/dingo/v3"
	log "github.com/sirupsen/logrus"
)

var container *dic.Container

func Execute() {
	config.Init()
	container, _ = dic.NewContainer(dingo.App)

	container.GetElastic().InstallMappings()
	container.GetSoftforkService().LoadSoftForks()

	indexer.LastBlockIndexed = getHeight()
	//location := getHeight()
	//target := uint64(2772750)
	//
	//for {
	//	if target < location {
	//		location = location - 400
	//		if err := container.GetRewinder().RewindToHeight(location); err != nil {
	//			log.WithError(err).Fatal("Failed to rewind index")
	//		}
	//	} else {
	//		break
	//	}
	//}

	if err := container.GetRewinder().RewindToHeight(indexer.LastBlockIndexed); err != nil {
		log.WithError(err).Fatal("Failed to rewind index")
	}

	if indexer.LastBlockIndexed != 0 {
		if block, err := container.GetBlockRepo().GetBlockByHeight(indexer.LastBlockIndexed); err != nil {
			log.WithError(err).Fatal("Failed to get block at height: ", indexer.LastBlockIndexed)
		} else {
			consensus, err := container.GetDaoConsensusRepo().GetConsensus()
			if err == nil {
				blockCycle := block.BlockCycle(consensus.BlocksPerVotingCycle, consensus.MinSumVotesPerVotingCycle)
				container.GetDaoProposalService().LoadVotingProposals(block, blockCycle)
				container.GetDaoPaymentRequestService().LoadVotingPaymentRequests(block, blockCycle)
			}
		}
	}

	// Bulk index the backlog
	container.GetIndexer().BulkIndex()

	// Subscribe to 0MQ
	container.GetSubscriber().Subscribe()
}

func getHeight() uint64 {
	if height, err := container.GetBlockRepo().GetHeight(); err != nil {
		log.WithError(err).Fatal("Failed to get block height")
	} else {
		if height == 2772753 {
			log.Fatal("DIE HERE")
		}
		if height >= uint64(config.Get().BulkIndexSize) {
			return height - uint64(config.Get().BulkIndexSize)
		}
	}

	return 0
}
