package repository

import (
	"github.com/gin-gonic/gin"
)

func (m *Repository) CheckConnectionDB(c *gin.Context) (map[string]interface{}, error) {
	var err error
	var result string

	data := make(map[string]interface{})

	query := `SELECT 1`
	_, err = m.DBConnection.Queryx(query)
	if err != nil {
		result = err.Error()
	} else {
		result = "OK"
	}

	data["DB Default"] = result
	return data, err
}

func (m *Repository) CheckConnectionRedis(c *gin.Context) (map[string]interface{}, error) {
	var err error

	data := make(map[string]interface{})

	result, err := m.RedisConnection.PING()
	if err != nil {
		return data, err
	}

	data["Redis Default"] = result
	return data, err
}
