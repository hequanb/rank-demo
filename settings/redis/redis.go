package redis

import (
	"boframe/settings"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var Cli *redis.Client

func Init(config *settings.RedisConfig) (err error) {
	Cli = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
		PoolSize: config.PoolSize,
	})
	
	_, err = Cli.Ping(context.Background()).Result()
	return err
}

func Close() {
	Cli.Close()
}

func IsErrNil(err error) bool {
	return errors.Is(err, redis.Nil)
}