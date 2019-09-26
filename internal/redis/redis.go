package redis

import (
	"fmt"
	"github.com/go-redis/redis"
)

type Redis struct {
	Client *redis.Client
}

func New(addr string, password string, db int) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	return &Redis{
		Client: client,
	}
}
