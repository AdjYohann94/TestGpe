package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/maxatome/go-testdeep/td"

	"gpe_project/internal/app/adapter/service"
)

func TestKinValidatorMiddlewareNotFoundPath(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(w)

	s := service.NewKinValidator()
	engine.Use(Kin(s))

	inputPayload := "{}"

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/noFound", strings.NewReader(inputPayload))
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(w, req)

	td.Cmp(t, w.Code, http.StatusNotFound)
	td.Cmp(t, w.Body.String(), `{"message":"route not found"}`)
}

func TestKinValidatorMiddlewareInputValidation(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(w)

	s := service.NewKinValidator()
	engine.Use(Kin(s))

	inputPayload := "{}"

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/test", strings.NewReader(inputPayload))
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(w, req)

	td.Cmp(t, w.Code, http.StatusUnprocessableEntity)
	td.Cmp(t, w.Body.String(), `{"message":"property \"name\" is missing"}`)
}
