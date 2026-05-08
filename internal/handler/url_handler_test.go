package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"snipqurl/internal/handler"
	"snipqurl/internal/model"
	"snipqurl/internal/service"

	"github.com/gin-gonic/gin"
)

type mockURLService struct {
	shortenFunc        func(url string) (*model.URL, error)
	getOriginalURLFunc func(code string) (*model.URL, error)
	generateQRFunc     func(url string) ([]byte, error)
}

func (m *mockURLService) Shorten(originalURL string) (*model.URL, error) {
	return m.shortenFunc(originalURL)
}

func (m *mockURLService) GetOriginalURL(code string) (*model.URL, error) {
	return m.getOriginalURLFunc(code)
}

func (m *mockURLService) GenerateQR(originalURL string) ([]byte, error) {
	return m.generateQRFunc(originalURL)
}

func TestURLHandler_Shorten(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		reqBody      map[string]string
		mockShorten  func(url string) (*model.URL, error)
		domainEnv    string
		expectedCode int
	}{
		{
			name:    "success",
			reqBody: map[string]string{"url": "https://google.com"},
			mockShorten: func(url string) (*model.URL, error) {
				return &model.URL{ShortCode: "abcdef12"}, nil
			},
			domainEnv:    "example.com",
			expectedCode: http.StatusOK,
		},
		{
			name:    "invalid url",
			reqBody: map[string]string{"url": "not-a-url"},
			mockShorten: func(url string) (*model.URL, error) {
				return nil, service.ErrInvalidURL
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:    "internal error",
			reqBody: map[string]string{"url": "https://google.com"},
			mockShorten: func(url string) (*model.URL, error) {
				return nil, errors.New("db error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("DOMAIN", tt.domainEnv)
			defer os.Unsetenv("DOMAIN")

			mockSvc := &mockURLService{shortenFunc: tt.mockShorten}
			h := handler.New(mockSvc)

			r := gin.New()
			r.POST("/shorten", h.Shorten)

			body, _ := json.Marshal(tt.reqBody)
			req, _ := http.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("expected status code %d, got %d", tt.expectedCode, w.Code)
			}
		})
	}
}

func TestURLHandler_Redirect(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		code         string
		mockGetURL   func(code string) (*model.URL, error)
		expectedCode int
	}{
		{
			name: "success",
			code: "abcdef12",
			mockGetURL: func(code string) (*model.URL, error) {
				return &model.URL{OriginalURL: "https://google.com"}, nil
			},
			expectedCode: http.StatusMovedPermanently,
		},
		{
			name: "not found",
			code: "unknown",
			mockGetURL: func(code string) (*model.URL, error) {
				return nil, errors.New("not found")
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockURLService{getOriginalURLFunc: tt.mockGetURL}
			h := handler.New(mockSvc)

			r := gin.New()
			r.GET("/:code", h.Redirect)

			req, _ := http.NewRequest(http.MethodGet, "/"+tt.code, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("expected status code %d, got %d", tt.expectedCode, w.Code)
			}
		})
	}
}
