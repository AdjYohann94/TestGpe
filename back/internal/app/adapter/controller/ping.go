package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingController struct {
}

func NewPingController(e *gin.Engine) *PingController {
	controller := &PingController{}

	e.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	return controller
}
