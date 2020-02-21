package conn

import (
	"log"

	"github.com/samaita/go-default/conf"
	"github.com/samaita/go-default/utils/sql/sqlx"
)

var DB DBCollection

type DBCollection struct {
	Default sqlx.DB
}

func InitDB() {
	db, err := sqlx.New(conf.AppConfig.DB.Driver, conf.AppConfig.DB.Path)
	if err != nil {
		log.Fatalf("[InitDB] Fatal: %v", err)
	}
	DB.Default = db
	log.Println("[InitDB] DB Running....")
}
