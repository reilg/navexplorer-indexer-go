package redis

import (
	"errors"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/config"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type Redis struct {
	client      *redis.Client
	reindexSize uint
}

var (
	ErrLastBlockIndexedRead  = errors.New("Failed to read last block indexed")
	ErrLastBlockIndexedWrite = errors.New("Failed to write last block indexed")
)

func NewRedis(addr string, password string, db int, reindexSize uint) *Redis {
	client := redis.NewClient(&redis.Options{Addr: addr, Password: password, DB: db})
	return &Redis{client, reindexSize}
}

func (r *Redis) Start() (uint64, error) {
	if !config.Get().Reindex {
		return r.GetLastBlockIndexed()
	}

	err := r.SetLastBlock(0)
	log.Info("Redis started")

	return 0, err
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
