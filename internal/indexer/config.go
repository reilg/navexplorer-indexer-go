package indexer

import (
	"github.com/go-redis/redis"
	"log"
	"strconv"
)

func (i *Indexer) getLastBlock() (uint64, error) {
	lastBlock, err := i.Redis.Client.Get("last_block_indexed").Result()
	if err != nil {
		if err == redis.Nil {
			log.Println("INFO: Last block not found")
			return 0, nil
		} else {
			return 0, err
		}
	}

	height, err := strconv.ParseUint(lastBlock, 10, 64)
	if err != nil || height < 0 {
		height = 0
	}

	return height, nil
}

func (i *Indexer) setLastBlock(height uint64) error {
	if err := i.Redis.Client.Set("last_block_indexed", height, 0).Err(); err != nil {
		log.Println(err)
		return ErrLastBlockWriteIndexed
	}

	return nil
}
