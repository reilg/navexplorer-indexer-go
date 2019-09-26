package indexer

import (
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/entity"
)

func CreateBlock(block navcoind.Block) entity.Block {
	return entity.Block{
		Hash:              block.Hash,
		Confirmations:     block.Confirmations,
		StrippedSize:      block.StrippedSize,
		Size:              block.Size,
		Weight:            block.Weight,
		Height:            block.Height,
		Version:           block.Version,
		VersionHex:        block.VersionHex,
		Merkleroot:        block.Merkleroot,
		Tx:                block.Tx,
		Time:              block.Time,
		MedianTime:        block.MedianTime,
		Nonce:             block.Nonce,
		Bits:              block.Bits,
		Difficulty:        block.Difficulty,
		Chainwork:         block.Chainwork,
		Previousblockhash: block.Previousblockhash,
		Nextblockhash:     block.Nextblockhash,
	}
}

func CreateBlockTransaction(tx navcoind.RawTransaction) entity.BlockTransaction {
	return entity.BlockTransaction{
		Hex:             tx.Hex,
		Txid:            tx.Txid,
		Hash:            tx.Hash,
		Size:            tx.Size,
		VSize:           tx.VSize,
		Version:         tx.Version,
		LockTime:        tx.LockTime,
		AnonDestination: tx.AnonDestination,
		Vin:             createVin(tx.Vin),
		Vout:            createVout(tx.Vout),
		BlockHash:       tx.BlockHash,
		Height:          tx.Height,
		Confirmations:   tx.Confirmations,
		Time:            tx.Time,
		BlockTime:       tx.BlockTime,
	}
}

func createVin(vins []navcoind.Vin) []entity.Vin {
	var inputs = make([]entity.Vin, 0)
	for _, i := range vins {
		input := entity.Vin{
			Coinbase:  i.Coinbase,
			ScriptSig: createScriptSig(i.ScriptSig),
			Value:     i.Value,
			ValueSat:  i.ValueSat,
			Address:   i.Address,
			Sequence:  i.Sequence,
		}
		if i.Txid != "" {
			input.Txid = &i.Txid
			input.Vout = &i.Vout
		}
		inputs = append(inputs, input)
	}

	return inputs
}

func createScriptSig(scriptSig navcoind.ScriptSig) *entity.ScriptSig {
	if scriptSig.Hex == "" && scriptSig.Asm == "" {
		return nil
	}

	return &entity.ScriptSig{
		Asm: scriptSig.Asm,
		Hex: scriptSig.Hex,
	}
}

func createVout(vouts []navcoind.Vout) []entity.Vout {
	var output = make([]entity.Vout, 0)
	for _, o := range vouts {
		output = append(output, entity.Vout{
			Value:    o.Value,
			ValueSat: o.ValueSat,
			N:        o.N,
			ScriptPubKey: entity.ScriptPubKey{
				Asm:       o.ScriptPubKey.Asm,
				Hex:       o.ScriptPubKey.Hex,
				ReqSigs:   o.ScriptPubKey.ReqSigs,
				Type:      o.ScriptPubKey.Type,
				Addresses: o.ScriptPubKey.Addresses,
			},
		})
	}

	return output
}
