package model

import (
	"fmt"
	"gorm.io/gorm"
)

// AutoMigrateModels uses the gorm migrate tool to perform auto migration in database from models declaration
// This could be used in dev mode only. Relations between models can be performed if linked models are created.
// Backward references must be disabled on first migrate to avoid recursion.
func AutoMigrateModels(db *gorm.DB) {
	err := db.AutoMigrate(&WorkCategory{}, &User{}, &Response{}, &Question{}, &Quiz{}, &ResponseChoice{}, &Barometer{}, &Notification{})
	if err != nil {
		panic(fmt.Sprintf("Could not auto migrate : %s", err.Error()))
	}
}
