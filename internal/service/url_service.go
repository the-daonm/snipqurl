package service

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"log"
	"net/url"

	"snipqurl/internal/model"
	"snipqurl/internal/repository"
)

type URLService interface {
	Shorten(originalURL string) (*model.URL, error)
	GetOrginalURL(code string) (*model.URL, error)
}

type urlService struct {
	repo repository.URLRepository
}

func New(repo *repository.URLRepository) URLService {
	return urlService{repo: repo}
}

func (s *urlService) Shorten(originalURL string) (*model.URL, error) {
	u, err := url.ParseRequestURI(originalURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return nil, fmt.Errorf("invalid url: %w", err)
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
		if err != nil {
			return nil, err
		}
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
