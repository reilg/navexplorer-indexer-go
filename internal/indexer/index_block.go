package indexer

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/events"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/entity"
	"github.com/gookit/event"
	log "github.com/sirupsen/logrus"
)

var (
	ErrTransactionNotFound = errors.New("Transaction not found")
)

func (i *Indexer) IndexBlocks() error {
	log.Info("Indexing all blocks")
	if i.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	lastBlock, err := i.getLastBlock()
	if err != nil {
		log.WithError(err).Error("Get last block")
		return err
	}

	if err := i.indexBlocks(lastBlock + 1); err != nil {
		if err != ErrOrphanBlockFound {
			log.WithError(err).Error("Orphan block found")
			return err
		}
		if err := i.RewindBy(10); err != nil {
			log.WithError(err).Error("Rewind blocks")
			return err
		}
	}

	return i.IndexBlocks()
}

func (i *Indexer) indexBlocks(height uint64) error {
	if err := i.IndexBlock(height); err != nil {
		return err
	}
	return i.indexBlocks(height + 1)
}

func (i *Indexer) IndexBlock(height uint64) error {
	hash, err := i.Navcoin.GetBlockHash(height)
	if err != nil {
		log.WithFields(log.Fields{"hash": hash, "height": height}).
			WithError(err).
			Error("Failed to GetBlockHash")
		return err
	}

	navBlock, err := i.Navcoin.GetBlock(hash)
	if err != nil {
		log.WithFields(log.Fields{"hash": hash, "height": height}).
			WithError(err).
			Error("Failed to GetBlock")
		return err
	}
	block := CreateBlock(navBlock)

	orphan, err := i.isOrphanBlock(block)
	if orphan == true {
		return ErrOrphanBlockFound
	}

	var txs = make([]entity.BlockTransaction, 0)
	for _, txHash := range block.Tx {
		rawTx, err := i.Navcoin.GetRawTransaction(txHash, true)
		if err != nil {
			log.WithFields(log.Fields{"hash": hash, "txHash": txHash, "height": height}).
				WithError(err).
				Error("Failed to GetRawTransaction")
			return err
		}

		txs = append(txs, CreateBlockTransaction(rawTx.(navcoind.RawTransaction)))
	}

	if err := i.applyInputs(&txs); err != nil {
		log.WithError(err).Error("Failed to apply inputs")
		return err
	}

	for idx, _ := range txs {
		applyType(&txs[idx], &txs)
		applyFees(&txs[idx], &block)
		applyStaking(&txs[idx], &block)
		applySpend(&txs[idx], &block)
		applyCFundPayout(&txs[idx], &block)
	}

	return i.persist(&txs, &block)
}

func (i *Indexer) applyInputs(txs *[]entity.BlockTransaction) error {
	log.Debug("*** APPLYING INPUTS ***")
	if len(*txs) == 0 {
		return nil
	}

	for idx, _ := range *txs {
		if len((*txs)[idx].Vin) == 0 {
			continue
		}
		vin := (*txs)[idx].Vin
		for vdx, _ := range vin {
			if vin[vdx].Vout == nil || vin[vdx].Txid == nil {
				continue
			}

			rawTx, err := i.Navcoin.GetRawTransaction(*vin[vdx].Txid, true)
			if err != nil {
				log.WithFields(log.Fields{"hash": *vin[vdx].Txid}).
					WithError(ErrTransactionNotFound).
					Fatal("Failed to get previous transaction")
			}
			prevTx := CreateBlockTransaction(rawTx.(navcoind.RawTransaction))

			if len(prevTx.Vout) <= *vin[vdx].Vout {
				log.WithFields(log.Fields{"index": vdx, "tx": prevTx.Hash}).Fatal("Vout does not exist")
			}

			previousOutput := prevTx.Vout[*vin[vdx].Vout]
			vin[vdx].Value = previousOutput.Value
			vin[vdx].ValueSat = previousOutput.ValueSat
			vin[vdx].Address = previousOutput.ScriptPubKey.Addresses[0]

			log.WithFields(log.Fields{"hash": prevTx.Hash}).Debug("Previous Transaction")
		}
	}

	return nil
}

func applyType(tx *entity.BlockTransaction, txs *[]entity.BlockTransaction) {
	var coinbase *entity.BlockTransaction
	for _, tx := range *txs {
		if tx.IsCoinbase() {
			coinbase = &tx
		}
	}
	log.Debug("*** APPLYING TYPE ***")
	if tx.IsCoinbase() {
		tx.MetaData.Type = string(entity.TxCoinbase)
	} else if tx.Vout.GetAmount() > tx.Vin.GetAmount() {
		if coinbase != nil && coinbase.Vout.HasOutputOfType(entity.VoutPoolStaking) {
			tx.MetaData.Type = string(entity.TxPoolStaking)
		}
		if tx.Vout.HasOutputOfType(entity.VoutColdStaking) {
			tx.MetaData.Type = string(entity.TxColdStaking)
		} else {
			tx.MetaData.Type = string(entity.TxStaking)
		}
	} else {
		tx.MetaData.Type = string(entity.TxSpend)
	}

	log.WithFields(log.Fields{"type": tx.MetaData.Type}).Debug("Transaction type")
}

func applyFees(tx *entity.BlockTransaction, block *entity.Block) {
	log.Debug("*** APPLYING FEES ***")
	if tx.IsSpend() {
		tx.MetaData.Fees = uint64(tx.Vin.GetAmount() - tx.Vout.GetAmount())
		block.MetaData.Fees += tx.MetaData.Fees
	}
	log.WithFields(log.Fields{"fees": tx.MetaData.Fees}).Debug("Transaction fees")
}

func applyStaking(tx *entity.BlockTransaction, block *entity.Block) {
	log.Debug("*** APPLYING STAKING ***")
	if tx.IsSpend() {
		return
	}

	if tx.IsAnyStaking() {
		log.Debug("Transaction is staking")
		if tx.Height >= 2761920 {
			log.Debug("Fixed stake reward")
			tx.MetaData.Stake = 2 // hard coded to 2 as static rewards arrived after block 2761920
			block.MetaData.Stake += tx.MetaData.Stake
		} else {
			log.Debug("Variable stake reward")
			tx.MetaData.Stake = uint64(tx.Vout.GetAmount() - tx.Vin.GetAmount())
			block.MetaData.Stake += tx.MetaData.Stake
		}
	} else if tx.IsCoinbase() {
		log.Debug("Transaction is coinbase")
		for _, o := range tx.Vout {
			if o.ScriptPubKey.Type == string(entity.VoutPubkeyhash) {
				tx.MetaData.Stake = o.ValueSat
				block.MetaData.Stake = o.ValueSat
			}
		}
	}

	voutsWithAddresses := tx.Vout.FilterWithAddresses()
	if len(voutsWithAddresses) != 0 {
		block.MetaData.StakedBy = voutsWithAddresses[0].ScriptPubKey.Addresses[0]
	}

	log.WithFields(log.Fields{"hash": tx.Hash, "stake": tx.MetaData.Stake}).Debug("Stake reward")
	log.WithFields(log.Fields{"hash": tx.Hash, "stakedBy": block.MetaData.StakedBy}).Debug("Stake by")
}

func applySpend(tx *entity.BlockTransaction, block *entity.Block) {
	log.Debug("*** APPLYING SPEND ***")
	if tx.MetaData.Type == string(entity.TxSpend) {
		tx.MetaData.Spend = tx.Vout.GetAmount()
		block.MetaData.Spend += tx.MetaData.Spend
		log.WithFields(log.Fields{"hash": tx.Hash, "spend": tx.Vout.GetAmount()}).Debug("Transaction spend")
	} else {
		log.Debug("Transaction is not a spend")
	}
}

func applyCFundPayout(tx *entity.BlockTransaction, block *entity.Block) {
	log.Debug("*** APPLYING CFUND PAYOUT ***")
	if !tx.IsCoinbase() {
		log.Debug("Only applies to coinbase TX")
		return
	}
	for _, o := range tx.Vout {
		if o.ScriptPubKey.Type == string(entity.VoutPubkeyhash) && tx.Version == 3 {
			block.MetaData.CFundPayout += o.ValueSat
		}
	}
	log.WithFields(log.Fields{"hash": tx.Hash, "payout": block.MetaData.CFundPayout}).Debug("Transaction cfund payout")
}

func (i *Indexer) persist(txs *[]entity.BlockTransaction, block *entity.Block) error {
	if err := i.persistBlockTransactions(txs, block); err != nil {
		log.WithError(err).Error("Failed to persist block transactions")
		return err
	}

	if err := i.persistBlock(*block); err != nil {
		log.WithError(err).Error("Failed to persist block")
		return err
	}

	b, _ := json.Marshal(block)
	log.Debug("")
	log.Debug("BLOCK: %s", string(b))

	t, _ := json.Marshal(txs)
	log.Debug("")
	log.Debug("TX: %s", string(t))
	log.Debug("")
	log.Debug("")

	if err := i.setLastBlock(block.Height); err != nil {
		log.WithError(err).Error("Failed to set last block indexed")
		return err
	}

	i.Elastic.Flush([]string{
		BlockIndex.Get(i.Network),
		BlockTransactionIndex.Get(i.Network),
	}...)

	go event.MustFire(string(events.EventBlockIndexed), event.M{"hash": block.Hash})

	log.WithFields(log.Fields{"height": block.Height}).Info("Block Indexed")
	return nil
}

func (i *Indexer) persistBlockTransactions(txs *[]entity.BlockTransaction, block *entity.Block) error {
	for _, tx := range *txs {
		_, err := i.Elastic.Client.
			Index().
			Index(BlockTransactionIndex.Get(i.Network)).
			Id(tx.Hash).
			BodyJson(tx).
			Do(context.Background())
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Indexer) persistBlock(block entity.Block) error {
	_, err := i.Elastic.Client.
		Index().
		Index(BlockIndex.Get(i.Network)).
		Id(block.Hash).
		BodyJson(block).
		Do(context.Background())
	return err
}
