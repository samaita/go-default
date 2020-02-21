package repository

import (
	"github.com/samaita/go-default/conn"
	src "github.com/samaita/go-default/src"
	"github.com/samaita/go-default/utils/redis/redigo"
	"github.com/samaita/go-default/utils/sql/sqlx"
)

type Repository struct {
	DBConnection    sqlx.DB
	RedisConnection redigo.Redis
	APIConnection   conn.APICollection
}

func NewRepository(conn sqlx.DB, redis redigo.Redis, api conn.APICollection) src.Repository {
	return &Repository{
		DBConnection:    conn,
		RedisConnection: redis,
		APIConnection:   api,
	}
}
