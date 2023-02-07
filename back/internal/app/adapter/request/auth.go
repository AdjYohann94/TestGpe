package request

import (
	"strings"

	"gpe_project/internal/app/adapter/postgresql/model"
	"gpe_project/internal/app/domain/valueobject"

	"github.com/gin-gonic/gin"
)

type AuthRequest struct{ *gin.Context }

type inputRegister struct {
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	WorkCategoryID *uint  `json:"workCategoryId"`
}

type inputLogin struct {
	Email         string `json:"email"`
	Password      string `json:"password"`
	StayConnected *bool  `json:"stayConnected"`
}

// GetValidatedRegisterPayload returns the user model from register information provided in the
// input payload.
func (r AuthRequest) GetValidatedRegisterPayload() (*model.User, error) {
	var req inputRegister
	if err := r.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	return inputRegisterToUserModel(req), nil
}

// GetValidatedLoginPayload returns the user model from login information provided in the
// input payload.
func (r AuthRequest) GetValidatedLoginPayload() (*valueobject.Login, error) {
	var req inputLogin
	if err := r.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	return inputLoginToLoginValueObject(req), nil
}

func inputRegisterToUserModel(req inputRegister) *model.User {
	user := &model.User{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		Password:       req.Password,
		Role:           model.RoleMember,
		WorkCategoryID: req.WorkCategoryID,
	}
	// TODO : temporary during demonstration
	if strings.Contains(user.Email, "admin") {
		user.Role = model.RoleAdmin
	}
	return user
}

func inputLoginToLoginValueObject(req inputLogin) *valueobject.Login {
	return &valueobject.Login{
		Email:         req.Email,
		Password:      req.Password,
		StayConnected: *req.StayConnected,
	}
}
