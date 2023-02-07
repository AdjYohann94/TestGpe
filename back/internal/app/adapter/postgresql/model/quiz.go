package model

import (
	"time"
)

type QuizStatus string
type QuestionType string

const (
	QuizDraft     QuizStatus = "draft"
	QuizPublished QuizStatus = "published"
	QuizClosed    QuizStatus = "closed"
)

const (
	QuestionTypeText     QuestionType = "text"
	QuestionTypeSingle   QuestionType = "single"
	QuestionTypeMultiple QuestionType = "multiple"
)

type Quiz struct {
	Model
	Name           string          `json:"name" gorm:"not null"`
	Description    string          `json:"description" gorm:"not null"`
	StartedAt      *time.Time      `json:"startedAt,omitempty"`
	ClosedAt       *time.Time      `json:"closedAt,omitempty"`
	Status         QuizStatus      `json:"status" gorm:"default:draft;not null"`
	WorkCategories []*WorkCategory `json:"workCategories" gorm:"many2many:quiz_work_categories;"`
	Questions      []*Question     `json:"questions,omitempty"`
	CreatorID      int             `json:"creatorId"`
	User           *User           `json:"user,omitempty" gorm:"foreignKey:CreatorID"`
}

type Question struct {
	Model
	Description     string            `json:"description,omitempty"`
	QuizID          uint              `json:"-"`
	Responses       []*Response       `json:"responses,omitempty"`
	QuestionType    QuestionType      `json:"questionType"`
	ResponseChoices []*ResponseChoice `json:"choices,omitempty"`
}

type ResponseChoice struct {
	Model
	QuestionID int          `json:"-"`
	Value      string       `json:"value"`
	Type       QuestionType `json:"type"`
	Question   Question     `json:"-"`
}

type Response struct {
	Model
	UserID           int            `json:"userId,omitempty"`
	QuestionID       int            `json:"questionId,omitempty"`
	ResponseChoiceID *int           `json:"responseChoiceId,omitempty"`
	Value            *string        `json:"value,omitempty"`
	User             User           `json:"-"`
	Question         Question       `json:"-"`
	ResponseChoice   ResponseChoice `json:"-"`
}

func QuestionsContainsID(w []*Question, id uint) bool {
	for i := range w {
		if w[i].ID == id {
			return true
		}
	}
	return false
}

func ResponseChoicesContainsID(w []*ResponseChoice, id uint) bool {
	for i := range w {
		if w[i].ID == id {
			return true
		}
	}
	return false
}

func QuestionByIDFromQuestions(questions []*Question, questionID uint) *Question {
	for i := range questions {
		if questions[i].ID == questionID {
			return questions[i]
		}
	}
	return nil
}
