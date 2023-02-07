package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gpe_project/internal/app/adapter/repository"
)

func ValidationError(err error) (int, any) {
	return http.StatusUnprocessableEntity, gin.H{"message": err.Error()}
}

func ProcessError(err error) (int, any) {
	if errors.Is(err, &repository.NotFoundError{}) {
		return http.StatusNotFound, gin.H{"message": err.Error()}
	}
	return http.StatusBadRequest, gin.H{"message": err.Error()}
}

func UnauthorizedError(err error) (int, any) {
	return http.StatusUnauthorized, gin.H{"message": err.Error()}
}

func SuccessMessage(message string) (int, any) {
	return http.StatusOK, gin.H{"message": message}
}
