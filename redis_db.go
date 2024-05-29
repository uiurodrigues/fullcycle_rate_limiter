package database

import (
	"time"

	"github.com/go-redis/redis"
)

var client *redis.Client

type RedisDB struct {
	client *redis.Client
}

func NewRedisDB(options *redis.Options) *RedisDB {
	client = redis.NewClient(options)
	return &RedisDB{client}
}

func (db *RedisDB) Get(key string) (int64, bool, error) {
	counter, err := db.client.Get(key).Int64()
	if err == redis.Nil {
		return 0, false, nil
	}

	return counter, true, nil
}

func (db *RedisDB) Set(key string, value int64, expiration time.Duration) error {
	err := db.client.Set(key, 1, expiration).Err()
	return err
}

func (db *RedisDB) Incr(key string) (int64, error) {
	counter, err := db.client.Incr(key).Result()
	if err != nil {
		return 0, err
	}
	return counter, nil
}

func (db *RedisDB) Block(key string, expiration time.Duration) error {
	blockedKey := key + "_blocked"
	err := db.client.Set(blockedKey, 1, expiration).Err()
	return err
}

func (db *RedisDB) IsBlocked(key string) (bool, error) {
	blockedKey := key + "_blocked"
	_, err := db.client.Get(blockedKey).Result()
	if err == redis.Nil {
		return false, nil
	}
	return true, err
}
