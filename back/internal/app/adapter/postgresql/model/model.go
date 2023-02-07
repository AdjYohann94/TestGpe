package model

import (
	"gorm.io/gorm"
	"time"
)

// Model provides a model base with id primary key and timestamp. CreatedAt and
// UpdatedAt fields are tracked and updating by gorm. The soft delete is handle by the
// DeleteAt property. This model redeclared this gorm.Model but with json tags for
// camel case output. The deleteAt property is not serialized.
type Model struct {
	ID        uint           `json:"id,omitempty" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
