package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"gpe_project/internal/app/adapter/middleware"
	"gpe_project/internal/app/adapter/postgresql/model"
	"gpe_project/internal/app/adapter/request"
	"gpe_project/internal/app/application/usecase"
)

type QuizController struct {
	usecase.Usecase
}

// NewQuizController boostrap quiz endpoints with their controllers.
// These routes are located under /quiz prefix.
func NewQuizController(e *gin.Engine, usecase usecase.Usecase) *QuizController {
	controller := &QuizController{usecase}

	group := e.Group("/quiz")
	group.GET("", controller.FindAll)
	group.GET("/stats", middleware.AdminOnly(usecase.Repositories), controller.FindAllStats)
	group.POST("", middleware.AdminOnly(usecase.Repositories), controller.Create)
	group.GET("/:id", controller.Find)
	group.GET("/:id/stats", middleware.AdminOnly(usecase.Repositories), controller.FindStats)
	group.DELETE("/:id", middleware.AdminOnly(usecase.Repositories), controller.Remove)
	group.POST("/:id/publish", middleware.AdminOnly(usecase.Repositories), controller.Publish)
	group.POST("/:id/close", middleware.AdminOnly(usecase.Repositories), controller.Close)
	group.POST("/:id/response", controller.AnswerQuiz)
	group.GET("/:id/response", middleware.AdminOnly(usecase.Repositories), controller.FindWithResponses)
	e.GET("/question/:id/response", middleware.AdminOnly(usecase.Repositories), controller.FindResponsesByQuestion)

	return controller
}

func (q *QuizController) AnswerQuiz(c *gin.Context) {
	claims, err := request.CommonRequest{Context: c}.GetValidatedClaimsHeaders()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	req, err := request.CommonRequest{Context: c}.GetValidatedResourceIdentifier()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	responsesInput, err := request.QuestRequest{Context: c}.GetValidatedQuizAnswersPayload()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	responses, err := q.QuizUsecase.AnswerQuiz(claims.UserID, req.ID, responsesInput)
	if err != nil {
		c.JSON(ProcessError(err))
		return
	}

	c.JSON(http.StatusCreated, responses)
}

func (q *QuizController) Close(c *gin.Context) {
	req, err := request.CommonRequest{Context: c}.GetValidatedResourceIdentifier()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}
	quiz, err := q.Repositories.QuizRepository.FindByIdOrFail(req.ID)
	if err != nil {
		c.JSON(ProcessError(err))
		return
	}

	switch quiz.Status {
	case model.QuizDraft, model.QuizPublished:
		if err := q.Repositories.QuizRepository.UpdateStatus(quiz, model.QuizClosed); err != nil {
			c.JSON(ProcessError(err))
			return
		}
		c.JSON(SuccessMessage("quiz closed"))
	case model.QuizClosed:
		c.JSON(ProcessError(fmt.Errorf("quiz already closed")))
	}
}

func (q *QuizController) Publish(c *gin.Context) {
	req, err := request.CommonRequest{Context: c}.GetValidatedResourceIdentifier()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	quiz, err := q.Repositories.QuizRepository.FindByIdOrFail(req.ID)
	if err != nil {
		c.JSON(ProcessError(err))
		return
	}

	switch quiz.Status {
	case model.QuizDraft:
		if err := q.Repositories.QuizRepository.UpdateStatus(quiz, model.QuizPublished); err != nil {
			c.JSON(ProcessError(err))
			return
		}
		c.JSON(SuccessMessage("quiz published successfully"))
	case model.QuizPublished:
		c.JSON(ProcessError(fmt.Errorf("quiz already published")))
	case model.QuizClosed:
		c.JSON(ProcessError(fmt.Errorf("quiz closed")))
	}
}

func (q *QuizController) Remove(c *gin.Context) {
	req, err := request.CommonRequest{Context: c}.GetValidatedResourceIdentifier()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	err = q.Repositories.QuizRepository.RemoveByID(req.ID)
	if err != nil {
		c.JSON(ProcessError(err))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (q *QuizController) Create(c *gin.Context) {
	claims, err := request.CommonRequest{Context: c}.GetValidatedClaimsHeaders()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	quiz, err := request.QuestRequest{Context: c}.GetValidatedCreateQuizPayload()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	quiz.CreatorID = int(claims.UserID)
	newQuiz, err := q.Repositories.QuizRepository.Create(*quiz)
	if err != nil {
		c.JSON(ProcessError(err))
		return
	}

	quizWithCategories, err := q.Repositories.QuizRepository.FindByIDFull(newQuiz.ID)
	if err != nil {
		c.JSON(ProcessError(err))
		return
	}

	c.JSON(http.StatusCreated, quizWithCategories)
}

func (q *QuizController) FindAll(c *gin.Context) {
	claims, err := request.CommonRequest{Context: c}.GetValidatedClaimsHeaders()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	quizzes, err := q.QuizUsecase.FindQuizzesWithHasResponded(claims.UserID)
	if err != nil {
		c.JSON(ProcessError(err))
		return
	}

	c.JSON(http.StatusOK, quizzes)
}

func (q *QuizController) FindAllStats(c *gin.Context) {
	stats, err := q.Repositories.QuizRepository.FindAllStats()
	if err != nil {
		c.JSON(ProcessError(err))
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (q *QuizController) Find(c *gin.Context) {
	req, err := request.CommonRequest{Context: c}.GetValidatedResourceIdentifier()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	claims, err := request.CommonRequest{Context: c}.GetValidatedClaimsHeaders()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	quiz, err := q.QuizUsecase.FindQuizByID(claims.UserID, req.ID)
	if err != nil {
		c.JSON(ProcessError(err))
		return
	}

	c.JSON(http.StatusOK, quiz)
}

func (q *QuizController) FindStats(c *gin.Context) {
	req, err := request.CommonRequest{Context: c}.GetValidatedResourceIdentifier()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	_, err = q.Repositories.QuizRepository.FindByIdOrFail(req.ID)
	if err != nil {
		c.JSON(ProcessError(err))
		return
	}

	stats, err := q.Repositories.QuizRepository.FindStatsForQuizID(req.ID)
	if err != nil {
		c.JSON(ProcessError(err))
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (q *QuizController) FindWithResponses(c *gin.Context) {
	req, err := request.CommonRequest{Context: c}.GetValidatedResourceIdentifier()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	quiz, err := q.Repositories.QuizRepository.FindByIDFull(req.ID)
	if err != nil {
		c.JSON(ProcessError(err))
		return
	}

	c.JSON(http.StatusOK, quiz)
}

func (q *QuizController) FindResponsesByQuestion(c *gin.Context) {
	req, err := request.CommonRequest{Context: c}.GetValidatedResourceIdentifier()
	if err != nil {
		c.JSON(ValidationError(err))
		return
	}

	question, err := q.Repositories.QuizRepository.FindQuestionWithAllResponses(req.ID)
	if err != nil {
		c.JSON(ProcessError(err))
		return
	}

	c.JSON(http.StatusOK, question)
}
