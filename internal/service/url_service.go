package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"

	"snipqurl/internal/model"
	"snipqurl/internal/repository"

	"github.com/skip2/go-qrcode"
)

type URLService interface {
	Shorten(originalURL string) (*model.URL, error)
	GetOriginalURL(code string) (*model.URL, error)
	GenerateQR(originalURL string) ([]byte, error)
}

type urlService struct {
	repo repository.URLRepository
}

func New(repo repository.URLRepository) URLService {
	return &urlService{repo: repo}
}

func (s *urlService) Shorten(originalURL string) (*model.URL, error) {
	err := validateURL(originalURL)
	if err != nil {
		return nil, err
	}

	var code string
	for {
		code = rand.Text()[:8]
		_, err := s.repo.FindByShortCode(code)
		if errors.Is(err, repository.ErrNotFound) {
			break
		}
		if err == nil {
			continue
		}
		return nil, err
	}

	newURL := &model.URL{
		OriginalURL: originalURL,
		ShortCode:   code,
	}
	err = s.repo.Save(newURL)
	if err != nil {
		return nil, fmt.Errorf("failed to save url: %w", err)
	}

	return newURL, nil
}

func (s *urlService) GetOriginalURL(code string) (*model.URL, error) {
	u, err := s.repo.FindByShortCode(code)
	if err != nil {
		return nil, err
	}

	err = s.repo.IncrementClick(code)
	if err != nil {
		log.Printf("fail increment click")
	}

	return u, nil
}

func (s *urlService) GenerateQR(originalURL string) ([]byte, error) {
	err := validateURL(originalURL)
	if err != nil {
		return nil, err
	}

	png, err := qrcode.Encode(originalURL, qrcode.Medium, 1024)
	if err != nil {
		return nil, err
	}

	return png, nil
}
