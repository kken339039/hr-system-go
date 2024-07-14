package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hr-system-go/app/plugins"
	"hr-system-go/app/plugins/env"
	"hr-system-go/app/plugins/logger"
	"reflect"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func init() {
	plugins.Registry = append(plugins.Registry, NewRedisStore)
}

var ErrDataIsNotPointer = errors.New("data is not pointer")

type RedisStore struct {
	env    *env.Env
	logger *logger.Logger
	rdb    *redis.Client
	ctx    context.Context
}

func NewRedisStore(env *env.Env, logger *logger.Logger) *RedisStore {
	ctx := context.Background()
	store := &RedisStore{
		env:    env,
		logger: logger,
		ctx:    ctx,
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

// func (s *RedisStore) Fetch(key string, valueType reflect.Type, expiration time.Duration, valueFunc func()) (interface{}, error) {
// 	existData, err := s.rdb.Get(key, valueType)
// 	if err == redis.Nil {
// 		value := valueFunc
// 		if err := s.Set(key, value, expiration); err != nil {
// 			return nil, err
// 		}
// 		return value, nil
// 	} else if err != nil {
// 		return nil, err
// 	} else {
// 		return existData, nil
// 	}
// }

func (s *RedisStore) Get(redisKey string, pointerData interface{}) error {
	if reflect.TypeOf(pointerData).Kind() != reflect.Ptr {
		return ErrDataIsNotPointer
	}

	data, err := s.rdb.Get(s.ctx, redisKey).Bytes()
	if err != nil {
		return err
	}

	//nolint:exhaustive,forcetypeassert
	switch reflect.TypeOf(pointerData).Elem().Kind() {
	case reflect.String:
		*pointerData.(*string) = string(data)
	case reflect.Int:
		intValue, err := strconv.Atoi(string(data))
		if err != nil {
			return err
		}
		*pointerData.(*int) = intValue
	case reflect.Int64:
		int64Value, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return err
		}
		*pointerData.(*int64) = int64Value
	case reflect.Float64:
		floatValue, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return err
		}
		*pointerData.(*float64) = floatValue
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(string(data))
		if err != nil {
			return err
		}
		*pointerData.(*bool) = boolValue
	case reflect.Slice:
		if reflect.TypeOf(pointerData).Elem().Elem().Kind() == reflect.Uint8 {
			*pointerData.(*[]byte) = data
		} else {
			if err := json.Unmarshal(data, pointerData); err != nil {
				return err
			}
		}
	default:
		if err := json.Unmarshal(data, pointerData); err != nil {
			return err
		}
	}

	return nil
}

func (s *RedisStore) Set(key string, value interface{}, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		s.logger.Error("Failed to Marshal value", zap.Error(err))
		return err
	}
	return s.rdb.Set(s.ctx, key, jsonValue, expiration).Err()
}

func (s *RedisStore) Delete(redisKey string) error {
	res, err := s.rdb.Del(s.ctx, redisKey).Result()
	if err != nil {
		s.logger.Error("Failed to Delete redis key", zap.Error(err))
		return err
	}
	s.logger.Info("keys deleted:", zap.Int64("numbers:", res))
	return nil
}

func (s *RedisStore) ClearAll() error {
	return s.rdb.FlushAll(s.ctx).Err()
}
