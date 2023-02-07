package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/maxatome/go-testdeep/td"
	"github.com/spf13/viper"
)

func TestAccessTokenAuthenticationWithoutAccessToken(t *testing.T) {
	viper.Reset()
	viper.Set("jwt.secret", "secret")

	w := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(w)

	engine.Use(AccessTokenAuthentication())

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/", http.NoBody)
	engine.ServeHTTP(w, req)

	td.Cmp(t, w.Code, http.StatusUnauthorized)
	td.Cmp(t, w.Body.String(), `{"message":"you're unauthorized due to No token value"}`)
}

func TestAccessTokenAuthenticationWithBadAccessToken(t *testing.T) {
	viper.Reset()
	viper.Set("jwt.secret", "secret")

	w := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(w)

	engine.Use(AccessTokenAuthentication())

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/", http.NoBody)
	req.Header.Set("X-Access-Token", "token")

	engine.ServeHTTP(w, req)

	td.Cmp(t, w.Code, http.StatusUnauthorized)
	td.Cmp(t, w.Body.String(), `{"message":"token contains an invalid number of segments"}`)
}

func TestRefreshTokenAuthenticationWithoutRefreshToken(t *testing.T) {
	viper.Reset()
	viper.Set("jwt.secret", "secret")

	w := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(w)

	engine.Use(RefreshTokenAuthentication())

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/", http.NoBody)
	engine.ServeHTTP(w, req)

	td.Cmp(t, w.Code, http.StatusUnauthorized)
	td.Cmp(t, w.Body.String(), `{"message":"you're unauthorized due to No token value"}`)
}

func TestRefreshTokenAuthenticationWithBadRefreshToken(t *testing.T) {
	viper.Reset()
	viper.Set("jwt.secret", "secret")

	w := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(w)

	engine.Use(RefreshTokenAuthentication())

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/", http.NoBody)
	req.Header.Set("X-Refresh-Token", "token")

	engine.ServeHTTP(w, req)

	td.Cmp(t, w.Code, http.StatusUnauthorized)
	td.Cmp(t, w.Body.String(), `{"message":"token contains an invalid number of segments"}`)
}
