package redis

import (
	"be-bwastartup/config"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Service interface{
	Save(key string, value interface{}, expire time.Duration) (string, error)
	Get(key string) (string, error)
}

type service struct {
	rdb *redis.Client
}

func NewService() *service {
	rdb, err := config.NewRedisClient()
	if err != nil {
		log.Fatal("Failed to connect to redis")
	}

	return &service{rdb: rdb}
}

var ctx = context.Background()

func (s *service) Save(key string, value interface{}, expire time.Duration) (string, error) {

	valueJSON, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	cacheValue, err := s.rdb.Set(ctx,key, valueJSON, expire).Result()
	if err != nil{
		return "", err
	}

	return cacheValue, err
}

func (s *service) Get(key string)  (string, error) {
	val, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	return  val, err
}
