package request

import (
	"github.com/gin-gonic/gin"

	"gpe_project/internal/app/adapter/postgresql/model"
)

type UserRequest struct{ *gin.Context }

type inputUpdateUser struct {
	FirstName      *string `json:"firstName"`
	LastName       *string `json:"lastName"`
	PhoneNumber    *string `json:"phoneNumber"`
	ZipCode        *string `json:"zipCode"`
	Address        *string `json:"address"`
	City           *string `json:"city"`
	WorkCategoryID *uint   `json:"workCategoryId"`
}

// GetValidatedUpdateUserPayload returns the user model from login information provided in the
// input payload to update user information.
func (r UserRequest) GetValidatedUpdateUserPayload() (*model.User, error) {
	var req inputUpdateUser
	if err := r.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	return inputUpdateUserToUserModel(req), nil
}

func inputUpdateUserToUserModel(req inputUpdateUser) *model.User {
	m := &model.User{
		PhoneNumber:    req.PhoneNumber,
		ZipCode:        req.ZipCode,
		Address:        req.Address,
		City:           req.City,
		WorkCategoryID: req.WorkCategoryID,
	}
	if req.FirstName != nil {
		m.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		m.LastName = *req.LastName
	}
	return m
}
