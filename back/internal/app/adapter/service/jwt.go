package service

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"

	"gpe_project/internal/app/adapter/postgresql/model"
)

// CreateAccessToken generate a signed string that contains user information
// with short validity.
func CreateAccessToken(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": strconv.Itoa(int(user.ID)),
		"type":   "access",
		"exp":    time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(getSecretKey())
	if err != nil {
		log.Println("[ERROR] cannot generate access token: ", err.Error())
		return "", fmt.Errorf("cannot generate access token")
	}

	return tokenString, nil
}

// CreateRefreshToken generate a signed string that contains uniq id
// to renew the access token when expired. The refresh token has a long
// validity expiration date.
func CreateRefreshToken(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": strconv.Itoa(int(user.ID)),
		"type":   "refresh",
		"exp":    time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString(getSecretKey())
	if err != nil {
		log.Println("[ERROR] cannot generate refresh token: ", err.Error())
		return "", fmt.Errorf("cannot generate refresh token")
	}

	return tokenString, nil
}

// ValidateJWTToken parse the jwtToken to ensure the token is valid against secret key and expiration date.
// Throw an error when token is not valid.
func ValidateJWTToken(jwtToken string) (*jwt.Token, error) {
	if jwtToken != "" {
		token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("you're unauthorized")
			}
			return getSecretKey(), nil
		})
		if err != nil {
			return nil, err
		}
		if token.Valid {
			return token, nil
		} else {
			return nil, errors.New("you're unauthorized due to invalid token")
		}
	} else {
		return nil, errors.New("you're unauthorized due to No token value")
	}
}

// getSecretKey returns the jwt secret phrase.
// Loaded using viper at KeyJwtSecret key.
func getSecretKey() []byte {
	secret := viper.GetString(KeyJwtSecret)
	return []byte(secret)
}

// ExtractJWTClaims returns a jwt.MapClaims that contains all claims for the
// jwt.Token. If token is not valid, no claims are returned.
func ExtractJWTClaims(token *jwt.Token) jwt.MapClaims {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims
	}

	return nil
}
