package usecase

import (
	"github.com/gin-gonic/gin"
)

func (m *Usecase) HealthCheck(c *gin.Context) (map[string]interface{}, error) {
	var err error
	var data, db, redis map[string]interface{}

	data = make(map[string]interface{})

	db, err = m.CheckConnectionDB(c)
	if err != nil {
		return data, err
	}

	redis, err = m.CheckConnectionRedis(c)
	if err != nil {
		return data, err
	}

	data["DB"] = db
	data["Redis"] = redis
	return data, err
}
