package model

import "time"

// Barometer defines a daily score for a barometer type. Unique constraint is set between the
// user, barometer type and day.
type Barometer struct {
	Model  `json:"-"`
	UserID uint       `json:"-" gorm:"not null"`
	Score  int        `json:"score" gorm:"not null"`
	Date   *time.Time `json:"date" gorm:"not null"`
	Type   string     `json:"type" gorm:"not null"`
	User   User       `json:"-"`
}
