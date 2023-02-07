package usecase

import (
	"fmt"

	"gpe_project/internal/app/domain"
	"gpe_project/internal/app/domain/valueobject"

	"gpe_project/internal/app/adapter/postgresql/model"
	"gpe_project/internal/app/adapter/repository"
)

type QuizUsecase struct {
	repository.Repositories
}

// FindQuizzesWithHasResponded calls FindQuizzes and add a responded proerty on each item
// to say of the connected user as at least one response for this quiz.
func (qu QuizUsecase) FindQuizzesWithHasResponded(userID uint) ([]domain.Quiz, error) {
	quizzes, err := qu.FindQuizzes(userID)
	if err != nil {
		return nil, err
	}

	// Iterate over each item to determine is user has responded
	output := make([]domain.Quiz, len(quizzes))
	for i, quiz := range quizzes {
		hasResponded, err := qu.QuizRepository.UserHasRespondToQuiz(userID, quiz)
		if err != nil {
			return nil, err
		}
		output[i] = domain.Quiz{
			Quiz:      quiz,
			Responded: hasResponded,
		}
	}

	return output, nil
}

// FindQuizzes will return all quizzes for the connected user. If the user is role admin, all quizzes
// are sent back. Else if user is only a member, return quiz for his work category and with status published or closed.
// If the user has no work category, send back an error.
// Responses and questions are not included if the output payload.
func (qu QuizUsecase) FindQuizzes(userID uint) (quizzes []*model.Quiz, err error) {
	userFound, err := qu.UserRepository.FindByIDOrFail(userID)
	if err != nil {
		return
	}

	if userFound.Role == model.RoleAdmin {
		return qu.QuizRepository.FindAll()
	}

	return qu.QuizRepository.FindAllAvailableForUserWorkCategory(*userFound.WorkCategoryID)
}

// FindQuizByID will return quiz by id.
// Only data for connected user is sent back (responses). If user is member, he cannot access the quiz if his
// work category is not in quiz target. Admin can retrieve details from all quizzes.
func (qu QuizUsecase) FindQuizByID(userID, quizID uint) (quiz *model.Quiz, err error) {
	userFound, err := qu.UserRepository.FindByIDOrFail(userID)
	if err != nil {
		return
	}

	if userFound.Role == model.RoleAdmin {
		return qu.QuizRepository.FindByIDFullForUserAdmin(quizID, userFound)
	}

	return qu.QuizRepository.FindByIDFullForUserMember(quizID, userFound)
}

// AnswerQuiz will save quiz response for the current user. THe quiz must be in status model.QuizPublished.
// This function will check questions references and choice responses. If question type is text, choices are
// not required but value yes. Once responses saved, they are sent back. The user with role member can answer a quiz
// only if he's on the quiz work category target. Admin can answer all quizzes.
func (qu QuizUsecase) AnswerQuiz(userID, quizId uint, responsesInput []*valueobject.Response) ([]*model.Response, error) {
	_, err := qu.UserRepository.FindByIDOrFail(userID)
	if err != nil {
		return nil, err
	}

	quiz, err := qu.FindQuizByID(userID, quizId)
	if err != nil {
		return nil, err
	}

	if quiz.Status != model.QuizPublished {
		return nil, fmt.Errorf("cannot answer an unpublished quiz")
	}

	var responses []*model.Response
	for _, response := range responsesInput {
		if !model.QuestionsContainsID(quiz.Questions, uint(response.QuestionId)) {
			return nil, fmt.Errorf("question with id %d not found in the Quiz %d", response.QuestionId, quiz.ID)
		}

		question := model.QuestionByIDFromQuestions(quiz.Questions, uint(response.QuestionId))
		if len(question.Responses) != 0 {
			return nil, fmt.Errorf("user has already answered this quiz")
		}

		res, err := BuildResponsesFromQuestionType(userID, question, response)
		if err != nil {
			return nil, err
		}
		responses = append(responses, res...)
	}

	if err := qu.Repositories.QuizRepository.CreateResponses(responses); err != nil {
		return nil, err
	}
	return responses, nil
}

// BuildResponsesFromQuestionType will build a response input from question and response value object.
// THis function check the question type and make some validations. If the question type is type model.QuestionTypeText,
// value is required and no response choices. If question type is model.QuestionTypeSingle or model.QuestionTypeMultiple
// choices are required. Choices are compared than available responses choices for the question.
func BuildResponsesFromQuestionType(
	userID uint,
	question *model.Question,
	response *valueobject.Response,
) ([]*model.Response, error) {
	if question.QuestionType == model.QuestionTypeText {
		if response.Value == "" || len(response.Choices) != 0 {
			return nil, fmt.Errorf("the question with id %d is of type text. the value field must be defined", question.ID)
		}
		return []*model.Response{{
			UserID:     int(userID),
			QuestionID: int(question.ID),
			Value:      &response.Value,
		}}, nil
	}
	if question.QuestionType == model.QuestionTypeMultiple || question.QuestionType == model.QuestionTypeSingle {
		if len(response.Choices) == 0 {
			return nil, fmt.Errorf("the question with id %d is of type multiple or single. the value choices must be passed with minimum one choice", question.ID)
		}
		if question.QuestionType == model.QuestionTypeSingle && len(response.Choices) > 1 {
			return nil, fmt.Errorf("the question with id %d is of type single. the value choices must be passed with maximum one choice", question.ID)
		}
		var res []*model.Response
		for _, responseChoice := range response.Choices {
			if !model.ResponseChoicesContainsID(question.ResponseChoices, uint(responseChoice)) {
				return nil, fmt.Errorf("the choice with id %d is not valid for the question with id %d", responseChoice, question.ID)
			}
			choice := responseChoice
			res = append(res, &model.Response{
				UserID:           int(userID),
				QuestionID:       int(question.ID),
				ResponseChoiceID: &choice,
			})
		}
		return res, nil
	}

	return nil, fmt.Errorf("invalid question type : %s", question.QuestionType)
}
