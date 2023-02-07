package model

const (
	RoleMember = "member"
	RoleAdmin  = "admin"
)

const (
	StatusActive = "active"
	StatusRevoke = "revoke"
)

type User struct {
	Model
	FirstName      string        `json:"firstName" gorm:"not null"`
	LastName       string        `json:"lastName" gorm:"not null"`
	Email          string        `json:"email" gorm:"uniqueIndex;not null"`
	Password       string        `json:"-" gorm:"not null"`
	Role           string        `json:"role,omitempty" gorm:"default:member;not null"`
	Status         string        `json:"status" gorm:"default:active"`
	PhoneNumber    *string       `json:"phoneNumber"`
	ZipCode        *string       `json:"zipCode"`
	Address        *string       `json:"address"`
	City           *string       `json:"city"`
	WorkCategoryID *uint         `json:"workCategoryId"`
	WorkCategory   *WorkCategory `json:"workCategory,omitempty"`
}

type WorkCategory struct {
	Model
	Name  string  `json:"name,omitempty" gorm:"unique;not null"`
	Users []*User `json:"users,omitempty"`
	Quiz  []*Quiz `json:"quiz,omitempty" gorm:"many2many:quiz_work_categories;"`
}

func WorkCategoriesContainsID(w []*WorkCategory, id uint) bool {
	for i := range w {
		if w[i].ID == id {
			return true
		}
	}
	return false
}
