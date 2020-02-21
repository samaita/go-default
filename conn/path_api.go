package conn

import (
	"log"

	"github.com/samaita/go-default/conf"
)

var API APICollection

type APICollection struct {
	Default string
}

func InitPath() {
	API.Default = conf.AppConfig.API.SampleAPI
	log.Println("[InitPath] Path Initiated....")
}
