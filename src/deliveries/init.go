package deliveries

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samaita/go-default/conn"
	"github.com/samaita/go-default/models"
	"github.com/samaita/go-default/src"
	repo "github.com/samaita/go-default/src/repository"
	"github.com/samaita/go-default/src/usecase"
	"github.com/samaita/go-default/utils"
)

var RP src.Repository
var UC src.Usecase
var AppHandler *Handler

type Handler struct {
	Usecase src.Usecase
}

func InitModule() {
	RP = repo.NewRepository(conn.DB.Default, conn.Redis.Default, conn.API)
	UC = usecase.NewUsecase(RP)
	AppHandler = NewHandler(UC)
}

func NewHandler(u src.Usecase) *Handler {
	return &Handler{Usecase: u}
}

// TODO: Move to middleware
func setContextTimeStart(c *gin.Context) {
	c.Set(utils.KeyTimeStart, time.Now())
}

// TODO: Move to middleware
func getInitialResponse() map[string]interface{} {
	d := make(map[string]interface{})
	d[models.FIELD_SUCCESS] = false
	return d
}

// TODO: Move to middleware
func getSuccessResponse(d map[string]interface{}) map[string]interface{} {
	d[models.FIELD_SUCCESS] = true
	return d
}
