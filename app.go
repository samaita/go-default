package main

import (
	"github.com/gin-gonic/gin"

	"github.com/samaita/go-default/conf"
	"github.com/samaita/go-default/conn"

	app "github.com/samaita/go-default/src/deliveries"
)

func init() {
	conf.LoadConfig()
	conn.InitDB()
	conn.InitRedis()

	app.InitModule()
}

func main() {
	// Release Gin
	// gin.SetMode(gin.ReleaseMode)
	// Middleware Setup
	router := gin.Default()

	// Assign API to handle deliveries
	health := router.Group("health")
	{
		health.GET("/check", app.HealthCheck)
	}

	// Running App
	router.Run(conf.AppConfig.Server.Port)
}
