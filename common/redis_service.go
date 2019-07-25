package common

import (
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"time"
)

/* from https://github.com/go-redis/redis/blob/master/command.go
type baseCmd struct {
	_args []interface{}
	err   error

	_readTimeout *time.Duration
}

type StringCmd struct {
	baseCmd

	val string
}

type StatusCmd struct {
	baseCmd

	val string
}

	Get(key string) *StringCmd
	Set(key string, value interface{}, expiration time.Duration) *StatusCmd

*/

// RedisIntf Interface to Redis commands
// All redis command to be called using this interface
type RedisIntf interface {
	Get(key string) (string, error)
	Set(key string, value interface{}, expiration time.Duration) error
}

// RedisService - Redis Pointer to redis
type RedisService struct {
	RedisClient *redis.Client
}

// NewRedisService get connection to redis and create a RedisService struct
func NewRedisService(redisOpt *RedisOptions) (*RedisService, error) {

	redisClient := redis.NewClient(&redis.Options{
		PoolSize:    10, // default
		IdleTimeout: 30 * time.Second,
		Addr:        redisOpt.Addr,
		Password:    "", // no password set
		DB:          0,  // use default DB
	})

	redisService := RedisService{}
	redisService.RedisClient = redisClient

	return &redisService, nil
}

// Get - call the Get method on the RedisClient
func (redis *RedisService) Get(key string) (string, error) {

	resp, err := redis.RedisClient.Get(key).Result()
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 268,
		}).Error(err)
	}

	return resp, err
}

// Set - Call the Set method on the Redis client
func (redis *RedisService) Set(key string, value interface{}, expiration time.Duration) error {

	err := redis.RedisClient.Set(key, value, 0).Err()
	if err != nil {
		log.WithFields(log.Fields{
			"msgnum": 265,
		}).Error(err)
		return err
	}

	return nil
}
