package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gpe_project/internal/app/adapter/middleware"
	"gpe_project/internal/app/adapter/request"
	"gpe_project/internal/app/application/usecase"
)

type UserController struct {
	usecase.Usecase
}

func NewUserController(e *gin.Engine, usecase usecase.Usecase) *UserController {
	controller := &UserController{usecase}

	e.GET("/me", middleware.AccessTokenAuthentication(), controller.Me)
	e.PUT("/me", middleware.AccessTokenAuthentication(), controller.UpdateUser)

	return controller
}

// Me controller will return the user connected information.
func (auth *UserController) Me(c *gin.Context) {
	claims, err := request.CommonRequest{Context: c}.GetValidatedClaimsHeaders()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	user, err := auth.UserUsecase.Me(claims.UserID)
	if err != nil {
		c.JSON(ProcessError(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser controller will update basic user information.
func (auth *UserController) UpdateUser(c *gin.Context) {
	claims, err := request.CommonRequest{Context: c}.GetValidatedClaimsHeaders()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	req, err := request.UserRequest{Context: c}.GetValidatedUpdateUserPayload()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	req.ID = claims.UserID
	user, err := auth.UserUsecase.UpdateUser(req)
	if err != nil {
		c.JSON(ProcessError(err))
		return
	}

	c.JSON(http.StatusOK, user)
}
