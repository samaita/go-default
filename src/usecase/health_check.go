package usecase

import (
	"github.com/gin-gonic/gin"
)

func (m *Usecase) CheckConnectionDB(c *gin.Context) (map[string]interface{}, error) {
	return m.repos.CheckConnectionDB(c)
}

func (m *Usecase) CheckConnectionRedis(c *gin.Context) (map[string]interface{}, error) {
	return m.repos.CheckConnectionRedis(c)
}
