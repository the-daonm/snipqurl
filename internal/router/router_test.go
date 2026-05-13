package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"snipqurl/internal/handler"
	"snipqurl/internal/router"
	"snipqurl/internal/service"

	"github.com/gin-gonic/gin"
)

type mockURLService struct {
	service.URLService
}

func TestRouter_Metrics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := handler.New(&mockURLService{})
	r := router.SetUp(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/metrics", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if w.Body.Len() == 0 {
		t.Error("expected non-empty body for /metrics")
	}
}
