package indexer

import "log"

func (i *Indexer) RewindBy(blocks uint64) error {
	log.Printf("INFO: Rewinding index by %d blocks", blocks)
	height, err := i.getLastBlock()
	if err != nil {
		return err
	}

	log.Printf("INFO: Rewinding index from %d by %d blocks", height, blocks)

	if height > 10 {
		height = height - blocks
	} else {
		height = 0
	}

	return i.setLastBlock(height)
}
