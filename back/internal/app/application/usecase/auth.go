package usecase

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"

	"gpe_project/internal/app/adapter/postgresql/model"
	"gpe_project/internal/app/adapter/repository"
	"gpe_project/internal/app/adapter/service"
	"gpe_project/internal/app/domain/valueobject"
)

type AuthUsecase struct {
	repository.Repositories
}

// Login will perform authentication, login and password must match an account information in the database
// If any user is found and password match the username, a new refresh uuid is generated and stored.
// AccessToken et Refresh Tokens are created and send back. The refresh token is generated only is the stay connected
// parameter is set to true.
func (u AuthUsecase) Login(loginInput *valueobject.Login) (accessToken, refreshToken string, err error) {
	userFound, err := u.UserRepository.FindByEmail(loginInput.Email)
	if err != nil {
		return
	}
	if userFound == nil {
		err = fmt.Errorf("invalid credentials")
		return
	}

	err = CompareHashPassword(loginInput.Password, userFound.Password)
	if err != nil {
		return
	}

	accessToken, err = service.CreateAccessToken(userFound)
	if err != nil {
		return
	}

	if loginInput.StayConnected {
		refreshToken, err = service.CreateRefreshToken(userFound)
		if err != nil {
			return
		}
	}

	return
}

// Refresh will check the user account status and if the account is valid. Generates a new access token
func (u AuthUsecase) Refresh(userID uint) (accessToken, refreshToken string, err error) {
	userFound, err := u.UserRepository.FindByIDOrFail(userID)
	if err != nil {
		return
	}

	if userFound.Status == model.StatusRevoke {
		err = fmt.Errorf("user is revoked")
		return
	}

	accessToken, err = service.CreateAccessToken(userFound)
	if err != nil {
		return
	}

	refreshToken, err = service.CreateRefreshToken(userFound)
	return
}

// HashPassword render the password hash using bcrypt from clear string
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Println("[ERROR] hash password has failed: ", err.Error())
		return "", fmt.Errorf("an error occured during password hash")
	}

	return string(bytes), nil
}

// CompareHashPassword check the hash password against clear input password.
// An error was returned is passwords not match
func CompareHashPassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return fmt.Errorf("invalid credentials")
	}

	return nil
}
