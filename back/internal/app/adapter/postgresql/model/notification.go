package model

type Level int

const (
	LevelMinimum Level = iota + 1
	LevelMedium
	LevelMaximum
)

type Notification struct {
	Model   `json:"-"`
	UserID  uint   `json:"-" gorm:"index,not null"`
	Read    bool   `json:"read" gorm:"index,default:false"`
	Message string `json:"message"`
	Level   Level  `json:"level" gorm:"default:1"`
	User    User   `json:"-"`
}
