package index

import (
	"encoding/json"
	"github.com/NavExplorer/navcoind-go"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/app/config"
	"github.com/dgraph-io/badger"
	"github.com/olliephillips/sett"
	"log"
	"strconv"
)

type indexer struct {
	navcoin *navcoind.Navcoind
	sett    *sett.Sett
}

func New() *indexer {
	navcoin, _ := navcoind.New(
		config.Get().Navcoind.Host,
		config.Get().Navcoind.Port,
		config.Get().Navcoind.User,
		config.Get().Navcoind.Password,
		false,
	)

	opts := sett.DefaultOptions
	opts.Dir = config.Get().Storage.Path
	opts.ValueDir = config.Get().Storage.Path

	return &indexer{navcoin: navcoin, sett: sett.Open(opts)}
}

func (i *indexer) Close() error {
	return i.sett.Close()
}

func (i *indexer) IndexAllBlocks() error {
	lastBlock, err := i.getLastBlockIndexed()
	if err != nil {
		return err
	}

	err = i.IndexBlock(lastBlock + 1)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return err
	}

	return i.IndexAllBlocks()
}

func (i *indexer) IndexBlock(height uint64) error {
	log.Printf("INFO: Indexing Block %d", height)

	hash, err := i.navcoin.GetBlockHash(height)
	if err != nil {
		return err
	}

	block, err := i.navcoin.GetBlock(hash)
	if err != nil {
		return err
	}

	return i.addBlock(block)
}

func (i * indexer) getLastBlockIndexed() (uint64, error) {
	lastBlock, err := i.sett.Table("meta").Get("last_block_indexed")
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return 0, nil
		}
		return 0, err
	}

	return strconv.ParseUint(lastBlock, 10, 64)
}

func (i * indexer) addBlock(block navcoind.Block) error {
	blockJson, err := json.Marshal(block)
	if err != nil {
		return err
	}

	err = i.sett.Table("block").Set(string(block.Height), string(blockJson))
	if err != nil {
		return err
	}

	err = i.sett.Table("meta").Set( "last_block_indexed", strconv.FormatUint(block.Height, 10))
	if err != nil {
		return err
	}

	return nil
}
