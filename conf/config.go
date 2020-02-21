package conf

import (
	"fmt"
	"log"
	"os"

	gcfg "gopkg.in/gcfg.v1"
)

var AppConfig Config

type Config struct {
	Server ServerConfig
	DB     DBConfig
	Redis  RedisConfig
	API    APIConfig
}

type ServerConfig struct {
	Port string
}

type DBConfig struct {
	Driver string
	Path   string
}

type RedisConfig struct {
	Path string
}

type APIConfig struct {
	SampleAPI string
}

func LoadConfig() {
	var err error
	var path string

	if os.Getenv("PVENV") == "development" || os.Getenv("PVENV") == "" {
		path = "./files/etc/go-default/main.development.ini"
	} else {
		path = fmt.Sprintf("/etc/go-default/main.%v.ini", os.Getenv("PVENV"))
	}

	if path != "" {
		log.Println("[LoadConfig] Path", path)
		err = gcfg.ReadFileInto(&AppConfig, path)
	}

	if err != nil {
		log.Fatalln("[LoadConfig]", err)
	}
}
