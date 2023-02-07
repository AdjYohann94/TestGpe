package domain

type BarometerStat struct {
	Entries        int     `json:"entries"`
	Score          float64 `json:"score"`
	Type           string  `json:"type"`
	WorkCategoryId int     `json:"workCategoryId,omitempty"`
	Year           int     `json:"year,omitempty"`
	Month          int     `json:"month,omitempty"`
	Day            int     `json:"day,omitempty"`
}

type QuizGlobalStats struct {
	OpenQuiz          int          `json:"openQuiz"`
	AvgQuestionsQuiz  float64      `json:"avgQuestionsQuiz"`
	Quizzes           []QuizStat   `json:"quizzes"`
	QuizCreatedByYear []QuizByYear `json:"quizCreatedByYear"`
}

type QuizStat struct {
	QuizId            int `json:"quizId,omitempty"`
	PeopleResponded   int `json:"peopleResponded"`
	PeopleConcerned   int `json:"peopleConcerned"`
	NumberOfQuestions int `json:"numberOfQuestions"`
}

type QuizByYear struct {
	Year   int `json:"year,omitempty"`
	Number int `json:"number,omitempty"`
}
