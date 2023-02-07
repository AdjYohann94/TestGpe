package usecase

import (
	"testing"

	"github.com/maxatome/go-testdeep/td"

	"gpe_project/internal/app/adapter/postgresql/model"
	"gpe_project/internal/app/domain/valueobject"
)

func TestBuildResponseWithBadResponseChoiceId(t *testing.T) {
	response, err := BuildResponsesFromQuestionType(
		1,
		&model.Question{
			Model:        model.Model{ID: 1},
			QuestionType: "single",
			ResponseChoices: []*model.ResponseChoice{{
				Model: model.Model{ID: 1},
			}},
		},
		&valueobject.Response{Choices: []int{2}},
	)
	td.CmpNil(t, response)
	td.Cmp(t, err.Error(), "the choice with id 2 is not valid for the question with id 1")
}

func TestBuildResponseWithNotEnoughResponseChoicesForTypeSingle(t *testing.T) {
	response, err := BuildResponsesFromQuestionType(
		1,
		&model.Question{
			Model:        model.Model{ID: 1},
			QuestionType: "single",
			ResponseChoices: []*model.ResponseChoice{{
				Model: model.Model{ID: 1},
			}},
		},
		&valueobject.Response{},
	)
	td.CmpNil(t, response)
	td.Cmp(t, err.Error(), "the question with id 1 is of type multiple or single. the value choices must be passed with minimum one choice")
}

func TestBuildResponseWithMoreThanOneResponseChoiceForTypeSingle(t *testing.T) {
	response, err := BuildResponsesFromQuestionType(
		1,
		&model.Question{
			Model:        model.Model{ID: 1},
			QuestionType: "single",
			ResponseChoices: []*model.ResponseChoice{{
				Model: model.Model{ID: 1},
			}},
		},
		&valueobject.Response{Choices: []int{1, 2}},
	)
	td.CmpNil(t, response)
	td.Cmp(t, err.Error(), "the question with id 1 is of type single. the value choices must be passed with maximum one choice")
}

func TestBuildResponseWithoutQuestionType(t *testing.T) {
	response, err := BuildResponsesFromQuestionType(
		1,
		&model.Question{
			Model:        model.Model{ID: 1},
			QuestionType: "",
		},
		&valueobject.Response{},
	)
	td.CmpNil(t, response)
	td.Cmp(t, err.Error(), "invalid question type : ")
}

func TestBuildResponseWithTypeTextAndNoValueDefined(t *testing.T) {
	response, err := BuildResponsesFromQuestionType(
		1,
		&model.Question{
			Model:        model.Model{ID: 1},
			QuestionType: "text",
		},
		&valueobject.Response{Value: ""},
	)
	td.CmpNil(t, response)
	td.Cmp(t, err.Error(), "the question with id 1 is of type text. the value field must be defined")
}

func TestBuildResponseWithTypeTextANdValueDefined(t *testing.T) {
	response, err := BuildResponsesFromQuestionType(
		1,
		&model.Question{
			Model:        model.Model{ID: 1},
			QuestionType: "text",
		},
		&valueobject.Response{
			QuestionId: 1,
			Choices:    nil,
			Value:      "some text",
		},
	)
	td.CmpNil(t, err)
	td.Cmp(t, response, []*model.Response{{
		UserID:           1,
		QuestionID:       1,
		ResponseChoiceID: nil,
		Value:            NewPointer("some text"),
	}})
}

func TestBuildResponseWithTypeSingle(t *testing.T) {
	response, err := BuildResponsesFromQuestionType(
		1,
		&model.Question{
			Model:           model.Model{ID: 1},
			QuestionType:    "single",
			ResponseChoices: []*model.ResponseChoice{{Model: model.Model{ID: 1}}},
		},
		&valueobject.Response{
			QuestionId: 1,
			Choices:    []int{1},
		},
	)
	td.CmpNil(t, err)
	td.Cmp(t, response, []*model.Response{{
		UserID:           1,
		QuestionID:       1,
		ResponseChoiceID: NewPointer(1),
	}})
}

func TestBuildResponseWithTypeMultiple(t *testing.T) {
	response, err := BuildResponsesFromQuestionType(
		1,
		&model.Question{
			Model:           model.Model{ID: 1},
			QuestionType:    "multiple",
			ResponseChoices: []*model.ResponseChoice{{Model: model.Model{ID: 1}}, {Model: model.Model{ID: 2}}},
		},
		&valueobject.Response{
			QuestionId: 1,
			Choices:    []int{1, 2},
		},
	)
	td.CmpNil(t, err)
	td.Cmp(t, response, []*model.Response{
		{
			UserID:           1,
			QuestionID:       1,
			ResponseChoiceID: NewPointer(1),
		},
		{
			UserID:           1,
			QuestionID:       1,
			ResponseChoiceID: NewPointer(2),
		}})
}

func NewPointer[T any](a T) *T {
	return &a
}
