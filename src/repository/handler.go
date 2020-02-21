package repository

import (
	"github.com/gin-gonic/gin"
)

func (m *Repository) HealthCheck(c *gin.Context) (map[string]interface{}, error) {
	var err error
	var data map[string]interface{}
	return data, err
}
