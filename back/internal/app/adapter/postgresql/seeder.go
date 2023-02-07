package postgresql

import (
	"gorm.io/gorm"

	"gpe_project/internal/app/adapter/postgresql/model"
	"gpe_project/internal/app/adapter/postgresql/seeds"
)

// Seeder struct is used only on test actions to inject, delete data in
// database.
type Seeder struct {
	db *gorm.DB
}

func NewSeeder(db *gorm.DB) Seeder {
	return Seeder{db: db}
}

func (s Seeder) AddUserNotifications(userID uint) {
	s.db.Create(seeds.GetNotificationSeedForUser(userID, 60))
}

func (s Seeder) PurgeNotifications() {
	s.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&model.Notification{})
}

func (s Seeder) AddBarometers(userID uint) {
	s.db.Create(seeds.GetBarometers(userID, 100))
}

func (s Seeder) PurgeBarometers() {
	s.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&model.Barometer{})
}
