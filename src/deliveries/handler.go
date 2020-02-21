package deliveries

import (
	"github.com/gin-gonic/gin"
	"github.com/samaita/go-default/utils"
)

func HealthCheck(c *gin.Context) {
	setContextTimeStart(c)
	AppHandler.HealthCheck(c)
}

func (m *Handler) HealthCheck(c *gin.Context) {
	var err error

	response := getInitialResponse()
	response, err = m.Usecase.HealthCheck(c)
	if err != nil {
		utils.HTTPInternalServerError(c, err.Error(), response)
		return
	}

	utils.HTTPSuccessResponse(c, response)
	return
}
