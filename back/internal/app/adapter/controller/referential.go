package controller

import (
	"net/http"

	"gpe_project/internal/app/application/usecase"

	"github.com/gin-gonic/gin"
)

type ReferentialController struct {
	usecase.Usecase
}

func NewReferentialController(e *gin.Engine, usecase usecase.Usecase) *ReferentialController {
	controller := &ReferentialController{usecase}

	e.GET("/work-category", controller.WorkCategories)

	return controller
}

func (r *ReferentialController) WorkCategories(c *gin.Context) {
	categories, err := r.Repositories.ReferentialRepository.FindAllWorkCategories()
	if err != nil {
		c.JSON(ProcessError(err))
		return
	}

	c.JSON(http.StatusOK, categories)
}
