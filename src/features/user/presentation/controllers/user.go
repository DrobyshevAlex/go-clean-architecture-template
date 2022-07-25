package controllers

import (
	"main/src/core/config"
	"main/src/core/http/controllers"
	"main/src/features/user/domain/usecases"
	"main/src/features/user/presentation/requests"
	"main/src/features/user/presentation/responses"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	*controllers.Controller

	userUsecase usecases.UserUsecase
}

func (c *UserController) Profile(ctx *gin.Context) {
	request := requests.GetProfileRequest{}
	if err := ctx.ShouldBindUri(&request); err != nil {
		c.SendError(ctx, err, 422)
		return
	}

	profileResponse, err := c.userUsecase.GetUserById(ctx, request.Id)
	if err != nil || profileResponse == nil {
		c.SendError(ctx, err, 500)
		return
	}

	response := responses.UserProfileResponse{
		FirstName: profileResponse.FirstName,
		LastName:  profileResponse.LastName,
	}

	ctx.JSON(200, response)

}

func NewUserController(
	userUsecase usecases.UserUsecase,
	config *config.Config,
) *UserController {
	return &UserController{
		userUsecase: userUsecase,
		Controller: &controllers.Controller{
			Config: config,
		},
	}
}
