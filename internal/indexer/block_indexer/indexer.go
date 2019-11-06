package block_indexer

import (
	"errors"
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/config"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/index"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/redis"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
	log "github.com/sirupsen/logrus"
)

type Indexer struct {
	elastic *index.Index
	cache   *redis.Redis
	navcoin *navcoind.Navcoind
}

var (
	ErrOrphanBlockFound = errors.New("Orphan block_indexer found")
)

func New(elastic *index.Index, cache *redis.Redis, navcoin *navcoind.Navcoind) *Indexer {
	return &Indexer{elastic: elastic, cache: cache, navcoin: navcoin}
}

func (i *Indexer) IndexBlocks() error {
	log.Info("Indexing all blocks")
	if config.Get().Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	lastBlock, err := i.cache.GetLastBlockIndexed()
	if err != nil {
		return err
	}

	if err := i.indexBlocks(lastBlock + 1); err != nil {
		if err != ErrOrphanBlockFound {
			log.WithError(err).Error(err)
			return err
		}

		_, err := i.cache.RewindBy(10)
		if err != nil {
			log.WithError(err).Error("Rewind blocks")
			return err
		}
	}

	return nil
}

func (i *Indexer) indexBlocks(height uint64) error {
	if err := i.IndexBlock(height); err != nil {
		return err
	}
	return i.indexBlocks(height + 1)
}

func (i *Indexer) IndexBlock(height uint64) error {
	//log.Info("Index Block @ height ", height)

	hash, err := i.navcoin.GetBlockHash(height)
	if err != nil {
		log.WithFields(log.Fields{"hash": hash, "height": height}).
			WithError(err).
			Error("Failed to GetBlockHash")
		return err
	}

	navBlock, err := i.navcoin.GetBlock(hash)
	if err != nil {
		log.WithFields(log.Fields{"hash": hash, "height": height}).WithError(err).Error("Failed to GetBlock")
		return err
	}
	block := CreateBlock(navBlock)

	orphan, _ := i.isOrphanBlock(block)
	if orphan == true {
		return ErrOrphanBlockFound
	}

	var txs = make([]explorer.BlockTransaction, 0)
	for _, txHash := range block.Tx {
		rawTx, err := i.navcoin.GetRawTransaction(txHash, true)
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

func (i *Indexer) applyInputs(txs *[]explorer.BlockTransaction) error {
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

			rawTx, err := i.navcoin.GetRawTransaction(*vin[vdx].Txid, true)
			if err != nil {
				log.WithFields(log.Fields{"hash": *vin[vdx].Txid}).WithError(err).Fatal("Failed to get previous transaction")
			}
			prevTx := CreateBlockTransaction(rawTx.(navcoind.RawTransaction))

			if len(prevTx.Vout) <= *vin[vdx].Vout {
				log.WithFields(log.Fields{"index": vdx, "tx": prevTx.Hash}).Fatal("Vout does not exist")
			}

			previousOutput := prevTx.Vout[*vin[vdx].Vout]
			vin[vdx].Value = previousOutput.Value
			vin[vdx].ValueSat = previousOutput.ValueSat
			vin[vdx].Addresses = previousOutput.ScriptPubKey.Addresses
			vin[vdx].PreviousOutput.Type = previousOutput.ScriptPubKey.Type
			vin[vdx].PreviousOutput.Height = prevTx.Height

			log.WithFields(log.Fields{"hash": prevTx.Hash}).Debug("Previous Transaction")
		}
	}

	return nil
}

func applyType(tx *explorer.BlockTransaction, txs *[]explorer.BlockTransaction) {
	var coinbase *explorer.BlockTransaction
	for _, tx := range *txs {
		if tx.IsCoinbase() {
			coinbase = &tx
			break
		}
	}

	if tx.IsCoinbase() {
		tx.Type = explorer.TxCoinbase
	} else if tx.Vout.GetAmount() > tx.Vin.GetAmount() {
		if coinbase != nil && coinbase.Vout.HasOutputOfType(explorer.VoutPoolStaking) {
			tx.Type = explorer.TxPoolStaking
		}
		if tx.Vout.HasOutputOfType(explorer.VoutColdStaking) {
			tx.Type = explorer.TxColdStaking
		} else {
			tx.Type = explorer.TxStaking
		}
	} else {
		tx.Type = explorer.TxSpend
	}
}

func applyFees(tx *explorer.BlockTransaction, block *explorer.Block) {
	if tx.IsSpend() {
		tx.Fees = tx.Vin.GetAmount() - tx.Vout.GetAmount()
		block.Fees += tx.Fees
	}
}

func applyStaking(tx *explorer.BlockTransaction, block *explorer.Block) {
	if tx.IsSpend() {
		return
	}

	if tx.IsAnyStaking() {
		if tx.Height >= 2761920 {
			tx.Stake = 2 // hard coded to 2 as static rewards arrived after block_indexer 2761920
			block.Stake += tx.Stake
		} else {
			tx.Stake = uint64(tx.Vout.GetAmount() - tx.Vin.GetAmount())
			block.Stake += tx.Stake
		}
	} else if tx.IsCoinbase() {
		for _, o := range tx.Vout {
			if o.ScriptPubKey.Type == explorer.VoutPubkeyhash {
				tx.Stake = o.ValueSat
				block.Stake = o.ValueSat
			}
		}
	}

	voutsWithAddresses := tx.Vout.FilterWithAddresses()
	if len(voutsWithAddresses) != 0 {
		block.StakedBy = voutsWithAddresses[0].ScriptPubKey.Addresses[0]
	}
}

func applySpend(tx *explorer.BlockTransaction, block *explorer.Block) {
	if tx.Type == explorer.TxSpend {
		tx.Spend = tx.Vout.GetAmount()
		block.Spend += tx.Spend
	}
}

func applyCFundPayout(tx *explorer.BlockTransaction, block *explorer.Block) {
	if !tx.IsCoinbase() {
		return
	}
	for _, o := range tx.Vout {
		if o.ScriptPubKey.Type == explorer.VoutPubkeyhash && tx.Version == 3 {
			block.CFundPayout += o.ValueSat
		}
	}
}
