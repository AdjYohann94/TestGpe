package service

import (
	"testing"

	"github.com/maxatome/go-testdeep/td"
	"github.com/spf13/viper"

	"gpe_project/internal/app/adapter/postgresql/model"
)

func TestCreateAccessToken(t *testing.T) {
	user := model.User{
		Model: model.Model{ID: 1},
	}
	viper.Reset()
	viper.Set(KeyJwtSecret, "secret")

	accessToken, err := CreateAccessToken(&user)
	td.CmpNoError(t, err)
	td.CmpNot(t, accessToken, "")
}

func TestCreateRefreshToken(t *testing.T) {
	user := model.User{
		Model: model.Model{ID: 1},
	}
	viper.Reset()
	viper.Set(KeyJwtSecret, "secret")

	refreshToken, err := CreateRefreshToken(&user)
	td.CmpNoError(t, err)
	td.CmpNot(t, refreshToken, "")
}

func TestValidateJWTToken(t *testing.T) {
	user := model.User{
		Model: model.Model{ID: 1},
	}
	viper.Reset()
	viper.Set(KeyJwtSecret, "secret")
	accessToken, err := CreateAccessToken(&user)
	td.CmpNoError(t, err)

	token, err := ValidateJWTToken(accessToken)
	td.CmpNoError(t, err)
	td.CmpNotNil(t, token)
}

func TestValidateJWTTokenWithEmptyToken(t *testing.T) {
	viper.Reset()
	viper.Set(KeyJwtSecret, "secret")

	token, err := ValidateJWTToken("")
	td.CmpNil(t, token)
	td.CmpString(t, err, "you're unauthorized due to No token value")
}

func TestValidateJWTTokenWithWrongToken(t *testing.T) {
	user := model.User{
		Model: model.Model{ID: 1},
	}
	viper.Reset()
	viper.Set(KeyJwtSecret, "secret")
	accessToken, err := CreateAccessToken(&user)
	td.CmpNoError(t, err)

	// Update secret to invalid token
	viper.Set(KeyJwtSecret, "other")

	token, err := ValidateJWTToken(accessToken)
	td.CmpNil(t, token)
	td.CmpString(t, err, "signature is invalid")
}

func TestExtractJWTClaims(t *testing.T) {
	user := model.User{
		Model: model.Model{ID: 1},
	}
	viper.Reset()
	viper.Set(KeyJwtSecret, "secret")
	accessToken, err := CreateAccessToken(&user)
	td.CmpNoError(t, err)

	token, err := ValidateJWTToken(accessToken)
	claims := ExtractJWTClaims(token)
	td.CmpNotNil(t, claims)
	td.CmpMap(t,
		map[string]interface{}(claims),
		map[string]interface{}{"userID": "1", "type": "access"},
		td.MapEntries{"exp": td.Ignore()})
}
