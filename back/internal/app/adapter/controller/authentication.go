package controller

import (
	"net/http"

	"gpe_project/internal/app/adapter/middleware"
	"gpe_project/internal/app/adapter/request"
	"gpe_project/internal/app/adapter/response"
	"gpe_project/internal/app/application/usecase"

	"github.com/gin-gonic/gin"
)

type AuthenticationController struct {
	usecase.Usecase
}

func NewAuthenticationController(e *gin.Engine, usecase usecase.Usecase) *AuthenticationController {
	controller := &AuthenticationController{usecase}
	group := e.Group("/auth")
	group.POST("/register", controller.Register)
	group.POST("/login", controller.Login)
	group.GET("/refresh", middleware.RefreshTokenAuthentication(), controller.Refresh)

	return controller
}

// Register controller will register the user by input the mail and password
func (auth *AuthenticationController) Register(c *gin.Context) {
	userInput, err := request.AuthRequest{Context: c}.GetValidatedRegisterPayload()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	_, err = auth.UserUsecase.CreateNewUser(*userInput)
	if err != nil {
		c.JSON(ProcessError(err))
		return
	}

	c.JSON(SuccessMessage("user registered"))
}

// Login controller will connect the user by input login password and generates access and refresh tokens.
func (auth *AuthenticationController) Login(c *gin.Context) {
	login, err := request.AuthRequest{Context: c}.GetValidatedLoginPayload()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	accessToken, refreshToken, err := auth.AuthUsecase.Login(login)
	if err != nil {
		c.JSON(UnauthorizedError(err))
		return
	}

	c.JSON(http.StatusOK, response.AuthenticationTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// Refresh controller will renew access and refresh tokens
func (auth *AuthenticationController) Refresh(c *gin.Context) {
	claims, err := request.CommonRequest{Context: c}.GetValidatedClaimsHeaders()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	accessToken, refreshToken, err := auth.AuthUsecase.Refresh(claims.UserID)
	if err != nil {
		c.JSON(ProcessError(err))
		return
	}

	c.JSON(http.StatusOK, response.AuthenticationTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
