package src

import (
	"github.com/gin-gonic/gin"
)

type Usecase interface {
	/* Handler */
	HealthCheck(c *gin.Context) (map[string]interface{}, error)

	CheckConnectionDB(c *gin.Context) (map[string]interface{}, error)
	CheckConnectionRedis(c *gin.Context) (map[string]interface{}, error)
}
