package request

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"gpe_project/internal/app/adapter/postgresql/model"
	"gpe_project/internal/app/domain/valueobject"
)

type QuestRequest struct{ Context *gin.Context }

type inputQuiz struct {
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	WorkCategories []uint          `json:"workCategories"`
	Questions      []inputQuestion `json:"questions"`
}

type inputQuestion struct {
	Description  string             `json:"description"`
	QuestionType model.QuestionType `json:"questionType"`
	Choices      []string           `json:"choices"`
}

type inputResponse struct {
	QuestionId int    `json:"questionId"`
	Choices    []int  `json:"choices"`
	Value      string `json:"value"`
}

func (request QuestRequest) GetValidatedQuizAnswersPayload() ([]*valueobject.Response, error) {
	var responsesInput []inputResponse
	if err := request.Context.ShouldBindJSON(&responsesInput); err != nil {
		return nil, err
	}

	var responses []*valueobject.Response
	for _, response := range responsesInput {
		responses = append(responses, inputResponseToResponseValueObject(response))
	}
	return responses, nil
}

func (request QuestRequest) GetValidatedCreateQuizPayload() (*model.Quiz, error) {
	var quizInput inputQuiz
	if err := request.Context.ShouldBindJSON(&quizInput); err != nil {
		return nil, err
	}

	questions := make([]*model.Question, len(quizInput.Questions))
	for i, question := range quizInput.Questions {
		questions[i] = &model.Question{
			Description:  question.Description,
			QuestionType: question.QuestionType,
		}

		if question.QuestionType == model.QuestionTypeText {
			continue
		} else {
			if len(question.Choices) < 2 {
				return nil, fmt.Errorf("at least two choices are required for question type %s", question.QuestionType)
			}
		}

		choices := make([]*model.ResponseChoice, len(question.Choices))
		for j, choice := range question.Choices {
			choices[j] = &model.ResponseChoice{
				Type:  question.QuestionType,
				Value: choice,
			}
		}

		questions[i].ResponseChoices = choices
	}

	categories := make([]*model.WorkCategory, len(quizInput.WorkCategories))
	for i, categoryID := range quizInput.WorkCategories {
		categories[i] = &model.WorkCategory{
			Model: model.Model{
				ID: categoryID,
			},
		}
	}

	return &model.Quiz{
		Name:           quizInput.Name,
		Description:    quizInput.Description,
		WorkCategories: categories,
		Questions:      questions,
	}, nil
}

func inputResponseToResponseValueObject(req inputResponse) *valueobject.Response {
	return &valueobject.Response{
		QuestionId: req.QuestionId,
		Choices:    req.Choices,
		Value:      req.Value,
	}
}
