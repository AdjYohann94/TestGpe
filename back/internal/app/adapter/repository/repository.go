package repository

import "gorm.io/gorm"

type Repositories struct {
	UserRepository         UserRepository
	QuizRepository         QuizRepository
	ReferentialRepository  ReferentialRepository
	BarometerRepository    BarometerRepository
	NotificationRepository NotificationRepository
}

// InitRepositories instantiate repositories once for the application lifecycle.
// Repositories could be linked to a database connection.
func InitRepositories(db *gorm.DB) Repositories {
	return Repositories{
		UserRepository{DB: db},
		QuizRepository{DB: db},
		ReferentialRepository{DB: db},
		BarometerRepository{DB: db},
		NotificationRepository{DB: db},
	}
}
