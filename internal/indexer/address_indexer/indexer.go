package address_indexer

import (
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
			i.indexAddressForTx(address, tx)
		}
	}
}

func (i *Indexer) indexAddressForTx(address string, tx explorer.BlockTransaction) {
	if tx.HasColdStakingInput(address) || tx.HasColdStakingOutput(address) {
		i.indexAddressForColdTx(address, tx)
	}

	addressTransaction := explorer.AddressTransaction{
		Hash:   address,
		Txid:   tx.Hash,
		Height: tx.Height,
		Time:   tx.Time,
		Cold:   false,
	}

	_, input := tx.Vin.GetAmountByAddress(address, addressTransaction.Cold)
	_, output := tx.Vout.GetAmountByAddress(address, addressTransaction.Cold)
	addressTransaction.Input = input
	addressTransaction.Output = output
	addressTransaction.Total = int64(addressTransaction.Output - addressTransaction.Input)

	bt, _ := json.Marshal(tx)
	if tx.IsStaking() {
		if tx.Vin.HasAddress(address) {
			addressTransaction.Type = explorer.TransferStake
		} else {
			addressTransaction.Type = explorer.TransferDelegateStake
		}
	} else if tx.IsCoinbase() {
		if tx.Version == 1 {
			// POW block_indexer
			addressTransaction.Type = explorer.TransferStake
		} else if tx.Version == 3 {
			addressTransaction.Type = explorer.TransferCommunityFundPayout
		} else {
			log.WithFields(log.Fields{"tx": string(bt)}).Fatal("Could not handle coinbase")
		}
	} else {
		if addressTransaction.Input > addressTransaction.Output {
			addressTransaction.Type = explorer.TransferSend
		} else {
			addressTransaction.Type = explorer.TransferReceive
		}
	}

	i.elastic.AddRequest(
		index.AddressTransactionIndex.Get(),
		fmt.Sprintf("%s-%s", address, tx.Hash),
		addressTransaction,
	)
}

func (i *Indexer) indexAddressForColdTx(address string, tx explorer.BlockTransaction) {
	addressTransaction := explorer.AddressTransaction{
		Hash:   address,
		Txid:   tx.Hash,
		Height: tx.Height,
		Time:   tx.Time,
		Cold:   true,
	}

	_, input := tx.Vin.GetAmountByAddress(address, addressTransaction.Cold)
	_, output := tx.Vout.GetAmountByAddress(address, addressTransaction.Cold)
	addressTransaction.Input = input
	addressTransaction.Output = output
	addressTransaction.Total = int64(addressTransaction.Output - addressTransaction.Input)

	bt, _ := json.Marshal(tx)
	if tx.IsSpend() {
		if addressTransaction.Input > addressTransaction.Output {
			addressTransaction.Type = explorer.TransferSend
		} else {
			addressTransaction.Type = explorer.TransferReceive
		}
	} else if tx.IsColdStaking() {
		if tx.Vin.HasAddress(address) {
			addressTransaction.Type = explorer.TransferStake
		} else {
			addressTransaction.Type = explorer.TransferDelegateStake
			btJson, _ := json.Marshal(bt)
			log.WithFields(log.Fields{"tx": string(btJson)}).Info("BlockTX")
			log.WithFields(log.Fields{"tx": string(bt)}).Fatal("Delegate staking recipient")
		}
	}

	i.elastic.AddRequest(
		index.AddressTransactionIndex.Get(),
		fmt.Sprintf("%s-%s-cold", address, tx.Hash),
		addressTransaction,
	)
}
