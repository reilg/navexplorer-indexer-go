package redis

import (
	"errors"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/config"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type Redis struct {
	client *redis.Client
}

var (
	ErrLastBlockIndexedRead  = errors.New("Failed to read last block_indexer indexed")
	ErrLastBlockIndexedWrite = errors.New("Failed to write last block_indexer indexed")
)

func New() *Redis {
	return &Redis{
		client: redis.NewClient(&redis.Options{
			Addr:     config.Get().Redis.Host,
			Password: config.Get().Redis.Password,
			DB:       config.Get().Redis.Db,
		}),
	}
}

func (r *Redis) Init() error {
	if !config.Get().Reindex {
		return r.RewindBy(10)
	}

	return r.SetLastBlock(0)
}

func (r *Redis) GetLastBlockIndexed() (uint64, error) {
	lastBlock, err := r.client.Get("last_block_indexed").Result()
	if err != nil {
		if err == redis.Nil {
			log.Info("INFO: Last block_indexer not found")
			return 0, nil
		} else {
			log.WithError(err).Error(ErrLastBlockIndexedRead.Error())
			return 0, ErrLastBlockIndexedRead
		}
	}

	height, err := strconv.ParseUint(lastBlock, 10, 64)
	if err != nil || height < 0 {
		height = 0
	}

	return height, nil
}

func (r *Redis) SetLastBlock(height uint64) error {
	if err := r.client.Set("last_block_indexed", height, 0).Err(); err != nil {
		log.WithError(err).Error(ErrLastBlockIndexedWrite.Error())
		return ErrLastBlockIndexedWrite
	}

	return nil
}

func (r *Redis) RewindBy(blocks uint64) error {
	height, err := r.GetLastBlockIndexed()
	if err != nil {
		return err
	}

	log.Infof("Rewinding last block indexed from %d by %d blocks", height, blocks)

	if height > 10 {
		height = height - blocks
	} else {
		height = 0
	}

	return r.SetLastBlock(height)
}
