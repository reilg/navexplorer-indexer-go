package indexer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/elastic"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	log "github.com/sirupsen/logrus"
)

type AddressIndexer struct {
	elastic *elastic.Elastic
	network string
}

func NewAddressIndexer(network string, es *elastic.Elastic) *AddressIndexer {
	return &AddressIndexer{network: network, elastic: es}
}

func (addressIndexer *AddressIndexer) IndexAddressesForBlock(hash string) {
	var block explorer.Block
	if err := addressIndexer.elastic.GetById(BlockIndex.Get(addressIndexer.network), hash, &block); err != nil {
		log.WithError(err).Fatal("Failed to get block")
	}

	for _, txHash := range block.Tx {
		var blockTransaction explorer.BlockTransaction
		if err := addressIndexer.elastic.GetById(BlockTransactionIndex.Get(addressIndexer.network), txHash, &blockTransaction); err != nil {
			log.WithError(err).Fatal("Failed to get block")
		}

		//log.WithFields(log.Fields{"height": block.Height}).Debug("Indexing Address Transactions")
		for _, hash := range blockTransaction.GetAllAddresses() {
			if err := addressIndexer.indexAddressForTx(hash, blockTransaction); err != nil {
				log.WithError(err).Fatal("Failed to index address transactions")
			}
		}
	}
}

func (addressIndexer *AddressIndexer) indexAddressForTx(hash string, blockTransaction explorer.BlockTransaction) error {
	addressTransaction := explorer.AddressTransaction{
		Hash:   hash,
		Txid:   blockTransaction.Hash,
		Height: blockTransaction.Height,
		Time:   blockTransaction.Time,
	}

	_, input := blockTransaction.Vin.GetAmountByAddress(hash)
	_, output := blockTransaction.Vout.GetAmountByAddress(hash)
	addressTransaction.Input = input
	addressTransaction.Output = output
	addressTransaction.Total = int64(addressTransaction.Output - addressTransaction.Input)

	bt, _ := json.Marshal(blockTransaction)
	if blockTransaction.IsPoolStaking() {
		log.WithFields(log.Fields{"tx": string(bt)}).Fatal("Could not handle pool staking")
	} else if blockTransaction.IsColdStaking() {
		log.WithFields(log.Fields{"tx": string(bt)}).Fatal("Could not handle cold staking")
	} else if blockTransaction.IsStaking() {
		if blockTransaction.Vin.HasAddress(hash) {
			addressTransaction.Type = explorer.TransferStake
		} else {
			btJson, _ := json.Marshal(bt)
			//	log.Debug("")
			log.WithFields(log.Fields{"tx": string(btJson)}).Info("BlockTX")
			log.WithFields(log.Fields{"tx": string(bt)}).Fatal("Delegate staking recipient")
		}
	} else if blockTransaction.IsCoinbase() {
		if blockTransaction.Version == 1 {
			// POW block
			addressTransaction.Type = explorer.TransferStake
		} else if blockTransaction.Version == 3 {
			log.WithFields(log.Fields{"tx": string(bt)}).Fatal("Could not handle cfund payout")
		} else {
			log.WithFields(log.Fields{"tx": string(bt)}).Fatal("Could not handle coinbase")
		}
	} else {
		if addressTransaction.Input < addressTransaction.Output {
			addressTransaction.Type = explorer.TransferSend
		} else {
			addressTransaction.Type = explorer.TransferReceive
		}
	}

	//if addressTransaction.Height == 293 && addressTransaction.Transfer[0].Type != explorer.TransferStake {
	//	b, _ := json.Marshal(addressTransaction)
	//	log.Debug("")
	//	log.Fatal("Address TRANSACTION: %s", string(b))
	//}

	_, err := addressIndexer.elastic.Client.
		Index().
		Index(AddressTransactionIndex.Get(addressIndexer.network)).
		Id(fmt.Sprintf("%s-%s", hash, blockTransaction.Hash)).
		BodyJson(addressTransaction).
		Do(context.Background())

	addressIndexer.elastic.Flush([]string{
		AddressTransactionIndex.Get(addressIndexer.network),
	}...)

	return err
}
