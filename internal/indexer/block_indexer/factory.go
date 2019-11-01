package block_indexer

import (
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/pkg/explorer"
)

func CreateBlock(block navcoind.Block) explorer.Block {
	return explorer.Block{
		RawBlock: explorer.RawBlock{
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
		},
	}
}

func CreateBlockTransaction(tx navcoind.RawTransaction) explorer.BlockTransaction {
	return explorer.BlockTransaction{
		RawBlockTransaction: explorer.RawBlockTransaction{
			Hex:             tx.Hex,
			Txid:            tx.Txid,
			Hash:            tx.Hash,
			Size:            tx.Size,
			VSize:           tx.VSize,
			Version:         tx.Version,
			LockTime:        tx.LockTime,
			AnonDestination: tx.AnonDestination,
			BlockHash:       tx.BlockHash,
			Height:          tx.Height,
			Confirmations:   tx.Confirmations,
			Time:            tx.Time,
			BlockTime:       tx.BlockTime,
		},
		Vin:  createVin(tx.Vin),
		Vout: createVout(tx.Vout),
	}
}

func createVin(vins []navcoind.Vin) []explorer.Vin {
	var inputs = make([]explorer.Vin, 0)
	for idx, _ := range vins {
		input := explorer.Vin{
			RawVin: explorer.RawVin{
				Coinbase: vins[idx].Coinbase,
				Sequence: vins[idx].Sequence,
			},
		}
		if vins[idx].Txid != "" {
			input.Txid = &vins[idx].Txid
			input.Vout = &vins[idx].Vout
		}
		inputs = append(inputs, input)
	}

	return inputs
}

func createVout(vouts []navcoind.Vout) []explorer.Vout {
	var output = make([]explorer.Vout, 0)
	for _, o := range vouts {
		output = append(output, explorer.Vout{
			RawVout: explorer.RawVout{
				Value:    o.Value,
				ValueSat: o.ValueSat,
				N:        o.N,
				ScriptPubKey: explorer.ScriptPubKey{
					Asm:       o.ScriptPubKey.Asm,
					Hex:       o.ScriptPubKey.Hex,
					ReqSigs:   o.ScriptPubKey.ReqSigs,
					Type:      explorer.VoutTypes[o.ScriptPubKey.Type],
					Addresses: o.ScriptPubKey.Addresses,
				},
			},
		})
	}

	return output
}
