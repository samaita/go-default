package conn

import (
	"log"

	"github.com/samaita/go-default/conf"
	"github.com/samaita/go-default/utils/redis/redigo"
)

var Redis RedisConnection
var defaultConfig *redigo.RedisConfig
var RedisPath RedisAddress

type RedisConnection struct {
	Default redigo.Redis
}

type RedisAddress struct {
	Path string
}

func InitRedis() {

	RedisPath = RedisAddress{
		Path: conf.AppConfig.Redis.Path,
	}

	Redis.Default = redigo.New(RedisPath.Path, defaultConfig)
	_, err := Redis.Default.PING()
	if err != nil {
		log.Fatalf("[InitRedis] Fatal: %v", err)
	}
	log.Println("[InitRedis] RedisMain Running....")
}
