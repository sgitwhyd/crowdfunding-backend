package config

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func NewRedisClient() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: 	 "localhost:6379",
		Password: "",
		DB: 	 0,
	});

	_, err := rdb.Ping(ctx).Result()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to redis")

	return rdb, nil

}