package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"time"

	"snipqurl/internal/metrics"
	"snipqurl/internal/model"
	"snipqurl/internal/repository"

	"github.com/skip2/go-qrcode"
)

var ErrExpired = errors.New("url has expired")
var ErrAliasTaken = errors.New("alias already taken")

type URLService interface {
	Shorten(originalURL string, alias string, expiresAt *time.Time) (*model.URL, error)
	GetOriginalURL(code string) (*model.URL, error)
	GenerateQR(originalURL string) ([]byte, error)
}

type urlService struct {
	repo repository.URLRepository
}

func New(repo repository.URLRepository) URLService {
	return &urlService{repo: repo}
}

func (s *urlService) Shorten(originalURL string, alias string, expiresAt *time.Time) (*model.URL, error) {
	normalizedURL, err := normalizeAndValidateURL(originalURL)
	if err != nil {
		return nil, err
	}

	var code string
	if alias != "" {
		code = alias
		_, err := s.repo.FindByShortCode(code)
		if err == nil {
			return nil, ErrAliasTaken
		}
		if !errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}
	} else {
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
	}

	newURL := &model.URL{
		OriginalURL: normalizedURL,
		ShortCode:   code,
		ExpiresAt:   expiresAt,
	}
	err = s.repo.Save(newURL)
	if err != nil {
		return nil, fmt.Errorf("failed to save url: %w", err)
	}

	metrics.URLsShortenedTotal.Inc()

	return newURL, nil
}

func (s *urlService) GetOriginalURL(code string) (*model.URL, error) {
	u, err := s.repo.FindByShortCode(code)
	if err != nil {
		return nil, err
	}

	if u.ExpiresAt != nil && time.Now().After(*u.ExpiresAt) {
		return nil, ErrExpired
	}

	err = s.repo.IncrementClick(code)
	if err != nil {
		log.Printf("fail increment click")
	}

	metrics.RedirectsTotal.Inc()

	return u, nil
}

func (s *urlService) GenerateQR(originalURL string) ([]byte, error) {
	normalizedURL, err := normalizeAndValidateURL(originalURL)
	if err != nil {
		return nil, err
	}

	png, err := qrcode.Encode(normalizedURL, qrcode.Medium, 1024)
	if err != nil {
		return nil, err
	}

	metrics.QRGeneratedTotal.Inc()

	return png, nil
}
