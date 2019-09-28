package address_indexer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/index"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	log "github.com/sirupsen/logrus"
)

type Indexer struct {
	elastic *index.Index
}

func New(elastic *index.Index) *Indexer {
	return &Indexer{elastic: elastic}
}

func (i *Indexer) IndexAddressesForTransactions(txs *[]explorer.BlockTransaction) {
	if len(*txs) == 0 {
		return
	}
	for _, tx := range *txs {
		for _, address := range tx.GetAllAddresses() {
			if err := i.indexAddressForTx(address, tx); err != nil {
				log.WithError(err).Fatal("Failed to index address transactions")
			}
		}
	}

	log.WithFields(log.Fields{"height": (*txs)[0].Height}).Info("Address Txs Indexed")
}

func (i *Indexer) indexAddressForTx(address string, tx explorer.BlockTransaction) error {
	addressTransaction := explorer.AddressTransaction{
		Hash:   address,
		Txid:   tx.Hash,
		Height: tx.Height,
		Time:   tx.Time,
	}

	_, input := tx.Vin.GetAmountByAddress(address)
	_, output := tx.Vout.GetAmountByAddress(address)
	addressTransaction.Input = input
	addressTransaction.Output = output
	addressTransaction.Total = int64(addressTransaction.Output - addressTransaction.Input)

	bt, _ := json.Marshal(tx)
	if tx.IsPoolStaking() {
		log.WithFields(log.Fields{"tx": string(bt)}).Fatal("Could not handle pool staking")
	} else if tx.IsColdStaking() {
		log.WithFields(log.Fields{"tx": string(bt)}).Fatal("Could not handle cold staking")
	} else if tx.IsStaking() {
		if tx.Vin.HasAddress(address) {
			addressTransaction.Type = explorer.TransferStake
		} else {
			btJson, _ := json.Marshal(bt)
			log.WithFields(log.Fields{"tx": string(btJson)}).Info("BlockTX")
			log.WithFields(log.Fields{"tx": string(bt)}).Fatal("Delegate staking recipient")
		}
	} else if tx.IsCoinbase() {
		if tx.Version == 1 {
			// POW block_indexer
			addressTransaction.Type = explorer.TransferStake
		} else if tx.Version == 3 {
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

	_, err := i.elastic.Client.
		Index().
		Index(index.AddressTransactionIndex.Get()).
		Id(fmt.Sprintf("%s-%s", address, tx.Hash)).
		BodyJson(addressTransaction).
		Do(context.Background())

	i.elastic.Flush([]string{
		index.AddressTransactionIndex.Get(),
	}...)

	return err
}
