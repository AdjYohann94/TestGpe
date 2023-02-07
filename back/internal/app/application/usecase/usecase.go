package usecase

import "gpe_project/internal/app/adapter/repository"

type Usecase struct {
	UserUsecase  UserUsecase
	AuthUsecase  AuthUsecase
	QuizUsecase  QuizUsecase
	Repositories repository.Repositories
}

// InitUsecase will create usecase struct once to reuse them during the application
// lifecycle because there is no need to create usecase for each request.
func InitUsecase(rep repository.Repositories) Usecase {
	return Usecase{
		UserUsecase:  UserUsecase{rep},
		AuthUsecase:  AuthUsecase{rep},
		QuizUsecase:  QuizUsecase{rep},
		Repositories: rep,
	}
}
