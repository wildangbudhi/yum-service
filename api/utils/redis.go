package utils

import (
	"strconv"

	"github.com/go-redis/redis/v8"
)

func NewRedisConnection(host, port, password, db string) (*redis.Client, error) {

	dbNumber, err := strconv.Atoi(db)

	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       dbNumber,
	})

	return rdb, nil
}
