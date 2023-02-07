package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/maxatome/go-testdeep/td"
)

func TestCorsMiddlewareRequestOptionsReturnsHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(w)

	engine.Use(CORSMiddleware())

	req, _ := http.NewRequestWithContext(ctx, http.MethodOptions, "/", http.NoBody)
	engine.ServeHTTP(w, req)

	td.Cmp(t, w.Code, http.StatusNoContent)
	td.Cmp(t, w.Header().Get("Access-Control-Allow-Origin"), "*")
	td.Cmp(t, w.Header().Get("Access-Control-Allow-Credentials"), "true")
	td.Cmp(t, w.Header().Get("Access-Control-Allow-Headers"), "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Access-Token, X-Refresh-Token")
	td.Cmp(t, w.Header().Get("Access-Control-Allow-Methods"), "POST,HEAD,PATCH, OPTIONS, GET, PUT, DELETE")
}
