package usecase

import (
	"fmt"

	"gpe_project/internal/app/adapter/postgresql/model"
	"gpe_project/internal/app/adapter/repository"
)

type UserUsecase struct {
	repository.Repositories
}

// CreateNewUser will create a new user in the user table. If the email is already in use
// The creation was aborted and throw an error.
func (u UserUsecase) CreateNewUser(user model.User) (userCreated *model.User, err error) {
	err = u.CheckIfUserExist(user)
	if err != nil {
		return nil, err
	}

	_, err = u.ReferentialRepository.FindByIdOrFail(*user.WorkCategoryID)
	if err != nil {
		return nil, err
	}

	user.Password, err = HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	userCreated, err = u.UserRepository.Create(user)
	return
}

// CheckIfUserExist will check if the user is already registered in the users table.
func (u UserUsecase) CheckIfUserExist(user model.User) error {
	userAccount, err := u.UserRepository.FindByEmail(user.Email)
	if err != nil {
		return err
	}

	if userAccount != nil {
		return fmt.Errorf("this email is already taken")
	}

	return nil
}

// Me will retrieve the current user connected
func (u UserUsecase) Me(userID uint) (userFound *model.User, err error) {
	userFound, err = u.UserRepository.FindByIDOrFail(userID)
	return
}

// UpdateUser will update the user information for connected user
func (u UserUsecase) UpdateUser(user *model.User) (*model.User, error) {
	_, err := u.UserRepository.FindByIDOrFail(user.ID)
	if err != nil {
		return nil, err
	}

	if user.WorkCategoryID != nil {
		_, err = u.ReferentialRepository.FindByIdOrFail(*user.WorkCategoryID)
		if err != nil {
			return nil, err
		}
	}

	return u.UserRepository.Update(user)
}
