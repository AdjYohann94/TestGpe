package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"gpe_project/internal/app/adapter/postgresql/model"
)

type UserRepository struct {
	DB *gorm.DB
}

func (ar UserRepository) Create(user model.User) (*model.User, error) {
	result := ar.DB.Create(&user)
	err := CatchCreateError(result.Error, result.RowsAffected)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ar UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	result := ar.DB.First(&user, "email = ?", email)
	err := CatchFindError(result.Error)
	if err != nil {
		return nil, err
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, nil
}

func (ar UserRepository) FindByEmailOrFail(email string) (*model.User, error) {
	var user model.User
	result := ar.DB.First(&user, "email = ?", email)
	err := CatchFindOrFailError(result.Error)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ar UserRepository) FindByIDOrFail(userID uint) (*model.User, error) {
	var user model.User
	result := ar.DB.First(&user, userID)
	err := CatchFindOrFailError(result.Error)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ar UserRepository) Update(user *model.User) (*model.User, error) {
	if user.WorkCategoryID != nil {
		result := ar.DB.Find(&[]model.WorkCategory{}, user.WorkCategoryID)
		if int(result.RowsAffected) == 0 {
			return nil, fmt.Errorf("some work categories does not exists")
		}
	}

	result := ar.DB.Omit("email", "password", "status", "role").Clauses(clause.Returning{}).Updates(&user)
	err := CatchUpdateError(result.Error, result.RowsAffected)
	if err != nil {
		return nil, err
	}

	return user, nil
}
