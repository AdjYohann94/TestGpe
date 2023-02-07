package repository

import (
	"fmt"
	"gorm.io/gorm"
	"time"

	"gpe_project/internal/app/adapter/postgresql/model"
	"gpe_project/internal/app/domain"
)

type QuizRepository struct {
	DB *gorm.DB
}

const peopleRespondToQuiz = `select count(distinct responses.user_id) as user_responded from responses inner join questions q on q.id = responses.question_id where quiz_id = ?`
const peopleConcernedByQuiz = `select count(*) as user_concerned from quiz_work_categories where quiz_id = ?`
const numberOfQuestionInQuiz = `select count(*) as number_of_questions from questions where quiz_id = ?`

const averageNumberOfQuestionInQuiz = `select avg(A.number_of_questions) as avg_questions_quiz from (select count(*) number_of_questions, quiz_id from questions group by quiz_id) as A`
const numberOfQuizCreatedByYear = `select count(*) as number, date_part('year', created_at) as year from quizzes where deleted_at is null group by year`
const numberOfOpenQuiz = `select count(*) as open_quiz from quizzes where status = 'published' and closed_at is null group by created_at`

const concernedRespondedQuiz = `select A.quiz_id, A.people_concerned, B.people_responded, C.number_of_questions from (select quiz_work_categories.quiz_id, count(*) as people_concerned from quizzes inner join quiz_work_categories on quizzes.id = quiz_work_categories.quiz_id where quizzes.deleted_at is null group by quiz_work_categories.quiz_id) as A left join (select quiz_id, count(distinct responses.user_id) people_responded from responses inner join questions q on q.id = responses.question_id group by q.quiz_id) as B on A.quiz_id = B.quiz_id left join (select count(*) number_of_questions, quiz_id from questions group by quiz_id) as C on A.quiz_id = C.quiz_id`

// FindAll returns all quizzes. Questions, responses and response choices
// are not sent back.
func (qr QuizRepository) FindAll() ([]*model.Quiz, error) {
	var quiz []*model.Quiz
	result := qr.DB.
		Preload("WorkCategories").
		Find(&quiz)
	err := CatchFindAllError(result.Error)
	if err != nil {
		return nil, err
	}

	return quiz, nil
}

func (qr QuizRepository) FindAllStats() (*domain.QuizGlobalStats, error) {
	var stats domain.QuizGlobalStats
	// Find number of people concerned of has responded
	result := qr.DB.Raw(concernedRespondedQuiz).Scan(&stats.Quizzes)
	err := CatchError(result.Error)
	if err != nil {
		return nil, err
	}

	// Find number of open quiz
	result = qr.DB.Raw(numberOfOpenQuiz).Scan(&stats.OpenQuiz)
	err = CatchError(result.Error)
	if err != nil {
		return nil, err
	}

	// Find number of quiz creation per year
	result = qr.DB.Raw(numberOfQuizCreatedByYear).Scan(&stats.QuizCreatedByYear)
	err = CatchError(result.Error)
	if err != nil {
		return nil, err
	}

	// Find average number of question in quiz
	result = qr.DB.Raw(averageNumberOfQuestionInQuiz).Scan(&stats.AvgQuestionsQuiz)
	err = CatchError(result.Error)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// UserHasRespondToQuiz returns true if the user respond to the quiz with at least one response.
func (qr QuizRepository) UserHasRespondToQuiz(userId uint, quiz *model.Quiz) (bool, error) {
	var count int64
	result := qr.DB.Preload("Question", "quiz_id = ?", quiz.ID).Find(&model.Response{}, "user_id = ?", userId).Count(&count)
	err := CatchFindAllError(result.Error)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindAllAvailableForUserWorkCategory returns quizzes for user where his work category is in
// quizzes work categories target and with status at least published. Questions, responses and response choices
// are not sent back.
func (qr QuizRepository) FindAllAvailableForUserWorkCategory(workCategoryID uint) ([]*model.Quiz, error) {
	var workCategory *model.WorkCategory
	result := qr.DB.
		Preload("Quiz.WorkCategories").
		Preload("Quiz", "status IN ?", []model.QuizStatus{model.QuizPublished, model.QuizClosed}).
		Where("id = ?", workCategoryID).
		First(&workCategory)
	err := CatchFindAllError(result.Error)
	if err != nil {
		return nil, err
	}

	return workCategory.Quiz, nil
}

// FindByIdOrFail returns the quiz by id without filter, if the quiz was not found, return
// not found error.
func (qr QuizRepository) FindByIdOrFail(id uint) (*model.Quiz, error) {
	var quiz *model.Quiz
	result := qr.DB.First(&quiz, id)
	err := CatchFindOrFailError(result.Error)
	if err != nil {
		return nil, err
	}

	return quiz, nil
}

// FindByIDFull returns the quiz by id with all responses of all users, questions and responses
// choices are included. Work categories also.
func (qr QuizRepository) FindByIDFull(id uint) (*model.Quiz, error) {
	var quiz *model.Quiz
	result := qr.DB.
		Preload("WorkCategories").
		Preload("Questions").
		Preload("Questions.Responses").
		Preload("Questions.ResponseChoices").
		First(&quiz, id)
	err := CatchFindOrFailError(result.Error)
	if err != nil {
		return nil, err
	}

	return quiz, nil
}

// FindQuestionWithAllResponses return the question by id with all responses of all users
func (qr QuizRepository) FindQuestionWithAllResponses(id uint) (*model.Question, error) {
	var question *model.Question
	result := qr.DB.
		Preload("ResponseChoices").
		Preload("Responses").
		First(&question, id)
	err := CatchFindAllError(result.Error)
	if err != nil {
		return nil, err
	}

	return question, nil
}

// FindByIDFullForUserMember returns the qui with user responses and questions with questions choices.
// The member can only view the quiz if his work category is on quiz work categories targets.
// If the requested quiz does not contain user work category, not found error is sent back.
func (qr QuizRepository) FindByIDFullForUserMember(quizID uint, user *model.User) (*model.Quiz, error) {
	var quiz *model.Quiz
	result := qr.DB.
		Preload("WorkCategories").
		Preload("Questions").
		Preload("Questions.Responses", qr.DB.Where(&model.Response{UserID: int(user.ID)})).
		Preload("Questions.ResponseChoices").
		First(&quiz, quizID)
	err := CatchFindOrFailError(result.Error)
	if err != nil {
		return nil, err
	}

	if model.WorkCategoriesContainsID(quiz.WorkCategories, *user.WorkCategoryID) {
		return quiz, nil
	}

	return nil, &NotFoundError{}
}

// FindByIDFullForUserAdmin returns the qui with user responses and questions with questions choices.
// The admin can view all quizzes.
func (qr QuizRepository) FindByIDFullForUserAdmin(quizID uint, user *model.User) (*model.Quiz, error) {
	var quiz *model.Quiz
	result := qr.DB.
		Preload("WorkCategories").
		Preload("Questions").
		Preload("Questions.Responses", qr.DB.Where(&model.Response{UserID: int(user.ID)})).
		Preload("Questions.ResponseChoices").
		First(&quiz, quizID)
	err := CatchFindOrFailError(result.Error)
	if err != nil {
		return nil, err
	}

	return quiz, nil
}

func (qr QuizRepository) FindStatsForQuizID(quizID uint) (*domain.QuizStat, error) {
	var stats domain.QuizStat
	// Find number of people respond to quiz
	result := qr.DB.Raw(peopleRespondToQuiz, quizID).Scan(&stats.PeopleResponded)
	err := CatchError(result.Error)
	if err != nil {
		return nil, err
	}

	// Find number of people concerned by the quiz
	result = qr.DB.Raw(peopleConcernedByQuiz, quizID).Scan(&stats.PeopleConcerned)
	err = CatchError(result.Error)
	if err != nil {
		return nil, err
	}

	// Get number of questions in the quiz
	result = qr.DB.Raw(numberOfQuestionInQuiz, quizID).Scan(&stats.NumberOfQuestions)
	err = CatchError(result.Error)
	if err != nil {
		return nil, err
	}

	stats.QuizId = int(quizID)

	return &stats, nil
}

// Create a new quiz. Associated work categories are validated before save.
// TODO : usecase to validate work category
func (qr QuizRepository) Create(quiz model.Quiz) (*model.Quiz, error) {
	err := qr.DB.Transaction(func(tx *gorm.DB) error {
		categoriesID := func() []uint {
			var id []uint
			for _, category := range quiz.WorkCategories {
				id = append(id, category.ID)
			}
			return id
		}()

		result := tx.Find(&[]model.WorkCategory{}, categoriesID)
		if int(result.RowsAffected) != len(quiz.WorkCategories) {
			return fmt.Errorf("some work categories does not exists")
		}

		result = tx.Create(&quiz)
		err := CatchCreateError(result.Error, result.RowsAffected)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		// Transaction aborted
		return nil, err
	}

	return &quiz, nil
}

// CreateResponses will save new quiz responses.
func (qr QuizRepository) CreateResponses(responses []*model.Response) error {
	result := qr.DB.Save(responses)
	err := CatchCreateError(result.Error, result.RowsAffected)
	return err
}

// RemoveByID soft deletes the quiz.
func (qr QuizRepository) RemoveByID(id uint) error {
	result := qr.DB.Delete(&model.Quiz{}, id)
	err := CatchDeleteError(result.Error, result.RowsAffected)
	if err != nil {
		return err
	}

	return nil
}

// UpdateStatus will update the quiz status. If status is closed, ClosedAt is filled with current time.
// If status is published, StartedAt is set to current time.
func (qr QuizRepository) UpdateStatus(quiz *model.Quiz, status model.QuizStatus) error {
	var result *gorm.DB
	if status == model.QuizClosed {
		result = qr.DB.Model(&quiz).Updates(map[string]interface{}{"Status": status, "ClosedAt": time.Now()})
	} else if status == model.QuizPublished {
		result = qr.DB.Model(&quiz).Updates(map[string]interface{}{"Status": status, "StartedAt": time.Now()})
	} else {
		return nil
	}
	err := CatchUpdateError(result.Error, result.RowsAffected)
	if err != nil {
		return err
	}

	return nil
}
