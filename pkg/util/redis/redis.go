package redis

import (
	"context"
	"github.com/lactobasilusprotectus/go-template/pkg/common/config"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

var (
	ctx = context.Background()
)

type Client struct {
	redis *redis.Client
}

func NewRedisClient(config config.RedisConfig) *Client {
	rdb := redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     config.Host,
		Username: config.Username,
		Password: config.Password,
		DB:       config.DB,
	})

	return &Client{
		redis: rdb,
	}
}

func (c *Client) Get(key string) (reply interface{}, err error) {
	result, err := c.redis.Get(ctx, key).Result()

	if err != nil {
		return nil, err
	}

	//close connection
	defer func(redis *redis.Client) {
		err = redis.Close()

		if err != nil {
			log.Fatalf("error closing redis connection: %v", err)
		}
	}(c.redis)

	return result, nil
}

func (c *Client) Set(key string, value interface{}, expireSeconds int) (err error) {
	if expireSeconds <= 0 {
		err = c.redis.Set(ctx, key, value, 0).Err()
	} else {
		err = c.redis.Set(ctx, key, value, time.Duration(expireSeconds)*time.Second).Err()
	}

	if err != nil {
		return err
	}

	//close connection
	defer func(redis *redis.Client) {
		err = redis.Close()

		if err != nil {
			log.Fatalf("error closing redis connection: %v", err)
		}
	}(c.redis)

	return nil
}
