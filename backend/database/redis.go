package database

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/sklyar-vlad/selfDev/internal/config"
)

func NewRedis(ctx context.Context, config config.ConfigDatabase) (*redis.Client, error) {
	opt, err := redis.ParseURL("redis://<user>:<pass>@localhost:6379/<db>")
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	return client, nil
}
