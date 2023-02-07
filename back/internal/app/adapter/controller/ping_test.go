package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/maxatome/go-testdeep/td"
)

func TestNewPingController(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(w)

	NewPingController(engine)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/ping", http.NoBody)
	engine.ServeHTTP(w, req)

	td.Cmp(t, w.Code, http.StatusOK)
	td.Cmp(t, w.Body.String(), `"pong"`)
}
