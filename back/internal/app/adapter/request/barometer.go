package request

import (
	"github.com/gin-gonic/gin"

	"gpe_project/internal/app/adapter/postgresql/model"
)

type BarometerRequest struct{ Context *gin.Context }

type inputBarometer struct {
	Type  string `json:"type"`
	Score int    `json:"score"`
}

type ScopeFilter struct {
	ScopeCategories bool   `form:"scopeCategories"`
	TimeScope       string `form:"timeScope"`
}

func (r BarometerRequest) GetValidatedBarometerPayload() (*model.Barometer, error) {
	var moodInput inputBarometer
	if err := r.Context.ShouldBindJSON(&moodInput); err != nil {
		return nil, err
	}

	return &model.Barometer{
		Score: moodInput.Score,
		Type:  moodInput.Type,
	}, nil
}

func (r BarometerRequest) GetValidatedFilters() (*ScopeFilter, error) {
	var req ScopeFilter
	if err := r.Context.ShouldBindQuery(&req); err != nil {
		return nil, err
	}

	return &req, nil
}
