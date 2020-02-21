package redigo

import (
	"sync"
	"time"

	redigo "github.com/garyburd/redigo/redis"
)

const (
	ErrNilRedis = "redigo: nil returned"
)

// Configuration for redis
type RedisConfig struct {
	MaxIdle   int
	MaxActive int
	MaxBatch  int
	Timeout   int
	Wait      bool
	Network   string
}

// Redis struct
type redis struct {
	Pool  *redigo.Pool
	mutex sync.Mutex
}

type Redis interface {
	PING() (string, error)
	TTL(key string) (time.Duration, error)
	EXISTS(key string) (bool, error)
	EXPIRE(key string, duration time.Duration) error
	INCR(key string) (int64, error)
	INCRBY(key string, val int) error

	SET(key string, value interface{}) error
	SETEX(key string, dur time.Duration, value interface{}) error
	GET(key string) (string, error)
	MGET(key ...string) ([]string, error)
	DEL(key string) error

	HSET(key string, field string, value string) error
	HSETEX(key string, duration time.Duration, field string, value string) error
	HMSET(key string, data map[string]interface{}) error
	HGET(key string, field string) (string, error)
	HGETP(keys []string, field string) (map[string]interface{}, error)
	HMGET(key string, field ...string) (map[string]string, error)
	HMGETP(keys []string, field ...string) (map[string]interface{}, error)
	HGETALL(key string) (map[string]string, error)
	HDEL(key string, fields ...string) error
	HEXISTS(key string, field string) (bool, error)

	LPUSH(key string, value interface{}) error
	LLEN(key string) (int64, error)
	LPUSHEX(key string, duration time.Duration, value interface{}) error
	LPUSHEXP(value map[string][]interface{}, duration time.Duration) error
	LREM(key string, count interface{}, value interface{}) (int64, error)
	LRANGE(key string, start interface{}, stop interface{}) ([]string, error)
	LRANGEP(keys []string, fields ...string) (map[string][]int64, error)
	LTRIM(key string, start interface{}, end interface{}) error
	LPOP(key string) (string, error)

	ZADD(key string, fields ...Z) error
	ZSCORE(key string, value interface{}) (int, error)
	ZREM(key string, fields ...string) error
	ZRANGE(key string, start int, end int) ([]string, error)
	ZRANGEBYSCORE(key string, opt ZRangeByScore) ([]string, error)

	SADD(key string, val ...interface{}) error
	SREM(key string, val ...interface{}) error
	SMEMBERS(key string) ([]string, error)
	SCARD(key string) (int64, error)
	SISMEMBERS(key string, field interface{}) (bool, error)

	RPUSH(key string, value ...string) (int64, error)
}

// New redis connection
func New(address string, conf *RedisConfig) Redis {
	if conf == nil {
		conf = &RedisConfig{
			MaxIdle:   10,
			MaxActive: 30,
			Timeout:   240,
			Wait:      true,
			Network:   "tcp",
		}
	}

	result := &redis{
		Pool: &redigo.Pool{
			MaxIdle:     conf.MaxIdle,
			MaxActive:   conf.MaxActive,
			IdleTimeout: time.Duration(conf.Timeout) * time.Second,
			Dial: func() (redigo.Conn, error) {
				return redigo.Dial(conf.Network, address)
			},
			TestOnBorrow: func(c redigo.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
			Wait: conf.Wait,
		},
	}

	return Redis(result)
}
