package controller

import (
	"net/http"

	"gpe_project/internal/app/adapter/postgresql/scope"
	"gpe_project/internal/app/adapter/request"
	"gpe_project/internal/app/application/usecase"

	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	usecase.Usecase
}

func NewNotificationController(e *gin.Engine, usecase usecase.Usecase) *NotificationController {
	controller := &NotificationController{usecase}

	group := e.Group("/notification")
	group.GET("", controller.UnreadNotifications)
	group.GET("/archive", controller.AllNotifications)
	group.POST("/read", controller.ReadNotifications)

	return controller
}

// UnreadNotifications returns all unread notifications for the current user
func (m *NotificationController) UnreadNotifications(c *gin.Context) {
	claims, err := request.CommonRequest{Context: c}.GetValidatedClaimsHeaders()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	filters, err := request.CommonRequest{Context: c}.GetValidatedCommonFilters()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	notifications, err := m.Repositories.NotificationRepository.RetrieveUnreadNotifications(
		claims.UserID, scope.BuildAndFilterInlineCondition(filters))
	if err != nil {
		c.JSON(ProcessError(err))
		return
	}
	c.JSON(http.StatusOK, notifications)
}

// AllNotifications returns all notifications for the current user
func (m *NotificationController) AllNotifications(c *gin.Context) {
	claims, err := request.CommonRequest{Context: c}.GetValidatedClaimsHeaders()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	filters, err := request.CommonRequest{Context: c}.GetValidatedCommonFilters()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	notifications, err := m.Repositories.NotificationRepository.RetrieveAllNotifications(
		claims.UserID, scope.BuildAndFilterInlineCondition(filters))
	if err != nil {
		c.JSON(ProcessError(err))
		return
	}
	c.JSON(http.StatusOK, notifications)
}

// ReadNotifications sets all notifications for the current user as read.
func (m *NotificationController) ReadNotifications(c *gin.Context) {
	claims, err := request.CommonRequest{Context: c}.GetValidatedClaimsHeaders()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	err = m.Repositories.NotificationRepository.ReadNotifications(claims.UserID)
	if err != nil {
		c.JSON(ProcessError(err))
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
