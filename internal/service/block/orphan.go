package block

import (
	"errors"
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	"github.com/getsentry/raven-go"
	log "github.com/sirupsen/logrus"
)

type OrphanService struct {
	navcoin *navcoind.Navcoind
}

func NewOrphanService(navcoin *navcoind.Navcoind) *OrphanService {
	return &OrphanService{navcoin}
}

var (
	ErrOrphanBlockFound = errors.New("Orphan block_indexer found")
)

func (o *OrphanService) IsOrphanBlock(block *explorer.Block) (bool, error) {
	if block.Height == 1 {
		return false, nil
	}

	previousBlockHash, err := o.navcoin.GetBlockHash(block.Height - 1)
	if err != nil {
		raven.CaptureError(err, nil)
		log.WithError(err).Fatal("OrphanService: Failed to get previous block hash")
	}

	orphan := previousBlockHash != block.Previousblockhash
	if orphan == true {
		raven.CaptureError(err, nil)
		log.WithFields(log.Fields{"height": block.Height, "hash": block.Hash}).Info("Orphan block found")
	}

	return orphan, nil
}
