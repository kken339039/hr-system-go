package redis

import (
	"context"
	"fmt"
	"hr-system-go/app/plugins"
	"hr-system-go/app/plugins/env"
	"hr-system-go/app/plugins/logger"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func init() {
	plugins.Registry = append(plugins.Registry, NewRedisStore)
}

type RedisStore struct {
	env    *env.Env
	logger *logger.Logger
	rdb    *redis.Client
}

func NewRedisStore(env *env.Env, logger *logger.Logger) *RedisStore {
	store := &RedisStore{
		env:    env,
		logger: logger,
	}
	return store
}

func (s *RedisStore) Connect(host string, port string, db int) {
	addr := fmt.Sprintf("%s:%s", host, port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       db,
	})

	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		s.logger.Error("Failed to connect Redis", zap.Error(err))
		panic(err)
	}

	s.rdb = rdb
	s.logger.Info(fmt.Sprintf("Connected to connect Redis %s", pong))
}
