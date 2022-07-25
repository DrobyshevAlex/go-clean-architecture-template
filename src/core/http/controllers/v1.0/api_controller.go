package v1_0

import (
	"main/src/core/http/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiController struct {
	*controllers.Controller
}

func (controller *ApiController) Version(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"version": "1.0"})
}

func NewApiController() *ApiController {
	return &ApiController{}
}
