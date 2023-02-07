package controller

import (
	"net/http"

	"gpe_project/internal/app/adapter/middleware"
	"gpe_project/internal/app/adapter/request"
	"gpe_project/internal/app/application/usecase"

	"github.com/gin-gonic/gin"
)

type BarometerController struct {
	usecase.Usecase
}

func NewBarometerController(e *gin.Engine, usecase usecase.Usecase) *BarometerController {
	controller := &BarometerController{usecase}

	group := e.Group("/barometer")
	group.POST("", controller.Define)
	group.GET("", controller.Barometers)
	group.GET("/stats", middleware.AdminOnly(usecase.Repositories), controller.Stats)

	return controller
}

// Define the mood for user. Mood is a daily value, the user define a score for his today's mood.
func (m *BarometerController) Define(c *gin.Context) {
	req, err := request.BarometerRequest{Context: c}.GetValidatedBarometerPayload()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	claims, err := request.CommonRequest{Context: c}.GetValidatedClaimsHeaders()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	req.UserID = claims.UserID
	barometer, err := m.Repositories.BarometerRepository.DefineForDay(req)
	if err != nil {
		c.JSON(ProcessError(err))
		return
	}
	c.JSON(http.StatusCreated, barometer)
}

// Barometers returns all barometers saved for the connected user for the current day.
func (m *BarometerController) Barometers(c *gin.Context) {
	claims, err := request.CommonRequest{Context: c}.GetValidatedClaimsHeaders()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	barometers, err := m.Repositories.BarometerRepository.RetrieveForDayUser(claims.UserID)
	if err != nil {
		c.JSON(ProcessError(err))
		return
	}
	c.JSON(http.StatusCreated, barometers)
}

// Stats returns all barometers stats.
func (m *BarometerController) Stats(c *gin.Context) {
	filter, err := request.BarometerRequest{Context: c}.GetValidatedFilters()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	stats, err := m.Repositories.BarometerRepository.GetStats(filter.ScopeCategories, filter.TimeScope)
	if err != nil {
		c.JSON(ProcessError(err))
		return
	}
	c.JSON(http.StatusOK, stats)
}
