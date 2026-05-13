package service_test

import (
	"testing"

	"snipqurl/internal/model"
	"snipqurl/internal/repository"
	"snipqurl/internal/service"
)

type mockURLRepository struct {
	saveFunc            func(url *model.URL) error
	findByShortCodeFunc func(code string) (*model.URL, error)
	incrementClickFunc  func(code string) error
	deleteExpiredFunc   func() (int64, error)
}

func (m *mockURLRepository) Save(url *model.URL) error {
	return m.saveFunc(url)
}

func (m *mockURLRepository) FindByShortCode(code string) (*model.URL, error) {
	return m.findByShortCodeFunc(code)
}

func (m *mockURLRepository) IncrementClick(code string) error {
	return m.incrementClickFunc(code)
}

func (m *mockURLRepository) DeleteExpired() (int64, error) {
	if m.deleteExpiredFunc != nil {
		return m.deleteExpiredFunc()
	}
	return 0, nil
}

func TestURLService_Shorten(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		alias       string
		mockFindBy  func(code string) (*model.URL, error)
		mockSave    func(url *model.URL) error
		expectError bool
	}{
		{
			name: "success",
			url:  "https://google.com",
			mockFindBy: func(code string) (*model.URL, error) {
				return nil, repository.ErrNotFound
			},
			mockSave: func(url *model.URL) error {
				return nil
			},
			expectError: false,
		},
		{
			name:        "invalid url",
			url:         "",
			expectError: true,
		},
		{
			name:  "alias taken",
			url:   "https://google.com",
			alias: "taken",
			mockFindBy: func(code string) (*model.URL, error) {
				return &model.URL{ShortCode: "taken"}, nil
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockURLRepository{
				findByShortCodeFunc: tt.mockFindBy,
				saveFunc:            tt.mockSave,
			}
			svc := service.New(mockRepo)

			_, err := svc.Shorten(tt.url, tt.alias, nil)
			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

func TestURLService_GetOriginalURL(t *testing.T) {
	tests := []struct {
		name          string
		code          string
		mockFindBy    func(code string) (*model.URL, error)
		mockIncrement func(code string) error
		expectError   bool
	}{
		{
			name: "success",
			code: "abcdef12",
			mockFindBy: func(code string) (*model.URL, error) {
				return &model.URL{OriginalURL: "https://google.com"}, nil
			},
			mockIncrement: func(code string) error {
				return nil
			},
			expectError: false,
		},
		{
			name: "not found",
			code: "unknown",
			mockFindBy: func(code string) (*model.URL, error) {
				return nil, repository.ErrNotFound
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockURLRepository{
				findByShortCodeFunc: tt.mockFindBy,
				incrementClickFunc:  tt.mockIncrement,
			}
			svc := service.New(mockRepo)

			_, err := svc.GetOriginalURL(tt.code)
			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}
