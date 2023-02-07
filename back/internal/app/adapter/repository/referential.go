package repository

import (
	"gorm.io/gorm"

	"gpe_project/internal/app/adapter/postgresql/model"
)

type ReferentialRepository struct {
	DB *gorm.DB
}

func (ar ReferentialRepository) FindAllWorkCategories() ([]*model.WorkCategory, error) {
	var categories []*model.WorkCategory
	result := ar.DB.
		Find(&categories)
	err := CatchFindAllError(result.Error)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (ar ReferentialRepository) FindByIdOrFail(id uint) (*model.WorkCategory, error) {
	var category *model.WorkCategory
	result := ar.DB.
		First(&category, id)
	err := CatchFindOrFailError(result.Error)
	if err != nil {
		return nil, err
	}

	return category, nil
}
